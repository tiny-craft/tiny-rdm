package convutil

import (
	"encoding/base64"
	strutil "tinyrdm/backend/utils/string"
)

type Base64Convert struct{}

func (Base64Convert) Encode(str string) (string, bool) {
	return base64.StdEncoding.EncodeToString([]byte(str)), true
}

func (Base64Convert) Decode(str string) (string, bool) {
	if decodedStr, err := base64.StdEncoding.DecodeString(str); err == nil {
		if s := string(decodedStr); !strutil.ContainsBinary(s) {
			return s, true
		}
	}
	return str, false
}
