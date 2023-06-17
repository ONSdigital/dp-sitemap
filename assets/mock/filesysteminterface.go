// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-sitemap/assets"
	"sync"
)

// Ensure, that FileSystemInterfaceMock does implement assets.FileSystemInterface.
// If this is not the case, regenerate this file with moq.
var _ assets.FileSystemInterface = &FileSystemInterfaceMock{}

// FileSystemInterfaceMock is a mock implementation of assets.FileSystemInterface.
//
//	func TestSomethingThatUsesFileSystemInterface(t *testing.T) {
//
//		// make and configure a mocked assets.FileSystemInterface
//		mockedFileSystemInterface := &FileSystemInterfaceMock{
//			GetFunc: func(contextMoqParam context.Context, path string) ([]byte, error) {
//				panic("mock out the Get method")
//			},
//		}
//
//		// use mockedFileSystemInterface in code that requires assets.FileSystemInterface
//		// and then make assertions.
//
//	}
type FileSystemInterfaceMock struct {
	// GetFunc mocks the Get method.
	GetFunc func(contextMoqParam context.Context, path string) ([]byte, error)

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Path is the path argument value.
			Path string
		}
	}
	lockGet sync.RWMutex
}

// Get calls GetFunc.
func (mock *FileSystemInterfaceMock) Get(contextMoqParam context.Context, path string) ([]byte, error) {
	if mock.GetFunc == nil {
		panic("FileSystemInterfaceMock.GetFunc: method is nil but FileSystemInterface.Get was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Path            string
	}{
		ContextMoqParam: contextMoqParam,
		Path:            path,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(contextMoqParam, path)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedFileSystemInterface.GetCalls())
func (mock *FileSystemInterfaceMock) GetCalls() []struct {
	ContextMoqParam context.Context
	Path            string
} {
	var calls []struct {
		ContextMoqParam context.Context
		Path            string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}
