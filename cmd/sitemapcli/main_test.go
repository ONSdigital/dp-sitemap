package main

import (
	"testing"

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
			fake_scroll:      true,
		}
		Convey("Then the args are valid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})

	//fakescroll is missing
	Convey("when fake scroll is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the fake_scroll is null", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeTrue) // Check if this is the requirement
		})
	})

	//fake scroll is false
	Convey("when fake scroll is false", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
			fake_scroll:      false,
		}
		Convey("Then the fake_scroll is false", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})

	//sitemap_index is missing
	Convey("when sitemap index missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
			fake_scroll:      true,
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)

			So(result, ShouldBeFalse) //
		})
	})
}

// func TestGenerateSitemap(t *testing.T) {
// 	t.Log("test")
// 	Convey("When config settings are given", func() {
// 		// Create a dummy config object
// 		cfg := &config.Config{
// 			OpenSearchConfig: config.OpenSearchConfig{
// 				APIURL: "http://localhost:9200",
// 			},
// 		}

// 		// Create dummy commandline flags
// 		commandline := &FlagFields{
// 			api_url:      "http://localhost:9200",
// 			zebedee_url:  "http://localhost:8082",
// 			sitemap_path: "/tmp/sitemap",
// 			fake_scroll:  true,
// 		}

// 		// Run the test using the Convey package
// 		convey.Convey("Generate sitemap", t, func() {
// 			convey.Convey("When given valid input", func() {

// 				convey.Convey("Then it should generate sitemap", func() {
// 					GenerateSitemap(cfg, commandline)
// 					expectedOutput := "sitemap generation job complete\n"
// 					convey.So(output, convey.ShouldEqual, expectedOutput)
// 				})
// 			})
// 		})
// 	})

//Check if return is needed for the function GenerateSitemap
