package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/event"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-net/v2/awsauth"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	robots_file_path string // path to the robots file
	api_url          string // elastic search api url
	sitemap_index    string // elastic search sitemap index
	scroll_timeout   string // elastic search scroll timeout
	scroll_size      int    // elastic search scroll size
	sitemap_path     string // path to the sitemap file
	zebedee_url      string // zebedee url
	fake_scroll      bool   // toggle to use or not the fake scroll implementation that replicates elastic search
	generate_sitemap bool   // generate the sitemap
	update_sitemap   bool   // updates the sitemap
}

func validConfig(flagFields *FlagFields) bool {
	log.Println("flag validation started")
	v := reflect.ValueOf(*flagFields)
	for i := 0; i < v.NumField(); i++ {
		flagTest := v.Field(i).String()
		if flagTest == "" {
			fmt.Println(v.Type().Field(i).Name + " is empty")
			return false
		} else {
			fmt.Println(v.Type().Field(i).Name + " = " + v.Field(i).String())
		}

	}
	log.Println("flag validation successful")
	return true
}

func validateCommandLines() (bool, *FlagFields) {
	commandline := FlagFields{}
	flag.StringVar(&commandline.robots_file_path, "robots-file-path", "test_robots.txt", "path to robots file")
	flag.StringVar(&commandline.sitemap_path, "sitemap-file-path", "test_sitemap", "path to sitemap file")
	flag.StringVar(&commandline.api_url, "api-url", "http://localhost", "elastic search api url")
	flag.StringVar(&commandline.zebedee_url, "zebedee-url", "http://localhost:8082", "zebedee url")
	flag.StringVar(&commandline.sitemap_index, "sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scroll_timeout, "scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scroll_size, "scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	flag.BoolVar(&commandline.fake_scroll, "enable-fake-scroll", true, "enable fake scroll")
	flag.BoolVar(&commandline.generate_sitemap, "generate-sitemap", false, "generate the sitemap")
	flag.BoolVar(&commandline.update_sitemap, "update-sitemap", false, "update the sitemap")

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
	valid, commandLine := validateCommandLines()
	if !valid {
		os.Exit(1)
	}

	cfg, err := config.Get()
	if err != nil {
		fmt.Println("Error retrieving config" + err.Error())
		os.Exit(1)
	}

	if commandLine.generate_sitemap {
		GenerateSitemap(cfg, commandLine)
	}

	if commandLine.update_sitemap {
		var scroll sitemap.Scroll
		if commandLine.fake_scroll {
			scroll = &FakeScroll{}
		} else {
			scroll = &sitemap.ElasticScroll{}
		}
		var store sitemap.FileStore
		if commandLine.fake_scroll {
			store = &sitemap.LocalStore{}
		} else {
			store = &sitemap.S3Store{}
		}
		zebedeeClient := zebedee.New(commandLine.zebedee_url)
		fetcher := sitemap.NewElasticFetcher(scroll, cfg, zebedeeClient)
		handler := event.NewContentPublishedHandler(store, zebedeeClient, cfg, fetcher)
		content, contentErr := getContent()
		fmt.Println(content)
		if contentErr != nil {
			fmt.Println("Failed to get event content from user:", err)
			return
		}

		err = handler.Handle(context.Background(), cfg, content)
		if err != nil {
			fmt.Println("Failed to handle event:", err)
			return
		}
	}

	GenerateRobotFile(cfg, commandLine)
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
		Addresses: []string{commandline.api_url},
		Transport: transport,
	})
	if err != nil {
		return
	}

	//Get zebedeeClient using arg -zebedee-url
	zebedeeClient := zebedee.New(commandline.zebedee_url)

	var scroll sitemap.Scroll
	if commandline.fake_scroll {
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
			config.English: commandline.sitemap_path + "_en",
			config.Welsh:   commandline.sitemap_path + "_cy",
		}),
	)

	//Generating sitemap
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
		config.English: commandline.robots_file_path,
	}

	store := &sitemap.LocalStore{}

	cfg.OpenSearchConfig.APIURL = commandline.api_url
	cfg.OpenSearchConfig.ScrollSize = commandline.scroll_size
	cfg.OpenSearchConfig.Signer = true

	body := robotFileWriter.GetRobotsFileBody(config.English, cfg.SitemapLocalFile)

	saveErr := store.SaveFile(commandline.robots_file_path, strings.NewReader(body))
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
