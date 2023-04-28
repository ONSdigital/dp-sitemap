package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-sitemap/sitemap"
)

type FakeScroll struct {
}

func NewFakeScroll() sitemap.Scroll {
	return &FakeScroll{}
}

func (fs *FakeScroll) StartScroll(ctx context.Context, result interface{}) error {
	res := fakeStartScroll()
	log.Println(res)
	return nil
}

func (fs *FakeScroll) GetScroll(ctx context.Context, id string, result interface{}) error {
	return nil
}

func fakeStartScroll() sitemap.ElasticResult {
	res := sitemap.ElasticResult{}

	for i := 0; i < 3; i++ {
		hit := sitemap.ElasticHit{
			Source: sitemap.ElasticHitSource{
				URI:         "https://localhost/test" + strconv.Itoa(i),
				ReleaseDate: time.Now(),
			},
		}
		res.Hits.Hits = append(res.Hits.Hits, hit)
	}

	return res
}
