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

func TestGeneratePublishingSitemap(t *testing.T) {
	store := &mock.FileStoreMock{}
	adder := &mock.AdderMock{}
	fetcher := &mock.FetcherMock{}
	fetcher.URLVersionsFunc = func(ctx context.Context, path, lastmod string) (*sitemap.URL, *sitemap.URL) {
		return &sitemap.URL{Loc: path, Lastmod: lastmod}, nil
	}
	Convey("When getting current sitemap returns an error", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			So(name, ShouldEqual, "sitemap.xml")
			return nil, errors.New("get file error")
		}

		g := sitemap.NewGenerator(
			sitemap.WithFileStore(store),
			sitemap.WithPublishingSitemapFile("sitemap.xml"),
		)

		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{})

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

		g := sitemap.NewGenerator(
			sitemap.WithFileStore(store),
			sitemap.WithPublishingSitemapFile("sitemap.xml"),
		)
		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to get current sitemap")
			So(err.Error(), ShouldContainSubstring, "get file error")
		})
	})
	Convey("When adder returns an error", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		adder.AddFunc = func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
			return "", 0, errors.New("adder error")
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithAdder(adder),
			sitemap.WithFileStore(store),
		)
		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to add to sitemap")
			So(err.Error(), ShouldContainSubstring, "adder error")
		})
	})

	Convey("When adder returns a non-existent file", t, func() {
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		adder.AddFunc = func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
			return "filename", 0, nil
		}
		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithAdder(adder),
			sitemap.WithFileStore(store),
		)
		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{})

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open publishing sitemap")
			So(err.Error(), ShouldContainSubstring, "no such file")
		})
	})

	Convey("When adder returns a file with known content", t, func() {
		store := &mock.FileStoreMock{}
		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}
		var tempFile string
		adder.AddFunc = func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
			So(url, ShouldResemble, &sitemap.URL{Loc: "a", Lastmod: "b"})
			file, err := os.CreateTemp("", "sitemap-incr")
			So(err, ShouldBeNil)
			_, err = file.WriteString("file content")
			So(err, ShouldBeNil)
			tempFile = file.Name()
			return tempFile, 1, nil
		}
		var uploadedFile string
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithAdder(adder),
			sitemap.WithFileStore(store),
			sitemap.WithPublishingSitemapFile("sitemap.xml"),
		)
		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{Loc: "a", Lastmod: "b"})

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
		adder.AddFunc = func(oldSitemap io.Reader, url *sitemap.URL) (string, int, error) {
			file, err := os.CreateTemp("", "sitemap-incr")
			So(err, ShouldBeNil)
			_, err = file.WriteString("file content")
			So(err, ShouldBeNil)
			tempFile = file.Name()
			return tempFile, 1, nil
		}
		var uploadedFile string
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "sitemap.xml")
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return errors.New("uploader error")
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithAdder(adder),
			sitemap.WithFileStore(store),
			sitemap.WithPublishingSitemapFile("sitemap.xml"),
		)
		err := g.MakePublishingSitemap(context.Background(), sitemap.URL{})

		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should pass correct file content to store", func() {
			So(uploadedFile, ShouldEqual, "file content")
		})
		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to save publishing sitemap file")
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
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "publishing-sitemap.xml")
			return nil
		}
		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithFileStore(store),
			sitemap.WithAdder(&sitemap.DefaultAdder{}),
		)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to fetch sitemap")
			So(err.Error(), ShouldContainSubstring, "fetcher error")
		})
	})

	Convey("When fetcher returns a non-existent file", t, func() {
		fetcher.GetFullSitemapFunc = func(ctx context.Context) (sitemap.Files, error) {
			return sitemap.Files{config.English: "filename"}, nil
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithFileStore(store),
			sitemap.WithAdder(&sitemap.DefaultAdder{}),
		)
		err := g.MakeFullSitemap(context.Background())

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
			_, err = file.WriteString("file content")
			So(err, ShouldBeNil)
			tempFile = file.Name()
			return sitemap.Files{config.English: tempFile}, nil
		}
		var uploadedFile string
		store := &mock.FileStoreMock{}
		store.SaveFileFunc = func(name string, reader io.Reader) error {
			body, err := io.ReadAll(reader)
			So(err, ShouldBeNil)
			uploadedFile = string(body)
			return nil
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithFileStore(store),
			sitemap.WithAdder(&sitemap.DefaultAdder{}),
		)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 2)
		})
		Convey("Generator should truncate publishing file and write the main file", func() {
			So(store.SaveFileCalls()[0].Name, ShouldEqual, "publishing-sitemap.xml")
			So(store.SaveFileCalls()[1].Name, ShouldEqual, "sitemap.xml")
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
			_, err = file.WriteString("file content")
			So(err, ShouldBeNil)
			tempFile = file.Name()
			return sitemap.Files{config.English: tempFile}, nil
		}
		store := &mock.FileStoreMock{}

		store.SaveFileFunc = func(name string, reader io.Reader) error {
			So(name, ShouldEqual, "publishing-sitemap.xml")
			return errors.New("uploader error")
		}

		g := sitemap.NewGenerator(
			sitemap.WithFetcher(fetcher),
			sitemap.WithFileStore(store),
			sitemap.WithAdder(&sitemap.DefaultAdder{}),
		)
		err := g.MakeFullSitemap(context.Background())

		Convey("Generator should call store", func() {
			So(store.SaveFileCalls(), ShouldHaveLength, 1)
		})
		Convey("Generator should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to save publishing sitemap file")
			So(err.Error(), ShouldContainSubstring, "uploader error")
		})
		Convey("Generator should remove the temporary file", func() {
			_, err := os.Stat(tempFile)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
}
