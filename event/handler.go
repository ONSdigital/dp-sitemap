package event

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"os"
	"time"
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
	urlEn, urlCy := h.fetcher.URLVersions(ctx, event.URI, "")
	if urlEn != nil {
		err1 := h.createSiteMap(ctx, event, "eng", "test_sitemap_en")
		if err1 != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// _, err = h.fetcher.GetZebedeeClient().GetPageDescription(ctx, "", "", "cy", event.URI)
	if urlCy != nil {
		err1 := h.createSiteMap(ctx, event, "cy", "test_sitemap_cy")
		if err1 != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return nil
}

func (h *ContentPublishedHandler) createSiteMap(ctx context.Context, event *ContentPublished, lang string, sitemapName string) error {

	currentSitemapName := sitemapName
	var tmpSitemapName string

	tmpSitemapName = h.generateTempSitemap(ctx, currentSitemapName, event, lang)

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

func (h *ContentPublishedHandler) generateTempSitemap(ctx context.Context, currentSitemapName string, event *ContentPublished, lang string) string {
	description, err := h.fetcher.GetZebedeeClient().GetPageDescription(ctx, "", "", lang, event.URI)
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

	releaseDate, err := time.Parse(time.RFC3339, description.Description.ReleaseDate)
	if err != nil {
		fmt.Println("Error parsing the release date", err)
		os.Exit(1)
	}
	url := h.fetcher.URLVersion(ctx, event.URI, releaseDate.Format("2006-01-02"), lang)

	var adder sitemap.DefaultAdder
	tmpSitemapName, _, err := adder.Add(currentSitemap, url)
	if err != nil {
		fmt.Println("Error creating temp sitemap file", err)
		os.Exit(1)
	}
	return tmpSitemapName
}
