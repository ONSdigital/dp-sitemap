package service

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-co-op/gocron"
)

type JobLimiter struct {
	concurrency int
	semaphore   chan struct{}
}

func NewJobLimiter(concurrency int) *JobLimiter {
	return &JobLimiter{
		concurrency: concurrency,
		semaphore:   make(chan struct{}, concurrency),
	}
}

func (l *JobLimiter) GetJob(j func(job gocron.Job)) func(job gocron.Job) {
	ctx := context.Background()
	return func(job gocron.Job) {
		select {
		case l.semaphore <- struct{}{}:
			log.Info(ctx, "limited job started", log.Data{"concurrency_limit": l.concurrency, "concurrency_now": len(l.semaphore)})
			j(job)
		default:
			log.Warn(ctx, "limited job skipped", log.Data{"concurrency_limit": l.concurrency, "concurrency_now": len(l.semaphore)})
		}
	}
}
