package service_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpEsMock "github.com/ONSdigital/dp-elasticsearch/v3/client/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-sitemap/clients"
	clientMock "github.com/ONSdigital/dp-sitemap/clients/mock"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/service"
	serviceMock "github.com/ONSdigital/dp-sitemap/service/mock"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	sitemapMock "github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	es710 "github.com/elastic/go-elasticsearch/v7"
	esapi710 "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
)

var (
	errKafkaConsumer = errors.New("Kafka consumer error")
	errHealthcheck   = errors.New("healthCheck error")
)

var funcDoGetKafkaConsumerErr = func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
	return nil, errKafkaConsumer
}

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

func TestRun(t *testing.T) {
	Convey("Having a set of mocked dependencies", t, func() {
		consumerMock := &kafkatest.IConsumerGroupMock{
			StartFunc:     func() error { return nil },
			LogErrorsFunc: func(ctx context.Context) {},
			CheckerFunc:   func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc:  func() *kafka.ConsumerGroupChannels { return &kafka.ConsumerGroupChannels{} },
		}

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}
		s3Mock := &sitemapMock.S3ClientMock{
			UploadFunc: func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, nil
			},
		}
		esMock := &dpEsMock.ClientMock{}
		esRawMock := &es710.Client{API: &esapi710.API{
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
		}}
		zebedeeMock := &clientMock.ZebedeeClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}
		funcDoGetKafkaConsumerOk := func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
			return consumerMock, nil
		}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetS3ClientFunc := func(cfg *config.S3Config) (sitemap.S3Client, error) {
			return s3Mock, nil
		}
		funcDoGetESClientFunc := func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
			return esMock, esRawMock, nil
		}
		funcDoGetZebedeeOk := func(cfg *config.Config) clients.ZebedeeClient {
			return zebedeeMock
		}

		Convey("Given that initialising Kafka consumer returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:    funcDoGetHTTPServerNil,
				DoGetKafkaConsumerFunc: funcDoGetKafkaConsumerErr,
				DoGetZebedeeClientFunc: funcDoGetZebedeeOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errKafkaConsumer)
				So(svcList.KafkaConsumer, ShouldBeFalse)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})

		Convey("Given that initialising healthcheck returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:    funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc:   funcDoGetHealthcheckErr,
				DoGetKafkaConsumerFunc: funcDoGetKafkaConsumerOk,
				DoGetS3ClientFunc:      funcDoGetS3ClientFunc,
				DoGetESClientsFunc:     funcDoGetESClientFunc,
				DoGetZebedeeClientFunc: funcDoGetZebedeeOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:    funcDoGetHTTPServer,
				DoGetHealthCheckFunc:   funcDoGetHealthcheckOk,
				DoGetKafkaConsumerFunc: funcDoGetKafkaConsumerOk,
				DoGetS3ClientFunc:      funcDoGetS3ClientFunc,
				DoGetESClientsFunc:     funcDoGetESClientFunc,
				DoGetZebedeeClientFunc: funcDoGetZebedeeOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeTrue)
			})

			Convey("The checkers are registered and the healthcheck and http server started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 3)
				So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "Kafka consumer")
				So(hcMock.AddCheckCalls()[1].Name, ShouldResemble, "Elasticsearch")
				So(hcMock.AddCheckCalls()[2].Name, ShouldResemble, "Zebedee client")
				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:")
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})

		Convey("Given that Checkers cannot be registered", func() {
			errAddheckFail := errors.New("Error(s) registering checkers for healthcheck")
			hcMockAddFail := &serviceMock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddheckFail },
				StartFunc:    func(ctx context.Context) {},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				DoGetKafkaConsumerFunc: funcDoGetKafkaConsumerOk,
				DoGetS3ClientFunc:      funcDoGetS3ClientFunc,
				DoGetESClientsFunc:     funcDoGetESClientFunc,
				DoGetZebedeeClientFunc: funcDoGetZebedeeOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 3)
				So(hcMockAddFail.AddCheckCalls()[0].Name, ShouldResemble, "Kafka consumer")
				So(hcMockAddFail.AddCheckCalls()[1].Name, ShouldResemble, "Elasticsearch")
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		hcStopped := false

		consumerMock := &kafkatest.IConsumerGroupMock{
			StartFunc:     func() error { return nil },
			LogErrorsFunc: func(ctx context.Context) {},
			StopFunc:      func() error { return nil },
			CloseFunc:     func(ctx context.Context, optFuncs ...kafka.OptFunc) error { return nil },
			CheckerFunc:   func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc:  func() *kafka.ConsumerGroupChannels { return &kafka.ConsumerGroupChannels{} },
		}

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		s3Mock := &sitemapMock.S3ClientMock{
			UploadFunc: func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, nil
			},
		}
		esMock := &dpEsMock.ClientMock{}
		esRawMock := &es710.Client{API: &esapi710.API{
			Scroll: func(o ...func(*esapi710.ScrollRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
			Search: func(o ...func(*esapi710.SearchRequest)) (*esapi710.Response, error) {
				return &esapi710.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				}, nil
			},
		}}
		zMock := clientMock.ZebedeeClientMock{
			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error { return nil },
			GetFileSizeFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error) {
				return zebedee.FileSize{Size: 1}, errors.New("no welsh content")
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
					return consumerMock, nil
				},
				DoGetS3ClientFunc: func(cfg *config.S3Config) (sitemap.S3Client, error) {
					return s3Mock, nil
				},
				DoGetESClientsFunc: func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
					return esMock, esRawMock, nil
				},
				DoGetZebedeeClientFunc: func(cfg *config.Config) clients.ZebedeeClient {
					return &zMock
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(consumerMock.CloseCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return failingserverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
					return consumerMock, nil
				},
				DoGetS3ClientFunc: func(cfg *config.S3Config) (sitemap.S3Client, error) {
					return s3Mock, nil
				},
				DoGetESClientsFunc: func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
					return esMock, esRawMock, nil
				},
				DoGetZebedeeClientFunc: func(cfg *config.Config) clients.ZebedeeClient {
					return &zMock
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
