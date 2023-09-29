package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ONSdigital/dp-sitemap/sitemap"
)

// FakeScroll is used to replicate elastic search
type FakeScroll struct {
}

func NewFakeScroll() sitemap.Scroll {
	return &FakeScroll{}
}

func (fs *FakeScroll) StartScroll(_ context.Context, result interface{}) error {
	return fakeStartScroll(result)
}

func (fs *FakeScroll) GetScroll(_ context.Context, _ string, _ interface{}) error {
	return nil
}

func fakeStartScroll(res interface{}) error {
	r, ok := res.(*sitemap.ElasticResult)
	if !ok {
		return fmt.Errorf("type assertion for %v failed", res)
	}

	// we have created a couple of test articles under $zebedee_root
	// so we can test the code locally and not in the sandbox environment
	hit := sitemap.ElasticHit{
		Source: sitemap.ElasticHitSource{
			URI:         "/economy/environmentalaccounts/articles/testarticle",
			ReleaseDate: time.Now(),
		},
	}
	r.Hits.Hits = append(r.Hits.Hits, hit)
	return nil
}
