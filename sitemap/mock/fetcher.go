// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"sync"
)

// Ensure, that FetcherMock does implement sitemap.Fetcher.
// If this is not the case, regenerate this file with moq.
var _ sitemap.Fetcher = &FetcherMock{}

// FetcherMock is a mock implementation of sitemap.Fetcher.
//
//	func TestSomethingThatUsesFetcher(t *testing.T) {
//
//		// make and configure a mocked sitemap.Fetcher
//		mockedFetcher := &FetcherMock{
//			GetFullSitemapFunc: func(ctx context.Context) (sitemap.Files, error) {
//				panic("mock out the GetFullSitemap method")
//			},
//			HasWelshContentFunc: func(ctx context.Context, path string) bool {
//				panic("mock out the HasWelshContent method")
//			},
//			URLVersionsFunc: func(ctx context.Context, path string, lastmod string) (sitemap.URL, *sitemap.URL) {
//				panic("mock out the URLVersions method")
//			},
//		}
//
//		// use mockedFetcher in code that requires sitemap.Fetcher
//		// and then make assertions.
//
//	}
type FetcherMock struct {
	// GetFullSitemapFunc mocks the GetFullSitemap method.
	GetFullSitemapFunc func(ctx context.Context) (sitemap.Files, error)

	// HasWelshContentFunc mocks the HasWelshContent method.
	HasWelshContentFunc func(ctx context.Context, path string) bool

	// URLVersionsFunc mocks the URLVersions method.
	URLVersionsFunc func(ctx context.Context, path string, lastmod string) (sitemap.URL, *sitemap.URL)

	// calls tracks calls to the methods.
	calls struct {
		// GetFullSitemap holds details about calls to the GetFullSitemap method.
		GetFullSitemap []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// HasWelshContent holds details about calls to the HasWelshContent method.
		HasWelshContent []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Path is the path argument value.
			Path string
		}
		// URLVersions holds details about calls to the URLVersions method.
		URLVersions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Path is the path argument value.
			Path string
			// Lastmod is the lastmod argument value.
			Lastmod string
		}
	}
	lockGetFullSitemap  sync.RWMutex
	lockHasWelshContent sync.RWMutex
	lockURLVersions     sync.RWMutex
}

// GetFullSitemap calls GetFullSitemapFunc.
func (mock *FetcherMock) GetFullSitemap(ctx context.Context) (sitemap.Files, error) {
	if mock.GetFullSitemapFunc == nil {
		panic("FetcherMock.GetFullSitemapFunc: method is nil but Fetcher.GetFullSitemap was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetFullSitemap.Lock()
	mock.calls.GetFullSitemap = append(mock.calls.GetFullSitemap, callInfo)
	mock.lockGetFullSitemap.Unlock()
	return mock.GetFullSitemapFunc(ctx)
}

// GetFullSitemapCalls gets all the calls that were made to GetFullSitemap.
// Check the length with:
//
//	len(mockedFetcher.GetFullSitemapCalls())
func (mock *FetcherMock) GetFullSitemapCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetFullSitemap.RLock()
	calls = mock.calls.GetFullSitemap
	mock.lockGetFullSitemap.RUnlock()
	return calls
}

// HasWelshContent calls HasWelshContentFunc.
func (mock *FetcherMock) HasWelshContent(ctx context.Context, path string) bool {
	if mock.HasWelshContentFunc == nil {
		panic("FetcherMock.HasWelshContentFunc: method is nil but Fetcher.HasWelshContent was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Path string
	}{
		Ctx:  ctx,
		Path: path,
	}
	mock.lockHasWelshContent.Lock()
	mock.calls.HasWelshContent = append(mock.calls.HasWelshContent, callInfo)
	mock.lockHasWelshContent.Unlock()
	return mock.HasWelshContentFunc(ctx, path)
}

// HasWelshContentCalls gets all the calls that were made to HasWelshContent.
// Check the length with:
//
//	len(mockedFetcher.HasWelshContentCalls())
func (mock *FetcherMock) HasWelshContentCalls() []struct {
	Ctx  context.Context
	Path string
} {
	var calls []struct {
		Ctx  context.Context
		Path string
	}
	mock.lockHasWelshContent.RLock()
	calls = mock.calls.HasWelshContent
	mock.lockHasWelshContent.RUnlock()
	return calls
}

// URLVersions calls URLVersionsFunc.
func (mock *FetcherMock) URLVersions(ctx context.Context, path string, lastmod string) (sitemap.URL, *sitemap.URL) {
	if mock.URLVersionsFunc == nil {
		panic("FetcherMock.URLVersionsFunc: method is nil but Fetcher.URLVersions was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Path    string
		Lastmod string
	}{
		Ctx:     ctx,
		Path:    path,
		Lastmod: lastmod,
	}
	mock.lockURLVersions.Lock()
	mock.calls.URLVersions = append(mock.calls.URLVersions, callInfo)
	mock.lockURLVersions.Unlock()
	return mock.URLVersionsFunc(ctx, path, lastmod)
}

// URLVersionsCalls gets all the calls that were made to URLVersions.
// Check the length with:
//
//	len(mockedFetcher.URLVersionsCalls())
func (mock *FetcherMock) URLVersionsCalls() []struct {
	Ctx     context.Context
	Path    string
	Lastmod string
} {
	var calls []struct {
		Ctx     context.Context
		Path    string
		Lastmod string
	}
	mock.lockURLVersions.RLock()
	calls = mock.calls.URLVersions
	mock.lockURLVersions.RUnlock()
	return calls
}
