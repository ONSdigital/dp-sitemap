// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"io"
	"sync"
)

// Ensure, that FileSaverMock does implement sitemap.FileSaver.
// If this is not the case, regenerate this file with moq.
var _ sitemap.FileSaver = &FileSaverMock{}

// FileSaverMock is a mock implementation of sitemap.FileSaver.
//
// 	func TestSomethingThatUsesFileSaver(t *testing.T) {
//
// 		// make and configure a mocked sitemap.FileSaver
// 		mockedFileSaver := &FileSaverMock{
// 			SaveFileFunc: func(body io.Reader) error {
// 				panic("mock out the SaveFile method")
// 			},
// 		}
//
// 		// use mockedFileSaver in code that requires sitemap.FileSaver
// 		// and then make assertions.
//
// 	}
type FileSaverMock struct {
	// SaveFileFunc mocks the SaveFile method.
	SaveFileFunc func(body io.Reader) error

	// calls tracks calls to the methods.
	calls struct {
		// SaveFile holds details about calls to the SaveFile method.
		SaveFile []struct {
			// Body is the body argument value.
			Body io.Reader
		}
	}
	lockSaveFile sync.RWMutex
}

// SaveFile calls SaveFileFunc.
func (mock *FileSaverMock) SaveFile(body io.Reader) error {
	if mock.SaveFileFunc == nil {
		panic("FileSaverMock.SaveFileFunc: method is nil but FileSaver.SaveFile was just called")
	}
	callInfo := struct {
		Body io.Reader
	}{
		Body: body,
	}
	mock.lockSaveFile.Lock()
	mock.calls.SaveFile = append(mock.calls.SaveFile, callInfo)
	mock.lockSaveFile.Unlock()
	return mock.SaveFileFunc(body)
}

// SaveFileCalls gets all the calls that were made to SaveFile.
// Check the length with:
//     len(mockedFileSaver.SaveFileCalls())
func (mock *FileSaverMock) SaveFileCalls() []struct {
	Body io.Reader
} {
	var calls []struct {
		Body io.Reader
	}
	mock.lockSaveFile.RLock()
	calls = mock.calls.SaveFile
	mock.lockSaveFile.RUnlock()
	return calls
}
