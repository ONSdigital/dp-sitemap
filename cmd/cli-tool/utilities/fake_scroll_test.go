package utilities

import (
	"testing"

	"github.com/ONSdigital/dp-sitemap/sitemap"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFakeStartScroll(t *testing.T) {
	Convey("When valid input", t, func() {
		testdata := sitemap.ElasticResult{}
		Convey("Then the are no errors", func() {
			result := fakeStartScroll(&testdata)
			So(result, ShouldBeEmpty)
		})
	})

	Convey("When not valid input", t, func() {
		testdata := sitemap.ElasticHit{}
		Convey("Then there is an error", func() {
			result := fakeStartScroll(&testdata)
			So(result, ShouldBeError)
		})
	})
}
