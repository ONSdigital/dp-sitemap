package sitemap

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

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
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			ID     string      `json:"_id"`
			Score  interface{} `json:"_score"`
			Source struct {
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
			} `json:"_source"`
			Sort []string `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

type URL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	Lastmod string   `xml:"lastmod"`
}

type ElasticFetcher struct {
	elastic *es710.Client
	cfg     *config.OpenSearchConfig
}

func NewElasticFetcher(elastic *es710.Client, cfg *config.OpenSearchConfig) *ElasticFetcher {
	return &ElasticFetcher{
		elastic: elastic,
		cfg:     cfg,
	}
}

func (f *ElasticFetcher) GetFullSitemap(ctx context.Context) (fileName string, err error) {
	file, err := os.CreateTemp("", "sitemap")
	if err != nil {
		return "", fmt.Errorf("failed to create sitemap file: %w", err)
	}
	fileName = file.Name()
	log.Info(ctx, "created sitemap file "+fileName)
	defer func() {
		file.Close()
		// clean up the temporary file if we're returning with an error
		if err != nil {
			removeErr := os.Remove(fileName)
			if removeErr != nil {
				log.Error(ctx, "failed to remove sitemap file", err)
				return
			}
			log.Info(ctx, "removed sitemap file "+fileName)
		}
	}()

	bufferedFile := bufio.NewWriter(file)
	defer bufferedFile.Flush()

	enc := xml.NewEncoder(bufferedFile)
	enc.Indent("", "	")

	_, err = bufferedFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	if err != nil {
		return fileName, fmt.Errorf("sitemap page xml header write error: %w", err)
	}

	var result ElasticResult
	err = f.StartScroll(ctx, &result)
	if err != nil {
		return fileName, fmt.Errorf("failed to start scroll: %w", err)
	}

	scrollID := result.ScrollID
	i := 0
	for len(result.Hits.Hits) > 0 {
		for _, hit := range result.Hits.Hits {
			i++
			err = enc.Encode(URL{
				Loc:     hit.Source.URI,
				Lastmod: hit.Source.ReleaseDate.Format("2006-01-02"),
			})
			if err != nil {
				return fileName, fmt.Errorf("sitemap page xml encode error: %w", err)
			}
		}
		result = ElasticResult{}
		err = f.GetScroll(ctx, scrollID, &result)
		if err != nil {
			return fileName, fmt.Errorf("failed to get scroll: %w", err)
		}
	}

	_, err = bufferedFile.WriteString(`</urlset>`)
	if err != nil {
		return fileName, fmt.Errorf("sitemap page xml footer write error: %w", err)
	}

	return fileName, nil
}

func (f *ElasticFetcher) StartScroll(ctx context.Context, result interface{}) error {
	res, err := f.elastic.Search(
		f.elastic.Search.WithIndex(f.cfg.SitemapIndex),
		f.elastic.Search.WithScroll(f.cfg.ScrollTimeout),
		f.elastic.Search.WithSize(f.cfg.ScrollSize),
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

	fmt.Print(res.Status() + "\n")

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (f *ElasticFetcher) GetScroll(ctx context.Context, id string, result interface{}) error {
	res, err := f.elastic.Scroll(
		f.elastic.Scroll.WithScroll(f.cfg.ScrollTimeout),
		f.elastic.Scroll.WithScrollID(id),
		f.elastic.Scroll.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Print(res.Status() + "\n")

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}
