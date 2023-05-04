package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ONSdigital/dp-sitemap/sitemap"
)

type FakeScroll struct {
}

func NewFakeScroll() sitemap.Scroll {
	return &FakeScroll{}
}

func (fs *FakeScroll) StartScroll(ctx context.Context, result interface{}) error {
	fakeStartScroll(result)
	return nil
}

func (fs *FakeScroll) GetScroll(ctx context.Context, id string, result interface{}) error {
	return nil
}

func fakeStartScroll(res interface{}) {
	r, ok := res.(*sitemap.ElasticResult)
	if !ok {
		fmt.Printf("Type assertion for %v failed.\n", res)
		os.Exit(1)
	}

	hit := sitemap.ElasticHit{
		Source: sitemap.ElasticHitSource{
			URI:         "/economy/environmentalaccounts/articles/testarticle",
			ReleaseDate: time.Now(),
		},
	}
	r.Hits.Hits = append(r.Hits.Hits, hit)
}
