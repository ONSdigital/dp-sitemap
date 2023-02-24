package steps

import (
	"context"
	"fmt"
	"os"
	"strings"

	assetmock "github.com/ONSdigital/dp-sitemap/assets/mock"
	"github.com/ONSdigital/dp-sitemap/robotseo"
	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^i have the following robot\.json:$`, iHaveTheFollowingRobotjson)
	ctx.Step(`^i invoke writejson with the sitemaps "([^"]*)"$`, c.iInvokeWritejsonWithTheSitemaps)
	ctx.Step(`^the content of the resulting robots file must be$`, c.theContentOfTheResultingRobotsFileMustBe)
}

func iHaveTheFollowingRobotjson(arg1 *godog.DocString) error {
	amock := assetmock.FileSystemInterfaceMock{
		GetFunc: func(contextMoqParam context.Context, path string) ([]byte, error) { return []byte(arg1.Content), nil },
	}
	robotseo.Init(&amock)
	return nil
}

func (c *Component) iInvokeWritejsonWithTheSitemaps(arg1 string) error {
	fw := robotseo.RobotFileWriter{}
	return fw.WriteRobotsFile(c.cfg, strings.Split(arg1, ","))
}

func (c *Component) theContentOfTheResultingRobotsFileMustBe(arg1 *godog.DocString) error {
	b, err := os.ReadFile(c.cfg.RobotsFilePath)
	if err != nil {
		return err
	}
	if strings.Compare(arg1.Content, string(b)) != 0 {
		return fmt.Errorf("robot file content mismatch actual [%s], expecting [%s]", string(b), arg1.Content)
	}
	return nil
}
