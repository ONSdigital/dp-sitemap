// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"sync"
)

// Ensure, that S3ClientMock does implement sitemap.S3Client.
// If this is not the case, regenerate this file with moq.
var _ sitemap.S3Client = &S3ClientMock{}

// S3ClientMock is a mock implementation of sitemap.S3Client.
//
// 	func TestSomethingThatUsesS3Client(t *testing.T) {
//
// 		// make and configure a mocked sitemap.S3Client
// 		mockedS3Client := &S3ClientMock{
// 			BucketNameFunc: func() string {
// 				panic("mock out the BucketName method")
// 			},
// 			GetFunc: func(key string) (io.ReadCloser, *int64, error) {
// 				panic("mock out the Get method")
// 			},
// 			UploadFunc: func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
// 				panic("mock out the Upload method")
// 			},
// 		}
//
// 		// use mockedS3Client in code that requires sitemap.S3Client
// 		// and then make assertions.
//
// 	}
type S3ClientMock struct {
	// BucketNameFunc mocks the BucketName method.
	BucketNameFunc func() string

	// GetFunc mocks the Get method.
	GetFunc func(key string) (io.ReadCloser, *int64, error)

	// UploadFunc mocks the Upload method.
	UploadFunc func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// BucketName holds details about calls to the BucketName method.
		BucketName []struct {
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Key is the key argument value.
			Key string
		}
		// Upload holds details about calls to the Upload method.
		Upload []struct {
			// Input is the input argument value.
			Input *s3manager.UploadInput
			// Options is the options argument value.
			Options []func(*s3manager.Uploader)
		}
	}
	lockBucketName sync.RWMutex
	lockGet        sync.RWMutex
	lockUpload     sync.RWMutex
}

// BucketName calls BucketNameFunc.
func (mock *S3ClientMock) BucketName() string {
	if mock.BucketNameFunc == nil {
		panic("S3ClientMock.BucketNameFunc: method is nil but S3Client.BucketName was just called")
	}
	callInfo := struct {
	}{}
	mock.lockBucketName.Lock()
	mock.calls.BucketName = append(mock.calls.BucketName, callInfo)
	mock.lockBucketName.Unlock()
	return mock.BucketNameFunc()
}

// BucketNameCalls gets all the calls that were made to BucketName.
// Check the length with:
//     len(mockedS3Client.BucketNameCalls())
func (mock *S3ClientMock) BucketNameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBucketName.RLock()
	calls = mock.calls.BucketName
	mock.lockBucketName.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *S3ClientMock) Get(key string) (io.ReadCloser, *int64, error) {
	if mock.GetFunc == nil {
		panic("S3ClientMock.GetFunc: method is nil but S3Client.Get was just called")
	}
	callInfo := struct {
		Key string
	}{
		Key: key,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(key)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedS3Client.GetCalls())
func (mock *S3ClientMock) GetCalls() []struct {
	Key string
} {
	var calls []struct {
		Key string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Upload calls UploadFunc.
func (mock *S3ClientMock) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if mock.UploadFunc == nil {
		panic("S3ClientMock.UploadFunc: method is nil but S3Client.Upload was just called")
	}
	callInfo := struct {
		Input   *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}{
		Input:   input,
		Options: options,
	}
	mock.lockUpload.Lock()
	mock.calls.Upload = append(mock.calls.Upload, callInfo)
	mock.lockUpload.Unlock()
	return mock.UploadFunc(input, options...)
}

// UploadCalls gets all the calls that were made to Upload.
// Check the length with:
//     len(mockedS3Client.UploadCalls())
func (mock *S3ClientMock) UploadCalls() []struct {
	Input   *s3manager.UploadInput
	Options []func(*s3manager.Uploader)
} {
	var calls []struct {
		Input   *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}
	mock.lockUpload.RLock()
	calls = mock.calls.Upload
	mock.lockUpload.RUnlock()
	return calls
}
