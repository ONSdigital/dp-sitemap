package robotseo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/exp/slices"
)

var robotList map[string]map[string]assets.SeoRobotModel

func Init(efs assets.FileSystemInterface) {
	robotList = map[string]map[string]assets.SeoRobotModel{}
	ctx := context.Background()
	for _, lang := range []string{config.English.String(), config.Welsh.String()} {
		fileName := "robot_" + lang + ".json"
		b, err := efs.Get(ctx, fileName)
		if err != nil {
			log.Error(ctx, "can't find "+fileName, err)
			panic("Can't find " + fileName)
		}

		rContent := map[string]assets.SeoRobotModel{}
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
	for _, lang := range []string{config.English.String(), config.Welsh.String()} {
		fileName := "robot_" + lang + ".json"
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
