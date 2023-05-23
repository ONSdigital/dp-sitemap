package event

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"os"
)

type ContentPublishedHandler struct {
	fileStore sitemap.FileStore
	fetcher   *sitemap.ElasticFetcher
}

func NewContentPublishedHandler(store sitemap.FileStore, fetcher *sitemap.ElasticFetcher) *ContentPublishedHandler {
	return &ContentPublishedHandler{
		fileStore: store,
		fetcher:   fetcher,
	}
}

// Handle takes a single event.
func (h *ContentPublishedHandler) Handle(ctx context.Context, cfg *config.Config, event *ContentPublished) (err error) {
	logData := log.Data{
		"eventContentPublished": event,
	}
	log.Info(ctx, "event handler called with event", logData)

	currentSitemapName := "test_sitemap_eng"
	var tmpSitemapName string

	tmpSitemapName = h.generateTempSitemap(ctx, currentSitemapName, event)

	currentSitemap, err := h.fileStore.CreateFile(currentSitemapName)
	if err != nil {
		fmt.Println("Error creating current sitemap", err)
		os.Exit(1)
	}
	defer currentSitemap.Close()

	tmpSitemap, err := h.fileStore.GetFile(tmpSitemapName)
	if err != nil {
		fmt.Println("Error opening temp sitemap", err)
		os.Exit(1)
	}
	defer tmpSitemap.Close()
	h.fileStore.CopyFile(tmpSitemap, currentSitemap)

	return nil
}

func (h *ContentPublishedHandler) generateTempSitemap(ctx context.Context, currentSitemapName string, event *ContentPublished) string {
	description, err := h.fetcher.GetZebedeeClient().GetPageDescription(ctx, "", "", "eng", event.URI)
	if err != nil {
		fmt.Println("Error getting page description", err)
		os.Exit(1)
	}

	currentSitemap, err := h.fileStore.GetFile(currentSitemapName)
	if err != nil {
		fmt.Println("Error opening current sitemap", err)
		os.Exit(1)
	}
	defer currentSitemap.Close()

	urlEn, _ := h.fetcher.URLVersions(ctx, event.URI, description.Description.ReleaseDate)

	var adder sitemap.DefaultAdder
	tmpSitemapName, _, err := adder.Add(currentSitemap, &urlEn)
	if err != nil {
		fmt.Println("Error creating temp sitemap file", err)
		os.Exit(1)
	}
	return tmpSitemapName
}
