package sitemap_test

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalStore(t *testing.T) {
	dir := os.TempDir()

	Convey("When an invalid file path is provided", t, func() {
		s := &sitemap.LocalStore{}
		err := s.SaveFile("", strings.NewReader("file content"))

		Convey("LocalSaver should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open a local file")
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})

	Convey("When an invalid file content is provided", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())

		s := &sitemap.LocalStore{}
		invalidBody := iotest.ErrReader(errors.New("invalid body"))
		err := s.SaveFile(randomFilename, invalidBody)
		defer os.Remove(randomFilename)

		Convey("LocalSaver should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to copy to a local file")
			So(err.Error(), ShouldContainSubstring, "invalid body")
		})
	})

	Convey("When local save is successful", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())

		s := &sitemap.LocalStore{}
		err := s.SaveFile(randomFilename, strings.NewReader("file content"))
		defer os.Remove(randomFilename)

		Convey("SaveFile should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Correct file should be available locally", func() {
			content, err := os.ReadFile(randomFilename)
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "file content")
		})
	})

	Convey("When local read fails", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())

		s := &sitemap.LocalStore{}
		body, err := s.GetFile(randomFilename)

		Convey("LocalStore should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Body should be empty", func() {
			content, err := io.ReadAll(body)
			So(err, ShouldBeNil)
			So(content, ShouldHaveLength, 0)
		})
	})

	Convey("When local read succeeds", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())
		err := os.WriteFile(randomFilename, []byte("file content"), 0o600)
		defer os.Remove(randomFilename)

		s := &sitemap.LocalStore{}
		body, err := s.GetFile(randomFilename)

		Convey("LocalStore should return no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("LocalStore should return correct file body", func() {
			content, err := io.ReadAll(body)
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "file content")
		})
	})
}
