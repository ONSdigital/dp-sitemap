package robotseo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-sitemap/assets"
	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/exp/slices"
)

var robotList map[string]assets.SeoRobotModel

func Init(efs assets.FileSystemInterface) {
	ctx := context.Background()
	b, err := efs.Get(ctx, "robot.json")
	if err != nil {
		log.Error(ctx, "can't find robot.json", err)
		panic("Can't find robot.json")
	}
	if err != nil {
		log.Error(ctx, "error reading robot.json", err)
		panic("Error reading robot.json")
	}

	robotList = map[string]assets.SeoRobotModel{}
	err = json.Unmarshal(b, &robotList)
	if err != nil {
		log.Error(ctx, "error reading robot.json", err)
		panic("Unable to read JSON")
	}

	// Validation
	// 1. Check there is at least 1 entry
	// 2. Check that same allow/deny dont exist for a user-agent
	if len(robotList) == 0 {
		log.Error(ctx, "no entry in robot.json", errors.New("robots.json cant be empty"))
		panic("robots.json cant be empty")
	}
	for ua, list := range robotList {
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
