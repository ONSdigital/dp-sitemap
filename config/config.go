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
	OutputFilePath             string        `envconfig:"OUTPUT_FILE_PATH"`
	KafkaConfig                KafkaConfig
	OpenSearchConfig           OpenSearchConfig
	S3Config                   S3Config
}

type S3Config struct {
	UploadBucketName string `envconfig:"S3_UPLOAD_BUCKET_NAME"`
	AwsRegion        string `envconfig:"S3_AWS_REGION"`
	LocalstackHost   string `envconfig:"S3_LOCALSTACK_HOST"`
}

type OpenSearchConfig struct {
	APIURL                string `envconfig:"OPENSEARCH_API_URL"`
	Signer                bool   `envconfig:"OPENSEARCH_SIGNER"`
	SignerFilename        string `envconfig:"OPENSEARCH_SIGNER_AWS_FILENAME"`
	SignerProfile         string `envconfig:"OPENSEARCH_SIGNER_AWS_PROFILE"`
	SignerRegion          string `envconfig:"OPENSEARCH_SIGNER_AWS_REGION"`
	SignerService         string `envconfig:"OPENSEARCH_SIGNER_AWS_SERVICE"`
	TLSInsecureSkipVerify bool   `envconfig:"OPENSEARCH_TLS_INSECURE_SKIP_VERIFY"`
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
		OutputFilePath:             "/tmp/helloworld.txt",
		KafkaConfig: KafkaConfig{
			Brokers:             []string{"localhost:9092"},
			Version:             "1.0.2",
			OffsetOldest:        true,
			NumWorkers:          1,
			ContentUpdatedGroup: "dp-sitemap",
			ContentUpdatedTopic: "content-updated",
		},
	}

	cfg.OpenSearchConfig = OpenSearchConfig{
		APIURL:                "http://localhost:11200",
		SignerFilename:        "",
		SignerProfile:         "",
		SignerRegion:          "eu-west-2",
		SignerService:         "es",
		Signer:                false,
		TLSInsecureSkipVerify: false,
	}

	return cfg, envconfig.Process("", cfg)
}
