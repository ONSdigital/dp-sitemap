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
