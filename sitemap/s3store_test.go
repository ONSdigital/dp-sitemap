package sitemap_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
)

func TestS3Store(t *testing.T) {
	bucket := "upload-bucket"
	fileKey := "sitemap-file-key"

	Convey("When s3 upload fails", t, func() {
		var uploadedFile string
		s3 := &mock.S3ClientMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, errors.New("uploader error")
		}
		s3.BucketNameFunc = func() string { return bucket }

		s := sitemap.NewS3Store(s3)
		err := s.SaveFile(fileKey, strings.NewReader("file content"))

		Convey("S3Store should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to upload file to s3")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("S3Store should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Store should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, bucket)
			So(*params.Key, ShouldEqual, fileKey)
		})
	})

	Convey("When s3 upload is successful", t, func() {
		var uploadedFile string
		s3 := &mock.S3ClientMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, nil
		}
		s3.BucketNameFunc = func() string { return bucket }

		s := sitemap.NewS3Store(s3)
		err := s.SaveFile(fileKey, strings.NewReader("file content"))

		Convey("SaveFile should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("S3Store should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Store should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, bucket)
			So(*params.Key, ShouldEqual, fileKey)
		})
	})

	Convey("When s3 get fails", t, func() {
		s3 := &mock.S3ClientMock{}
		s3.GetFunc = func(key string) (io.ReadCloser, *int64, error) { return nil, nil, errors.New("s3 get error") }

		s := sitemap.NewS3Store(s3)
		body, err := s.GetFile(fileKey)

		Convey("S3Store should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Body should be empty", func() {
			content, err := io.ReadAll(body)
			So(err, ShouldBeNil)
			So(content, ShouldHaveLength, 0)
		})
		Convey("S3Store should call s3 get", func() {
			So(s3.GetCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Store should pass correct file key to s3 get", func() {
			So(s3.GetCalls()[0].Key, ShouldEqual, fileKey)
		})
	})

	Convey("When s3 get succeeds", t, func() {
		s3 := &mock.S3ClientMock{}
		s3.GetFunc = func(key string) (io.ReadCloser, *int64, error) {
			return io.NopCloser(strings.NewReader("file content")), nil, nil
		}

		s := sitemap.NewS3Store(s3)
		body, err := s.GetFile(fileKey)

		Convey("S3Store should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("S3Store should return correct file body", func() {
			content, err := io.ReadAll(body)
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "file content")
		})
		Convey("S3Store should call s3 get", func() {
			So(s3.GetCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Store should pass correct file key to s3 get", func() {
			So(s3.GetCalls()[0].Key, ShouldEqual, fileKey)
		})
	})
}
