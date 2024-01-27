package convutil

import (
	"encoding/hex"
	"strings"
)

type HexConvert struct{}

func (HexConvert) Encode(str string) (string, bool) {
	hexStrArr := strings.Split(str, "\\x")
	hexStr := strings.Join(hexStrArr, "")
	if decodeStr, err := hex.DecodeString(hexStr); err == nil {
		return string(decodeStr), true
	}

	return str, false
}

func (HexConvert) Decode(str string) (string, bool) {
	decodeStr := hex.EncodeToString([]byte(str))
	decodeStr = strings.ToUpper(decodeStr)
	var resultStr strings.Builder
	for i := 0; i < len(decodeStr); i += 2 {
		resultStr.WriteString("\\x")
		resultStr.WriteString(decodeStr[i : i+2])
	}
	return resultStr.String(), true
}
