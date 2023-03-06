package sitemap_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	. "github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
)

func TestS3Saver(t *testing.T) {
	bucket := "upload-bucket"
	fileKey := map[Language]string{English: "sitemap-file-key", Welsh: "sitemap-file-key-cy"}

	Convey("When s3 upload fails", t, func() {
		var uploadedFile string
		s3 := &mock.S3UploaderMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, errors.New("uploader error")
		}

		s := sitemap.NewS3Saver(s3, bucket, fileKey)
		err := s.SaveFile(English, strings.NewReader("file content"))

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to upload file to s3")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("S3Saver should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Saver should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, bucket)
			So(*params.Key, ShouldEqual, English.String())
		})
	})

	Convey("When s3 upload is successful", t, func() {
		var uploadedFile string
		s3 := &mock.S3UploaderMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, nil
		}

		s := sitemap.NewS3Saver(s3, bucket, fileKey)
		err := s.SaveFile(English, strings.NewReader("file content"))

		Convey("SaveFile should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("S3Saver should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("S3Saver should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, bucket)
			So(*params.Key, ShouldEqual, English.String())
		})
	})
}
