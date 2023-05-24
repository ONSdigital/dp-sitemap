package sitemap_test

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	zcMock "github.com/ONSdigital/dp-sitemap/clients/mock"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
	esapi710 "github.com/elastic/go-elasticsearch/v7/esapi"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFetcher(t *testing.T) {
	cfg := &config.Config{}
	zc := zcMock.ZebedeeClientMock{
		CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
		GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
			return zebedee.FileSize{Size: 1}, errors.New("no welsh content")
		},
	}

	Convey("When elastic start scroll returns an error", t, func() {
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return nil, errors.New("start scroll error")
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zc)
		filename, err := f.GetFullSitemap(context.Background())

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to start scroll")
			So(err.Error(), ShouldContainSubstring, "start scroll error")
		})
		Convey("Temporary sitemap file should be created and then cleaned up", func() {
			So(filename[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filename[config.English])
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
	Convey("When elastic start scroll returns no hits", t, func() {
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zc)
		filenames, err := f.GetFullSitemap(context.Background())
		defer func() {
			for _, fl := range filenames {
				os.Remove(fl)
			}
		}()

		Convey("Fetcher should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filenames[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filenames[config.English])
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be a valid xml and include no urls", func() {
			sitemapContent, err := os.ReadFile(filenames[config.English])
			So(err, ShouldBeNil)
			So(string(sitemapContent), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">

</urlset>`)
		})
	})
	Convey("When elastic start scroll returns hits", t, func() {
		var receivedScrollID string
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_1",
										"release_date": "2014-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_2",
										"release_date": "2023-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`)),
				}, nil
			},
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				req := esapi710.ScrollRequest{}
				for _, f := range o {
					f(&req)
				}
				receivedScrollID = req.ScrollID
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zc)
		filenames, err := f.GetFullSitemap(context.Background())
		defer func() {
			for _, fl := range filenames {
				os.Remove(fl)
			}
		}()

		Convey("Fetcher should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Correct scroll ID should be passed", func() {
			So(receivedScrollID, ShouldEqual, "scroll_id_1")
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filenames[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filenames[config.English])
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be valid and include all received urls", func() {
			sitemapContent, err := os.ReadFile(filenames[config.English])
			So(err, ShouldBeNil)
			So(string(sitemapContent), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
<url>
  <loc>uri_1</loc>
  <lastmod>2014-12-10</lastmod>
</url>
<url>
  <loc>uri_2</loc>
  <lastmod>2023-03-31</lastmod>
</url>
</urlset>`)
		})
	})
	Convey("When subsequent scroll returns error", t, func() {
		var receivedScrollID string
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_1",
										"release_date": "2014-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_2",
										"release_date": "2023-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`)),
				}, nil
			},
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				req := esapi710.ScrollRequest{}
				for _, f := range o {
					f(&req)
				}
				receivedScrollID = req.ScrollID
				return nil, errors.New("subsequent scroll error")
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zc)
		filename, err := f.GetFullSitemap(context.Background())

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to get scroll")
			So(err.Error(), ShouldContainSubstring, "subsequent scroll error")
		})
		Convey("Correct scroll ID should be passed", func() {
			So(receivedScrollID, ShouldEqual, "scroll_id_1")
		})
		Convey("Temporary sitemap file should be created and then cleaned up", func() {
			So(filename[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filename[config.English])
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
	Convey("When subsequent scrolls return hits", t, func() {
		var receivedScrollID string
		subsequentScrollsLeft := 2
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_1",
										"release_date": "2014-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_2",
										"release_date": "2023-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`)),
				}, nil
			},
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				req := esapi710.ScrollRequest{}
				for _, f := range o {
					f(&req)
				}
				receivedScrollID = req.ScrollID
				body := io.NopCloser(strings.NewReader("{}"))
				if subsequentScrollsLeft > 0 {
					body = io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_3",
										"release_date": "2015-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_4",
										"release_date": "2024-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`))
				}
				subsequentScrollsLeft--
				return &esapi710.Response{
					Body: body,
				}, nil
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zc)
		filenames, err := f.GetFullSitemap(context.Background())
		defer func() {
			for _, fl := range filenames {
				os.Remove(fl)
			}
		}()

		Convey("Fetcher should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Correct scroll ID should be passed", func() {
			So(receivedScrollID, ShouldEqual, "scroll_id_1")
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filenames[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filenames[config.English])
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be valid and include all received urls", func() {
			sitemapContent, err := os.ReadFile(filenames[config.English])
			So(err, ShouldBeNil)
			So(string(sitemapContent), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
<url>
  <loc>uri_1</loc>
  <lastmod>2014-12-10</lastmod>
</url>
<url>
  <loc>uri_2</loc>
  <lastmod>2023-03-31</lastmod>
</url>
<url>
  <loc>uri_3</loc>
  <lastmod>2015-12-10</lastmod>
</url>
<url>
  <loc>uri_4</loc>
  <lastmod>2024-03-31</lastmod>
</url>
<url>
  <loc>uri_3</loc>
  <lastmod>2015-12-10</lastmod>
</url>
<url>
  <loc>uri_4</loc>
  <lastmod>2024-03-31</lastmod>
</url>
</urlset>`)
		})
	})
	Convey("When subsequent scrolls return hits (with welsh content)", t, func() {
		var receivedScrollID string
		subsequentScrollsLeft := 2
		zcWithWelsh := zcMock.ZebedeeClientMock{
			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
			GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
				return zebedee.FileSize{Size: 1}, nil
			},
		}
		esMock := &es710.Client{API: &esapi710.API{
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_1",
										"release_date": "2014-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_2",
										"release_date": "2023-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`)),
				}, nil
			},
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				req := esapi710.ScrollRequest{}
				for _, f := range o {
					f(&req)
				}
				receivedScrollID = req.ScrollID
				body := io.NopCloser(strings.NewReader("{}"))
				if subsequentScrollsLeft > 0 {
					body = io.NopCloser(strings.NewReader(`
					{
						"_scroll_id": "scroll_id_1",
						"hits": {
							"hits": [
								{
									"_source": {
										"uri": "uri_3",
										"release_date": "2015-12-10T00:00:00.000Z"
									}
								},
								{
									"_source": {
										"uri": "uri_4",
										"release_date": "2024-03-31T00:00:00.000Z"
									}
								}
							]
						}
					}
					`))
				}
				subsequentScrollsLeft--
				return &esapi710.Response{
					Body: body,
				}, nil
			},
		}}

		scroller := sitemap.NewElasticScroll(esMock, cfg)

		f := sitemap.NewElasticFetcher(scroller, cfg, &zcWithWelsh)
		filenames, err := f.GetFullSitemap(context.Background())
		defer func() {
			for _, fl := range filenames {
				os.Remove(fl)
			}
		}()

		Convey("Fetcher should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Correct scroll ID should be passed", func() {
			So(receivedScrollID, ShouldEqual, "scroll_id_1")
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filenames[config.English], ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filenames[config.English])
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be valid and include all received urls", func() {
			sitemapContent, err := os.ReadFile(filenames[config.English])
			So(err, ShouldBeNil)
			So(string(sitemapContent), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
<url>
  <loc>uri_1</loc>
  <lastmod>2014-12-10</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_1</href>
  </xhtml:link>
</url>
<url>
  <loc>uri_2</loc>
  <lastmod>2023-03-31</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_2</href>
  </xhtml:link>
</url>
<url>
  <loc>uri_3</loc>
  <lastmod>2015-12-10</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_3</href>
  </xhtml:link>
</url>
<url>
  <loc>uri_4</loc>
  <lastmod>2024-03-31</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_4</href>
  </xhtml:link>
</url>
<url>
  <loc>uri_3</loc>
  <lastmod>2015-12-10</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_3</href>
  </xhtml:link>
</url>
<url>
  <loc>uri_4</loc>
  <lastmod>2024-03-31</lastmod>
  <xhtml:link>
    <rel>alternate</rel>
    <hreflang>cy</hreflang>
    <href>uri_4</href>
  </xhtml:link>
</url>
</urlset>`)
		})
	})
}
