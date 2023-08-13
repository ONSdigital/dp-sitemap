package sitemap

import (
	"encoding/xml"
	"github.com/ONSdigital/dp-sitemap/config"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadStaticSitemap(t *testing.T) {
	Convey("given we have static sitemap file", t, func() {
		oldSitemapName := "test_sitemap_en"
		staticSitemapName := "sitemap_en_test.json"
		Convey("when loading english static sitemap", func() {
			store := LocalStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(oldSitemapName, staticSitemapName, cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the file should exist", func() {
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("given we have static sitemap file sitemap_cy.json", t, func() {
		oldSitemapName := "test_sitemap_cy"
		staticSitemapName := "sitemap_cy_test.json"
		Convey("when loading welsh static sitemap", func() {
			store := LocalStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(oldSitemapName, staticSitemapName, cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the file should exist", func() {
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldBeNil)
			})
		})
	})
}

func expectedUrlSetEnglish() *UrlsetReader {
	return &UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL: []URLReader{
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "cy",
					Link:    "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				},
			},
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle2",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "cy",
					Link:    "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle2",
				},
			},
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "cy",
					Link:    "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				},
			},
		},
	}
}

func expectedUrlSetWelsh() *UrlsetReader {
	return &UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL: []URLReader{
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "en",
					Link:    "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				},
			},
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle2",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "en",
					Link:    "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle2",
				},
			},
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				Lastmod: "01-01-2023",
				Alternate: &AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "en",
					Link:    "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				},
			},
		},
	}
}
