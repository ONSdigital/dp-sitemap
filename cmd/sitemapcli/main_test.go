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
		}
		Convey("Then the args are valid", func() {
			result := validConfig(&testdata)
			So(result, ShouldBeTrue)
		})
	})
	Convey("when sitemap_index is missing", t, func() {
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
	Convey("when robots_file_path is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "",
			api_url:          "test.api.url",
			sitemap_index:    "2",
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
	Convey("when api_url is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "",
			sitemap_index:    "2",
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
	Convey("when scroll_timeout is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when zebedee_url is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "200",
			scroll_size:      2,
			zebedee_url:      "",
			sitemap_path:     "/path",
		}
		Convey("Then the args are invalid", func() {

			result := validConfig(&testdata)
			So(result, ShouldBeFalse)
		})
	})
	Convey("when sitemap_path is missing", t, func() {
		testdata := FlagFields{
			robots_file_path: "robot_file.txt",
			api_url:          "test.api.url",
			sitemap_index:    "2",
			scroll_timeout:   "",
			scroll_size:      2,
			zebedee_url:      "test.zeebedee.url",
			sitemap_path:     "",
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
