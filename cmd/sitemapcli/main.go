package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-net/v2/awsauth"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/cmd"
	"github.com/ONSdigital/dp-sitemap/cmd/utilities"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

// FlagFields Config represents service configurations for dp-sitemap
type FlagFields struct {
	robotsFilePath  string // path to the robots file
	apiURL          string // elastic search api url
	sitemapIndex    string // elastic search sitemap index
	scrollTimeout   string // elastic search scroll timeout
	scrollSize      int    // elastic search scroll size
	sitemapPath     string // path to the sitemap file
	zebedeeURL      string // zebedee url
	fakeScroll      bool   // toggle to use or not the fake scroll implementation that replicates elastic search
	generateSitemap bool   // generate the sitemap
	updateSitemap   bool   // updates the sitemap
}

func validConfig(flagFields *FlagFields) bool {
	log.Println("flag validation started")
	v := reflect.ValueOf(*flagFields)
	for i := 0; i < v.NumField(); i++ {
		flagTest := v.Field(i).String()
		if flagTest == "" {
			fmt.Println(v.Type().Field(i).Name + " is empty")
			return false
		}
		fmt.Println(v.Type().Field(i).Name + " = " + v.Field(i).String())
	}
	log.Println("flag validation successful")
	return true
}

func validateCommandLines() (bool, *FlagFields) {
	commandline := FlagFields{}
	flag.StringVar(&commandline.robotsFilePath, "robots-file-path", "test_robots.txt", "path to robots file")
	flag.StringVar(&commandline.sitemapPath, "sitemap-file-path", "test_sitemap", "path to sitemap file")
	flag.StringVar(&commandline.apiURL, "api-url", "http://localhost", "elastic search api url")
	flag.StringVar(&commandline.zebedeeURL, "zebedee-url", "http://localhost:8082", "zebedee url")
	flag.StringVar(&commandline.sitemapIndex, "sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scrollTimeout, "scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scrollSize, "scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	flag.BoolVar(&commandline.fakeScroll, "enable-fake-scroll", true, "enable fake scroll")
	flag.BoolVar(&commandline.generateSitemap, "generate-sitemap", false, "generate the sitemap")
	flag.BoolVar(&commandline.updateSitemap, "update-sitemap", false, "update the sitemap")

	flag.Parse()
	if !validConfig(&commandline) {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
			fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
			flag.PrintDefaults()
		}
		return false, nil
	}
	return true, &commandline
}
func main() {
	cmd.Execute()
}

func GenerateSitemap(cfg *config.Config, commandline *FlagFields) {
	// Create local file store
	store := &sitemap.LocalStore{}

	// Get ElasticSearch Clients
	var transport http.RoundTripper = dphttp.DefaultTransport
	if cfg.OpenSearchConfig.Signer {
		// add SignerRegion,SignerService
		var err error
		transport, err = awsauth.NewAWSSignerRoundTripper(cfg.OpenSearchConfig.APIURL, cfg.OpenSearchConfig.SignerFilename, cfg.OpenSearchConfig.SignerRegion, cfg.OpenSearchConfig.SignerService, awsauth.Options{TlsInsecureSkipVerify: cfg.OpenSearchConfig.TLSInsecureSkipVerify})
		if err != nil {
			fmt.Printf("failed to save file")
			return
		}
	}

	// Get rawClient using arg -api-url
	rawClient, err := es710.NewClient(es710.Config{
		Addresses: []string{commandline.apiURL},
		Transport: transport,
	})
	if err != nil {
		return
	}

	// Get zebedeeClient using arg -zebedee-url
	zebedeeClient := zebedee.New(commandline.zebedeeURL)

	var scroll sitemap.Scroll
	if commandline.fakeScroll {
		scroll = utilities.NewFakeScroll()
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
			config.English: commandline.sitemapPath + "_en",
			config.Welsh:   commandline.sitemapPath + "_cy",
		}),
	)

	// Generating sitemap
	genErr := generator.MakeFullSitemap(context.Background())
	if genErr != nil {
		fmt.Println("Error writing sitemap file", genErr.Error())
		return
	}
	fmt.Println("sitemap generation job complete")
}
func GenerateRobotFile(cfg *config.Config, commandline *FlagFields) {
	robotseo.Init(assets.NewFromEmbeddedFilesystem())
	robotFileWriter := robotseo.RobotFileWriter{}
	cfg.RobotsFilePath = map[config.Language]string{
		config.English: commandline.robotsFilePath,
	}

	store := &sitemap.LocalStore{}

	cfg.OpenSearchConfig.APIURL = commandline.apiURL
	cfg.OpenSearchConfig.ScrollSize = commandline.scrollSize
	cfg.OpenSearchConfig.Signer = true

	body := robotFileWriter.GetRobotsFileBody(config.English, cfg.SitemapLocalFile)

	saveErr := store.SaveFile(commandline.robotsFilePath, strings.NewReader(body))
	if saveErr != nil {
		fmt.Println("failed to save file")
		return
	}
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

// Delete this func?
// func mainold() {
// 	valid, commandLine := validateCommandLines()
// 	if !valid {
// 		os.Exit(1)
// 	}

// 	cfg, err := config.Get()
// 	if err != nil {
// 		fmt.Println("Error retrieving config" + err.Error())
// 		os.Exit(1)
// 	}

// 	if commandLine.generate_sitemap {
// 		GenerateSitemap(cfg, commandLine)
// 	}

// 	if commandLine.update_sitemap {
// 		var scroll sitemap.Scroll
// 		if commandLine.fake_scroll {
// 			scroll = &utilities.FakeScroll{}
// 		} else {
// 			scroll = &sitemap.ElasticScroll{}
// 		}
// 		var store sitemap.FileStore
// 		if commandLine.fake_scroll {
// 			store = &sitemap.LocalStore{}
// 		} else {
// 			store = &sitemap.S3Store{}
// 		}
// 		zebedeeClient := zebedee.New(commandLine.zebedee_url)
// 		fetcher := sitemap.NewElasticFetcher(scroll, cfg, zebedeeClient)
// 		handler := event.NewContentPublishedHandler(store, zebedeeClient, cfg, fetcher)
// 		content, contentErr := getContent()
// 		fmt.Println(content)
// 		if contentErr != nil {
// 			fmt.Println("Failed to get event content from user:", err)
// 			return
// 		}

// 		err = handler.Handle(context.Background(), cfg, content)
// 		if err != nil {
// 			fmt.Println("Failed to handle event:", err)
// 			return
// 		}
// 	}

// 	GenerateRobotFile(cfg, commandLine)
// }
