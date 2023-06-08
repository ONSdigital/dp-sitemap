package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ONSdigital/dp-sitemap/sitemap"
)

type FakeScroll struct {
}

func NewFakeScroll() sitemap.Scroll {
	return &FakeScroll{}
}

func (fs *FakeScroll) StartScroll(ctx context.Context, result interface{}) error {
	return fakeStartScroll(result)
}

func (fs *FakeScroll) GetScroll(ctx context.Context, id string, result interface{}) error {
	return nil
}

func fakeStartScroll(res interface{}) error {
	r, ok := res.(*sitemap.ElasticResult)
	if !ok {
		return fmt.Errorf("type assertion for %v failed", res)
	}

	hit := sitemap.ElasticHit{
		Source: sitemap.ElasticHitSource{
			URI:         "/economy/environmentalaccounts/articles/testarticle",
			ReleaseDate: time.Now(),
		},
	}
	r.Hits.Hits = append(r.Hits.Hits, hit)
	return nil
}
