package event

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"net/url"
	"os"
)

type ContentPublishedHandler struct {
	fileStore     sitemap.FileStore
	zebedeeClient clients.ZebedeeClient
	config        *config.Config
}

func NewContentPublishedHandler(store sitemap.FileStore, client clients.ZebedeeClient, cfg *config.Config) *ContentPublishedHandler {
	return &ContentPublishedHandler{
		fileStore:     store,
		zebedeeClient: client,
		config:        cfg,
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
		fmt.Println(err)
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
	description, err := h.zebedeeClient.GetPageDescription(ctx, "", "", "eng", event.URI)
	if err != nil {
		fmt.Println("Error getting page description", err)
		os.Exit(1)
	}

	enLoc, _ := url.JoinPath(h.config.DpOnsURLHostNameEn, event.URI)
	var url = &sitemap.URL{Loc: enLoc, Lastmod: description.Description.ReleaseDate}
	currentSitemap, err := h.fileStore.GetFile(currentSitemapName)
	if err != nil {
		fmt.Println("Error opening current sitemap", err)
		os.Exit(1)
	}
	defer currentSitemap.Close()

	var adder sitemap.DefaultAdder
	tmpSitemapName, _, err := adder.Add(currentSitemap, url)
	if err != nil {
		fmt.Println("Error creating temp sitemap file", err)
		os.Exit(1)
	}
	return tmpSitemapName
}
