package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	"github.com/ONSdigital/dp-sitemap/schema"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-sitemap"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	pConfig := &kafka.ProducerConfig{
		KafkaVersion:      &cfg.KafkaConfig.Version,
		Topic:             cfg.KafkaConfig.ContentUpdatedTopic,
		BrokerAddrs:       cfg.KafkaConfig.Brokers,
		MinBrokersHealthy: &cfg.KafkaConfig.NumWorkers,
	}
	if cfg.KafkaConfig.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaConfig.SecCACerts,
			cfg.KafkaConfig.SecClientCert,
			cfg.KafkaConfig.SecClientKey,
			cfg.KafkaConfig.SecSkipVerify,
		)
	}
	kafkaProducer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		log.Error(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.KafkaConfig.ContentUpdatedTopic})
		os.Exit(1)
	}

	// kafka error logging go-routines
	kafkaProducer.LogErrors(ctx)

	time.Sleep(500 * time.Millisecond)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		e := scanEvent(scanner)
		log.Info(ctx, "sending hello-called event", log.Data{"helloCalledEvent": e})

		bytes, kafkaErr := schema.ContentPublishedEvent.Marshal(e)
		if kafkaErr != nil {
			log.Error(ctx, "hello-called event error", kafkaErr)
			os.Exit(1)
		}

		// Send bytes to Output channel, after calling Initialise just in case it is not initialised.
		err := kafkaProducer.Initialise(ctx)
		if err != nil {
			log.Error(ctx, "error initialising kafka producer", err)
			return
		}
		kafkaProducer.Channels().Output <- bytes
	}
}

// scanEvent creates a HelloCalled event according to the user input
func scanEvent(scanner *bufio.Scanner) *event.ContentPublished {
	fmt.Println("--- [Send Kafka ContentPublished] ---")

	fmt.Println("Press enter to send message")
	fmt.Printf("$ ")
	scanner.Scan()
	scanner.Text()

	return &event.ContentPublished{
		URI:          "test-uri",
		DataType:     "thedatatype",
		CollectionID: "thecollectionid",
		JobID:        "thejobId",
		SearchIndex:  "theSearchIndex",
		TraceID:      "theTraceId",
	}
}
