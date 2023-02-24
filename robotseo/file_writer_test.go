package robotseo

import (
	"testing"

	asset "github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRobotFileWriter_WriteRobotsFile(t *testing.T) {
	var expectedError error
	r := &RobotFileWriter{}
	cfg, _ := config.Get()

	Convey("no file path provided", t, func() {
		robotList = map[string]asset.SeoRobotModel{}
		expectedError = ErrNoRobotsFilePath
		cfg.RobotsFilePath = ""
		So(r.WriteRobotsFile(cfg, []string{}), ShouldEqual, expectedError)
	})

	Convey("no robots body", t, func() {
		robotList = map[string]asset.SeoRobotModel{}
		expectedError = ErrNoRobotsBody
		cfg.RobotsFilePath = "/tmp/dp_robot.txt"
		So(r.WriteRobotsFile(cfg, []string{}), ShouldEqual, expectedError)
	})
}
