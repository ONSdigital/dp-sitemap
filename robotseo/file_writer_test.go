package robotseo

import (
	"testing"

	asset "github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRobotFileWriter_WriteRobotsFile(t *testing.T) {
	robotList = map[string]map[string]asset.SeoRobotModel{}
	var expectedError error
	r := &RobotFileWriter{}
	cfg, _ := config.Get()

	Convey("no file path provided", t, func() {
		robotList[config.English.String()] = map[string]asset.SeoRobotModel{}
		expectedError = ErrNoRobotsFilePath
		cfg.RobotsFilePath[config.English.String()] = ""
		So(r.WriteRobotsFile(cfg, map[string]string{}), ShouldEqual, expectedError)
	})

	Convey("no robots body", t, func() {
		robotList[config.English.String()] = map[string]asset.SeoRobotModel{}
		expectedError = ErrNoRobotsBody
		cfg.RobotsFilePath[config.English.String()] = "/tmp/dp_robot.txt"
		So(r.WriteRobotsFile(cfg, map[string]string{}), ShouldEqual, expectedError)
	})
}
