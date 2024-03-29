package service

import (
	"context"
	"strings"
	"time"

	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/event"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const schedulerTagFullSitemap = "full-sitemap"

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	server          HTTPServer
	router          *mux.Router
	serviceList     *ExternalServiceList
	healthCheck     HealthChecker
	consumer        kafka.IConsumerGroup
	shutdownTimeout time.Duration
	scheduler       *gocron.Scheduler
	esClient        dpEsClient.Client
	s3Client        sitemap.S3Client
}

// Run the service
func Run(ctx context.Context, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	// Get HTTP Server with collectionID checkHeader middleware
	r := mux.NewRouter()
	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get Kafka consumer
	consumer, err := serviceList.GetKafkaConsumer(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumer", err)
		return nil, err
	}

	// Get S3 Client
	s3Client, err := serviceList.GetS3Client(cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise s3 client", err)
		return nil, err
	}

	// Get ElasticSearch Clients
	esClient, esRawClient, err := serviceList.GetESClient(ctx, cfg)
	if err != nil {
		log.Error(ctx, "Failed to create dp-elasticsearch clients", err)
		return nil, err
	}

	zebedeeClient := serviceList.GetZebedee(cfg)

	var (
		store                 sitemap.FileStore
		fullSitemapFiles      sitemap.Files
		publishingSitemapFile string
	)
	switch cfg.SitemapSaveLocation {
	case "s3":
		store = sitemap.NewS3Store(
			s3Client,
		)
		fullSitemapFiles = cfg.S3Config.SitemapFileKey
		publishingSitemapFile = cfg.S3Config.PublishingSitemapFileKey

	default:
		store = &sitemap.LocalStore{}
		fullSitemapFiles = cfg.SitemapLocalFile
		publishingSitemapFile = cfg.PublishingSitemapLocalFile
	}

	scroll := sitemap.NewElasticScroll(esRawClient, cfg)
	fetcher := sitemap.NewElasticFetcher(scroll, cfg, zebedeeClient)
	handler := event.NewContentPublishedHandler(store, zebedeeClient, cfg, fetcher)

	// Event Handler for Kafka Consumer
	event.Consume(ctx, consumer, handler, cfg)

	if consumerStartErr := consumer.Start(); consumerStartErr != nil {
		log.Fatal(ctx, "error starting the consumer", consumerStartErr)
		return nil, consumerStartErr
	}

	// Kafka error logging go-routine
	consumer.LogErrors(ctx)

	// Get HealthCheck
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err = registerCheckers(ctx, hc, consumer, esClient, zebedeeClient); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if serveErr := s.ListenAndServe(); serveErr != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonModeAll()

	runSitemapGeneration := func() {
		log.Info(ctx, "full sitemap generation from callback start")
		jobs, findErr := scheduler.FindJobsByTag(schedulerTagFullSitemap)
		if findErr != nil {
			log.Error(ctx, "failed to find full sitemap generation job", findErr)
			return
		}
		if len(jobs) != 1 {
			log.Error(ctx, "unexpected jobs found", findErr)
			return
		}
		if jobs[0].IsRunning() {
			log.Info(ctx, "full sitemap generation from callback - job is already running")
			return
		}
		runErr := scheduler.RunByTag(schedulerTagFullSitemap)
		if runErr != nil {
			log.Error(ctx, "failed to run full sitemap generation job", runErr)
			return
		}
		log.Info(ctx, "full sitemap generation from callback complete")
	}

	scroller := sitemap.NewElasticScroll(esRawClient, cfg)

	generator := sitemap.NewGenerator(
		sitemap.WithFetcher(sitemap.NewElasticFetcher(
			scroller,
			cfg,
			zebedeeClient,
		)),
		sitemap.WithAdder(&sitemap.DefaultAdder{}),
		sitemap.WithFileStore(store),
		sitemap.WithFullSitemapFiles(fullSitemapFiles),
		sitemap.WithPublishingSitemapFile(publishingSitemapFile),
		sitemap.WithPublishingSitemapMaxSize(cfg.PublishingSitemapMaxSize, runSitemapGeneration),
	)

	robotFileWriter := robotseo.RobotFileWriter{}

	generateSitemapJob := func(job gocron.Job) {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.SitemapGenerationTimeout)
		defer cancel()
		log.Info(ctx, "sitemap generation job start", log.Data{"last_run": job.LastRun(), "next_run": job.NextRun(), "run_count": job.RunCount()})
		genErr := generator.MakeFullSitemap(ctx)
		if genErr != nil {
			log.Error(ctx, "failed to generate sitemap", genErr)
			return
		}
		log.Info(ctx, "sitemap generation job complete", log.Data{"last_run": job.LastRun(), "next_run": job.NextRun(), "run_count": job.RunCount()})

		// write robots file
		// TODO: pass sitemap file path (once URL is known)
		for _, lang := range []config.Language{config.English, config.Welsh} {
			body := robotFileWriter.GetRobotsFileBody(lang, cfg.SitemapLocalFile)
			saveErr := store.SaveFile(cfg.RobotsFilePath[lang], strings.NewReader(body))
			if saveErr != nil {
				log.Error(ctx, "failed to save file", saveErr)
				return
			}
		}
		log.Info(ctx, "wrote robots file")
	}

	_, err = scheduler.Every(cfg.SitemapGenerationFrequency).Tag(schedulerTagFullSitemap).DoWithJobDetails(generateSitemapJob)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run scheduler")
	}
	scheduler.StartAsync()

	return &Service{
		server:          s,
		router:          r,
		serviceList:     serviceList,
		healthCheck:     hc,
		consumer:        consumer,
		shutdownTimeout: cfg.GracefulShutdownTimeout,
		scheduler:       scheduler,
		esClient:        esClient,
		s3Client:        s3Client,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.shutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutdown gracefully closes up
	var gracefulShutdown bool

	go func() {
		defer cancel()
		var hasShutdownError bool

		// stop healthcheck, as it depends on everything else
		if svc.serviceList.HealthCheck {
			svc.healthCheck.Stop()
		}
		// stop the scheduler
		log.Info(ctx, "stopping scheduler")
		svc.scheduler.Stop()
		if !svc.scheduler.IsRunning() {
			log.Info(ctx, "stopped scheduler")
		}

		// If kafka consumer exists, stop listening to it.
		// This will automatically stop the event consumer loops and no more messages will be processed.
		// The kafka consumer will be closed after the service shuts down.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "stopping kafka consumer listener")
			if err := svc.consumer.Stop(); err != nil {
				log.Error(ctx, "error stopping kafka consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped kafka consumer listener")
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// If kafka consumer exists, close it.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "closing kafka consumer")
			if err := svc.consumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing kafka consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed kafka consumer")
		}

		if !hasShutdownError {
			gracefulShutdown = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	if !gracefulShutdown {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	hc HealthChecker,
	consumer kafka.IConsumerGroup,
	esClient dpEsClient.Client,
	zebedeeClient clients.ZebedeeClient,
) error {
	hasErrors := false

	if err := hc.AddCheck("Kafka consumer", consumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for Kafka", err)
	}

	if err := hc.AddCheck("Elasticsearch", esClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error creating elasticsearch health check", err)
	}

	if err := hc.AddCheck("Zebedee client", zebedeeClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for ZebedeeClient", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
