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
	var result strings.Builder
	total := len(str)
	for i := 0; i < total; i += 8 {
		b, _ := strconv.ParseInt(str[i:i+8], 2, 8)
		result.WriteByte(byte(b))
	}
	return result.String(), true
}

func (BinaryConvert) Decode(str string) (string, bool) {
	var binary strings.Builder
	for _, char := range str {
		binary.WriteString(fmt.Sprintf("%08b", int(char)))
	}
	return binary.String(), true
}
