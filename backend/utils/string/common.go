package strutil

import "unicode/utf8"

func containsBinary(str string) bool {
	//buf := []byte(str)
	//size := 0
	//for start := 0; start < len(buf); start += size {
	//	var r rune
	//	if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
	//		return true
	//	}
	//}

	if !utf8.ValidString(str) {
		return true
	}
	return false
}
