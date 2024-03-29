// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	dpEsClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-sitemap/clients"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/service"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	es710 "github.com/elastic/go-elasticsearch/v7"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
// 	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
// 		// make and configure a mocked service.Initialiser
// 		mockedInitialiser := &InitialiserMock{
// 			DoGetESClientsFunc: func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
// 				panic("mock out the DoGetESClients method")
// 			},
// 			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
// 				panic("mock out the DoGetHTTPServer method")
// 			},
// 			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
// 				panic("mock out the DoGetHealthCheck method")
// 			},
// 			DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
// 				panic("mock out the DoGetKafkaConsumer method")
// 			},
// 			DoGetS3ClientFunc: func(cfg *config.S3Config) (sitemap.S3Client, error) {
// 				panic("mock out the DoGetS3Client method")
// 			},
// 			DoGetZebedeeClientFunc: func(cfg *config.Config) clients.ZebedeeClient {
// 				panic("mock out the DoGetZebedeeClient method")
// 			},
// 		}
//
// 		// use mockedInitialiser in code that requires service.Initialiser
// 		// and then make assertions.
//
// 	}
type InitialiserMock struct {
	// DoGetESClientsFunc mocks the DoGetESClients method.
	DoGetESClientsFunc func(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error)

	// DoGetHTTPServerFunc mocks the DoGetHTTPServer method.
	DoGetHTTPServerFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// DoGetHealthCheckFunc mocks the DoGetHealthCheck method.
	DoGetHealthCheckFunc func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error)

	// DoGetKafkaConsumerFunc mocks the DoGetKafkaConsumer method.
	DoGetKafkaConsumerFunc func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error)

	// DoGetS3ClientFunc mocks the DoGetS3Client method.
	DoGetS3ClientFunc func(cfg *config.S3Config) (sitemap.S3Client, error)

	// DoGetZebedeeClientFunc mocks the DoGetZebedeeClient method.
	DoGetZebedeeClientFunc func(cfg *config.Config) clients.ZebedeeClient

	// calls tracks calls to the methods.
	calls struct {
		// DoGetESClients holds details about calls to the DoGetESClients method.
		DoGetESClients []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.OpenSearchConfig
		}
		// DoGetHTTPServer holds details about calls to the DoGetHTTPServer method.
		DoGetHTTPServer []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// DoGetHealthCheck holds details about calls to the DoGetHealthCheck method.
		DoGetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BuildTime is the buildTime argument value.
			BuildTime string
			// GitCommit is the gitCommit argument value.
			GitCommit string
			// Version is the version argument value.
			Version string
		}
		// DoGetKafkaConsumer holds details about calls to the DoGetKafkaConsumer method.
		DoGetKafkaConsumer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// KafkaCfg is the kafkaCfg argument value.
			KafkaCfg *config.KafkaConfig
		}
		// DoGetS3Client holds details about calls to the DoGetS3Client method.
		DoGetS3Client []struct {
			// Cfg is the cfg argument value.
			Cfg *config.S3Config
		}
		// DoGetZebedeeClient holds details about calls to the DoGetZebedeeClient method.
		DoGetZebedeeClient []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
	}
	lockDoGetESClients     sync.RWMutex
	lockDoGetHTTPServer    sync.RWMutex
	lockDoGetHealthCheck   sync.RWMutex
	lockDoGetKafkaConsumer sync.RWMutex
	lockDoGetS3Client      sync.RWMutex
	lockDoGetZebedeeClient sync.RWMutex
}

// DoGetESClients calls DoGetESClientsFunc.
func (mock *InitialiserMock) DoGetESClients(ctx context.Context, cfg *config.OpenSearchConfig) (dpEsClient.Client, *es710.Client, error) {
	if mock.DoGetESClientsFunc == nil {
		panic("InitialiserMock.DoGetESClientsFunc: method is nil but Initialiser.DoGetESClients was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.OpenSearchConfig
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetESClients.Lock()
	mock.calls.DoGetESClients = append(mock.calls.DoGetESClients, callInfo)
	mock.lockDoGetESClients.Unlock()
	return mock.DoGetESClientsFunc(ctx, cfg)
}

// DoGetESClientsCalls gets all the calls that were made to DoGetESClients.
// Check the length with:
//     len(mockedInitialiser.DoGetESClientsCalls())
func (mock *InitialiserMock) DoGetESClientsCalls() []struct {
	Ctx context.Context
	Cfg *config.OpenSearchConfig
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.OpenSearchConfig
	}
	mock.lockDoGetESClients.RLock()
	calls = mock.calls.DoGetESClients
	mock.lockDoGetESClients.RUnlock()
	return calls
}

// DoGetHTTPServer calls DoGetHTTPServerFunc.
func (mock *InitialiserMock) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.DoGetHTTPServerFunc == nil {
		panic("InitialiserMock.DoGetHTTPServerFunc: method is nil but Initialiser.DoGetHTTPServer was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockDoGetHTTPServer.Lock()
	mock.calls.DoGetHTTPServer = append(mock.calls.DoGetHTTPServer, callInfo)
	mock.lockDoGetHTTPServer.Unlock()
	return mock.DoGetHTTPServerFunc(bindAddr, router)
}

// DoGetHTTPServerCalls gets all the calls that were made to DoGetHTTPServer.
// Check the length with:
//     len(mockedInitialiser.DoGetHTTPServerCalls())
func (mock *InitialiserMock) DoGetHTTPServerCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockDoGetHTTPServer.RLock()
	calls = mock.calls.DoGetHTTPServer
	mock.lockDoGetHTTPServer.RUnlock()
	return calls
}

// DoGetHealthCheck calls DoGetHealthCheckFunc.
func (mock *InitialiserMock) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	if mock.DoGetHealthCheckFunc == nil {
		panic("InitialiserMock.DoGetHealthCheckFunc: method is nil but Initialiser.DoGetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}{
		Cfg:       cfg,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		Version:   version,
	}
	mock.lockDoGetHealthCheck.Lock()
	mock.calls.DoGetHealthCheck = append(mock.calls.DoGetHealthCheck, callInfo)
	mock.lockDoGetHealthCheck.Unlock()
	return mock.DoGetHealthCheckFunc(cfg, buildTime, gitCommit, version)
}

// DoGetHealthCheckCalls gets all the calls that were made to DoGetHealthCheck.
// Check the length with:
//     len(mockedInitialiser.DoGetHealthCheckCalls())
func (mock *InitialiserMock) DoGetHealthCheckCalls() []struct {
	Cfg       *config.Config
	BuildTime string
	GitCommit string
	Version   string
} {
	var calls []struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}
	mock.lockDoGetHealthCheck.RLock()
	calls = mock.calls.DoGetHealthCheck
	mock.lockDoGetHealthCheck.RUnlock()
	return calls
}

// DoGetKafkaConsumer calls DoGetKafkaConsumerFunc.
func (mock *InitialiserMock) DoGetKafkaConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
	if mock.DoGetKafkaConsumerFunc == nil {
		panic("InitialiserMock.DoGetKafkaConsumerFunc: method is nil but Initialiser.DoGetKafkaConsumer was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		KafkaCfg *config.KafkaConfig
	}{
		Ctx:      ctx,
		KafkaCfg: kafkaCfg,
	}
	mock.lockDoGetKafkaConsumer.Lock()
	mock.calls.DoGetKafkaConsumer = append(mock.calls.DoGetKafkaConsumer, callInfo)
	mock.lockDoGetKafkaConsumer.Unlock()
	return mock.DoGetKafkaConsumerFunc(ctx, kafkaCfg)
}

// DoGetKafkaConsumerCalls gets all the calls that were made to DoGetKafkaConsumer.
// Check the length with:
//     len(mockedInitialiser.DoGetKafkaConsumerCalls())
func (mock *InitialiserMock) DoGetKafkaConsumerCalls() []struct {
	Ctx      context.Context
	KafkaCfg *config.KafkaConfig
} {
	var calls []struct {
		Ctx      context.Context
		KafkaCfg *config.KafkaConfig
	}
	mock.lockDoGetKafkaConsumer.RLock()
	calls = mock.calls.DoGetKafkaConsumer
	mock.lockDoGetKafkaConsumer.RUnlock()
	return calls
}

// DoGetS3Client calls DoGetS3ClientFunc.
func (mock *InitialiserMock) DoGetS3Client(cfg *config.S3Config) (sitemap.S3Client, error) {
	if mock.DoGetS3ClientFunc == nil {
		panic("InitialiserMock.DoGetS3ClientFunc: method is nil but Initialiser.DoGetS3Client was just called")
	}
	callInfo := struct {
		Cfg *config.S3Config
	}{
		Cfg: cfg,
	}
	mock.lockDoGetS3Client.Lock()
	mock.calls.DoGetS3Client = append(mock.calls.DoGetS3Client, callInfo)
	mock.lockDoGetS3Client.Unlock()
	return mock.DoGetS3ClientFunc(cfg)
}

// DoGetS3ClientCalls gets all the calls that were made to DoGetS3Client.
// Check the length with:
//     len(mockedInitialiser.DoGetS3ClientCalls())
func (mock *InitialiserMock) DoGetS3ClientCalls() []struct {
	Cfg *config.S3Config
} {
	var calls []struct {
		Cfg *config.S3Config
	}
	mock.lockDoGetS3Client.RLock()
	calls = mock.calls.DoGetS3Client
	mock.lockDoGetS3Client.RUnlock()
	return calls
}

// DoGetZebedeeClient calls DoGetZebedeeClientFunc.
func (mock *InitialiserMock) DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient {
	if mock.DoGetZebedeeClientFunc == nil {
		panic("InitialiserMock.DoGetZebedeeClientFunc: method is nil but Initialiser.DoGetZebedeeClient was just called")
	}
	callInfo := struct {
		Cfg *config.Config
	}{
		Cfg: cfg,
	}
	mock.lockDoGetZebedeeClient.Lock()
	mock.calls.DoGetZebedeeClient = append(mock.calls.DoGetZebedeeClient, callInfo)
	mock.lockDoGetZebedeeClient.Unlock()
	return mock.DoGetZebedeeClientFunc(cfg)
}

// DoGetZebedeeClientCalls gets all the calls that were made to DoGetZebedeeClient.
// Check the length with:
//     len(mockedInitialiser.DoGetZebedeeClientCalls())
func (mock *InitialiserMock) DoGetZebedeeClientCalls() []struct {
	Cfg *config.Config
} {
	var calls []struct {
		Cfg *config.Config
	}
	mock.lockDoGetZebedeeClient.RLock()
	calls = mock.calls.DoGetZebedeeClient
	mock.lockDoGetZebedeeClient.RUnlock()
	return calls
}
