package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
)

var apiurl = flag.String("api-url", "OPENSEARCH_API_URL", "")

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Printf("Hello, %s!\n", strings.ToLower(*apiurl))

	robotseo.Init(assets.NewFromEmbeddedFilesystem())
	robotFileWriter := robotseo.RobotFileWriter{}

	cfg, err := config.Get()
	if err != nil {
		fmt.Println("Error retrieving config" + err.Error())
		os.Exit(1)
	}
	if wErr := robotFileWriter.WriteRobotsFile(cfg, []string{}); wErr != nil {
		fmt.Println("Error writing robot files", wErr.Error())
		return
	}

}
