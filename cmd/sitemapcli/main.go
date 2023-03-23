package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	//"testing"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
)

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	lang             string
	robots_file_path string
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
func test() bool {
	fmt.Println("flagtest start..") //put log entry")

	flagfields := FlagFields{
		lang:             "lang",
		robots_file_path: "/temp/robotfile.txt",
	}
	v := reflect.ValueOf(flagfields)
	//for i := 1; i < v.NumField(); i++ {
	for i := 1; i < 2; i++ {
		flagtest := v.Field(i).String()
		fmt.Println("flagtest is " + flagtest) //put log entry
		if flagtest == "" {
			fmt.Println("flagtest is empty")
			return false
		} else if flagtest == "robots-file-path" {
			fmt.Println("flagtest is robots-file-path" + flagtest)
			return true
		} else if flagtest == "lang" {
			fmt.Println("flagtest is lang" + flagtest)
			//check with language options
			return true
		} else {
			fmt.Println("flag field is incorrect)" + flagtest)
			return false
		}

	}
	fmt.Println("flagtest end..") //put log entry")

	return true

}
func (f *FlagFields) walkThru() {

	//v := reflect.ValueOf(f).Elem()
	fmt.Println(reflect.ValueOf(f).Elem())
	for i := 0; i < v.NumField(); i++ {
		found := false
		flag.Visit(func(f *flag.Flag) {
			if f.Value.String() == name {
				found = true
			}
		})
		return found
	}

	return nil
}

var robotsfilepath = flag.String("robots-file-path", "", "robot.txt_PATH")
var lang = flag.String("lang", "eng", "eng")

func validateFlag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f. == name {
			found = true
		}
	})
	return found
}

// config struct
// func validateFlagLooper(config string) {

// }

func main() {

	test()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	walkThru(flag.Flag)
	//fmt.Println(reflect.TypeOf(os.Args))
	//fmt.Println(walkThru(flag))
	//argList := reflect.
	// if !validateFlag("robots-file-path") {
	// 	fmt.Println("robot.txt file path not found")
	// 	flag.PrintDefaults()
	// 	return
	// }
	// if !validateFlag("lang") {
	// 	fmt.Println("language is not definsssed")
	// 	flag.PrintDefaults()
	// 	return
	// }
	fmt.Printf("Hello, %s!\n", strings.ToLower(*robotsfilepath))

	//Creating dp_robot_file.txt
	robotseo.Init(assets.NewFromEmbeddedFilesystem())
	robotFileWriter := robotseo.RobotFileWriter{}

	cfg, err := config.Get()
	if err != nil {
		fmt.Println("Error retrieving config" + err.Error())
		os.Exit(1)
	}
	cfg.RobotsFilePath = *robotsfilepath

	if wErr := robotFileWriter.WriteRobotsFile(cfg, []string{}); wErr != nil {
		fmt.Println("Error writing robot files", wErr.Error())
		return
	}

	// //Generate sitemap
	// ctx, cancel := context.WithTimeout(context.Background(), cfg.SitemapGenerationTimeout)

	// generator := sitemap.NewGenerator(
	// 	sitemap.NewElasticFetcher(
	// 		config.Get()
	// 	),
	// 	saver,
	// )

	// genErr := generator.MakeFullSitemap(ctx)
	// if genErr != nil {
	// 	log.Error(ctx, "failed to generate sitemap", genErr)
	// 	return
	// }
	// log.Info(ctx, "sitemap generation job complete", log.Data{"last_run": job.LastRun(), "next_run": job.NextRun(), "run_count": job.RunCount()})

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
