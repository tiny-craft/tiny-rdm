package convutil

import (
	"bytes"
	"encoding/json"
	"strings"
)

type JsonConvert struct{}

func (JsonConvert) Enable() bool {
	return true
}

func (JsonConvert) Decode(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if (strings.HasPrefix(trimedStr, "{") && strings.HasSuffix(trimedStr, "}")) ||
		(strings.HasPrefix(trimedStr, "[") && strings.HasSuffix(trimedStr, "]")) {
		var out bytes.Buffer
		if err := json.Indent(&out, []byte(trimedStr), "", "  "); err == nil {
			return out.String(), true
		}
	}
	return str, false
}

func (JsonConvert) Encode(str string) (string, bool) {
	var dst bytes.Buffer
	if err := json.Compact(&dst, []byte(str)); err != nil {
		return str, false
	}
	return dst.String(), true
}
