package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"github.com/ONSdigital/dp-sitemap/event"
	"io"
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

type StaticURLEn struct {
	URL         string `json:"url"`
	ReleaseDate string `json:"releaseDate"`
	Cy          bool   `json:"cy"`
}

type StaticURLCy struct {
	URL         string `json:"url"`
	ReleaseDate string `json:"releaseDate"`
	En          bool   `json:"en"`
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
		err = loadStaticSitemap(context.Background(), cfg, "test_sitemap_en", &sitemap.LocalStore{}, config.English)
		if err != nil {
			fmt.Println("Failed to load english static sitemap:", err)
			return
		}
		err = loadStaticSitemap(context.Background(), cfg, "test_sitemap_cy", &sitemap.LocalStore{}, config.Welsh)
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

func loadStaticSitemap(ctx context.Context, cfg *config.Config, sitemapName string, store sitemap.FileStore, lang config.Language) error {
	efs := assets.NewFromEmbeddedFilesystem()

	staticSitemaps := make(map[config.Language]string)
	staticSitemaps[config.English] = "sitemap/sitemap_en.json"
	staticSitemaps[config.Welsh] = "sitemap/sitemap_cy.json"

	b, err := efs.Get(ctx, staticSitemaps[lang])
	if err != nil {
		panic("can't find file " + staticSitemaps[lang])
	}

	var contentEn []StaticURLEn
	var contentCy []StaticURLCy

	if lang == config.English {
		err = json.Unmarshal(b, &contentEn)
		if err != nil {
			return fmt.Errorf("unable to read json: %w", err)
		}
	} else {
		err = json.Unmarshal(b, &contentCy)
		if err != nil {
			return fmt.Errorf("unable to read json: %w", err)
		}
	}

	// get the old sitemap
	oldSitemapFile, err := store.GetFile(sitemapName)
	if err != nil {
		return fmt.Errorf("unable to get file: %w", err)
	}
	defer oldSitemapFile.Close()

	sitemapReader := sitemap.UrlsetReader{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}
	decoder := xml.NewDecoder(oldSitemapFile)
	err = decoder.Decode(&sitemapReader)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to decode old sitemap: %w", err)
		}
	}

	// move old sitemap urls to new sitemap
	sitemapWriter := sitemap.Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}

	for _, url := range sitemapReader.URL {
		var u sitemap.URL
		u.Loc = url.Loc
		u.Lastmod = url.Lastmod
		u.Alternate = &sitemap.AlternateURL{}
		if url.Alternate != nil {
			u.Alternate.Rel = url.Alternate.Rel
			u.Alternate.Link = url.Alternate.Link
			u.Alternate.Lang = url.Alternate.Lang
		}
		sitemapWriter.URL = append(sitemapWriter.URL, u)
	}

	// range through static content
	if lang == config.English {
		for _, item := range contentEn {
			var u sitemap.URL
			u.Loc = cfg.DpOnsURLHostNameEn + item.URL
			u.Lastmod = item.ReleaseDate
			u.Alternate = &sitemap.AlternateURL{}
			if item.Cy == true {
				u.Alternate.Rel = "alternate"
				u.Alternate.Link = cfg.DpOnsURLHostNameCy + item.URL
				u.Alternate.Lang = "cy"
			}
			sitemapWriter.URL = append(sitemapWriter.URL, u)
		}
	} else {
		for _, item := range contentCy {
			var u sitemap.URL
			u.Loc = cfg.DpOnsURLHostNameEn + item.URL
			u.Lastmod = item.ReleaseDate
			u.Alternate = &sitemap.AlternateURL{}
			if item.En == true {
				u.Alternate.Rel = "alternate"
				u.Alternate.Link = cfg.DpOnsURLHostNameEn + item.URL
				u.Alternate.Lang = "en"
			}
			sitemapWriter.URL = append(sitemapWriter.URL, u)
		}
	}

	marshaledContent, err := xml.MarshalIndent(sitemapWriter, "", "  ")
	if err != nil {
		return err
	}
	header := []byte(xml.Header)
	header = append(header, marshaledContent...)
	reader := bytes.NewReader(header)
	err = store.SaveFile(sitemapName, reader)
	if err != nil {
		return err
	}
	return nil
}
