package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/filesaver.go -pkg mock . FileSaver
//go:generate moq -out mock/fetcher.go -pkg mock . Fetcher

type FileSaver interface {
	SaveFile(body io.Reader) error
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) (string, error)
}

type Generator struct {
	fetcher Fetcher
	saver   FileSaver
}

func NewGenerator(fetcher Fetcher, saver FileSaver) *Generator {
	return &Generator{
		fetcher: fetcher,
		saver:   saver,
	}
}

func (g *Generator) MakeFullSitemap(ctx context.Context) error {
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

	err = g.saver.SaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to save sitemap file: %w", err)
	}
	return nil
}
