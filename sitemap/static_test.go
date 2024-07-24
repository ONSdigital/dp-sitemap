package sitemap

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadStaticSitemap(t *testing.T) {
	Convey("given we have static sitemap file", t, func() {
		oldSitemapName := "test_sitemap_en"
		staticSitemapName := "sitemap_en.json"
		Convey("when loading english static sitemap", func() {
			store := LocalStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(cfg, oldSitemapName, staticSitemapName, cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the file should exist", func() {
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldBeNil)
			})
			Convey("And when we delete it, it should not exist", func() {
				err = os.Remove(oldSitemapName)
				So(err, ShouldBeNil)
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("given we have static sitemap file sitemap_cy.json", t, func() {
		oldSitemapName := "test_sitemap_cy"
		staticSitemapName := "sitemap_cy.json"
		Convey("when loading welsh static sitemap", func() {
			store := LocalStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(cfg, oldSitemapName, staticSitemapName, cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the file should exist", func() {
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldBeNil)
			})
			Convey("And when we delete it, it should not exist", func() {
				err = os.Remove(oldSitemapName)
				So(err, ShouldBeNil)
				_, err = os.Stat(oldSitemapName)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func expectedURLSetEnglish() *UrlsetReader {
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

func expectedURLSetWelsh() *UrlsetReader {
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
