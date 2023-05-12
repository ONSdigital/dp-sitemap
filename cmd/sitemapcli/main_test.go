package main

import (
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	// Define test cases
	Convey("when all the args filled", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "1000",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the args are valid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})
	//sitemap_index is missing
	Convey("when some args are missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
}

func TestGenerateSitemap(t *testing.T) {
	Convey("When config settings are given", func() {
		// Create a dummy config object
		cfg := &config.Config{
			OpenSearchConfig: config.OpenSearchConfig{
				APIURL: "http://localhost:9200",
			},
		}

		// Create dummy commandline flags
		commandline := &FlagFields{
			api_url:      "http://localhost:9200",
			zebedee_url:  "http://localhost:8082",
			sitemap_path: "/tmp/sitemap",
			fake_scroll:  true,
		}
		// Call the GenerateSitemap function
		GenerateSitemap(cfg, commandline)

		Convey("Then sitemap generation job complete", func() {
			So()
		})
	})

}
