package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/filesaver.go -pkg mock . FileSaver
//go:generate moq -out mock/fetcher.go -pkg mock . Fetcher

type FileSaver interface {
	SaveFile(lang config.Language, body io.Reader) error
	UploadFiles(paths []string) error
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) ([]string, error)
	HasWelshContent(ctx context.Context, path string) bool
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
	fileNames, err := g.fetcher.GetFullSitemap(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer func() {
		for _, fl := range fileNames {
			err = os.Remove(fl)
			if err != nil {
				log.Error(ctx, "failed to remove temporary sitemap file "+fl, err)
				return
			}
			log.Info(ctx, "removed temporary sitemap file "+fl)
		}
	}()

	for _, fl := range fileNames {
		lang := config.English
		if strings.Contains(fl, tempSitemapFileCy) {
			lang = config.Welsh
		}
		file, err := os.Open(fl)
		if err != nil {
			return fmt.Errorf("failed to open sitemap: %w", err)
		}

		err = g.saver.SaveFile(lang, file)
		file.Close()
		if err != nil {
			return fmt.Errorf("failed to save sitemap file: %w", err)
		}
	}
	return nil
}
