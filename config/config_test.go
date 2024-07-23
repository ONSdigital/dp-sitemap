package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.BindAddr, ShouldEqual, "localhost:")
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.KafkaConfig.Brokers, ShouldHaveLength, 1)
				So(cfg.KafkaConfig.Brokers[0], ShouldEqual, "localhost:9092")
				So(cfg.KafkaConfig.Version, ShouldEqual, "1.0.2")
				So(cfg.KafkaConfig.SecProtocol, ShouldEqual, "")
				So(cfg.KafkaConfig.NumWorkers, ShouldEqual, 1)
				So(cfg.KafkaConfig.ContentUpdatedGroup, ShouldEqual, "dp-sitemap")
				So(cfg.KafkaConfig.ContentUpdatedTopic, ShouldEqual, "content-updated")
				So(cfg.OpenSearchConfig.ElasticSearchURL, ShouldEqual, "http://localhost:11200")
				So(cfg.OpenSearchConfig.ElasticSearchIndex, ShouldEqual, "ons")
				So(cfg.OpenSearchConfig.ScrollTimeout, ShouldEqual, time.Minute)
				So(cfg.OpenSearchConfig.ScrollSize, ShouldEqual, 10000)
				So(cfg.OpenSearchConfig.DebugFirstPageOnly, ShouldEqual, false)
				So(cfg.OpenSearchConfig.SignerFilename, ShouldEqual, "")
				So(cfg.OpenSearchConfig.SignerProfile, ShouldEqual, "")
				So(cfg.OpenSearchConfig.SignerRegion, ShouldEqual, "eu-west-2")
				So(cfg.OpenSearchConfig.SignerService, ShouldEqual, "es")
				So(cfg.OpenSearchConfig.TLSInsecureSkipVerify, ShouldEqual, false)
				So(cfg.OpenSearchConfig.ElasticSearchURL, ShouldEqual, "http://localhost:11200")
				So(cfg.ZebedeeURL, ShouldEqual, "http://localhost:8082")
				So(cfg.DpOnsURLHostNameEn, ShouldEqual, "https://dp.aws.onsdigital.uk/")
				So(cfg.DpOnsURLHostNameCy, ShouldEqual, "https://cy.dp.aws.onsdigital.uk/")
				So(cfg.SitemapSaveLocation, ShouldEqual, "local")
				So(cfg.SitemapLocalFile[English], ShouldEqual, "/tmp/dp-sitemap-en.xml")
				So(cfg.SitemapLocalFile[Welsh], ShouldEqual, "/tmp/dp-sitemap-cy.xml")
				So(cfg.PublishingSitemapLocalFile, ShouldEqual, "/tmp/dp-publishing-sitemap.xml")
				So(cfg.PublishingSitemapMaxSize, ShouldEqual, 500)
				So(cfg.S3Config.SitemapFileKey[English], ShouldEqual, "sitemap-en")
				So(cfg.S3Config.SitemapFileKey[Welsh], ShouldEqual, "sitemap-cy")
				So(cfg.S3Config.PublishingSitemapFileKey, ShouldEqual, "publishing-sitemap")
				So(cfg.RobotsFilePath, ShouldNotBeEmpty)
				So(cfg.Debug, ShouldBeTrue)
			})

			Convey("Then a second call to config should return the same config", func() {
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
