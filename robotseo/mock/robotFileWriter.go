// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"sync"
)

// Ensure, that RobotFileWriterInterfaceMock does implement robotseo.RobotFileWriterInterface.
// If this is not the case, regenerate this file with moq.
var _ robotseo.RobotFileWriterInterface = &RobotFileWriterInterfaceMock{}

// RobotFileWriterInterfaceMock is a mock implementation of robotseo.RobotFileWriterInterface.
//
// 	func TestSomethingThatUsesRobotFileWriterInterface(t *testing.T) {
//
// 		// make and configure a mocked robotseo.RobotFileWriterInterface
// 		mockedRobotFileWriterInterface := &RobotFileWriterInterfaceMock{
// 			GetRobotsFileBodyFunc: func() string {
// 				panic("mock out the GetRobotsFileBody method")
// 			},
// 			WriteRobotsFileFunc: func(cfg *config.Config, sitemaps []string) error {
// 				panic("mock out the WriteRobotsFile method")
// 			},
// 		}
//
// 		// use mockedRobotFileWriterInterface in code that requires robotseo.RobotFileWriterInterface
// 		// and then make assertions.
//
// 	}
type RobotFileWriterInterfaceMock struct {
	// GetRobotsFileBodyFunc mocks the GetRobotsFileBody method.
	GetRobotsFileBodyFunc func() string

	// WriteRobotsFileFunc mocks the WriteRobotsFile method.
	WriteRobotsFileFunc func(cfg *config.Config, sitemaps []string) error

	// calls tracks calls to the methods.
	calls struct {
		// GetRobotsFileBody holds details about calls to the GetRobotsFileBody method.
		GetRobotsFileBody []struct {
		}
		// WriteRobotsFile holds details about calls to the WriteRobotsFile method.
		WriteRobotsFile []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// Sitemaps is the sitemaps argument value.
			Sitemaps []string
		}
	}
	lockGetRobotsFileBody sync.RWMutex
	lockWriteRobotsFile   sync.RWMutex
}

// GetRobotsFileBody calls GetRobotsFileBodyFunc.
func (mock *RobotFileWriterInterfaceMock) GetRobotsFileBody() string {
	if mock.GetRobotsFileBodyFunc == nil {
		panic("RobotFileWriterInterfaceMock.GetRobotsFileBodyFunc: method is nil but RobotFileWriterInterface.GetRobotsFileBody was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetRobotsFileBody.Lock()
	mock.calls.GetRobotsFileBody = append(mock.calls.GetRobotsFileBody, callInfo)
	mock.lockGetRobotsFileBody.Unlock()
	return mock.GetRobotsFileBodyFunc()
}

// GetRobotsFileBodyCalls gets all the calls that were made to GetRobotsFileBody.
// Check the length with:
//     len(mockedRobotFileWriterInterface.GetRobotsFileBodyCalls())
func (mock *RobotFileWriterInterfaceMock) GetRobotsFileBodyCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetRobotsFileBody.RLock()
	calls = mock.calls.GetRobotsFileBody
	mock.lockGetRobotsFileBody.RUnlock()
	return calls
}

// WriteRobotsFile calls WriteRobotsFileFunc.
func (mock *RobotFileWriterInterfaceMock) WriteRobotsFile(cfg *config.Config, sitemaps []string) error {
	if mock.WriteRobotsFileFunc == nil {
		panic("RobotFileWriterInterfaceMock.WriteRobotsFileFunc: method is nil but RobotFileWriterInterface.WriteRobotsFile was just called")
	}
	callInfo := struct {
		Cfg      *config.Config
		Sitemaps []string
	}{
		Cfg:      cfg,
		Sitemaps: sitemaps,
	}
	mock.lockWriteRobotsFile.Lock()
	mock.calls.WriteRobotsFile = append(mock.calls.WriteRobotsFile, callInfo)
	mock.lockWriteRobotsFile.Unlock()
	return mock.WriteRobotsFileFunc(cfg, sitemaps)
}

// WriteRobotsFileCalls gets all the calls that were made to WriteRobotsFile.
// Check the length with:
//     len(mockedRobotFileWriterInterface.WriteRobotsFileCalls())
func (mock *RobotFileWriterInterfaceMock) WriteRobotsFileCalls() []struct {
	Cfg      *config.Config
	Sitemaps []string
} {
	var calls []struct {
		Cfg      *config.Config
		Sitemaps []string
	}
	mock.lockWriteRobotsFile.RLock()
	calls = mock.calls.WriteRobotsFile
	mock.lockWriteRobotsFile.RUnlock()
	return calls
}
