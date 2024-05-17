package convutil

import (
	"strings"
	strutil "tinyrdm/backend/utils/string"
)

type JsonConvert struct{}

func (JsonConvert) Enable() bool {
	return true
}

func (JsonConvert) Decode(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if (strings.HasPrefix(trimedStr, "{") && strings.HasSuffix(trimedStr, "}")) ||
		(strings.HasPrefix(trimedStr, "[") && strings.HasSuffix(trimedStr, "]")) {
		return strutil.JSONBeautify(trimedStr, "  "), true
	}
	return str, false
}

func (JsonConvert) Encode(str string) (string, bool) {
	return strutil.JSONMinify(str), true
}
