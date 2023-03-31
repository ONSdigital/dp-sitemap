package service

import (
	"context"
	"net/http"

	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
)

//go:generate moq -out mock/initialiser.go -pkg mock . Initialiser
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/healthCheck.go -pkg mock . HealthChecker

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error)
	DoGetKafkaConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error)
	DoGetS3Client(cfg *config.S3Config) (sitemap.S3Client, error)
	DoGetESClients(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error)
	DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}

// EventConsumer defines the required methods from event Consumer
type EventConsumer interface {
	Close(ctx context.Context) (err error)
}
