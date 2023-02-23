package steps

import (
	"context"
	"net/http"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/service"
	"github.com/ONSdigital/dp-sitemap/service/mock"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"

	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpEsMock "github.com/ONSdigital/dp-elasticsearch/v3/client/mocks"
)

type Component struct {
	componenttest.ErrorFeature
	serviceList   *service.ExternalServiceList
	KafkaConsumer kafka.IConsumerGroup
	killChannel   chan os.Signal
	apiFeature    *componenttest.APIFeature
	errorChan     chan error
	svc           *service.Service
	cfg           *config.Config
}

func NewComponent() *Component {

	c := &Component{errorChan: make(chan error)}

	consumer := kafkatest.NewMessageConsumer(false)
	consumer.CheckerFunc = funcCheck
	consumer.StartFunc = func() error { return nil }
	consumer.LogErrorsFunc = func(ctx context.Context) {}
	c.KafkaConsumer = consumer

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	c.cfg = cfg

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc: c.DoGetConsumer,
		DoGetHealthCheckFunc:   c.DoGetHealthCheck,
		DoGetHTTPServerFunc:    c.DoGetHTTPServer,
		DoGetS3ClientFunc: func(ctx context.Context, cfg *config.S3Config) (sitemap.S3Uploader, error) {
			return nil, nil
		},
		DoGetESClientsFunc: func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
			return &dpEsMock.ClientMock{}, &es710.Client{}, nil
		},
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {
	os.Remove(c.cfg.RobotsFilePath)
}

func (c *Component) Reset() {
	os.Remove(c.cfg.RobotsFilePath)
}

func (c *Component) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
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
