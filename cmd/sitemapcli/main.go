package main

import (
	"context"
	"flag"
	"fmt"
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

func validateCommandLines() bool {

	commandline := FlagFields{}
	flag.StringVar(&commandline.robots_file_path, "robots-file-path", "robot_file.txt", "robotfile.txt")
	flag.StringVar(&commandline.sitemap_path, "sitemap-file-path", "sitemap.xml", "sitemap.xml")
	flag.StringVar(&commandline.api_url, "api-url", "", "")
	flag.StringVar(&commandline.zebedee_url, "zebedee-url", "", "")
	flag.StringVar(&commandline.sitemap_index, "sitemap-index", "", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scroll_timeout, "scroll-timeout", "", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scroll_size, "scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")

	flag.Parse()
	if !validConfig(&commandline) {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
			fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
			flag.PrintDefaults()
		}
		return false
	} else {
		return true
	}
}
func main() {
	//validate commandline - reduce this part as well
	commandline := FlagFields{}
	flag.StringVar(&commandline.robots_file_path, "robots-file-path", "robot_file.txt", "robotfile.txt")
	flag.StringVar(&commandline.sitemap_path, "sitemap-file-path", "sitemap.xml", "sitemap.xml")
	flag.StringVar(&commandline.api_url, "api-url", "", "")
	flag.StringVar(&commandline.zebedee_url, "zebedee-url", "", "")
	flag.StringVar(&commandline.sitemap_index, "sitemap-index", "", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scroll_timeout, "scroll-timeout", "", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scroll_size, "scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")

	flag.Parse()
	if !validConfig(&commandline) {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
			fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
			flag.PrintDefaults()
		}
		return
	}
	cfg, err := config.Get()

	if err != nil {
		fmt.Println("Error retrieving config" + err.Error())
		os.Exit(1)
	} else {

		GenerateSitemap(cfg, err, commandline)

		//create robot.txt file
		GenerateRobotFile(cfg, commandline)
	}

}

func GenerateSitemap(cfg *config.Config, err error, commandline FlagFields) {
	//Create local file store
	store := &sitemap.LocalStore{}

	//Get ElasticSearch Clients
	var transport http.RoundTripper = dphttp.DefaultTransport
	if cfg.OpenSearchConfig.Signer {
		//add SignerRegion,SignerService
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

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			rawClient,
			cfg,
			zebedeeClient,
		)),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
		sitemap.WithFileStore(store),
		sitemap.WithFullSitemapFiles(map[config.Language]string{
			config.English: commandline.sitemap_path + "_eng",
			config.Welsh:   commandline.sitemap_path + "_welsh",
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

func GenerateRobotFile(cfg *config.Config, commandline FlagFields) {

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
