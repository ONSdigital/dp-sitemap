package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	assetmock "github.com/ONSdigital/dp-sitemap/assets/mock"
	zcMock "github.com/ONSdigital/dp-sitemap/clients/mock"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cucumber/godog"
	es710 "github.com/elastic/go-elasticsearch/v7"
	esapi710 "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^i have the following robot\.json:$`, iHaveTheFollowingRobotjson)
	ctx.Step(`^i invoke writejson with the sitemap "([^"]*)"$`, c.iInvokeWritejsonWithTheSitemap)
	ctx.Step(`^the content of the resulting robots file must be$`, c.theContentOfTheResultingRobotsFileMustBe)
	ctx.Step(`^I generate a local sitemap$`, c.iGenerateLocalSitemap)
	ctx.Step(`^I index the following URLs:$`, c.iIndexTheFollowingURLs)
	ctx.Step(`^the content of the resulting sitemap should be$`, c.theContentOfTheResultingSitemapShouldBe)
	ctx.Step(`^I generate S3 sitemap$`, c.iGenerateS3Sitemap)
	ctx.Step(`^the content of the S3 sitemap should be$`, c.theContentOfTheS3SitemapShouldBe)
	ctx.Step(`^I add a URL "([^"]*)" dated "([^"]*)" to sitemap "([^"]*)"$`, c.iAddAURLDatedToSitemap)
	ctx.Step(`^Sitemap "([^"]*)" looks like the following:$`, c.sitemapLooksLikeTheFollowing)
	ctx.Step(`^Sitemap "([^"]*)" doesn\'t exist yet$`, c.sitemapDoesntExistYet)
	ctx.Step(`^the new content of the sitemap "([^"]*)" should be$`, c.theNewContentOfTheSitemapShouldBe)
	ctx.Step(`^URL "([^"]*)" has Welsh version$`, c.uRLHasWelshVersion)
}

func (c *Component) uRLHasWelshVersion(url string) error {
	c.welshVersion[url+"/data_cy.json"] = true
	return nil
}

func (c *Component) iAddAURLDatedToSitemap(url, date, sitemapID string) error {
	zc := zcMock.ZebedeeClientMock{
		CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
		GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
			if c.welshVersion[uri] {
				return zebedee.FileSize{Size: 1}, nil
			}
			return zebedee.FileSize{Size: 1}, errors.New("no welsh content")
		},
	}

	es := sitemap.NewElasticScroll(c.EsClient, c.cfg)

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			es,
			c.cfg,
			&zc,
		)),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
		sitemap.WithFileStore(&sitemap.LocalStore{}),
		sitemap.WithPublishingSitemapFile(c.files[sitemapID]),
	)
	err := generator.MakePublishingSitemap(context.Background(), sitemap.URL{Loc: url, Lastmod: date})
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) sitemapLooksLikeTheFollowing(sitemapID string, body *godog.DocString) error {
	file, err := os.CreateTemp("", "sitemap-component-test")
	if err != nil {
		return fmt.Errorf("failed to create sitemap file: %w", err)
	}
	err = os.WriteFile(file.Name(), []byte(body.Content), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write to sitemap file: %w", err)
	}
	c.files[sitemapID] = file.Name()
	return nil
}

func (c *Component) sitemapDoesntExistYet(sitemapID string) error {
	c.files[sitemapID] = path.Join(os.TempDir(), "sitemap-test-"+uuid.NewString())
	return nil
}

func (c *Component) theNewContentOfTheSitemapShouldBe(sitemapID string, body *godog.DocString) error {
	file, err := os.ReadFile(c.files[sitemapID])
	if err != nil {
		return fmt.Errorf("failed to read from sitemap file: %w", err)
	}
	if strings.Compare(body.Content, string(file)) != 0 {
		return fmt.Errorf("sitemap file content mismatch actual [%s], expecting [%s]", string(file), body.Content)
	}
	return nil
}

func (c *Component) iGenerateLocalSitemap() error {
	hits, err := c.indexSearchHits()
	if err != nil {
		return err
	}
	zc := zcMock.ZebedeeClientMock{
		CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
		GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
			return zebedee.FileSize{Size: 1}, errors.New("no welsh content")
		},
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

	es := sitemap.NewElasticScroll(c.EsClient, c.cfg)

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			es,
			c.cfg,
			&zc,
		)),
		sitemap.WithFileStore(&sitemap.LocalStore{}),
		sitemap.WithFullSitemapFiles(c.cfg.SitemapLocalFile),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
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
	b, err := os.ReadFile(c.cfg.SitemapLocalFile[config.English])
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
	zc := zcMock.ZebedeeClientMock{
		CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
		GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
			return zebedee.FileSize{Size: 1}, errors.New("no welsh content")
		},
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

	s3uploader := &mock.S3ClientMock{}
	s3uploader.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
		body, readErr := io.ReadAll(input.Body)
		if readErr != nil {
			return nil, readErr
		}
		c.S3UploadedSitemap[*input.Key] = string(body)
		return nil, nil
	}
	s3uploader.BucketNameFunc = func() string { return c.cfg.S3Config.UploadBucketName }

	es := sitemap.NewElasticScroll(c.EsClient, c.cfg)

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			es,
			c.cfg,
			&zc,
		)),
		sitemap.WithFileStore(sitemap.NewS3Store(s3uploader)),
		sitemap.WithFullSitemapFiles(c.cfg.S3Config.SitemapFileKey),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
	)
	err = generator.MakeFullSitemap(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) theContentOfTheS3SitemapShouldBe(arg1 *godog.DocString) error {
	fileName := c.cfg.S3Config.SitemapFileKey[config.English]
	if strings.Compare(arg1.Content, c.S3UploadedSitemap[fileName]) != 0 {
		return fmt.Errorf("s3 sitemap file content mismatch actual [%s], expecting [%s]", c.S3UploadedSitemap[fileName], arg1.Content)
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

func (c *Component) iInvokeWritejsonWithTheSitemap(arg1 string) error {
	fw := robotseo.RobotFileWriter{}
	body := fw.GetRobotsFileBody(config.English, map[config.Language]string{config.English: arg1})
	err := os.WriteFile(c.cfg.RobotsFilePath[config.English], []byte(body), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write to robots file: %w", err)
	}
	return err
}

func (c *Component) theContentOfTheResultingRobotsFileMustBe(arg1 *godog.DocString) error {
	b, err := os.ReadFile(c.cfg.RobotsFilePath[config.English])
	if err != nil {
		return err
	}
	if strings.Compare(arg1.Content, string(b)) != 0 {
		return fmt.Errorf("robot file content mismatch actual [%s], expecting [%s]", string(b), arg1.Content)
	}
	return nil
}
