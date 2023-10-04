package features

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRobotFile(t *testing.T) {
	Convey("Given we want to get the robot file from the embedded fs", t, func() {
		Convey("When we call GetRobotFile func for the english robot file", func() {
			fileName := "robot_en.json"
			_, err := GetRobotFile(fileName)
			Convey("Than there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When we call GetRobotFile func for the welsh robot file", func() {
			fileName := "robot_cy.json"
			_, err := GetRobotFile(fileName)
			Convey("Than there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When the file does not exist in the embedded fs", func() {
			fileName := "robot_gr.json"
			_, err := GetRobotFile(fileName)
			Convey("Than there should be an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
