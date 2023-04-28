package sitemap

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ONSdigital/dp-sitemap/config"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

type ElasticScroll struct {
	elastic *es710.Client
	cfg     *config.Config
}

func NewElasticScroll(elastic *es710.Client, cfg *config.Config) *ElasticScroll {
	return &ElasticScroll{
		elastic: elastic,
		cfg:     cfg,
	}
}

func (f *ElasticScroll) GetScroll(ctx context.Context, id string, result interface{}) error {
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

func (f *ElasticScroll) StartScroll(ctx context.Context, result interface{}) error {
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
