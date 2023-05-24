package sitemap

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
)

type ElasticResult struct {
	ScrollID string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{}  `json:"max_score"`
		Hits     []ElasticHit `json:"hits"`
	} `json:"hits"`
}

type ElasticHit struct {
	Index  string           `json:"_index"`
	Type   string           `json:"_type"`
	ID     string           `json:"_id"`
	Score  interface{}      `json:"_score"`
	Source ElasticHitSource `json:"_source"`
	Sort   []string         `json:"sort"`
}
type ElasticHitSource struct {
	Type            string      `json:"type"`
	URI             string      `json:"uri"`
	JobID           string      `json:"job_id"`
	SearchIndex     string      `json:"search_index"`
	Cdid            string      `json:"cdid"`
	DatasetID       string      `json:"dataset_id"`
	Edition         string      `json:"edition"`
	Keywords        []string    `json:"keywords"`
	MetaDescription string      `json:"meta_description"`
	ReleaseDate     time.Time   `json:"release_date"`
	Summary         string      `json:"summary"`
	Title           string      `json:"title"`
	Topics          interface{} `json:"topics"`
	Cancelled       bool        `json:"cancelled"`
	Finalised       bool        `json:"finalised"`
	Published       bool        `json:"published"`
	CanonicalTopic  string      `json:"canonical_topic"`
}
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Xhtml   string   `xml:"xmlns:xhtml,attr"`
	URL     []URL    `xml:"url"`
}

type URL struct {
	XMLName   xml.Name      `xml:"url"`
	Loc       string        `xml:"loc"`
	Lastmod   string        `xml:"lastmod"`
	Alternate *AlternateURL `xml:"xhtml:link,omitempty"`
}

type AlternateURL struct {
	XMLName xml.Name `xml:"xhtml:link,omitempty"`
	Rel     string   `xml:"rel,omitempty"`
	Lang    string   `xml:"hreflang,omitempty"`
	Link    string   `xml:"href,omitempty"`
}

type UrlsetReader struct {
	XMLName xml.Name    `xml:"urlset"`
	Xmlns   string      `xml:"xmlns,attr"`
	Xhtml   string      `xml:"xmlns:xhtml,attr"`
	URL     []URLReader `xml:"url"`
}

type URLReader struct {
	XMLName   xml.Name            `xml:"url"`
	Loc       string              `xml:"loc"`
	Lastmod   string              `xml:"lastmod"`
	Alternate *AlternateURLReader `xml:"link,omitempty"`
}

type AlternateURLReader struct {
	XMLName xml.Name `xml:"link,omitempty"`
	Rel     string   `xml:"rel,omitempty"`
	Lang    string   `xml:"hreflang,omitempty"`
	Link    string   `xml:"href,omitempty"`
}

type ElasticFetcher struct {
	scroll  Scroll
	cfg     *config.Config
	zClient clients.ZebedeeClient
}

func NewElasticFetcher(scroll Scroll, cfg *config.Config, zc clients.ZebedeeClient) *ElasticFetcher {
	return &ElasticFetcher{
		scroll:  scroll,
		cfg:     cfg,
		zClient: zc,
	}
}

func (f *ElasticFetcher) HasWelshContent(ctx context.Context, path string) bool {
	welshPath := path + "/data_cy.json"
	log.Info(ctx, "Checking welsh content for "+welshPath)
	_, err := f.zClient.GetFileSize(ctx, "", "", config.Welsh.String(), welshPath)
	return err == nil
}

func (f *ElasticFetcher) URLVersions(ctx context.Context, path, lastmod string) (en URL, cy *URL) {
	enLoc, _ := url.JoinPath(f.cfg.DpOnsURLHostNameEn, path)
	en = URL{
		Loc:     enLoc,
		Lastmod: lastmod,
	}
	if f.HasWelshContent(ctx, path) {
		// has equivalent welsh content
		// 1. Add the welsh as and alternate in english sitemap
		// 2. Add the english link as alternate in the welsh sitemap
		cyLoc, _ := url.JoinPath(f.cfg.DpOnsURLHostNameCy, path)
		en.Alternate = &AlternateURL{
			Rel:  "alternate",
			Lang: config.Welsh.String(),
			Link: cyLoc,
		}
		cy = &URL{
			Loc:     cyLoc,
			Lastmod: lastmod,
			Alternate: &AlternateURL{
				Rel:  "alternate",
				Lang: config.English.String(),
				Link: enLoc,
			},
		}
	}
	return
}

var (
	tempSitemapFileEn = "sitemap_en"
	tempSitemapFileCy = "sitemap_cy"
)

func (f *ElasticFetcher) GetFullSitemap(ctx context.Context) (fileNames Files, err error) {
	fileEn, err := os.CreateTemp("", tempSitemapFileEn)
	if err != nil {
		return fileNames, fmt.Errorf("failed to create sitemap_en file: %w", err)
	}
	fileCy, err := os.CreateTemp("", tempSitemapFileCy)
	if err != nil {
		return fileNames, fmt.Errorf("failed to create sitemap_cy file: %w", err)
	}
	fileNameEn := fileEn.Name()
	fileNameCy := fileCy.Name()

	fileNames = Files{config.English: fileNameEn, config.Welsh: fileNameCy}
	files := map[config.Language]*os.File{config.English: fileEn, config.Welsh: fileCy}

	log.Info(ctx, "created sitemap files "+fileNameEn+", "+fileNameCy)
	defer func() {
		for i, fl := range files {
			fl.Close()
			// clean up the temporary file if we're returning with an error
			if err != nil {
				removeErr := os.Remove(fileNames[i])
				if removeErr != nil {
					log.Error(ctx, "failed to remove sitemap file "+fileNames[i], err)
					return
				}
				log.Info(ctx, "removed sitemap file "+fileNames[i])
			}
		}
	}()

	bufferedFileEn := bufio.NewWriter(fileEn)
	defer bufferedFileEn.Flush()
	encEn := xml.NewEncoder(bufferedFileEn)
	encEn.Indent("", "  ")

	bufferedFileCy := bufio.NewWriter(fileCy)
	defer bufferedFileCy.Flush()
	encCy := xml.NewEncoder(bufferedFileCy)
	encCy.Indent("", "  ")

	sitemapHdContent := xml.Header + `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n"
	_, err = bufferedFileEn.WriteString(sitemapHdContent)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_en page xml header write error: %w", err)
	}
	_, err = bufferedFileCy.WriteString(sitemapHdContent)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_cy page xml header write error: %w", err)
	}

	var result ElasticResult
	err = f.scroll.StartScroll(ctx, &result)
	if err != nil {
		return fileNames, fmt.Errorf("failed to start scroll: %w", err)
	}

	scrollID := result.ScrollID
	for len(result.Hits.Hits) > 0 {
		for i := range result.Hits.Hits {
			urlEn, urlCy := f.URLVersions(
				ctx,
				result.Hits.Hits[i].Source.URI,
				result.Hits.Hits[i].Source.ReleaseDate.Format("2006-01-02"),
			)

			err = encEn.Encode(urlEn)
			if err != nil {
				return fileNames, fmt.Errorf("sitemap_en page xml encode error: %w", err)
			}
			if urlCy != nil {
				err = encCy.Encode(urlCy)
				if err != nil {
					return fileNames, fmt.Errorf("sitemap_cy page xml encode error: %w", err)
				}
			}
		}
		result = ElasticResult{}
		err = f.scroll.GetScroll(ctx, scrollID, &result)
		if err != nil {
			return fileNames, fmt.Errorf("failed to get scroll: %w", err)
		}
	}

	_, err = bufferedFileEn.WriteString("\n" + `</urlset>`)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_en page xml footer write error: %w", err)
	}
	_, err = bufferedFileCy.WriteString("\n" + `</urlset>`)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_cy page xml footer write error: %w", err)
	}

	return fileNames, nil
}

func (f *ElasticFetcher) GetZebedeeClient() clients.ZebedeeClient {
	return f.zClient
}

func (f *ElasticFetcher) GetConfig() *config.Config {
	return f.cfg
}
