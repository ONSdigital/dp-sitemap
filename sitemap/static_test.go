package sitemap

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/ONSdigital/dp-sitemap/config"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestStore struct {
}

var buf = new(bytes.Buffer)

func (s *TestStore) SaveFile(name string, body io.Reader) error {
	io.Copy(buf, body)
	return nil
}
func (s *TestStore) GetFile(name string) (body io.ReadCloser, err error) { return nil, nil }
func (s *TestStore) CopyFile(src io.Reader, dest io.Writer) error        { return nil }
func (s *TestStore) CreateFile(name string) (io.ReadWriteCloser, error)  { return nil, nil }
func TestLoadStaticSitemap(t *testing.T) {
	Convey("given we have static sitemap file sitemap_en.json", t, func() {

		Convey("when loading english static sitemap", func() {
			store := TestStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en_test.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &store)
			urlSet := &UrlsetReader{}
			xml.Unmarshal(buf.Bytes(), urlSet)
			buf.Reset()
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
			store := TestStore{}
			cfg, _ := config.Get()
			err := LoadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy_test.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &store)
			urlSet := &UrlsetReader{}
			xml.Unmarshal(buf.Bytes(), urlSet)
			buf.Reset()
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("We should have the correct content loaded", func() {
				So(urlSet, ShouldResemble, expectedUrlSetWelsh())
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
