package steps

import (
	"context"
	"log"
	"net/http"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/service"
	"github.com/ONSdigital/dp-sitemap/service/mock"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/cucumber/godog"
	es710 "github.com/elastic/go-elasticsearch/v7"

	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpEsMock "github.com/ONSdigital/dp-elasticsearch/v3/client/mocks"
)

type Component struct {
	componenttest.ErrorFeature
	serviceList       *service.ExternalServiceList
	KafkaConsumer     kafka.IConsumerGroup
	EsClient          *es710.Client
	EsIndex           *godog.Table
	S3UploadedSitemap map[string]string
	killChannel       chan os.Signal
	apiFeature        *componenttest.APIFeature
	errorChan         chan error
	svc               *service.Service
	cfg               *config.Config
	files             map[string]string
	welshVersion      map[string]bool
}

func testCfg() config.Config {
	return config.Config{
		KafkaConfig: config.KafkaConfig{
			Brokers:             []string{"localhost:9092", "localhost:9093"},
			ContentUpdatedGroup: "dp-sitemap",
			ContentUpdatedTopic: "content-updated",
		},
	}
}

func NewComponent(ctx context.Context) *Component {
	c := &Component{errorChan: make(chan error)}

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	kafkaOffset := kafka.OffsetOldest
	consumer, err := kafka.NewConsumerGroup(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:  cfg.KafkaConfig.Brokers,
			Topic:        cfg.KafkaConfig.ContentUpdatedTopic,
			GroupName:    cfg.KafkaConfig.ContentUpdatedGroup,
			KafkaVersion: &cfg.KafkaConfig.Version,
			Offset:       &kafkaOffset,
		},
	)
	if err != nil {
		return nil
	}

	c.KafkaConsumer = consumer

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc: c.DoGetConsumer,
		DoGetHealthCheckFunc:   c.DoGetHealthCheck,
		DoGetHTTPServerFunc:    c.DoGetHTTPServer,
		DoGetS3ClientFunc: func(cfg *config.S3Config) (sitemap.S3Client, error) {
			return nil, nil
		},
		DoGetESClientsFunc: func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
			return &dpEsMock.ClientMock{}, &es710.Client{}, nil
		},
	}

	c.serviceList = service.NewServiceList(initMock)

	c.files = make(map[string]string)
	c.S3UploadedSitemap = make(map[string]string)
	c.welshVersion = make(map[string]bool)

	return c
}

func (c *Component) Reset() {
	c.CleanFile(c.cfg.RobotsFilePath[config.English])
	c.CleanFile(c.cfg.RobotsFilePath[config.Welsh])
	for _, file := range c.files {
		c.CleanFile(file)
	}
}

func (c *Component) CleanFile(file string) {
	_, err := os.Stat(file)
	if err != nil {
		// nothing to do
		return
	}
	err = os.Remove(file)
	if err != nil {
		log.Fatal("failed to clean up file: " + err.Error())
	}
}

func (c *Component) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	return dphttp.NewServer(bindAddr, router)
}

func (c *Component) DoGetConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafkaConsumer kafka.IConsumerGroup, err error) {
	return c.KafkaConsumer, nil
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}
