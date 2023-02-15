package assets

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestRobotFileValidity(t *testing.T) {
	t.Parallel()
	b, err := os.ReadFile("robot/robot.json")
	assert.Nil(t, err, "unable to read")
	err = fastjson.Validate(string(b))
	assert.Nil(t, err, "json file is invalid")
}

// To detect any json definition issues (that are not caught by above test-case)
// esp: duplicates
func TestMarshalUnmarshalRoundtrip(t *testing.T) {
	robot := map[string]SeoRobotModel{}
	cOne := &bytes.Buffer{}
	cTwo := &bytes.Buffer{}
	b, err := os.ReadFile("robot/robot.json")
	assert.Nil(t, err, "unable to read json file")
	err = json.Unmarshal(b, &robot)
	assert.Nil(t, err, "json file unmarshal error")
	by, err := json.Marshal(&robot)
	assert.Nil(t, err, "json file marshal error")

	json.Compact(cOne, b)
	json.Compact(cTwo, by)

	assert.Equal(t, len(cOne.String()), len(cTwo.String()))
}
