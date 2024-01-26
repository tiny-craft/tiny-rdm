package strutil

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"regexp"
	"strconv"
	"strings"
	"tinyrdm/backend/types"
)

// ConvertTo convert string to specified type
// @param decodeType empty string indicates automatic detection
// @param formatType empty string indicates automatic detection
func ConvertTo(str, decodeType, formatType string) (value, resultDecode, resultFormat string) {
	if len(str) <= 0 {
		// empty content
		if len(formatType) <= 0 {
			resultFormat = types.FORMAT_RAW
		} else {
			resultFormat = formatType
		}
		if len(decodeType) <= 0 {
			resultDecode = types.DECODE_NONE
		} else {
			resultDecode = decodeType
		}
		return
	}

	// decode first
	value, resultDecode = decodeWith(str, decodeType)
	// then format content
	value, resultFormat = viewAs(value, formatType)
	return
}

func decodeWith(str, decodeType string) (value, resultDecode string) {
	if len(decodeType) > 0 {
		switch decodeType {
		case types.DECODE_NONE:
			value = str
			resultDecode = decodeType
			return

		case types.DECODE_BASE64:
			if base64Str, ok := decodeBase64(str); ok {
				value = base64Str
			} else {
				value = str
			}
			resultDecode = decodeType
			return

		case types.DECODE_GZIP:
			if gzipStr, ok := decodeGZip(str); ok {
				value = gzipStr
			} else {
				value = str
			}
			resultDecode = decodeType
			return

		case types.DECODE_DEFLATE:
			if falteStr, ok := decodeDeflate(str); ok {
				value = falteStr
			} else {
				value = str
			}
			resultDecode = decodeType
			return

		case types.DECODE_ZSTD:
			if zstdStr, ok := decodeZStd(str); ok {
				value = zstdStr
			} else {
				value = str
			}
			resultDecode = decodeType
			return

		case types.DECODE_BROTLI:
			if brotliStr, ok := decodeBrotli(str); ok {
				value = brotliStr
			} else {
				value = str
			}
			resultDecode = decodeType
			return
		}
	}

	return autoDecode(str)
}

// attempt try possible decode method
// if no decode is possible, it will return the origin string value and "none" decode type
func autoDecode(str string) (value, resultDecode string) {
	if len(str) > 0 {
		// pure digit content may incorrect regard as some encoded type, skip decode
		if match, _ := regexp.MatchString(`^\d+$`, str); !match {
			var ok bool
			if len(str)%4 == 0 && len(str) >= 12 && !isSameChar(str) {
				if value, ok = decodeBase64(str); ok {
					resultDecode = types.DECODE_BASE64
					return
				}
			}

			if value, ok = decodeGZip(str); ok {
				resultDecode = types.DECODE_GZIP
				return
			}

			// FIXME: skip decompress with deflate due to incorrect format checking
			//if value, ok = decodeDeflate(str); ok {
			//	resultDecode = types.DECODE_DEFLATE
			//	return
			//}

			if value, ok = decodeZStd(str); ok {
				resultDecode = types.DECODE_ZSTD
				return
			}

			// FIXME: skip decompress with brotli due to incorrect format checking
			//if value, ok = decodeBrotli(str); ok {
			//	resultDecode = types.DECODE_BROTLI
			//	return
			//}

			if value, ok = decodeMsgpack(str); ok {
				resultDecode = types.DECODE_MSGPACK
				return
			}
		}
	}

	value = str
	resultDecode = types.DECODE_NONE
	return
}

func viewAs(str, formatType string) (value, resultFormat string) {
	if len(formatType) > 0 {
		switch formatType {
		case types.FORMAT_RAW, types.FORMAT_YAML, types.FORMAT_XML:
			value = str
			resultFormat = formatType
			return

		case types.FORMAT_JSON:
			if jsonStr, ok := decodeJson(str); ok {
				value = jsonStr
			} else {
				value = str
			}
			resultFormat = formatType
			return

		case types.FORMAT_HEX:
			if hexStr, ok := decodeToHex(str); ok {
				value = hexStr
			} else {
				value = str
			}
			resultFormat = formatType
			return

		case types.FORMAT_BINARY:
			if binStr, ok := decodeBinary(str); ok {
				value = binStr
			} else {
				value = str
			}
			resultFormat = formatType
			return
		}
	}

	return autoViewAs(str)
}

// attempt automatic convert to possible types
// if no conversion is possible, it will return the origin string value and "plain text" type
func autoViewAs(str string) (value, resultFormat string) {
	if len(str) > 0 {
		var ok bool
		if value, ok = decodeJson(str); ok {
			resultFormat = types.FORMAT_JSON
			return
		}

		if containsBinary(str) {
			if value, ok = decodeToHex(str); ok {
				resultFormat = types.FORMAT_HEX
				return
			}
		}
	}

	value = str
	resultFormat = types.FORMAT_RAW
	return
}

func decodeJson(str string) (string, bool) {
	trimedStr := strings.TrimSpace(str)
	if (strings.HasPrefix(trimedStr, "{") && strings.HasSuffix(trimedStr, "}")) ||
		(strings.HasPrefix(trimedStr, "[") && strings.HasSuffix(trimedStr, "]")) {
		var out bytes.Buffer
		if err := json.Indent(&out, []byte(trimedStr), "", "  "); err == nil {
			return out.String(), true
		}
	}
	return str, false
}

func decodeBase64(str string) (string, bool) {
	if decodedStr, err := base64.StdEncoding.DecodeString(str); err == nil {
		if s := string(decodedStr); !containsBinary(s) {
			return s, true
		}
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

func decodeToHex(str string) (string, bool) {
	decodeStr := hex.EncodeToString([]byte(str))
	var resultStr strings.Builder
	for i := 0; i < len(decodeStr); i += 2 {
		resultStr.WriteString("\\x")
		resultStr.WriteString(strings.ToUpper(decodeStr[i : i+2]))
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

func decodeZStd(str string) (string, bool) {
	if reader, err := zstd.NewReader(strings.NewReader(str)); err == nil {
		defer reader.Close()
		if decompressed, err := io.ReadAll(reader); err == nil {
			return string(decompressed), true
		}
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

func decodeMsgpack(str string) (string, bool) {
	var decodedStr string
	if err := msgpack.Unmarshal([]byte(str), &decodedStr); err == nil {
		return decodedStr, true
	}

	var decodedObj map[string]any
	if err := msgpack.Unmarshal([]byte(str), &decodedObj); err == nil {
		if b, err := json.Marshal(decodedObj); err == nil {
			return string(b), true
		}
	}

	return str, false
}

func SaveAs(str, format, decode string) (value string, err error) {
	value = str
	switch format {
	case types.FORMAT_JSON:
		if jsonStr, ok := encodeJson(str); ok {
			value = jsonStr
		} else {
			err = errors.New("invalid json data")
			return
		}

	case types.FORMAT_HEX:
		if hexStr, ok := encodeHex(str); ok {
			value = hexStr
		} else {
			err = errors.New("invalid hex data")
			return
		}

	case types.FORMAT_BINARY:
		if binStr, ok := encodeBinary(str); ok {
			value = binStr
		} else {
			err = errors.New("invalid binary data")
			return
		}
	}

	switch decode {
	case types.DECODE_NONE:
		return

	case types.DECODE_BASE64:
		value, _ = encodeBase64(value)
		return

	case types.DECODE_GZIP:
		if gzipStr, ok := encodeGZip(str); ok {
			value = gzipStr
		} else {
			err = errors.New("fail to build gzip")
		}
		return

	case types.DECODE_DEFLATE:
		if deflateStr, ok := encodeDeflate(str); ok {
			value = deflateStr
		} else {
			err = errors.New("fail to build deflate")
		}
		return

	case types.DECODE_ZSTD:
		if zstdStr, ok := encodeZStd(str); ok {
			value = zstdStr
		} else {
			err = errors.New("fail to build zstd")
		}
		return

	case types.DECODE_BROTLI:
		if brotliStr, ok := encodeBrotli(str); ok {
			value = brotliStr
		} else {
			err = errors.New("fail to build brotli")
		}
		return

	case types.DECODE_MSGPACK:
		if msgpackStr, ok := encodeMsgpack(str); ok {
			value = msgpackStr
		} else {
			err = errors.New("fail to build msgpack")
		}
		return
	}
	return str, nil
}

func encodeJson(str string) (string, bool) {
	var dst bytes.Buffer
	if err := json.Compact(&dst, []byte(str)); err != nil {
		return str, false
	}
	return dst.String(), true
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

func encodeZStd(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer, err := zstd.NewWriter(&buf)
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
	if zstdStr, err := compress([]byte(str)); err == nil {
		return zstdStr, true
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

func encodeMsgpack(str string) (string, bool) {
	var err error
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	if err = enc.EncodeString(str); err == nil {
		return buf.String(), true
	}

	return str, false
}
