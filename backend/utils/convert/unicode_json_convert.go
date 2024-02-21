package convutil

import (
	"bytes"
	"encoding/json"
	"strings"
)

type UnicodeJsonConvert struct{}

func (UnicodeJsonConvert) Enable() bool {
	return true
}

func (UnicodeJsonConvert) Decode(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if strings.HasPrefix(trimedStr, "{") && strings.HasSuffix(trimedStr, "}") {
		var obj map[string]any
		if err := json.Unmarshal([]byte(trimedStr), &obj); err == nil {
			var out []byte
			if out, err = json.MarshalIndent(obj, "", " "); err == nil {
				return string(out), true
			}
		}
	} else if strings.HasPrefix(trimedStr, "[") && strings.HasSuffix(trimedStr, "]") {
		var arr []any
		if err := json.Unmarshal([]byte(trimedStr), &arr); err == nil {
			var out []byte
			if out, err = json.MarshalIndent(arr, "", " "); err == nil {
				return string(out), true
			}
		}
	}
	return str, false
}

func (UnicodeJsonConvert) Encode(str string) (string, bool) {
	var dst bytes.Buffer
	if err := json.Compact(&dst, []byte(str)); err != nil {
		return str, false
	}
	return dst.String(), true
}
