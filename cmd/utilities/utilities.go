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
	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	RobotsFilePath string // path to the robots file
	ApiUrl         string // elastic search api url
	SitemapIndex   string // elastic search sitemap index
	ScrollTimeout  string // elastic search scroll timeout
	ScrollSize     int    // elastic search scroll size
	SitemapPath    string // path to the sitemap file
	ZebedeeUrl     string // zebedee url
	FakeScroll     bool   // toggle to use or not the fake scroll implementation that replicates elastic search
}

func GenerateSitemap(cfg *config.Config, commandline *FlagFields) {
	//Create local file store
	store := &sitemap.LocalStore{}

	//Get ElasticSearch Clients
	var transport http.RoundTripper = dphttp.DefaultTransport
	if cfg.OpenSearchConfig.Signer {
		//add SignerRegion,SignerService
		var err error
		transport, err = awsauth.NewAWSSignerRoundTripper(cfg.OpenSearchConfig.APIURL, cfg.OpenSearchConfig.SignerFilename, cfg.OpenSearchConfig.SignerRegion, cfg.OpenSearchConfig.SignerService, awsauth.Options{TlsInsecureSkipVerify: cfg.OpenSearchConfig.TLSInsecureSkipVerify})
		if err != nil {
			fmt.Printf("failed to save file")
			return
		}
	}

	//Get rawClient using arg -api-url
	rawClient, err := es710.NewClient(es710.Config{
		Addresses: []string{commandline.ApiUrl},
		Transport: transport,
	})
	if err != nil {
		return
	}

	//Get zebedeeClient using arg -zebedee-url
	zebedeeClient := zebedee.New(commandline.ZebedeeUrl)

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
			config.English: commandline.SitemapPath + "_en",
			config.Welsh:   commandline.SitemapPath + "_cy",
		}),
	)

	//Generating sitemap
	genErr := generator.MakeFullSitemap(context.Background())
	if genErr != nil {
		fmt.Println("Error writing sitemap file", genErr.Error())
		return
	}

	GenerateRobotFile(cfg, commandline)

	fmt.Println("sitemap generation job complete")
}

func GenerateRobotFile(cfg *config.Config, commandline *FlagFields) {

	robotseo.Init(assets.NewFromEmbeddedFilesystem())
	robotFileWriter := robotseo.RobotFileWriter{}
	cfg.RobotsFilePath = map[config.Language]string{
		config.English: commandline.RobotsFilePath,
	}

	store := &sitemap.LocalStore{}

	cfg.OpenSearchConfig.APIURL = commandline.ApiUrl
	cfg.OpenSearchConfig.ScrollSize = commandline.ScrollSize
	cfg.OpenSearchConfig.Signer = true

	body := robotFileWriter.GetRobotsFileBody(config.English, cfg.SitemapLocalFile)

	saveErr := store.SaveFile(commandline.RobotsFilePath, strings.NewReader(body))
	if saveErr != nil {
		fmt.Println("failed to save file")
		return
	}
}

func UpdateSitemap(cfg *config.Config, commandLine *FlagFields) {

	var scroll sitemap.Scroll
	if commandLine.FakeScroll {
		scroll = &FakeScroll{}
	} else {
		scroll = &sitemap.ElasticScroll{}
	}
	var store sitemap.FileStore
	if commandLine.FakeScroll {
		store = &sitemap.LocalStore{}
	} else {
		store = &sitemap.S3Store{}
	}
	zebedeeClient := zebedee.New(commandLine.ZebedeeUrl)
	fetcher := sitemap.NewElasticFetcher(scroll, cfg, zebedeeClient)
	handler := event.NewContentPublishedHandler(store, zebedeeClient, cfg, fetcher)
	content, contentErr := getContent()
	if contentErr != nil {
		fmt.Println("Failed to get event content from user:", contentErr)
		return
	}

	err := handler.Handle(context.Background(), cfg, content)
	if err != nil {
		fmt.Println("Failed to handle event:", err)
		return
	}

	GenerateRobotFile(cfg, commandLine)
	fmt.Println("sitemap update job complete")

}

func getContent() (*event.ContentPublished, error) {
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
