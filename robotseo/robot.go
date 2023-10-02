package robotseo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/features"
	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/exp/slices"
)

var robotList map[config.Language]map[string]SeoRobotModel

func Init(pathToRobotFile string) {
	robotList = map[config.Language]map[string]SeoRobotModel{}
	ctx := context.Background()
	var b []byte
	var err error
	var fileName string

	if !strings.HasSuffix(pathToRobotFile, "/") {
		pathToRobotFile += "/"
	}
	for _, lang := range []config.Language{config.English, config.Welsh} {
		fileName = "robot_" + lang.String() + ".json"

		// if pathToRobotFile is empty (the default) we get the robot file for the component tests
		// otherwise we get the robot file from the local file store and the path specified by pathToRobotFile
		if pathToRobotFile == "" {
			b, err = features.GetRobotFile(fileName)
		} else {
			b, err = os.ReadFile(pathToRobotFile + fileName)
		}

		if err != nil {
			log.Error(ctx, "can't find "+fileName, err)
			panic("Can't find " + fileName)
		}

		rContent := map[string]SeoRobotModel{}
		err = json.Unmarshal(b, &rContent)
		if err != nil {
			log.Error(ctx, "error reading "+fileName, err)
			panic("Unable to read JSON")
		}
		robotList[lang] = rContent
	}

	// Validation
	// 1. Check there is at least 1 entry
	// 2. Check that same allow/deny dont exist for a user-agent
	for _, lang := range []config.Language{config.English, config.Welsh} {
		fileName := "robot_" + lang.String() + ".json"
		rList := robotList[lang]
		if len(rList) == 0 {
			log.Error(ctx, "no entry in "+fileName, errors.New(fileName+" cant be empty"))
			panic(fileName + " cant be empty")
		}
		for ua, list := range rList {
			if len(list.AllowList) == 0 || len(list.DenyList) == 0 {
				continue
			}
			for _, allow := range list.AllowList {
				if slices.Contains(list.DenyList, allow) {
					panic(fmt.Sprintf("user agent [%s], contains [%s] in both allow and deny", ua, allow))
				}
			}
		}
	}
}
