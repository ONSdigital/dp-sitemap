// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"io"
	"sync"
)

// Ensure, that AdderMock does implement sitemap.Adder.
// If this is not the case, regenerate this file with moq.
var _ sitemap.Adder = &AdderMock{}

// AdderMock is a mock implementation of sitemap.Adder.
//
//	func TestSomethingThatUsesAdder(t *testing.T) {
//
//		// make and configure a mocked sitemap.Adder
//		mockedAdder := &AdderMock{
//			AddFunc: func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
//				panic("mock out the Add method")
//			},
//		}
//
//		// use mockedAdder in code that requires sitemap.Adder
//		// and then make assertions.
//
//	}
type AdderMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error)

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			// OldSitemap is the oldSitemap argument value.
			OldSitemap io.Reader
			// URL is the url argument value.
			URL *sitemap.URL
		}
	}
	lockAdd sync.RWMutex
}

// Add calls AddFunc.
func (mock *AdderMock) Add(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
	if mock.AddFunc == nil {
		panic("AdderMock.AddFunc: method is nil but Adder.Add was just called")
	}
	callInfo := struct {
		OldSitemap io.Reader
		URL        *sitemap.URL
	}{
		OldSitemap: oldSitemap,
		URL:        url,
	}
	mock.lockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	mock.lockAdd.Unlock()
	return mock.AddFunc(oldSitemap, url)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//
//	len(mockedAdder.AddCalls())
func (mock *AdderMock) AddCalls() []struct {
	OldSitemap io.Reader
	URL        *sitemap.URL
} {
	var calls []struct {
		OldSitemap io.Reader
		URL        *sitemap.URL
	}
	mock.lockAdd.RLock()
	calls = mock.calls.Add
	mock.lockAdd.RUnlock()
	return calls
}
