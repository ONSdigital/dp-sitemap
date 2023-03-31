package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"

	//"testing"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
)

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	lang             string
	robots_file_path string
	api_url          string
	sitemap_index    string
	scroll_timeout   string
	scroll_size      int
}

// test function for FlagFields
// func test() bool {
// 	fmt.Println("flagtest start..") //put log entry")

// 	flagfields := FlagFields{
// 		lang:             "lang",
// 		robots_file_path: "/temp/robotfile.txt",
// 	}
// 	v := reflect.ValueOf(flagfields)
// 	for i := 0; i < v.NumField(); i++ {
// 		flagtest := v.Field(i).String()
// 		fmt.Println("flagtest is " + flagtest) //put log entry
// 		if flagtest == "" {
// 			fmt.Println("flagtest is empty")
// 			return false
// 		} else if flagtest == "robots-file-path" {
// 			fmt.Println("flagtest is robots-file-path")
// 			return true
// 		} else if flagtest == "lang" {
// 			fmt.Println("flagtest is lang")
// 			//check with language options
// 			return true
// 		} else {
// 			fmt.Println("flagtest is lang")
// 			return false
// 		}

// 	}
// 	fmt.Println("flagtest end..") //put log entry")

// 	return true

// }

// test function for FlagFields
func validConfig(flagfields *FlagFields) bool {
	fmt.Println("flagtest start..") //put log entry")

	v := reflect.ValueOf(*flagfields)
	//for i := 1; i < v.NumField(); i++ {
	for i := 0; i < v.NumField(); i++ {
		flagtest := v.Field(i).String()
		fmt.Println("flagtest is " + flagtest) //put log entry
		if flagtest == "" {
			fmt.Println(v.Type().Field(i).Name + " is empty")
			//logger
			return false
		} else {
			fmt.Println(v.Type().Field(i).Name + " = " + v.Field(i).String())
		}

	}
	fmt.Println("flagtest end..") //put log entry

	return true

}

// func (f *FlagFields) walkThru() error {

// 	v := reflect.ValueOf(f).Elem()
// 	fmt.Println(v)
// 	for i := 0; i < v.NumField(); i++ {
// 		found := false
// 		fmt.Println(f.robots_file_path)
// 		flag.Lookup(*&f.lang)
// 		flag.Visit(func(f *flag.Flag) {
// 			fmt.Println("heree")
// 			if f.Value.String() == *lang {
// 				return true
// 			}
// 		})
// 		return found
// 	}
// 	fmt.Println("walkthru end")
// 	return nil
// }

// func validateFlag(name string) bool {
// 	found := false
// 	flag.Visit(func(f *flag.Flag) {
// 		if f.Value.String() == *lang {
// 			found = true
// 		}
// 	})
// 	return found
// }

// config struct
// func validateFlagLooper(config string) {

// }

func main() {

	commandline := FlagFields{}
	flag.StringVar(&commandline.robots_file_path, "robots-file-path", "", "robotfile.txt")
	flag.StringVar(&commandline.lang, "lang", "eng", "eng")
	flag.StringVar(&commandline.api_url, "api-url", "", "OPENSEARCH_API_URL")
	flag.StringVar(&commandline.sitemap_index, "sitemap-index", "", "OPENSEARCH_SITEMAP_INDEX")
	flag.StringVar(&commandline.scroll_timeout, "scroll-timeout", "", "OPENSEARCH_SCROLL_TIMEOUT")
	flag.IntVar(&commandline.scroll_size, "scroll-size", 5, "OPENSEARCH_SCROLL_TIMEOUT")

	flag.Parse()
	if !validConfig(&commandline) {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
			fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
			flag.PrintDefaults()
		}
		return
	}

	//Creating dp_robot_file.txt
	robotseo.Init(assets.NewFromEmbeddedFilesystem())
	robotFileWriter := robotseo.RobotFileWriter{}

	cfg, err := config.Get()
	if err != nil {
		fmt.Println("Error retrieving config" + err.Error())
		os.Exit(1)
	}
	cfg.RobotsFilePath = commandline.robots_file_path

	if wErr := robotFileWriter.WriteRobotsFile(cfg, []string{}); wErr != nil {
		fmt.Println("Error writing robot files", wErr.Error())
		return
	}

	//Generate sitemap
	//ctx, cancel := context.WithTimeout(context.Background(), cfg.SitemapGenerationTimeout)

	saver := &mock.FileSaverMock{}
	fetcher := &mock.FetcherMock{}
	fetcher.HasWelshContentFunc = func(ctx context.Context, path string) bool { return false }

	// fetcher.GetFullSitemapFunc = func(ctx context.Context) ([]string, error) {
	// 	return []string{""}, errors.New("fetcher error")
	// }

	generator := sitemap.NewGenerator(fetcher, saver)

	genErr := generator.MakeFullSitemap(context.Background())
	if genErr != nil {
		fmt.Println("Error writing robot files", genErr.Error())
		return
	}
	fmt.Println("sitemap generation job complete")
	//log.Info(ctx, "sitemap generation job complete", log.Data{"last_run": job.LastRun(), "next_run": job.NextRun(), "run_count": job.RunCount()})

}

// func TestFlagFields_walkThru(t *testing.T) {
// 	f := FlagFields{value: ""}
// 	err := f.walkThru()
// 	if err == nil {
// 		t.Error("Expected error but got nil")
// 	}

// 	f2 := FlagFields{value: "some value"}
// 	err = f2.walkThru()
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// }
