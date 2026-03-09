package convutil

import (
	"fmt"
	"strconv"
	"strings"
)

type BinaryConvert struct{}

func (BinaryConvert) Enable() bool {
	return true
}

func (BinaryConvert) Encode(str string) (string, bool) {
	total := len(str)
	if total%8 != 0 {
		return str, false
	}
	var result strings.Builder
	for i := 0; i < total; i += 8 {
		b, err := strconv.ParseUint(str[i:i+8], 2, 8)
		if err != nil {
			return str, false
		}
		result.WriteByte(byte(b))
	}
	return result.String(), true
}

func (BinaryConvert) Decode(str string) (string, bool) {
	var binary strings.Builder
	for _, char := range []byte(str) {
		binary.WriteString(fmt.Sprintf("%08b", int(char)))
	}
	return binary.String(), true
}
