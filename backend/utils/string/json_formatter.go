package strutil

import (
	"strings"
	"unicode"
)

// Convert from https://github.com/ObuchiYuki/SwiftJSONFormatter

// ArrayIterator defines the iterator for an array
type ArrayIterator[T any] struct {
	array []T
	head  int
}

// NewArrayIterator initializes a new ArrayIterator with the given array
func NewArrayIterator[T any](array []T) *ArrayIterator[T] {
	return &ArrayIterator[T]{
		array: array,
		head:  -1,
	}
}

// HasNext returns true if there are more elements to iterate over
func (it *ArrayIterator[T]) HasNext() bool {
	return it.head+1 < len(it.array)
}

// PeekNext returns the next element without advancing the iterator
func (it *ArrayIterator[T]) PeekNext() *T {
	if it.head+1 < len(it.array) {
		return &it.array[it.head+1]
	}
	return nil
}

// Next returns the next element and advances the iterator
func (it *ArrayIterator[T]) Next() *T {
	defer func() {
		it.head++
	}()
	return it.PeekNext()
}

// JSONBeautify formats a JSON string with indentation
func JSONBeautify(value string, indent string) string {
	if len(indent) <= 0 {
		indent = "    "
	}
	return format(value, indent, "\n", " ")
}

// JSONMinify formats a JSON string by removing all unnecessary whitespace
func JSONMinify(value string) string {
	return format(value, "", "", "")
}

// format applies the specified formatting to a JSON string
func format(value string, indent string, newLine string, separator string) string {
	var formatted strings.Builder
	chars := NewArrayIterator([]rune(value))
	indentLevel := 0

	for chars.HasNext() {
		if char := chars.Next(); char != nil {
			switch *char {
			case '{', '[':
				formatted.WriteRune(*char)
				consumeWhitespaces(chars)
				peeked := chars.PeekNext()
				if peeked != nil && (*peeked == '}' || *peeked == ']') {
					chars.Next()
					formatted.WriteRune(*peeked)
				} else {
					indentLevel++
					formatted.WriteString(newLine)
					formatted.WriteString(strings.Repeat(indent, indentLevel))
				}
			case '}', ']':
				indentLevel--
				formatted.WriteString(newLine)
				formatted.WriteString(strings.Repeat(indent, max(0, indentLevel)))
				formatted.WriteRune(*char)
			case '"':
				str := consumeString(chars)
				//str = convertUnicodeString(str)
				formatted.WriteString(str)
			case ',':
				consumeWhitespaces(chars)
				formatted.WriteRune(',')
				peeked := chars.PeekNext()
				if peeked != nil && *peeked != '}' && *peeked != ']' {
					formatted.WriteString(newLine)
					formatted.WriteString(strings.Repeat(indent, max(0, indentLevel)))
				}
			case ':':
				formatted.WriteString(":" + separator)
			default:
				if !unicode.IsSpace(*char) {
					formatted.WriteRune(*char)
				}
			}
		}
	}

	return formatted.String()
}

// consumeWhitespaces advances the iterator past any whitespace characters
func consumeWhitespaces(iter *ArrayIterator[rune]) {
	for iter.HasNext() {
		if peeked := iter.PeekNext(); peeked != nil && unicode.IsSpace(*peeked) {
			iter.Next()
		} else {
			break
		}
	}
}

// consumeString consumes a JSON string value from the iterator
func consumeString(iter *ArrayIterator[rune]) string {
	var sb strings.Builder
	sb.WriteRune('"')
	escaping := false

	for iter.HasNext() {
		if char := iter.Next(); char != nil {
			if *char == '\n' {
				return sb.String() // Unterminated string
			}

			sb.WriteRune(*char)

			if escaping {
				escaping = false
			} else {
				if *char == '\\' {
					escaping = true
				}
				if *char == '"' {
					break
				}
			}
		}
	}

	return sb.String()
}

func convertUnicodeString(str string) string {
	// TODO: quote UTF-16 characters
	//if len(str) > 2 {
	//	if unqStr, err := strconv.Unquote(str); err == nil {
	//		return strconv.Quote(unqStr)
	//	}
	//}
	return str
}
