package robotseo

import "embed"

//go:embed robot/robot_en.json
//go:embed robot/robot_cy.json

var folder embed.FS

func GetRobotFile(filename string) ([]byte, error) {
	file, err := folder.ReadFile("robot/" + filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
