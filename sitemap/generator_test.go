package sitemap_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	saver := &mock.FileSaverMock{}
	fetcher := &mock.FetcherMock{}

	Convey("When fetcher returns an error", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (string, error) {
			return "", errors.New("fetcher error")
		}

		g := sitemap.NewGenerator(fetcher, saver)
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

		g := sitemap.NewGenerator(fetcher, saver)
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
		saver := &mock.FileSaverMock{}
		saver.SaveFileFunc = func(reader io.Reader) error {
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil
		}

		g := sitemap.NewGenerator(fetcher, saver)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Generator should call saver", func() {
			So(saver.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to saver", func() {
			So(uploadedFile, ShouldEqual, "file content")
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
		saver := &mock.FileSaverMock{}

		saver.SaveFileFunc = func(reader io.Reader) error {
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return errors.New("uploader error")
		}

		g := sitemap.NewGenerator(fetcher, saver)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should call saver", func() {
			So(saver.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to saver", func() {
			So(uploadedFile, ShouldEqual, "file content")
		})
		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to save sitemap file")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
}
