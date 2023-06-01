package event

import (
	"context"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"time"
)

type ContentPublishedHandler struct {
	fileStore     sitemap.FileStore
	zebedeeClient clients.ZebedeeClient
	config        *config.Config
	fetcher       sitemap.Fetcher
}

func NewContentPublishedHandler(store sitemap.FileStore, client clients.ZebedeeClient, cfg *config.Config, fetcher sitemap.Fetcher) *ContentPublishedHandler {
	return &ContentPublishedHandler{
		fileStore:     store,
		zebedeeClient: client,
		config:        cfg,
		fetcher:       fetcher,
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
		err = h.createSiteMap(ctx, event, "eng", "test_sitemap_en")
		if err != nil {
			return err
		}
	}

	// _, err = h.fetcher.GetZebedeeClient().GetPageDescription(ctx, "", "", "cy", event.URI)
	if urlCy != nil {
		err = h.createSiteMap(ctx, event, "cy", "test_sitemap_cy")
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *ContentPublishedHandler) createSiteMap(ctx context.Context, event *ContentPublished, lang string, sitemapName string) error {

	currentSitemapName := sitemapName
	var tmpSitemapName string

	tmpSitemapName, err := h.generateTempSitemap(ctx, currentSitemapName, event, lang)
	if err != nil {
		return err
	}

	currentSitemap, err := h.fileStore.CreateFile(currentSitemapName)
	if err != nil {
		log.Error(ctx, "Error creating current sitemap", err)
		return err
	}
	defer currentSitemap.Close()

	tmpSitemap, err := h.fileStore.GetFile(tmpSitemapName)
	if err != nil {
		log.Error(ctx, "Error opening temp sitemap", err)
		return err
	}
	defer tmpSitemap.Close()

	h.fileStore.CopyFile(tmpSitemap, currentSitemap)

	return nil
}

func (h *ContentPublishedHandler) generateTempSitemap(ctx context.Context, currentSitemapName string, event *ContentPublished, lang string) (string, error) {
	description, err := h.zebedeeClient.GetPageDescription(ctx, "", "", lang, event.URI)
	if err != nil {
		log.Error(ctx, "Error getting page description", err)
		return "", err
	}

	currentSitemap, err := h.fileStore.GetFile(currentSitemapName)
	if err != nil {
		log.Error(ctx, "Error opening current sitemap", err)
		return "", err
	}
	defer currentSitemap.Close()

	releaseDate, err := time.Parse(time.RFC3339, description.Description.ReleaseDate)
	if err != nil {
		log.Error(ctx, "Error parsing the release date", err)
		return "", err
	}
	url := h.fetcher.URLVersion(ctx, event.URI, releaseDate.Format("2006-01-02"), lang)

	var adder sitemap.DefaultAdder
	tmpSitemapName, _, err := adder.Add(currentSitemap, url)
	if err != nil {
		log.Error(ctx, "Error creating temp sitemap file", err)
		return "", err
	}

	return tmpSitemapName, nil
}
