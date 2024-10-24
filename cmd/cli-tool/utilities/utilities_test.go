package utilities

import (
	"errors"
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	defaultTestDataType = "abc"
)

func TestCreateCliSitemapGenerator(t *testing.T) {
	Convey("Given valid config and command line flags/Fake scroll is True", t, func() {
		cfg, _ := config.Get()
		commandline := &FlagFields{
			RobotsFilePath:     "robot_file.txt",
			ElasticSearchURL:   "http://localhost",
			ScrollTimeout:      "1000",
			ScrollSize:         2,
			ZebedeeURL:         "http://localhost:8082",
			SitemapPath:        "test_sitemap",
			FakeScroll:         true,
			ElasticSearchIndex: "1",
		}
		cfg.OpenSearchConfig.Signer = true
		So(cfg.OpenSearchConfig.Signer, ShouldBeTrue)

		Convey("When OpenSearchConfig.Signer is true and no errors", func() {
			generator, err := createCliSitemapGenerator(cfg, commandline)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
				So(generator, ShouldNotBeNil)
			})
		})

		Convey("When Fakescroll is true", func() {
			So(commandline.FakeScroll, ShouldBeTrue)
			generator, err := createCliSitemapGenerator(cfg, commandline)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
				So(generator, ShouldNotBeNil)
			})
		})

		commandline.FakeScroll = false
		So(commandline.FakeScroll, ShouldBeFalse)

		Convey("When Fakescroll is false", func() {
			generator, err := createCliSitemapGenerator(cfg, commandline)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
				So(generator, ShouldNotBeNil)
			})
		})
	})
}

func TestUpdateSitemap(t *testing.T) {
	Convey("Given a valid config and command line flags", t, func() {
		cfg := &config.Config{}
		commandLine := &FlagFields{}

		Convey("When FakeScroll is true", func() {
			commandLine.FakeScroll = true

			Convey("When getContent does not return an error", func() {
				getContent = func() (*event.ContentPublished, error) {
					var cont event.ContentPublished
					cont.URI = "1"
					cont.CollectionID = "1"
					cont.DataType = defaultTestDataType
					cont.JobID = "1"
					cont.SearchIndex = "2"
					cont.TraceID = "1"
					return &cont, nil
				}
				err := UpdateSitemap(cfg, commandLine)
				Convey("UpdateSitemap should not return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("When getContent returns an error", func() {
				getContent = func() (*event.ContentPublished, error) {
					return nil, errors.New("Error")
				}
				err := UpdateSitemap(cfg, commandLine)

				Convey("UpdateSitemap should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When FakeScroll is false", func() {
			commandLine.FakeScroll = false

			Convey("When getContent does not return an error", func() {
				getContent = func() (*event.ContentPublished, error) {
					var cont event.ContentPublished
					cont.URI = "1"
					cont.CollectionID = "1"
					cont.DataType = defaultTestDataType
					cont.JobID = "1"
					cont.SearchIndex = "2"
					cont.TraceID = "1"
					return &cont, nil
				}
				err := UpdateSitemap(cfg, commandLine)
				Convey("UpdateSitemap should not return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("When getContent returns an error", func() {
				getContent = func() (*event.ContentPublished, error) {
					return nil, errors.New("Error")
				}
				err := UpdateSitemap(cfg, commandLine)

				Convey("UpdateSitemap should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
