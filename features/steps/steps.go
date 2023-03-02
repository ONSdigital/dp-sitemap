package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	assetmock "github.com/ONSdigital/dp-sitemap/assets/mock"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cucumber/godog"
	es710 "github.com/elastic/go-elasticsearch/v7"
	esapi710 "github.com/elastic/go-elasticsearch/v7/esapi"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^i have the following robot\.json:$`, iHaveTheFollowingRobotjson)
	ctx.Step(`^i invoke writejson with the sitemaps "([^"]*)"$`, c.iInvokeWritejsonWithTheSitemaps)
	ctx.Step(`^the content of the resulting robots file must be$`, c.theContentOfTheResultingRobotsFileMustBe)
	ctx.Step(`^I generate a local sitemap$`, c.iGenerateLocalSitemap)
	ctx.Step(`^I index the following URLs:$`, c.iIndexTheFollowingURLs)
	ctx.Step(`^the content of the resulting sitemap should be$`, c.theContentOfTheResultingSitemapShouldBe)
	ctx.Step(`^I generate S3 sitemap$`, c.iGenerateS3Sitemap)
	ctx.Step(`^the content of the S3 sitemap should be$`, c.theContentOfTheS3SitemapShouldBe)
}

func (c *Component) iGenerateLocalSitemap() error {
	hits, err := c.indexSearchHits()
	if err != nil {
		return err
	}
	c.EsClient = &es710.Client{API: &esapi710.API{
		Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
			return &esapi710.Response{
				Body: io.NopCloser(strings.NewReader(hits)),
			}, nil
		},
		Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
			return &esapi710.Response{
				Body: io.NopCloser(strings.NewReader("{}")),
			}, nil
		},
	}}

	generator := sitemap.NewGenerator(
		sitemap.NewElasticFetcher(
			c.EsClient,
			&c.cfg.OpenSearchConfig,
		),
		sitemap.NewLocalSaver(c.cfg.SitemapLocalFile),
	)
	err = generator.MakeFullSitemap(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) indexSearchHits() (string, error) {
	hits := sitemap.ElasticResult{
		ScrollID: "scroll1",
	}
	for _, url := range c.EsIndex.Rows {
		releaseDate, err := time.Parse("2006-01-02", url.Cells[1].Value)
		if err != nil {
			return "", err
		}
		hits.Hits.Hits = append(hits.Hits.Hits, sitemap.ElasticHit{
			Source: sitemap.ElasticHitSource{
				URI:         url.Cells[0].Value,
				ReleaseDate: releaseDate,
			},
		})
	}
	jsonHits, err := json.Marshal(hits)
	if err != nil {
		return "", err
	}

	return string(jsonHits), nil
}

func (c *Component) iIndexTheFollowingURLs(urls *godog.Table) error {
	c.EsIndex = urls
	return nil
}

func (c *Component) theContentOfTheResultingSitemapShouldBe(arg1 *godog.DocString) error {
	b, err := os.ReadFile(c.cfg.SitemapLocalFile)
	if err != nil {
		return err
	}
	if strings.Compare(arg1.Content, string(b)) != 0 {
		return fmt.Errorf("sitemap file content mismatch actual [%s], expecting [%s]", string(b), arg1.Content)
	}
	return nil
}

func (c *Component) iGenerateS3Sitemap() error {
	hits, err := c.indexSearchHits()
	if err != nil {
		return err
	}
	c.EsClient = &es710.Client{API: &esapi710.API{
		Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
			return &esapi710.Response{
				Body: io.NopCloser(strings.NewReader(hits)),
			}, nil
		},
		Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
			return &esapi710.Response{
				Body: io.NopCloser(strings.NewReader("{}")),
			}, nil
		},
	}}

	s3uploader := &mock.S3UploaderMock{}
	s3uploader.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
		body, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, err
		}
		c.S3UploadedSitemap = string(body)
		return nil, nil
	}

	generator := sitemap.NewGenerator(
		sitemap.NewElasticFetcher(
			c.EsClient,
			&c.cfg.OpenSearchConfig,
		),
		sitemap.NewS3Saver(s3uploader, c.cfg.S3Config.UploadBucketName, c.cfg.S3Config.SitemapFileKey),
	)
	err = generator.MakeFullSitemap(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) theContentOfTheS3SitemapShouldBe(arg1 *godog.DocString) error {
	if strings.Compare(arg1.Content, c.S3UploadedSitemap) != 0 {
		return fmt.Errorf("s3 sitemap file content mismatch actual [%s], expecting [%s]", c.S3UploadedSitemap, arg1.Content)
	}
	return nil
}

func iHaveTheFollowingRobotjson(arg1 *godog.DocString) error {
	amock := assetmock.FileSystemInterfaceMock{
		GetFunc: func(contextMoqParam context.Context, path string) ([]byte, error) { return []byte(arg1.Content), nil },
	}
	robotseo.Init(&amock)
	return nil
}

func (c *Component) iInvokeWritejsonWithTheSitemaps(arg1 string) error {
	fw := robotseo.RobotFileWriter{}
	return fw.WriteRobotsFile(c.cfg, strings.Split(arg1, ","))
}

func (c *Component) theContentOfTheResultingRobotsFileMustBe(arg1 *godog.DocString) error {
	b, err := os.ReadFile(c.cfg.RobotsFilePath)
	if err != nil {
		return err
	}
	if strings.Compare(arg1.Content, string(b)) != 0 {
		return fmt.Errorf("robot file content mismatch actual [%s], expecting [%s]", string(b), arg1.Content)
	}
	return nil
}
