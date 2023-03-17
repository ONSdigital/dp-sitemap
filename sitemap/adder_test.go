package sitemap_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-sitemap/sitemap"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAdder(t *testing.T) {
	Convey("When xml decode returns an error", t, func() {
		oldSitemap := strings.NewReader("<<<")

		a := &sitemap.DefaultAdder{}
		filename, err := a.Add(oldSitemap, sitemap.URL{})

		Convey("Adder should return correct error", func() {
			So(err.Error(), ShouldContainSubstring, "failed to decode old sitemap")
			So(err.Error(), ShouldContainSubstring, "XML syntax error")
		})
		Convey("Temporary sitemap file should be created and then cleaned up", func() {
			So(filename, ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filename)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "no such file or directory")
		})
	})
	Convey("When old sitemap is empty", t, func() {
		oldSitemap := strings.NewReader("")

		a := &sitemap.DefaultAdder{}
		filename, err := a.Add(oldSitemap, sitemap.URL{Loc: "a", Lastmod: "b"})
		defer func() {
			removeErr := os.Remove(filename)
			So(removeErr, ShouldBeNil)
		}()

		Convey("Adder should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filename, ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filename)
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be valid and include the received url", func() {
			sitemapContent, err := os.ReadFile(filename)
			So(err, ShouldBeNil)
			So(string(sitemapContent), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>a</loc>
    <lastmod>b</lastmod>
  </url>
</urlset>`)
		})
	})
	Convey("When old sitemap contains urls", t, func() {
		oldSitemap := strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?>
		<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
		  <url>
			<loc>a</loc>
			<lastmod>b</lastmod>
		  </url>
		  <url>
		  <loc>c</loc>
		  <lastmod>d</lastmod>
		</url>
		</urlset>`)

		a := &sitemap.DefaultAdder{}
		filename, err := a.Add(oldSitemap, sitemap.URL{Loc: "e", Lastmod: "f"})
		defer func() {
			removeErr := os.Remove(filename)
			So(removeErr, ShouldBeNil)
		}()

		Convey("Adder should return with no error", func() {
			So(err, ShouldBeNil)
		})
		Convey("Temporary sitemap file should be created and available", func() {
			So(filename, ShouldContainSubstring, "sitemap")
			_, err := os.Stat(filename)
			So(err, ShouldBeNil)
		})
		Convey("Sitemap should be valid and include both old and new urls", func() {
			sitemap, err := os.ReadFile(filename)
			So(err, ShouldBeNil)
			So(string(sitemap), ShouldEqual, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>a</loc>
    <lastmod>b</lastmod>
  </url>
  <url>
    <loc>c</loc>
    <lastmod>d</lastmod>
  </url>
  <url>
    <loc>e</loc>
    <lastmod>f</lastmod>
  </url>
</urlset>`)
		})
	})
}
