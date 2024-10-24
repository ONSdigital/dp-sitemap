package utilities

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-net/v2/awsauth"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

func createCliSitemapGenerator(cfg *config.Config, commandline *FlagFields) (*sitemap.Generator, error) {
	store := &sitemap.LocalStore{}

	var transport http.RoundTripper = dphttp.DefaultTransport

	if cfg.OpenSearchConfig.Signer {
		var err error
		transport, err = awsauth.NewAWSSignerRoundTripper(cfg.OpenSearchConfig.ElasticSearchURL, cfg.OpenSearchConfig.SignerFilename, cfg.OpenSearchConfig.SignerRegion, cfg.OpenSearchConfig.SignerService, awsauth.Options{TlsInsecureSkipVerify: cfg.OpenSearchConfig.TLSInsecureSkipVerify})
		if err != nil {
			fmt.Printf("failed to save file")
			return nil, err
		}
	}

	rawClient, err := es710.NewClient(es710.Config{
		Addresses: []string{commandline.ElasticSearchURL},
		Transport: transport,
	})
	if err != nil {
		return nil, err
	}

	// Get zebedeeClient using arg -zebedee-url
	zebedeeClient := zebedee.New(commandline.ZebedeeURL)

	var scroll sitemap.Scroll
	if commandline.FakeScroll {
		scroll = NewFakeScroll()
	} else {
		scroll = sitemap.NewElasticScroll(
			rawClient,
			cfg,
		)
	}

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			scroll,
			cfg,
			zebedeeClient,
		)),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
		sitemap.WithFileStore(store),
		sitemap.WithFullSitemapFiles(map[config.Language]string{
			config.English: commandline.SitemapPath + "_en.xml",
			config.Welsh:   commandline.SitemapPath + "_cy.xml",
		}),
	)

	return generator, nil
}

func GenerateSitemap(cfg *config.Config, commandline *FlagFields) {
	generator, err := createCliSitemapGenerator(cfg, commandline)
	if err != nil {
		fmt.Println("Error creating sitemap generator", err.Error())
		os.Exit(1)
	}

	// Generating sitemap
	genErr := generator.MakeFullSitemap(context.Background())
	if genErr != nil {
		fmt.Println("Error writing sitemap file", genErr.Error())
		os.Exit(1)
	}
	fmt.Println("sitemap generation job complete")
}

func GenerateRobotFile(cfg *config.Config, commandline *FlagFields) {
	robotseo.Init(commandline.RobotsFilePathReader)
	robotFileWriter := robotseo.RobotFileWriter{}
	cfg.RobotsFilePath = map[config.Language]string{
		config.English: commandline.RobotsFilePath,
	}

	store := &sitemap.LocalStore{}

	cfg.OpenSearchConfig.ElasticSearchURL = commandline.ElasticSearchURL
	cfg.OpenSearchConfig.ScrollSize = commandline.ScrollSize
	cfg.OpenSearchConfig.Signer = true

	body := robotFileWriter.GetRobotsFileBody(config.English, cfg.SitemapLocalFile)

	saveErr := store.SaveFile(commandline.RobotsFilePath, strings.NewReader(body))
	if saveErr != nil {
		fmt.Println("failed to save file")
		return
	}
	fmt.Println("robot file creation successful")
}

func UpdateSitemap(cfg *config.Config, commandLine *FlagFields) error {
	var scroll sitemap.Scroll
	if commandLine.FakeScroll {
		scroll = &FakeScroll{}
	} else {
		scroll = &sitemap.ElasticScroll{}
	}
	var store sitemap.FileStore
	if cfg.SitemapSaveLocation == "local" {
		store = &sitemap.LocalStore{}
	} else {
		store = &sitemap.S3Store{}
	}
	zebedeeClient := zebedee.New(commandLine.ZebedeeURL)
	fetcher := sitemap.NewElasticFetcher(scroll, cfg, zebedeeClient)
	handler := event.NewContentPublishedHandler(store, zebedeeClient, cfg, fetcher)
	content, contentErr := getContent()
	if contentErr != nil {
		fmt.Println("Failed to get event content from user:", contentErr)
		return contentErr
	}

	err := handler.Handle(context.Background(), cfg, content)
	if err != nil {
		fmt.Println("Failed to handle event:", err)
		return err
	}
	fmt.Println("sitemap update job complete")
	return nil
}

func LoadStaticSitemap(cfg *config.Config, commandLine *FlagFields) {
	if !strings.HasSuffix(commandLine.SitemapPathReader, "/") {
		commandLine.SitemapPathReader += "/"
	}
	var err error
	err = sitemap.LoadStaticSitemap(cfg, commandLine.SitemapPath+"_en", commandLine.SitemapPathReader+"sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &sitemap.LocalStore{})
	if err != nil {
		fmt.Println("Failed to load english static sitemap:", err)
		return
	}
	err = sitemap.LoadStaticSitemap(cfg, commandLine.SitemapPath+"_cy", commandLine.SitemapPathReader+"sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &sitemap.LocalStore{})
	if err != nil {
		fmt.Println("Failed to load welsh static sitemap:", err)
		return
	}
}

var getContent = func() (*event.ContentPublished, error) {
	content := &event.ContentPublished{}
	fmt.Print("Please enter URI: ")
	text, err := getText()
	if err != nil {
		return nil, err
	}
	content.URI = *text
	fmt.Print("Please enter Data Type: ")
	text, err = getText()
	if err != nil {
		return nil, err
	}
	content.DataType = *text
	fmt.Print("Please enter Collection ID: ")
	text, err = getText()
	if err != nil {
		return nil, err
	}
	content.CollectionID = *text
	fmt.Print("Please enter Job ID: ")
	text, err = getText()
	if err != nil {
		return nil, err
	}
	content.JobID = *text
	fmt.Print("Please enter Search Index: ")
	text, err = getText()
	if err != nil {
		return nil, err
	}
	content.SearchIndex = *text
	fmt.Print("Please enter Trace ID: ")
	text, err = getText()
	if err != nil {
		return nil, err
	}
	content.TraceID = *text
	return content, nil
}

func getText() (*string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// we convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return &text, nil
}
