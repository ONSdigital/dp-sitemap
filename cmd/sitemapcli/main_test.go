package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/ONSdigital/dp-sitemap/assets"
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
		err := sitemap.LoadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
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
		err := sitemap.LoadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
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
	Convey("Given we have static sitemap file sitemap_en.json", t, func() {

		Convey("When loading english static sitemap", func() {
			store := mock.FileStoreMock{}
			buf := new(bytes.Buffer)
			store.SaveFileFunc = func(name string, body io.Reader) error {
				io.Copy(buf, body)
				return nil
			}
			cfg, _ := config.Get()
			err := sitemap.LoadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
			urlSet := &sitemap.UrlsetReader{}
			xml.Unmarshal(buf.Bytes(), urlSet)

			Convey("Than there should be no error", func() {
				So(err, ShouldBeNil)
				Convey("And we should have the correct content loaded", func() {
					So(urlSet, ShouldResemble, expectedUrlSetEnglish())
				})
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
			err := sitemap.LoadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
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
	staticSitemapName := "sitemap_en.json"
	cfg, _ := config.Get()
	efs := assets.NewFromEmbeddedFilesystem()

	b, err := efs.Get(context.Background(), assets.Sitemap, staticSitemapName)
	if err != nil {
		panic("can't find file " + staticSitemapName)
	}

	var content []sitemap.StaticURL

	err = json.Unmarshal(b, &content)
	if err != nil {
		return nil
	}

	// move old sitemap urls to new sitemap
	sitemapReader := sitemap.UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml:   "",
	}

	// range through static content
	for _, contentItem := range content {
		var newURL sitemap.URLReader
		newURL.XMLName = xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"}
		newURL.Loc = cfg.DpOnsURLHostNameEn + contentItem.URL
		newURL.Lastmod = contentItem.ReleaseDate
		newURL.Alternate = &sitemap.AlternateURLReader{}
		if contentItem.HasAltLang == true {
			newURL.Alternate.XMLName = xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"}
			newURL.Alternate.Rel = "alternate"
			newURL.Alternate.Link = cfg.DpOnsURLHostNameCy + contentItem.URL
			newURL.Alternate.Lang = "cy"
		}
		sitemapReader.URL = append(sitemapReader.URL, newURL)
	}
	return &sitemapReader
}

func expectedUrlSetWelsh() *sitemap.UrlsetReader {
	staticSitemapName := "sitemap_cy.json"
	cfg, _ := config.Get()
	efs := assets.NewFromEmbeddedFilesystem()

	b, err := efs.Get(context.Background(), assets.Sitemap, staticSitemapName)
	if err != nil {
		panic("can't find file " + staticSitemapName)
	}

	var content []sitemap.StaticURL

	err = json.Unmarshal(b, &content)
	if err != nil {
		return nil
	}

	// move old sitemap urls to new sitemap
	sitemapReader := sitemap.UrlsetReader{
		XMLName: xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml:   "",
	}

	// range through static content
	for _, contentItem := range content {
		var newURL sitemap.URLReader
		newURL.XMLName = xml.Name{Space: "http://www.sitemaps.org/schemas/sitemap/0.9", Local: "url"}
		newURL.Loc = cfg.DpOnsURLHostNameCy + contentItem.URL
		newURL.Lastmod = contentItem.ReleaseDate
		newURL.Alternate = &sitemap.AlternateURLReader{}
		if contentItem.HasAltLang == true {
			newURL.Alternate.XMLName = xml.Name{Space: "http://www.w3.org/1999/xhtml", Local: "link"}
			newURL.Alternate.Rel = "alternate"
			newURL.Alternate.Link = cfg.DpOnsURLHostNameEn + contentItem.URL
			newURL.Alternate.Lang = "en"
		}
		sitemapReader.URL = append(sitemapReader.URL, newURL)
	}
	return &sitemapReader
}

func expectedUrlSetWelshV1() *sitemap.UrlsetReader {
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
