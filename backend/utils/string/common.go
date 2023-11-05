package strutil

import (
	"unicode"
)

func containsBinary(str string) bool {
	//size := 0
	//for start := 0; start < len(buf); start += size {
	//	var r rune
	//	if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
	//		return true
	//	}
	//}
	rs := []rune(str)
	for _, r := range rs {
		if !unicode.IsPrint(r) && r != '\n' {
			return true
		}
	}
	return false
}
