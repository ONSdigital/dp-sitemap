package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/filestore.go -pkg mock . FileStore
//go:generate moq -out mock/fetcher.go -pkg mock . Fetcher
//go:generate moq -out mock/adder.go -pkg mock . Adder

type FileStore interface {
	SaveFile(name string, body io.Reader) error
	GetFile(name string) (body io.ReadCloser, err error)
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) (string, error)
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
	defer currentSitemap.Close()

	fileName, err := g.adder.Add(currentSitemap, url)
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
	defer file.Close()

	err = g.store.SaveFile(name, file)
	if err != nil {
		return fmt.Errorf("failed to save incremental sitemap file: %w", err)
	}
	return nil
}

func (g *Generator) MakeFullSitemap(ctx context.Context, name string) error {
	fileName, err := g.fetcher.GetFullSitemap(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch sitemap: %w", err)
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
		return fmt.Errorf("failed to open sitemap: %w", err)
	}
	defer file.Close()

	err = g.store.SaveFile(name, file)
	if err != nil {
		return fmt.Errorf("failed to save sitemap file: %w", err)
	}
	return nil
}
