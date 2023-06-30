package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	// Define test cases
	Convey("when all the args filled", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "1000",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the args are valid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})
	//sitemap_index is missing
	Convey("when some args are missing", t, func() {
		tetestdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&tetestdata)
			So(result, ShouldBeFalse)
		})
	})
}

func TestLoadStaticSitemap(t *testing.T) {
	expectedUrlset := expectedUrlSetEnglish()
	Convey("when loading english static sitemap", t, func() {
		store := mock.FileStoreMock{}
		buf := new(bytes.Buffer)
		store.SaveFileFunc = func(name string, body io.Reader) error {
			io.Copy(buf, body)
			return nil
		}
		cfg, _ := config.Get()
		err := loadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
		urlSet := &sitemap.UrlsetReader{}
		xml.Unmarshal(buf.Bytes(), urlSet)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("We should have the correct content loaded", func() {
			So(urlSet, ShouldResemble, expectedUrlset)
		})
	})

	Convey("when loading welsh static sitemap", t, func() {
		store := mock.FileStoreMock{}
		buf := new(bytes.Buffer)
		store.SaveFileFunc = func(name string, body io.Reader) error {
			io.Copy(buf, body)
			return nil
		}
		cfg, _ := config.Get()
		err := loadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
		urlSet := &sitemap.UrlsetReader{}
		xml.Unmarshal(buf.Bytes(), urlSet)
		expectedUrlset = expectedUrlSetWelsh()
		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("We should have the correct content loaded", func() {
			So(urlSet, ShouldResemble, expectedUrlset)
		})
	})
}

func TestLoadStaticSitemapV2(t *testing.T) {
	Convey("given we have static sitemap file sitemap_en.json", t, func() {

		Convey("when loading english static sitemap", func() {
			store := mock.FileStoreMock{}
			buf := new(bytes.Buffer)
			store.SaveFileFunc = func(name string, body io.Reader) error {
				io.Copy(buf, body)
				return nil
			}
			cfg, _ := config.Get()
			err := loadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
			urlSet := &sitemap.UrlsetReader{}
			xml.Unmarshal(buf.Bytes(), urlSet)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("We should have the correct content loaded", func() {
				So(urlSet, ShouldResemble, expectedUrlSetEnglish())
			})
		})
	})

	Convey("given we have static sitemap file sitemap_cy.json", t, func() {

		Convey("when loading welsh static sitemap", func() {
			store := mock.FileStoreMock{}
			buf := new(bytes.Buffer)
			store.SaveFileFunc = func(name string, body io.Reader) error {
				io.Copy(buf, body)
				return nil
			}
			cfg, _ := config.Get()
			err := loadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
			urlSet := &sitemap.UrlsetReader{}
			xml.Unmarshal(buf.Bytes(), urlSet)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("We should have the correct content loaded", func() {
				So(urlSet, ShouldResemble, expectedUrlSetWelsh())
			})
		})
	})
}

func expectedUrlSetEnglish() *sitemap.UrlsetReader {
	return &sitemap.UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL: []sitemap.URLReader{
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				Lastmod: "01-01-2023",
				Alternate: &sitemap.AlternateURLReader{
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
				Alternate: &sitemap.AlternateURLReader{
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
				Alternate: &sitemap.AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "cy",
					Link:    "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				},
			},
		},
	}
}

func expectedUrlSetWelsh() *sitemap.UrlsetReader {
	return &sitemap.UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL: []sitemap.URLReader{
			{
				XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"},
				Loc:     "https://cy.dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle1",
				Lastmod: "01-01-2023",
				Alternate: &sitemap.AlternateURLReader{
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
				Alternate: &sitemap.AlternateURLReader{
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
				Alternate: &sitemap.AlternateURLReader{
					XMLName: xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"},
					Rel:     "alternate",
					Lang:    "en",
					Link:    "https://dp.aws.onsdigital.uk/economy/environmentalaccounts/articles/testarticle3",
				},
			},
		},
	}
}
