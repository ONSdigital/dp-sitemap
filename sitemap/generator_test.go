package sitemap_test

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateIncrementalSitemap(t *testing.T) {
	store := &mock.FileStoreMock{}
	adder := &mock.AdderMock{}
	fetcher := &mock.FetcherMock{}
	fetcher.URLVersionsFunc = func(ctx context.Context, path, lastmod string) (sitemap.URL, *sitemap.URL) {
		return sitemap.URL{Loc: path, Lastmod: lastmod}, nil
	}
	Convey("When getting current sitemap returns an error", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			So(name, ShouldEqual, "sitemap.xml")
			return nil, errors.New("get file error")
		}

		g := sitemap.NewGenerator(nil, nil, store)
		err := g.MakeIncrementalSitemap(context.Background(), "sitemap.xml", sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to get current sitemap")
			So(err.Error(), ShouldContainSubstring, "get file error")
		})
	})
	Convey("When getting current sitemap returns an error", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			So(name, ShouldEqual, "sitemap.xml")
			return nil, errors.New("get file error")
		}

		g := sitemap.NewGenerator(nil, nil, store)
		err := g.MakeIncrementalSitemap(context.Background(), "sitemap.xml", sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to get current sitemap")
			So(err.Error(), ShouldContainSubstring, "get file error")
		})
	})
	Convey("When adder returns an error", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		adder.AddFunc = func(oldSitemap io.Reader, url sitemap.URL) (string, error) {
			return "", errors.New("adder error")
		}

		g := sitemap.NewGenerator(fetcher, adder, store)
		err := g.MakeIncrementalSitemap(context.Background(), "", sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to add to sitemap")
			So(err.Error(), ShouldContainSubstring, "adder error")
		})
	})

	Convey("When adder returns a non-existent file", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		adder.AddFunc = func(oldSitemap io.Reader, url sitemap.URL) (string, error) {
			return "filename", nil
		}

		g := sitemap.NewGenerator(fetcher, adder, store)
		err := g.MakeIncrementalSitemap(context.Background(), "", sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open incremental sitemap")
			So(err.Error(), ShouldContainSubstring, "no such file")
		})
	})

	Convey("When adder returns a file with known content", t, func() {
		store := &mock.FileStoreMock{}
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		var tempFile string
		adder.AddFunc = func(oldSitemap io.Reader, url sitemap.URL) (string, error) {
			So(url, ShouldResemble, sitemap.URL{Loc: "a", Lastmod: "b"})
			file, err := os.CreateTemp("", "sitemap-incr")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return tempFile, nil
		}
		var uploadedFile string
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil
		}

		g := sitemap.NewGenerator(fetcher, adder, store)
		err := g.MakeIncrementalSitemap(context.Background(), "sitemap.xml", sitemap.URL{Loc: "a", Lastmod: "b"})

		Convey("Generator should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to store", func() {
			So(uploadedFile, ShouldEqual, "file content")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
	Convey("When save file returns with an error", t, func() {
		store := &mock.FileStoreMock{}
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		var tempFile string
		adder.AddFunc = func(oldSitemap io.Reader, url sitemap.URL) (string, error) {
			file, err := os.CreateTemp("", "sitemap-incr")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return tempFile, nil
		}
		var uploadedFile string
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return errors.New("uploader error")
		}

		g := sitemap.NewGenerator(fetcher, adder, store)
		err := g.MakeIncrementalSitemap(context.Background(), "sitemap.xml", sitemap.URL{})

		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to store", func() {
			So(uploadedFile, ShouldEqual, "file content")
		})
		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to save incremental sitemap file")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
}

func TestGenerateFullSitemap(t *testing.T) {
	store := &mock.FileStoreMock{}
	fetcher := &mock.FetcherMock{}

	fetcher.HasWelshContentFunc = func(ctx context.Context, path string) bool { return false }
	Convey("When fetcher returns an error", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (sitemap.Files, error) {
			return nil, errors.New("fetcher error")
		}

		g := sitemap.NewGenerator(fetcher, nil, store)
		err := g.MakeFullSitemap(context.Background(), nil)

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to fetch sitemap")
			So(err.Error(), ShouldContainSubstring, "fetcher error")
		})
	})

	Convey("When fetcher returns a non-existent file", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (sitemap.Files, error) {
			return sitemap.Files{config.English: "filename"}, nil
		}

		g := sitemap.NewGenerator(fetcher, nil, store)
		err := g.MakeFullSitemap(context.Background(), nil)

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open sitemap")
			So(err.Error(), ShouldContainSubstring, "no such file")
		})
	})

	Convey("When fetcher returns a file with known content", t, func() {
		var tempFile string
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (sitemap.Files, error) {
			file, err := os.CreateTemp("", "sitemap")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return sitemap.Files{config.English: tempFile}, nil
		}
		var uploadedFile string
		store := &mock.FileStoreMock{}
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil
		}

		g := sitemap.NewGenerator(fetcher, nil, store)
		err := g.MakeFullSitemap(context.Background(), sitemap.Files{config.English: "sitemap.xml"})

		Convey("Generator should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to store", func() {
			So(uploadedFile, ShouldEqual, "file content")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})

	Convey("When save file returns with an error", t, func() {
		var tempFile string
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (sitemap.Files, error) {
			file, err := os.CreateTemp("", "sitemap")
			So(err, ShouldBeNil)
			file.WriteString("file content")
			tempFile = file.Name()
			return sitemap.Files{config.English: tempFile}, nil
		}
		var uploadedFile string
		store := &mock.FileStoreMock{}

		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return errors.New("uploader error")
		}

		g := sitemap.NewGenerator(fetcher, nil, store)
		err := g.MakeFullSitemap(context.Background(), sitemap.Files{config.English: "sitemap.xml"})

		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to store", func() {
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
