package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/event"
	"net/http"
	"os"
	"reflect"
	"strings"

	//"testing"

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
	robots_file_path string
	api_url          string
	sitemap_index    string
	scroll_timeout   string
	scroll_size      int
	sitemap_path     string
	zebedee_url      string
	fake_scroll      bool
}

type StaticURL struct {
	URL         string `json:"url"`
	ReleaseDate string `json:"releaseDate"`
	HasAltLang  bool   `json:"hasAltLang"`
}

// test function for FlagFields
func validConfig(flagfields *FlagFields) bool {
	fmt.Println("flag validation started..") //put log entry

	v := reflect.ValueOf(*flagfields)
	for i := 0; i < v.NumField(); i++ {
		flagtest := v.Field(i).String()
		if flagtest == "" {
			fmt.Println(v.Type().Field(i).Name + " is empty")
			return false
		} else {
			fmt.Println(v.Type().Field(i).Name + " = " + v.Field(i).String())
		}

	}
	fmt.Println("flagtest validation succesfull..") //put log entry

	return true

}

func validateCommandLines() (bool, *FlagFields) {

	commandline := FlagFields{}
	flag.StringVar(&commandline.robots_file_path, "robots-file-path", "test_robots.txt", "robotfile.txt")
	flag.StringVar(&commandline.sitemap_path, "sitemap-file-path", "test_sitemap", "sitemap.xml")
	flag.StringVar(&commandline.api_url, "api-url", "http://localhost", "")
	flag.StringVar(&commandline.zebedee_url, "zebedee-url", "http://localhost:8082", "")
	flag.StringVar(&commandline.sitemap_index, "sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scroll_timeout, "scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scroll_size, "scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	flag.BoolVar(&commandline.fake_scroll, "enable-fake-scroll", true, "enable fake scroll")

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

	choice, err := menu()
	if err != nil {
		fmt.Println("Error retrieving user choice:", err)
		os.Exit(1)
	}
	switch choice {
	case 1:
		GenerateSitemap(cfg, commandLine)
	case 2:
		zebedeeClient := zebedee.New(commandLine.zebedee_url)
		fetcher := sitemap.NewElasticFetcher(&FakeScroll{}, cfg, zebedeeClient)
		handler := event.NewContentPublishedHandler(&sitemap.LocalStore{}, zebedeeClient, cfg, fetcher)
		content := &event.ContentPublished{
			URI:          "economy/environmentalaccounts/articles/testarticle3",
			DataType:     "theDateType",
			CollectionID: "theCollectionId",
			JobID:        "theJobId",
			SearchIndex:  "theSearchIndex",
			TraceID:      "theTraceId",
		}

		err := handler.Handle(context.Background(), cfg, content)
		if err != nil {
			fmt.Println("Failed to handle event:", err)
			return
		}
	case 3:
		err = loadStaticSitemap(context.Background(), "test_sitemap_en", "sitemap_en.json", cfg.DpOnsURLHostNameEn, cfg.DpOnsURLHostNameCy, "cy", &sitemap.LocalStore{})
		if err != nil {
			fmt.Println("Failed to load english static sitemap:", err)
			return
		}
		err = loadStaticSitemap(context.Background(), "test_sitemap_cy", "sitemap_cy.json", cfg.DpOnsURLHostNameCy, cfg.DpOnsURLHostNameEn, "en", &sitemap.LocalStore{})
		if err != nil {
			fmt.Println("Failed to load welsh static sitemap:", err)
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
		fmt.Println("Error writing sitemappp file", genErr.Error())
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

	body := robotFileWriter.GetRobotsFileBody(config.Language(config.English), cfg.SitemapLocalFile)
	fmt.Printf("robot file path is: %v", commandline.robots_file_path)

	saveErr := store.SaveFile(commandline.robots_file_path, strings.NewReader(body))

	if saveErr != nil {
		fmt.Println("failed to save file")
		return
	}
}

func menu() (int, error) {
	var i = 0
	for i < 1 || i > 3 {
		fmt.Println("*** Menu ***")
		fmt.Println("1. Generate sitemap")
		fmt.Println("2. Update sitemap")
		fmt.Println("3. Load static sitemap")
		fmt.Print("Choice: ")
		if _, err := fmt.Scan(&i); err != nil {
			return 0, err
		}
	}
	return i, nil
}

func loadStaticSitemap(ctx context.Context, oldSitemapName, staticSitemapName, DpOnsURLHostName, DpOnsURLHostNameAlt, altLang string, store sitemap.FileStore) error {
	efs := assets.NewFromEmbeddedFilesystem()

	b, err := efs.Get(ctx, assets.Sitemap, staticSitemapName)
	if err != nil {
		panic("can't find file " + staticSitemapName)
	}

	var content []StaticURL

	err = json.Unmarshal(b, &content)
	if err != nil {
		return fmt.Errorf("unable to read json: %w", err)
	}

	// move old sitemap urls to new sitemap
	sitemapWriter := sitemap.Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}

	// range through static content
	for _, item := range content {
		var newURL sitemap.URL
		newURL.Loc = DpOnsURLHostName + item.URL
		newURL.Lastmod = item.ReleaseDate
		newURL.Alternate = &sitemap.AlternateURL{}
		if item.HasAltLang == true {
			newURL.Alternate.Rel = "alternate"
			newURL.Alternate.Link = DpOnsURLHostNameAlt + item.URL
			newURL.Alternate.Lang = altLang
		}
		sitemapWriter.URL = append(sitemapWriter.URL, newURL)
	}

	marshaledContent, err := xml.MarshalIndent(sitemapWriter, "", "  ")
	if err != nil {
		return err
	}
	header := []byte(xml.Header)
	header = append(header, marshaledContent...)
	reader := bytes.NewReader(header)
	err = store.SaveFile(oldSitemapName, reader)
	if err != nil {
		return err
	}
	return nil
}
