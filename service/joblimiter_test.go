package service_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ONSdigital/dp-sitemap/service"
	"github.com/go-co-op/gocron"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJobLimiter(t *testing.T) {
	Convey("given a job is limited to 5 concurrency", t, func() {
		var counter atomic.Int32
		longJob := func(job gocron.Job) {
			counter.Add(1)
			time.Sleep(time.Second)
		}
		limiter := service.NewJobLimiter(5)
		limitedJob := limiter.GetJob(longJob)

		Convey("when the job is run 100 times in parallel", func() {
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					limitedJob(gocron.Job{})
					wg.Done()
				}()
			}
			wg.Wait()
			Convey("the job should only run 5 times", func() {
				So(counter.Load(), ShouldEqual, 5)
			})
		})
	})
	Convey("given a job is limited to 1 concurrency", t, func() {
		var counter atomic.Int32
		longJob := func(job gocron.Job) {
			counter.Add(1)
			time.Sleep(time.Second)
		}
		limiter := service.NewJobLimiter(1)
		limitedJob := limiter.GetJob(longJob)

		Convey("when the job is run 100 times in parallel", func() {
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					limitedJob(gocron.Job{})
					wg.Done()
				}()
			}
			wg.Wait()
			Convey("the job should only run once", func() {
				So(counter.Load(), ShouldEqual, 1)
			})
		})
	})
}
