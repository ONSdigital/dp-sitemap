package sitemap_test

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"
	"testing/iotest"

	. "github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalSaver(t *testing.T) {
	dir := os.TempDir()

	Convey("When an invalid file path is provided", t, func() {
		s := sitemap.NewLocalSaver(map[Language]string{English: ""})
		err := s.SaveFile(English, strings.NewReader("file content"))

		Convey("LocalSaver should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to open a local file")
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})

	Convey("When an invalid file content is provided", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())

		s := sitemap.NewLocalSaver(map[Language]string{English: randomFilename})
		invalidBody := iotest.ErrReader(errors.New("invalid body"))
		err := s.SaveFile(English, invalidBody)
		defer os.Remove(randomFilename)

		Convey("LocalSaver should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to copy to a local file")
			So(err.Error(), ShouldContainSubstring, "invalid body")
		})
	})

	Convey("When local save is successful", t, func() {
		randomFilename := path.Join(dir, "sitemap-test-"+uuid.NewString())

		s := sitemap.NewLocalSaver(map[Language]string{English: randomFilename})
		err := s.SaveFile(English, strings.NewReader("file content"))
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
}
