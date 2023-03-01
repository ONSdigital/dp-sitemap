package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-sitemap
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	SitemapGenerationFrequency time.Duration `envconfig:"SITEMAP_GENERATION_FREQUENCY"`
	SitemapGenerationTimeout   time.Duration `envconfig:"SITEMAP_GENERATION_TIMEOUT"`
	RobotsFilePath             string        `envconfig:"ROBOTS_FILE_PATH"`
	KafkaConfig                KafkaConfig
	OpenSearchConfig           OpenSearchConfig
	SitemapSaveLocation        string `envconfig:"SITEMAP_SAVE_LOCATION"` // "local" or "s3", default: "local"
	SitemapLocalFile           string `envconfig:"SITEMAP_LOCAL_FILE"`
	S3Config                   S3Config
	ZebedeeURL                 string `envconfig:"ZEBEDEE_URL"`
}

type S3Config struct {
	UploadBucketName string `envconfig:"S3_UPLOAD_BUCKET_NAME"`
	SitemapFileKey   string `envconfig:"S3_SITEMAP_FILE_KEY"`
	AwsRegion        string `envconfig:"S3_AWS_REGION"`
	LocalstackHost   string `envconfig:"S3_LOCALSTACK_HOST"`
}

type OpenSearchConfig struct {
	APIURL                string        `envconfig:"OPENSEARCH_API_URL"`
	SitemapIndex          string        `envconfig:"OPENSEARCH_SITEMAP_INDEX"`
	ScrollTimeout         time.Duration `envconfig:"OPENSEARCH_SCROLL_TIMEOUT"`
	ScrollSize            int           `envconfig:"OPENSEARCH_SCROLL_TIMEOUT"`
	Signer                bool          `envconfig:"OPENSEARCH_SIGNER"`
	SignerFilename        string        `envconfig:"OPENSEARCH_SIGNER_AWS_FILENAME"`
	SignerProfile         string        `envconfig:"OPENSEARCH_SIGNER_AWS_PROFILE"`
	SignerRegion          string        `envconfig:"OPENSEARCH_SIGNER_AWS_REGION"`
	SignerService         string        `envconfig:"OPENSEARCH_SIGNER_AWS_SERVICE"`
	TLSInsecureSkipVerify bool          `envconfig:"OPENSEARCH_TLS_INSECURE_SKIP_VERIFY"`
}

// KafkaConfig contains the config required to connect to Kafka
// TODO: change "hello-called" to your topic (config field name, env var name, default value later)
type KafkaConfig struct {
	Brokers             []string `envconfig:"KAFKA_ADDR"`
	Version             string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest        bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	SecProtocol         string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts          string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientKey        string   `envconfig:"KAFKA_SEC_CLIENT_KEY"    json:"-"`
	SecClientCert       string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecSkipVerify       bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	NumWorkers          int      `envconfig:"KAFKA_NUM_WORKERS"`
	ContentUpdatedGroup string   `envconfig:"KAFKA_CONTENT_UPDATED_GROUP"`
	ContentUpdatedTopic string   `envconfig:"KAFKA_CONTENT_UPDATED_TOPIC"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		SitemapGenerationFrequency: time.Hour,
		SitemapGenerationTimeout:   10 * time.Minute,
		RobotsFilePath:             "./dp_robot_file.txt",
		KafkaConfig: KafkaConfig{
			Brokers:             []string{"localhost:9092"},
			Version:             "1.0.2",
			OffsetOldest:        true,
			NumWorkers:          1,
			ContentUpdatedGroup: "dp-sitemap",
			ContentUpdatedTopic: "content-updated",
		},
		SitemapSaveLocation: "local",
		SitemapLocalFile:    "/tmp/dp-sitemap.xml",
		ZebedeeURL:          "http://localhost:8082",
	}

	cfg.OpenSearchConfig = OpenSearchConfig{
		APIURL:                "http://localhost:11200",
		SitemapIndex:          "ons",
		ScrollTimeout:         time.Minute,
		ScrollSize:            10000,
		SignerFilename:        "",
		SignerProfile:         "",
		SignerRegion:          "eu-west-2",
		SignerService:         "es",
		Signer:                false,
		TLSInsecureSkipVerify: false,
	}

	return cfg, envconfig.Process("", cfg)
}
