package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"

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
	SaveFiles(paths []string) error
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) (Files, error)
	HasWelshContent(ctx context.Context, path string) bool
	URLVersions(ctx context.Context, path string, lastmod string) (en URL, cy *URL)
}
type Adder interface {
	Add(oldSitemap io.Reader, url URL) (string, error)
}

type Generator struct {
	fetcher Fetcher
	adder   Adder
	store   FileStore
}

func NewGenerator(fetcher Fetcher, adder Adder, store FileStore) *Generator {
	return &Generator{
		fetcher: fetcher,
		adder:   adder,
		store:   store,
	}
}

func (g *Generator) MakeIncrementalSitemap(ctx context.Context, name string, url URL) error {
	currentSitemap, err := g.store.GetFile(name)
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

	return g.AppendURL(ctx, currentSitemap, urlEn, name)
}

func (g *Generator) AppendURL(ctx context.Context, sitemap io.ReadCloser, url URL, destination string) error {
	fileName, err := g.adder.Add(sitemap, url)
	if err != nil {
		return fmt.Errorf("failed to add to sitemap: %w", err)
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
		return fmt.Errorf("failed to open incremental sitemap: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Error(ctx, "failed to close incremental sitemap file", closeErr)
		}
	}()

	err = g.store.SaveFile(destination, file)
	if err != nil {
		return fmt.Errorf("failed to save incremental sitemap file: %w", err)
	}
	return nil
}

func (g *Generator) MakeFullSitemap(ctx context.Context, fileNames Files) error {
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

		err = g.store.SaveFile(fileNames[lang], file)
		file.Close()
		if err != nil {
			return fmt.Errorf("failed to save sitemap file: %w", err)
		}
	}
	return nil
}
