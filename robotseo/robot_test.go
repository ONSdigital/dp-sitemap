package robotseo

import (
	"context"
	"errors"
	"testing"

	asset "github.com/ONSdigital/dp-sitemap/assets"
	mockassets "github.com/ONSdigital/dp-sitemap/assets/mock"
	"github.com/ONSdigital/dp-sitemap/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	var shouldError bool
	var returnBytes []byte
	var called bool

	fsMock := mockassets.FileSystemInterfaceMock{GetFunc: func(contextMoqParam context.Context, path string) ([]byte, error) {
		called = true
		if shouldError {
			return nil, errors.New("error")
		}
		return returnBytes, nil
	}}

	Convey("Init calls asset function and panics on error", t, func() {
		shouldError = true
		So(func() { Init(&fsMock) }, ShouldPanicWith, "Can't find robot_en.json")
		So(called, ShouldBeTrue)
	})

	Convey("bad data in json panics on error", t, func() {
		shouldError = false
		returnBytes = []byte(`a,b`)
		So(func() { Init(&fsMock) }, ShouldPanicWith, "Unable to read JSON")
		So(called, ShouldBeTrue)
	})

	Convey("bad data in json panics on error", t, func() {
		shouldError = false
		returnBytes = []byte(`{}`)
		So(func() { Init(&fsMock) }, ShouldPanicWith, "robot_en.json cant be empty")
		So(called, ShouldBeTrue)
	})

	Convey("happy path scenario", t, func() {
		shouldError = false
		returnBytes = []byte(`{
			"Googlebot": {
			  "AllowList": ["/googlecontent"],
			  "DenyList":  ["/googlecontent"]
			}
		}`)
		So(func() { Init(&fsMock) }, ShouldPanicWith, "user agent [Googlebot], contains [/googlecontent] in both allow and deny")
		So(called, ShouldBeTrue)
	})

	Convey("happy path scenario", t, func() {
		shouldError = false
		returnBytes = []byte(`{
			"Googlebot": {
			  "AllowList": ["/googlecontent"],
			  "DenyList":  ["/googleprivate"]
			},
			"Bingbot": {
				"AllowList": ["/bingcontent"],
				"DenyList":  ["/bingprivate"]
			  },
			  "*": {
				"AllowList": ["/"],
				"DenyList":  ["/private"]
			  }
		}`)
		So(func() { Init(&fsMock) }, ShouldNotPanic)
		So(called, ShouldBeTrue)
		So(len(robotList[config.English]), ShouldEqual, 3)
	})
}

func TestGetRobotsFileBody(t *testing.T) {
	robotList = map[config.Language]map[string]asset.SeoRobotModel{}
	var expectedRobotsBody string
	r := RobotFileWriter{}

	Convey("no robots data", t, func() {
		robotList[config.English] = map[string]asset.SeoRobotModel{}
		expectedRobotsBody = ""
		So(r.GetRobotsFileBody(config.English, map[config.Language]string{}), ShouldEqual, expectedRobotsBody)
	})

	Convey("simple allow/deny with one user-agent", t, func() {
		robotList[config.English] = map[string]asset.SeoRobotModel{
			"GoogleBot": {AllowList: []string{"/googleallow"}, DenyList: []string{"/googledeny"}}}
		expectedRobotsBody = `
User-agent: GoogleBot
Allow: /googleallow
Disallow: /googledeny
`
		So(r.GetRobotsFileBody(config.English, map[config.Language]string{}), ShouldEqual, expectedRobotsBody)
	})

	Convey("multiple allow/deny with one user-agent", t, func() {
		robotList[config.English] = map[string]asset.SeoRobotModel{
			"GoogleBot": {AllowList: []string{"/googleallow1", "/googleallow2"}, DenyList: []string{"/googledeny1", "/googledeny2"}}}
		expectedRobotsBody = `
User-agent: GoogleBot
Allow: /googleallow1
Allow: /googleallow2
Disallow: /googledeny1
Disallow: /googledeny2
`
		So(r.GetRobotsFileBody(config.English, map[config.Language]string{}), ShouldEqual, expectedRobotsBody)
	})

	Convey("multiple allow/deny with multiple user-agents", t, func() {
		robotList[config.English] = map[string]asset.SeoRobotModel{
			"BingBot":   {AllowList: []string{"/bingallow1", "/bingallow2"}, DenyList: []string{"/bingdeny1", "/bingdeny2"}},
			"GoogleBot": {AllowList: []string{"/googleallow1", "/googleallow2"}, DenyList: []string{"/googledeny1", "/googledeny2"}}}
		bot1 := `
User-agent: BingBot
Allow: /bingallow1
Allow: /bingallow2
Disallow: /bingdeny1
Disallow: /bingdeny2
`
		bot2 := `
User-agent: GoogleBot
Allow: /googleallow1
Allow: /googleallow2
Disallow: /googledeny1
Disallow: /googledeny2
`
		robotFile := r.GetRobotsFileBody(config.English, map[config.Language]string{})
		So(robotFile, ShouldContainSubstring, bot1)
		So(robotFile, ShouldContainSubstring, bot2)
	})
}
