package sitemap

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type StaticURL struct {
	URL         string `json:"url"`
	ReleaseDate string `json:"releaseDate"`
	HasAltLang  bool   `json:"hasAltLang"`
}

func LoadStaticSitemap(oldSitemapName, staticSitemapName, dpOnsURLHostName, dpOnsURLHostNameAlt, altLang string, store FileStore) error {
	b, err := GetStaticSitemap(staticSitemapName)
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
