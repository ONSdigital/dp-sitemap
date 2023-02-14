package robotseo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/exp/slices"
)

var robotList map[string]SeoRobotModel

func Init(asset func(name string) ([]byte, error)) {
	ctx := context.Background()
	b, err := asset("robot/robot.json")
	if err != nil {
		log.Error(ctx, "can't find robot.json", err)
		panic("Can't find robot.json")
	}

	robotList = map[string]SeoRobotModel{}
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

func GetRobotsFileBody() string {
	robot := strings.Builder{}
	for k, v := range robotList {
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
