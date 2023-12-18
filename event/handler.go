package event

import (
	"context"

	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
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
	pageInfo, err := h.fetcher.GetPageInfo(ctx, event.URI)
	if err != nil {
		log.Error(ctx, "error getting page information", err, log.Data{"uri": event.URI})
		return err
	}

	if pageInfo.URLs[config.English] != nil {
		err = h.createSiteMap(ctx, config.English, cfg.SitemapLocalFile[config.English], pageInfo)
		if err != nil {
			return err
		}
	}

	if pageInfo.URLs[config.Welsh] != nil {
		err = h.createSiteMap(ctx, config.Welsh, cfg.SitemapLocalFile[config.Welsh], pageInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *ContentPublishedHandler) createSiteMap(ctx context.Context, lang config.Language, sitemapName string, pageInfo *sitemap.PageInfo) error {
	currentSitemapName := sitemapName
	var tmpSitemapName string

	tmpSitemapName, err := h.generateTempSitemap(ctx, currentSitemapName, lang, pageInfo)
	if err != nil {
		return err
	}

	currentSitemap, err := h.fileStore.CreateFile(currentSitemapName)
	if err != nil {
		log.Error(ctx, "error creating current sitemap", err)
		return err
	}
	defer currentSitemap.Close()

	tmpSitemap, err := h.fileStore.GetFile(tmpSitemapName)
	if err != nil {
		log.Error(ctx, "error opening temp sitemap", err)
		return err
	}
	defer tmpSitemap.Close()

	err = h.fileStore.CopyFile(tmpSitemap, currentSitemap)
	if err != nil {
		log.Error(ctx, "error copying file", err)
		return err
	}

	err = h.fileStore.DeleteFile(tmpSitemapName)
	if err != nil {
		log.Error(ctx, "error deleting temp sitemap", err)
	}

	return nil
}

func (h *ContentPublishedHandler) generateTempSitemap(ctx context.Context, currentSitemapName string, lang config.Language, pageInfo *sitemap.PageInfo) (string, error) {
	currentSitemap, err := h.fileStore.GetFile(currentSitemapName)
	if err != nil {
		log.Error(ctx, "error opening current sitemap", err)
		return "", err
	}
	defer currentSitemap.Close()

	var adder sitemap.DefaultAdder
	tmpSitemapName, _, err := adder.Add(currentSitemap, pageInfo.URLs[lang])
	if err != nil {
		log.Error(ctx, "error creating temp sitemap file", err)
		return "", err
	}

	return tmpSitemapName, nil
}
