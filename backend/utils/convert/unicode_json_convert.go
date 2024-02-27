package convutil

import (
	"bytes"
	"encoding/json"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type UnicodeJsonConvert struct{}

func (UnicodeJsonConvert) Enable() bool {
	return true
}

func (UnicodeJsonConvert) Decode(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if (strings.HasPrefix(trimedStr, "{") && strings.HasSuffix(trimedStr, "}")) ||
		(strings.HasPrefix(trimedStr, "[") && strings.HasSuffix(trimedStr, "]")) {
		var out bytes.Buffer
		if err := json.Indent(&out, []byte(trimedStr), "", "  "); err == nil {
			if quoteStr, ok := UnquoteUnicodeJson([]byte(trimedStr)); ok {
				return string(quoteStr), true
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

func UnquoteUnicodeJson(s []byte) ([]byte, bool) {
	var unquoted bytes.Buffer
	r := 0
	ls := len(s)
	for r < ls {
		c := s[r]
		offset := 1
		if c == '"' {
			// find next '"'
			for ; r+offset < ls; offset++ {
				if s[r+offset] == '"' && s[r+offset-1] != '\\' {
					offset += 1
					if ub, ok := unquoteBytes(s[r : r+offset]); ok {
						unquoted.WriteByte('"')
						unquoted.Write(ub)
						unquoted.WriteByte('"')
					} else {
						return nil, false
					}
					break
				}
			}
			// can not find close '"' until reach to the end of content
			if r+offset >= ls {
				return nil, false
			}
		} else {
			unquoted.WriteByte(c)
		}
		r += offset
	}
	return unquoted.Bytes(), true
}

func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}
