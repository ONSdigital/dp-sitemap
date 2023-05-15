package event

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"os"
	"time"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
)

type ContentPublishedHandler struct {
	FileStore sitemap.FileStore
}

// Handle takes a single event.
func (h *ContentPublishedHandler) Handle(ctx context.Context, cfg *config.Config, event *ContentPublished) (err error) {
	logData := log.Data{
		"eventContentPublished": event,
	}
	log.Info(ctx, "event handler called with event", logData)

	var adder sitemap.DefaultAdder
	currentSitemapName := "test_sitemap_eng"
	var tmpSitemapName string

	generateTempSitemap := func() {

		y, m, d := time.Now().Date()
		date := fmt.Sprintf("%d-%d-%d", y, m, d)
		var url = &sitemap.URL{Loc: event.URI, Lastmod: date}
		currentSitemap, err := h.FileStore.GetFile(currentSitemapName)
		if err != nil {
			fmt.Println("Error opening current sitemap", err)
			os.Exit(1)
		}
		defer currentSitemap.Close()

		tmpSitemapName, _, err = adder.Add(currentSitemap, url)
		if err != nil {
			fmt.Println("Error creating temp sitemap file", err)
			os.Exit(1)
		}
		fmt.Println("temp sitemap file name:", tmpSitemapName)
	}

	generateTempSitemap()

	currentSitemap, err := h.FileStore.CreateFile(currentSitemapName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer currentSitemap.Close()

	tmpSitemap, err := os.Open(tmpSitemapName)
	if err != nil {
		fmt.Println("open temp sitemap", err)
		os.Exit(1)
	}
	defer tmpSitemap.Close()
	h.FileStore.CopyFile(tmpSitemap, currentSitemap)

	return nil
}
