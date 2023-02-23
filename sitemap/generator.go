package sitemap

import (
	"context"
	"fmt"
	"os"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out mock/s3uploader.go -pkg mock . S3Uploader
//go:generate moq -out mock/fetcher.go -pkg mock . Fetcher

type S3Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type Fetcher interface {
	GetFullSitemap(ctx context.Context) (string, error)
}

type Generator struct {
	fetcher  Fetcher
	uploader S3Uploader
	cfg      *config.S3Config
}

func NewGenerator(fetcher Fetcher, uploader S3Uploader, cfg *config.S3Config) *Generator {
	return &Generator{
		fetcher:  fetcher,
		uploader: uploader,
		cfg:      cfg,
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

	_, err = g.uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: &g.cfg.UploadBucketName,
		Key:    &g.cfg.SitemapFileKey,
	})
	if err != nil {
		return fmt.Errorf("failed to upload sitemap: %w", err)
	}
	return nil
}
