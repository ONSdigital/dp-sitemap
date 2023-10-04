package features

import "embed"

//go:embed steps/robot/robot_en.json
//go:embed steps/robot/robot_cy.json

var folder embed.FS

// GetRobotFile Gets the robot file from the embedded files
func GetRobotFile(filename string) ([]byte, error) {
	file, err := folder.ReadFile("steps/robot/" + filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
