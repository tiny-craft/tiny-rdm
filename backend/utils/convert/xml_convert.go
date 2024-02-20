package convutil

import (
	"encoding/xml"
	"strings"
)

type XmlConvert struct{}

func (XmlConvert) Enable() bool {
	return true
}

func (XmlConvert) Encode(str string) (string, bool) {
	return str, true
}

func (XmlConvert) Decode(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if !strings.HasPrefix(trimedStr, "<") && !strings.HasSuffix(trimedStr, ">") {
		return str, false
	}
	var obj any
	err := xml.Unmarshal([]byte(trimedStr), &obj)
	return str, err == nil
}
