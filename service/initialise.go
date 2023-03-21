package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	dpEs "github.com/ONSdigital/dp-elasticsearch/v3"
	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-net/v2/awsauth"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	es710 "github.com/elastic/go-elasticsearch/v7"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck   bool
	KafkaConsumer bool
	S3Client      bool
	ESClient      bool
	ZebedeeClient bool
	Init          Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:   false,
		KafkaConsumer: false,
		S3Client:      false,
		Init:          initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server and sets the Server flag to true
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetKafkaConsumer creates a Kafka consumer and sets the consumer flag to true
func (e *ExternalServiceList) GetKafkaConsumer(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, error) {
	consumer, err := e.Init.DoGetKafkaConsumer(ctx, &cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}
	e.KafkaConsumer = true
	return consumer, nil
}

// GetS3Client creates an S3Client and sets the S3Client flag to true
func (e *ExternalServiceList) GetS3Client(cfg *config.Config) (sitemap.S3Client, error) {
	consumer, err := e.Init.DoGetS3Client(&cfg.S3Config)
	if err != nil {
		return nil, err
	}
	e.S3Client = true
	return consumer, nil
}

// GetESClient creates an ESClient and sets the ESClient flag to true
func (e *ExternalServiceList) GetESClient(ctx context.Context, cfg *config.Config) (dpClient dpEsClient.Client, rawClient *es710.Client, err error) {
	dpClient, rawClient, err = e.Init.DoGetESClients(ctx, &cfg.OpenSearchConfig)
	if err != nil {
		return nil, nil, err
	}
	e.ESClient = true
	return dpClient, rawClient, nil
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// GetZebedee return zebedee client
func (e *ExternalServiceList) GetZebedee(cfg *config.Config) clients.ZebedeeClient {
	zebedeeClient := e.Init.DoGetZebedeeClient(cfg)
	e.ZebedeeClient = true
	return zebedeeClient
}

// DoGetZebedeeClient gets and initialises the Zebedee Client
func (e *Init) DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient {
	zebedeeClient := zebedee.New(cfg.ZebedeeURL)
	return zebedeeClient
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetKafkaConsumer returns a Kafka Consumer group
func (e *Init) DoGetKafkaConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (dpkafka.IConsumerGroup, error) {
	kafkaOffset := dpkafka.OffsetNewest
	if kafkaCfg.OffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}
	cgConfig := &dpkafka.ConsumerGroupConfig{
		KafkaVersion: &kafkaCfg.Version,
		Offset:       &kafkaOffset,
		Topic:        kafkaCfg.ContentUpdatedTopic,
		GroupName:    kafkaCfg.ContentUpdatedGroup,
		BrokerAddrs:  kafkaCfg.Brokers,
	}
	if kafkaCfg.SecProtocol == config.KafkaTLSProtocolFlag {
		cgConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			kafkaCfg.SecCACerts,
			kafkaCfg.SecClientCert,
			kafkaCfg.SecClientKey,
			kafkaCfg.SecSkipVerify,
		)
	}
	kafkaConsumer, err := dpkafka.NewConsumerGroup(
		ctx,
		cgConfig,
	)
	if err != nil {
		return nil, err
	}

	return kafkaConsumer, nil
}

// DoGetS3Client returns a S3Client
func (e *Init) DoGetS3Client(cfg *config.S3Config) (sitemap.S3Client, error) {
	if cfg.LocalstackHost != "" {
		s, err := session.NewSession(&aws.Config{
			Endpoint:         aws.String(cfg.LocalstackHost),
			Region:           aws.String(cfg.AwsRegion),
			S3ForcePathStyle: aws.Bool(true),
			Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		})
		if err != nil {
			return nil, err
		}

		return dps3.NewClientWithSession(cfg.UploadBucketName, s), nil
	}

	s3Client, err := dps3.NewClient(cfg.AwsRegion, cfg.UploadBucketName)
	if err != nil {
		return nil, err
	}
	return s3Client, nil
}

// DoGetS3Clients returns a DP and raw Elastic clients
func (e *Init) DoGetESClients(ctx context.Context, cfg *config.OpenSearchConfig) (esClient dpEsClient.Client, rawClient *es710.Client, err error) {
	var transport http.RoundTripper = dphttp.DefaultTransport
	if cfg.Signer {
		transport, err = awsauth.NewAWSSignerRoundTripper(cfg.SignerFilename, cfg.SignerProfile, cfg.SignerRegion, cfg.SignerService, awsauth.Options{TlsInsecureSkipVerify: cfg.TLSInsecureSkipVerify})
		if err != nil {
			log.Error(ctx, "failed to create aws auth roundtripper", err)
			return nil, nil, err
		}
	}

	esClient, err = dpEs.NewClient(dpEsClient.Config{
		ClientLib: dpEsClient.GoElasticV710,
		Address:   cfg.APIURL,
		Transport: transport,
	})
	if err != nil {
		return nil, nil, err
	}

	rawClient, err = es710.NewClient(es710.Config{
		Addresses: []string{cfg.APIURL},
		Transport: transport,
	})
	if err != nil {
		return nil, nil, err
	}

	return esClient, rawClient, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}
