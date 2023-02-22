package sitemap_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	s3 := &mock.S3UploaderMock{}
	fetcher := &mock.FetcherMock{}
	cfg := &config.S3Config{
		UploadBucketName: "upload-bucket",
		SitemapFileKey:   "sitemap-file-key",
	}

	Convey("When fetcher returns an error", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (string, error) {
			return "", errors.New("fetcher error")
		}

		g := sitemap.NewGenerator(fetcher, s3, cfg)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to fetch sitemap")
			So(err.Error(), ShouldContainSubstring, "fetcher error")
		})
	})

	Convey("When fetcher returns a non-existent file", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (string, error) {
			return "filename", nil
		}

		g := sitemap.NewGenerator(fetcher, s3, cfg)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open sitemap")
			So(err.Error(), ShouldContainSubstring, "no such file")
		})
	})

	Convey("When fetcher returns a file with known content", t, func() {
		var tempFile string
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (string, error) {
			file, err := os.CreateTemp("", "sitemap")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return tempFile, nil
		}
		var uploadedFile string
		s3 := &mock.S3UploaderMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, nil
		}

		g := sitemap.NewGenerator(fetcher, s3, cfg)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Generator should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, cfg.UploadBucketName)
			So(*params.Key, ShouldEqual, cfg.SitemapFileKey)
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})

	Convey("When uploader returns with an error", t, func() {
		var tempFile string
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (string, error) {
			file, err := os.CreateTemp("", "sitemap")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return tempFile, nil
		}
		var uploadedFile string
		s3 := &mock.S3UploaderMock{}
		s3.UploadFunc = func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
			body, err := io.ReadAll(input.Body)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil, errors.New("uploader error")
		}

		g := sitemap.NewGenerator(fetcher, s3, cfg)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should call uploader", func() {
			So(s3.UploadCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct params to uploader", func() {
			params := s3.UploadCalls()[0].Input

			So(uploadedFile, ShouldEqual, "file content")
			So(*params.Bucket, ShouldEqual, cfg.UploadBucketName)
			So(*params.Key, ShouldEqual, cfg.SitemapFileKey)
		})
		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to upload sitemap")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
}
