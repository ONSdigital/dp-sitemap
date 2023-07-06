package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidConfig(t *testing.T) {
	// Define test cases
	Convey("when all the args filled", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "test.api.url",
			sitemapIndex:   "2",
			scrollTimeout:  "1000",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "/path",
		}
		Convey("Then the args are valid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})
	Convey("when sitemap_index is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "test.api.url",
			sitemapIndex:   "",
			scrollTimeout:  "200",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when robots_file_path is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "",
			apiURL:         "test.api.url",
			sitemapIndex:   "2",
			scrollTimeout:  "200",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when api_url is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "",
			sitemapIndex:   "2",
			scrollTimeout:  "200",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when scroll_timeout is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "test.api.url",
			sitemapIndex:   "2",
			scrollTimeout:  "",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when zebedee_url is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "test.api.url",
			sitemapIndex:   "2",
			scrollTimeout:  "200",
			scrollSize:     2,
			zebedeeURL:     "",
			sitemapPath:    "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when sitemap_path is missing", t, func() {
		testdata := FlagFields{
			robotsFilePath: "robot_file.txt",
			apiURL:         "test.api.url",
			sitemapIndex:   "2",
			scrollTimeout:  "",
			scrollSize:     2,
			zebedeeURL:     "test.zeebedee.url",
			sitemapPath:    "",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when all args are missing", t, func() {
		testdata := FlagFields{}
		Convey("Then the args are invalid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
}
