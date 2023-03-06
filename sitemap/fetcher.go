package sitemap

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
	es710 "github.com/elastic/go-elasticsearch/v7"
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

type ElasticFetcher struct {
	elastic *es710.Client
	cfg     *config.Config
	zClient clients.ZebedeeClient
}

func NewElasticFetcher(elastic *es710.Client, cfg *config.Config, zc clients.ZebedeeClient) *ElasticFetcher {
	return &ElasticFetcher{
		elastic: elastic,
		cfg:     cfg,
		zClient: zc,
	}
}

func (f *ElasticFetcher) HasWelshContent(ctx context.Context, path string) bool {
	welshPath := path + "/data_cy.json"
	log.Info(ctx, "Checking welsh content for "+welshPath)
	_, err := f.zClient.GetFileSize(ctx, "", "", "cy", welshPath)
	return err == nil
}

var (
	tempSitemapFileEn = "sitemap_en"
	tempSitemapFileCy = "sitemap_cy"
)

func (f *ElasticFetcher) GetFullSitemap(ctx context.Context) (fileNames []string, err error) {
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

	fileNames = []string{fileNameEn, fileNameCy}
	files := []*os.File{fileEn, fileCy}

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

	sitemapHdContent := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	_, err = bufferedFileEn.WriteString(sitemapHdContent)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_en page xml header write error: %w", err)
	}
	_, err = bufferedFileCy.WriteString(sitemapHdContent)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_cy page xml header write error: %w", err)
	}

	var result ElasticResult
	err = f.StartScroll(ctx, &result)
	if err != nil {
		return fileNames, fmt.Errorf("failed to start scroll: %w", err)
	}

	scrollID := result.ScrollID
	for len(result.Hits.Hits) > 0 {
		for i := range result.Hits.Hits {
			uriEn, _ := url.JoinPath(f.cfg.DpOnsURLHostNameEn, result.Hits.Hits[i].Source.URI)

			if f.HasWelshContent(ctx, result.Hits.Hits[i].Source.URI) {
				// has equivalent welsh content
				// 1. Add the welsh as and alternate in english sitemap
				// 2. Add the english link as alternate in the welsh sitemap
				uriCy, _ := url.JoinPath(f.cfg.DpOnsURLHostNameCy, result.Hits.Hits[i].Source.URI)

				err = encEn.Encode(URL{
					Loc:     uriEn,
					Lastmod: result.Hits.Hits[i].Source.ReleaseDate.Format("2006-01-02"),
					Alternate: &AlternateURL{
						Rel:  "alternate",
						Lang: "cy",
						Link: uriCy,
					},
				})
				if err != nil {
					return fileNames, fmt.Errorf("sitemap_en page xml encode error: %w", err)
				}

				err = encCy.Encode(URL{
					Loc:     uriCy,
					Lastmod: result.Hits.Hits[i].Source.ReleaseDate.Format("2006-01-02"),
					Alternate: &AlternateURL{
						Rel:  "alternate",
						Lang: "en",
						Link: uriEn,
					},
				})
				if err != nil {
					return fileNames, fmt.Errorf("sitemap_cy page xml encode error: %w", err)
				}
			} else { // no welsh content
				err = encEn.Encode(URL{
					Loc:     uriEn,
					Lastmod: result.Hits.Hits[i].Source.ReleaseDate.Format("2006-01-02"),
				})
				if err != nil {
					return fileNames, fmt.Errorf("sitemap_en page xml encode error: %w", err)
				}
			}
		}
		result = ElasticResult{}
		err = f.GetScroll(ctx, scrollID, &result)
		if err != nil {
			return fileNames, fmt.Errorf("failed to get scroll: %w", err)
		}
	}

	_, err = bufferedFileEn.WriteString(`</urlset>`)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_en page xml footer write error: %w", err)
	}
	_, err = bufferedFileCy.WriteString(`</urlset>`)
	if err != nil {
		return fileNames, fmt.Errorf("sitemap_cy page xml footer write error: %w", err)
	}

	return fileNames, nil
}

func (f *ElasticFetcher) StartScroll(ctx context.Context, result interface{}) error {
	res, err := f.elastic.Search(
		f.elastic.Search.WithIndex(f.cfg.OpenSearchConfig.SitemapIndex),
		f.elastic.Search.WithScroll(f.cfg.OpenSearchConfig.ScrollTimeout),
		f.elastic.Search.WithSize(f.cfg.OpenSearchConfig.ScrollSize),
		f.elastic.Search.WithContext(ctx),
		f.elastic.Search.WithBody(strings.NewReader(`
		{
			"query": {
				"match_all": {}
			},
			"sort": [
				{"_id": "asc"}
			]
		}`),
		),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (f *ElasticFetcher) GetScroll(ctx context.Context, id string, result interface{}) error {
	res, err := f.elastic.Scroll(
		f.elastic.Scroll.WithScroll(f.cfg.OpenSearchConfig.ScrollTimeout),
		f.elastic.Scroll.WithScrollID(id),
		f.elastic.Scroll.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}
