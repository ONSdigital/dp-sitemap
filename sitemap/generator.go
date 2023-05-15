package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/filestore.go -pkg mock . FileStore
//go:generate moq -out mock/fetcher.go -pkg mock . Fetcher
//go:generate moq -out mock/adder.go -pkg mock . Adder

type Files map[config.Language]string

type FileStore interface {
	SaveFile(name string, body io.Reader) error
	GetFile(name string) (body io.ReadCloser, err error)
	CopyFile(src io.Reader, dest io.Writer) error
	CreateFile(name string) (io.ReadWriteCloser, error)
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) (Files, error)
	HasWelshContent(ctx context.Context, path string) bool
	URLVersions(ctx context.Context, path string, lastmod string) (en URL, cy *URL)
}
type Adder interface {
	Add(oldSitemap io.Reader, url *URL) (file string, size int, err error)
}

type Generator struct {
	fetcher               Fetcher
	adder                 Adder
	store                 FileStore
	publishingSitemapMx   sync.Mutex
	maxSize               int
	maxSizeCallback       func()
	fullSitemapFiles      Files
	publishingSitemapFile string
}
type GeneratorOptions func(*Generator) *Generator

func NewGenerator(opts ...GeneratorOptions) *Generator {
	g := &Generator{
		fullSitemapFiles:      Files{config.English: "sitemap.xml"},
		publishingSitemapFile: "publishing-sitemap.xml",
	}
	for _, opt := range opts {
		g = opt(g)
	}
	return g
}

func WithFetcher(f Fetcher) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.fetcher = f
		return g
	}
}

func WithAdder(a Adder) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.adder = a
		return g
	}
}

func WithFileStore(s FileStore) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.store = s
		return g
	}
}

func WithFullSitemapFiles(f Files) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.fullSitemapFiles = f
		return g
	}
}

func WithPublishingSitemapFile(f string) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.publishingSitemapFile = f
		return g
	}
}

func WithPublishingSitemapMaxSize(size int, callback func()) GeneratorOptions {
	return func(g *Generator) *Generator {
		g.maxSize = size
		g.maxSizeCallback = callback
		return g
	}
}

func (g *Generator) MakePublishingSitemap(ctx context.Context, url URL) error {
	g.publishingSitemapMx.Lock()
	defer g.publishingSitemapMx.Unlock()

	currentSitemap, err := g.store.GetFile(g.publishingSitemapFile)
	if err != nil {
		return fmt.Errorf("failed to get current sitemap: %w", err)
	}
	defer func() {
		closeErr := currentSitemap.Close()
		if closeErr != nil {
			log.Error(ctx, "failed to close current sitemap file", closeErr)
		}
	}()

	urlEn, _ := g.fetcher.URLVersions(
		ctx,
		url.Loc,
		url.Lastmod,
	)

	size, err := g.AppendURL(ctx, currentSitemap, &urlEn, g.publishingSitemapFile)
	if err != nil {
		return err
	}

	if g.maxSize > 0 && size > g.maxSize {
		go g.maxSizeCallback()
	}

	return nil
}

func (g *Generator) TruncatePublishingSitemap(ctx context.Context) error {
	g.publishingSitemapMx.Lock()
	defer g.publishingSitemapMx.Unlock()

	_, err := g.AppendURL(ctx, io.NopCloser(strings.NewReader("")), nil, g.publishingSitemapFile)
	return err
}

func (g *Generator) AppendURL(ctx context.Context, sitemap io.ReadCloser, url *URL, destination string) (int, error) {
	fileName, size, err := g.adder.Add(sitemap, url)
	if err != nil {
		return 0, fmt.Errorf("failed to add to sitemap: %w", err)
	}
	defer func() {
		err = os.Remove(fileName)
		if err != nil {
			log.Error(ctx, "failed to remove temporary sitemap file "+fileName, err)
			return
		}
		log.Info(ctx, "removed temporary sitemap file "+fileName)
	}()

	file, err := os.Open(fileName)
	if err != nil {
		return 0, fmt.Errorf("failed to open publishing sitemap: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Error(ctx, "failed to close publishing sitemap file", closeErr)
		}
	}()

	err = g.store.SaveFile(destination, file)
	if err != nil {
		return 0, fmt.Errorf("failed to save publishing sitemap file: %w", err)
	}
	return size, nil
}

func (g *Generator) MakeFullSitemap(ctx context.Context) error {
	// first truncate the publishing sitemap as all URLs that are
	// currently there will be automatically included in the full sitemap
	err := g.TruncatePublishingSitemap(ctx)
	if err != nil {
		return fmt.Errorf("failed to truncate publishing sitemap: %w", err)
	}

	sitemaps, err := g.fetcher.GetFullSitemap(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer func() {
		for _, fl := range sitemaps {
			err = os.Remove(fl)
			if err != nil {
				log.Error(ctx, "failed to remove temporary sitemap file "+fl, err)
				return
			}
			log.Info(ctx, "removed temporary sitemap file "+fl)
		}
	}()

	for lang, fl := range sitemaps {
		file, err := os.Open(fl)
		if err != nil {
			return fmt.Errorf("failed to open sitemap: %w", err)
		}

		err = g.store.SaveFile(g.fullSitemapFiles[lang], file)
		file.Close()
		if err != nil {
			return fmt.Errorf("failed to save sitemap file: %w", err)
		}
	}
	return nil
}
