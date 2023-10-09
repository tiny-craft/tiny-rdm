package strutil

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"strconv"
	"strings"
	"tinyrdm/backend/types"
	"unicode/utf8"
)

// ConvertTo convert string to specified type
// @param targetType  empty string indicates automatic detection of the string type
func ConvertTo(str, targetType string) (value, resultType string) {
	if len(str) <= 0 {
		// empty content
		if len(targetType) <= 0 {
			resultType = types.PLAIN_TEXT
		} else {
			resultType = targetType
		}
		return
	}

	switch targetType {
	case types.PLAIN_TEXT:
		value = str
		resultType = targetType
		return

	case types.JSON:
		value, _ = decodeJson(str)
		resultType = targetType
		return

	case types.BASE64_TEXT, types.BASE64_JSON:
		if base64Str, ok := decodeBase64(str); ok {
			if targetType == types.BASE64_JSON {
				value, _ = decodeJson(base64Str)
			} else {
				value = base64Str
			}
		} else {
			value = str
		}
		resultType = targetType
		return

	case types.HEX:
		if hexStr, ok := decodeHex(str); ok {
			value = hexStr
		} else {
			value = str
		}
		resultType = targetType
		return

	case types.BINARY:
		if binStr, ok := decodeBinary(str); ok {
			value = binStr
		} else {
			value = str
		}
		resultType = targetType
		return

	case types.GZIP, types.GZIP_JSON:
		if gzipStr, ok := decodeGZip(str); ok {
			if targetType == types.GZIP_JSON {
				value, _ = decodeJson(gzipStr)
			} else {
				value = gzipStr
			}
		} else {
			value = str
		}
		resultType = targetType
		return

	case types.DEFLATE, types.DEFLATE_JSON:
		if deflateStr, ok := decodeDeflate(str); ok {
			if targetType == types.DEFLATE_JSON {
				value, _ = decodeJson(deflateStr)
			} else {
				value = deflateStr
			}
		} else {
			value = str
		}
		resultType = targetType
		return

	case types.BROTLI, types.BROTLI_JSON:
		if brotliStr, ok := decodeBrotli(str); ok {
			if targetType == types.BROTLI_JSON {
				value, _ = decodeJson(brotliStr)
			} else {
				value = brotliStr
			}
		} else {
			value = str
		}
		resultType = targetType
		return
	}

	// type isn't specified or unknown, try to automatically detect and return converted value
	return autoToType(str)
}

// attempt automatic convert to possible types
// if no conversion is possible, it will return the origin string value and "plain text" type
func autoToType(str string) (value, resultType string) {
	if len(str) > 0 {
		var ok bool
		if value, ok = decodeJson(str); ok {
			resultType = types.JSON
			return
		}

		if value, ok = decodeBase64(str); ok {
			if value, ok = decodeJson(value); ok {
				resultType = types.BASE64_JSON
				return
			}
			resultType = types.BASE64_TEXT
			return
		}

		if value, ok = decodeGZip(str); ok {
			if value, ok = decodeJson(value); ok {
				resultType = types.GZIP_JSON
				return
			}
			resultType = types.GZIP
			return
		}

		if value, ok = decodeDeflate(str); ok {
			if value, ok = decodeJson(value); ok {
				resultType = types.DEFLATE_JSON
				return
			}
			resultType = types.DEFLATE
			return
		}

		if value, ok = decodeBrotli(str); ok {
			if value, ok = decodeJson(value); ok {
				resultType = types.BROTLI_JSON
				return
			}
			resultType = types.BROTLI
			return
		}

		if isBinary(str) {
			if value, ok = decodeHex(str); ok {
				resultType = types.HEX
				return
			}
		}
	}

	value = str
	resultType = types.PLAIN_TEXT
	return
}

func isBinary(str string) bool {
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

func decodeJson(str string) (string, bool) {
	var data any
	if (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
		(strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) {
		if err := json.Unmarshal([]byte(str), &data); err == nil {
			var jsonByte []byte
			if jsonByte, err = json.MarshalIndent(data, "", "  "); err == nil {
				return string(jsonByte), true
			}
		}
	}
	return str, false
}

func decodeBase64(str string) (string, bool) {
	if decodedStr, err := base64.StdEncoding.DecodeString(str); err == nil {
		return string(decodedStr), true
	}
	return str, false
}

func decodeBinary(str string) (string, bool) {
	var binary strings.Builder
	for _, char := range str {
		binary.WriteString(fmt.Sprintf("%08b", int(char)))
	}
	return binary.String(), true
}

func decodeHex(str string) (string, bool) {
	decodeStr := hex.EncodeToString([]byte(str))
	var resultStr strings.Builder
	for i := 0; i < len(decodeStr); i += 2 {
		resultStr.WriteString("\\x")
		resultStr.WriteString(decodeStr[i : i+2])
	}
	return resultStr.String(), true
}

func decodeGZip(str string) (string, bool) {
	if reader, err := gzip.NewReader(strings.NewReader(str)); err == nil {
		defer reader.Close()
		var decompressed []byte
		if decompressed, err = io.ReadAll(reader); err == nil {
			return string(decompressed), true
		}
	}
	return str, false
}

func decodeDeflate(str string) (string, bool) {
	reader := flate.NewReader(strings.NewReader(str))
	defer reader.Close()
	if decompressed, err := io.ReadAll(reader); err == nil {
		return string(decompressed), true
	}
	return str, false
}

func decodeBrotli(str string) (string, bool) {
	reader := brotli.NewReader(strings.NewReader(str))
	if decompressed, err := io.ReadAll(reader); err == nil {
		return string(decompressed), true
	}
	return str, false
}

func SaveAs(str, targetType string) (value string, err error) {
	switch targetType {
	case types.PLAIN_TEXT:
		return str, nil

	case types.BASE64_TEXT:
		base64Str, _ := encodeBase64(str)
		return base64Str, nil

	case types.HEX:
		hexStr, _ := encodeHex(str)
		return hexStr, nil

	case types.BINARY:
		binStr, _ := encodeBinary(str)
		return binStr, nil

	case types.JSON, types.BASE64_JSON, types.GZIP_JSON, types.DEFLATE_JSON, types.BROTLI_JSON:
		if jsonStr, ok := encodeJson(str); ok {
			switch targetType {
			case types.BASE64_JSON:
				base64Str, _ := encodeBase64(jsonStr)
				return base64Str, nil
			case types.GZIP_JSON:
				gzipStr, _ := encodeGZip(jsonStr)
				return gzipStr, nil
			case types.DEFLATE_JSON:
				deflateStr, _ := encodeDeflate(jsonStr)
				return deflateStr, nil
			case types.BROTLI_JSON:
				brotliStr, _ := encodeBrotli(jsonStr)
				return brotliStr, nil
			default:
				return jsonStr, nil
			}
		} else {
			return str, errors.New("invalid json")
		}

	case types.GZIP:
		if gzipStr, ok := encodeGZip(str); ok {
			return gzipStr, nil
		} else {
			return str, errors.New("fail to build gzip data")
		}

	case types.DEFLATE:
		if deflateStr, ok := encodeDeflate(str); ok {
			return deflateStr, nil
		} else {
			return str, errors.New("fail to build deflate data")
		}

	case types.BROTLI:
		if brotliStr, ok := encodeBrotli(str); ok {
			return brotliStr, nil
		} else {
			return str, errors.New("fail to build brotli data")
		}
	}
	return str, errors.New("fail to save with unknown error")
}

func encodeJson(str string) (string, bool) {
	var data any
	if (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
		(strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) {
		if err := json.Unmarshal([]byte(str), &data); err == nil {
			var jsonByte []byte
			if jsonByte, err = json.Marshal(data); err == nil {
				return string(jsonByte), true
			}
		}
	}
	return str, false
}

func encodeBase64(str string) (string, bool) {
	return base64.StdEncoding.EncodeToString([]byte(str)), true
}

func encodeBinary(str string) (string, bool) {
	var result strings.Builder
	total := len(str)
	for i := 0; i < total; i += 8 {
		b, _ := strconv.ParseInt(str[i:i+8], 2, 8)
		result.WriteByte(byte(b))
	}
	return result.String(), true
}

func encodeHex(str string) (string, bool) {
	hexStrArr := strings.Split(str, "\\x")
	hexStr := strings.Join(hexStrArr, "")
	if decodeStr, err := hex.DecodeString(hexStr); err == nil {
		return string(decodeStr), true
	}

	return str, false
}

func encodeGZip(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write([]byte(str)); err != nil {
			writer.Close()
			return "", err
		}
		writer.Close()
		return string(buf.Bytes()), nil
	}

	if gzipStr, err := compress([]byte(str)); err == nil {
		return gzipStr, true
	}
	return str, false
}

func encodeDeflate(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
		if err != nil {
			return "", err
		}
		if _, err = writer.Write([]byte(str)); err != nil {
			writer.Close()
			return "", err
		}
		writer.Close()
		return string(buf.Bytes()), nil
	}
	if deflateStr, err := compress([]byte(str)); err == nil {
		return deflateStr, true
	}
	return str, false
}

func encodeBrotli(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer := brotli.NewWriter(&buf)
		if _, err := writer.Write([]byte(str)); err != nil {
			writer.Close()
			return "", err
		}
		writer.Close()
		return string(buf.Bytes()), nil
	}
	if brotliStr, err := compress([]byte(str)); err == nil {
		return brotliStr, true
	}
	return str, false
}
