package sitemap

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
)

type StaticURL struct {
	URL         string `json:"url"`
	ReleaseDate string `json:"releaseDate"`
	HasAltLang  bool   `json:"hasAltLang"`
}

func LoadStaticSitemap(ctx context.Context, oldSitemapName, staticSitemapName, dpOnsURLHostName, dpOnsURLHostNameAlt, altLang string, store FileStore) error {
	efs := assets.NewFromEmbeddedFilesystem()

	b, err := efs.Get(ctx, assets.Sitemap, staticSitemapName)
	if err != nil {
		panic("can't find file " + staticSitemapName)
	}

	var content []StaticURL

	err = json.Unmarshal(b, &content)
	if err != nil {
		return fmt.Errorf("unable to read json: %w", err)
	}

	// move old sitemap urls to new sitemap
	sitemapWriter := Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}

	// range through static content
	for _, item := range content {
		var newURL URL
		newURL.Loc = dpOnsURLHostName + item.URL
		newURL.Lastmod = item.ReleaseDate
		newURL.Alternate = &AlternateURL{}
		if item.HasAltLang {
			newURL.Alternate.Rel = "alternate"
			newURL.Alternate.Link = dpOnsURLHostNameAlt + item.URL
			newURL.Alternate.Lang = altLang
		}
		sitemapWriter.URL = append(sitemapWriter.URL, newURL)
	}

	marshaledContent, err := xml.MarshalIndent(sitemapWriter, "", "  ")
	if err != nil {
		return err
	}
	header := []byte(xml.Header)
	header = append(header, marshaledContent...)
	reader := bytes.NewReader(header)
	err = store.SaveFile(oldSitemapName, reader)
	if err != nil {
		return err
	}
	return nil
}

type StaticSitemapConfig struct {
	Lang                  config.Language // the main language
	HostName              string          // the main language host name url
	StaticSitemapFileName string          // the static sitemap file name
	SitemapFileName       string          // the sitemap file name
	AlternateHostName     string          // the alternate language host name url
	AlternateLang         config.Language // the alternate language
}

func GetConfigAsArray(cfg *config.Config) []StaticSitemapConfig {
	return []StaticSitemapConfig{
		{
			Lang:                  config.English,
			HostName:              cfg.DpOnsURLHostNameEn,
			StaticSitemapFileName: "sitemap_en.json",
			SitemapFileName:       "test_sitemap_en",
			AlternateHostName:     cfg.DpOnsURLHostNameCy,
			AlternateLang:         config.Welsh,
		},
		{
			Lang:                  config.Welsh,
			HostName:              cfg.DpOnsURLHostNameCy,
			StaticSitemapFileName: "sitemap_cy.json",
			SitemapFileName:       "test_sitemap_cy",
			AlternateHostName:     cfg.DpOnsURLHostNameEn,
			AlternateLang:         config.English,
		},
	}
}
