package robotseo

import (
	"strings"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/pkg/errors"
)

//go:generate moq -out mock/robotFileWriter.go -pkg mock . RobotFileWriterInterface
type RobotFileWriterInterface interface {
	GetRobotsFileBody(lang string, sitemap map[config.Language]string) string
}
type RobotFileWriter struct {
}

var (
	ErrNoRobotsBody     = errors.New("no robots body")
	ErrNoRobotsFilePath = errors.New("no robots file path given")
)

func (r *RobotFileWriter) GetRobotsFileBody(lang config.Language, sitemap map[config.Language]string) string {
	robot := strings.Builder{}
	for k, v := range robotList[lang] {
		robot.WriteString("\nUser-agent: " + k)
		for _, allow := range v.AllowList {
			robot.WriteString("\nAllow: " + allow)
		}
		for _, deny := range v.DenyList {
			robot.WriteString("\nDisallow: " + deny)
		}
		if sm, ok := sitemap[lang]; ok {
			robot.WriteString("\n\nsitemap: " + sm)
		}
		robot.WriteString("\n")
	}
	return robot.String()
}
