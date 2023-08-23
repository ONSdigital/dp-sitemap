package main

import (
	"context"
	"github.com/ONSdigital/dp-sitemap/global"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-sitemap/cmd"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/ONSdigital/dp-sitemap/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const serviceName = "dp-sitemap"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()
	cmdErr := cmd.GetRootCommand().Execute()
	if cmdErr != nil {
		log.Error(ctx, "error initialising kafka producer", cmdErr)
		return
	}

	if err := run(ctx); err != nil {
		log.Error(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	robotseo.Init(global.CmdFlagFields.RobotsFilePathReader)

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)
	svcList := service.NewServiceList(&service.Init{})
	svc, err := service.Run(ctx, svcList, BuildTime, GitCommit, Version, svcErrors)
	if err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		log.Error(ctx, "service error received", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return svc.Close(ctx)
}
