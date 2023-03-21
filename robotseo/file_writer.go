package robotseo

import (
	"os"
	"strings"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/pkg/errors"
)

//go:generate moq -out mock/robotFileWriter.go -pkg mock . RobotFileWriterInterface
type RobotFileWriterInterface interface {
	WriteRobotsFile(cfg *config.Config, sitemaps []string) error
	GetRobotsFileBody() string
}
type RobotFileWriter struct {
}

var (
	ErrNoRobotsBody     = errors.New("no robots body")
	ErrNoRobotsFilePath = errors.New("no robots file path given")
)

func (r *RobotFileWriter) WriteRobotsFile(cfg *config.Config, sitemap map[string]string) error {
	for _, lang := range []string{"en", "cy"} {
		if cfg.RobotsFilePath[lang] == "" {
			return ErrNoRobotsFilePath
		}

		robotFile := strings.Builder{}
		body := r.GetRobotsFileBody(lang)
		if body == "" {
			return ErrNoRobotsBody
		}
		_, err := robotFile.WriteString(body)
		if err != nil {
			return errors.Wrap(err, "error writing to buffer")
		}

		robotFile.WriteString("\n")
		if sm, ok := sitemap[lang]; ok {
			robotFile.WriteString("sitemap: " + sm + "\n")
		}

		if err := os.WriteFile(cfg.RobotsFilePath[lang], []byte(robotFile.String()), 0600); err != nil {
			return errors.Wrap(err, "error writing to file")
		}
	}

	return nil
}

func (r *RobotFileWriter) GetRobotsFileBody(lang string) string {
	robot := strings.Builder{}
	for k, v := range robotList[lang] {
		robot.WriteString("\nUser-agent: " + k)
		for _, allow := range v.AllowList {
			robot.WriteString("\nAllow: " + allow)
		}
		for _, deny := range v.DenyList {
			robot.WriteString("\nDisallow: " + deny)
		}
		robot.WriteString("\n")
	}
	return robot.String()
}
