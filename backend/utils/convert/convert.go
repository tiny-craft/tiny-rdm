package convutil

import (
	"errors"
	"regexp"
	"tinyrdm/backend/types"
	strutil "tinyrdm/backend/utils/string"
)

type DataConvert interface {
	Encode(string) (string, bool)
	Decode(string) (string, bool)
}

var (
	jsonConv    JsonConvert
	base64Conv  Base64Convert
	binaryConv  BinaryConvert
	hexConv     HexConvert
	gzipConv    GZipConvert
	deflateConv DeflateConvert
	zstdConv    ZStdConvert
	brotliConv  BrotliConvert
	msgpackConv MsgpackConvert
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

		case types.DECODE_BASE64:
			if base64Str, ok := base64Conv.Decode(str); ok {
				value = base64Str
			} else {
				value = str
			}

		case types.DECODE_GZIP:
			if gzipStr, ok := gzipConv.Decode(str); ok {
				value = gzipStr
			} else {
				value = str
			}

		case types.DECODE_DEFLATE:
			if falteStr, ok := deflateConv.Decode(str); ok {
				value = falteStr
			} else {
				value = str
			}

		case types.DECODE_ZSTD:
			if zstdStr, ok := zstdConv.Decode(str); ok {
				value = zstdStr
			} else {
				value = str
			}

		case types.DECODE_BROTLI:
			if brotliStr, ok := brotliConv.Decode(str); ok {
				value = brotliStr
			} else {
				value = str
			}

		case types.DECODE_MSGPACK:
			if msgpackStr, ok := msgpackConv.Decode(str); ok {
				value = msgpackStr
			} else {
				value = str
			}
		}
		resultDecode = decodeType
		return
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
			if len(str)%4 == 0 && len(str) >= 12 && !strutil.IsSameChar(str) {
				if value, ok = base64Conv.Decode(str); ok {
					resultDecode = types.DECODE_BASE64
					return
				}
			}

			if value, ok = gzipConv.Decode(str); ok {
				resultDecode = types.DECODE_GZIP
				return
			}

			// FIXME: skip decompress with deflate due to incorrect format checking
			//if value, ok = decodeDeflate(str); ok {
			//	resultDecode = types.DECODE_DEFLATE
			//	return
			//}

			if value, ok = zstdConv.Decode(str); ok {
				resultDecode = types.DECODE_ZSTD
				return
			}

			// FIXME: skip decompress with brotli due to incorrect format checking
			//if value, ok = decodeBrotli(str); ok {
			//	resultDecode = types.DECODE_BROTLI
			//	return
			//}

			if value, ok = msgpackConv.Decode(str); ok {
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

		case types.FORMAT_JSON:
			if jsonStr, ok := jsonConv.Decode(str); ok {
				value = jsonStr
			} else {
				value = str
			}

		case types.FORMAT_HEX:
			if hexStr, ok := hexConv.Decode(str); ok {
				value = hexStr
			} else {
				value = str
			}

		case types.FORMAT_BINARY:
			if binStr, ok := binaryConv.Decode(str); ok {
				value = binStr
			} else {
				value = str
			}
		}
		resultFormat = formatType
		return
	}

	return autoViewAs(str)
}

// attempt automatic convert to possible types
// if no conversion is possible, it will return the origin string value and "plain text" type
func autoViewAs(str string) (value, resultFormat string) {
	if len(str) > 0 {
		var ok bool
		if value, ok = jsonConv.Decode(str); ok {
			resultFormat = types.FORMAT_JSON
			return
		}

		if strutil.ContainsBinary(str) {
			if value, ok = hexConv.Decode(str); ok {
				resultFormat = types.FORMAT_HEX
				return
			}
		}
	}

	value = str
	resultFormat = types.FORMAT_RAW
	return
}

func SaveAs(str, format, decode string) (value string, err error) {
	value = str
	switch format {
	case types.FORMAT_JSON:
		if jsonStr, ok := jsonConv.Encode(str); ok {
			value = jsonStr
		} else {
			err = errors.New("invalid json data")
			return
		}

	case types.FORMAT_HEX:
		if hexStr, ok := hexConv.Encode(str); ok {
			value = hexStr
		} else {
			err = errors.New("invalid hex data")
			return
		}

	case types.FORMAT_BINARY:
		if binStr, ok := binaryConv.Encode(str); ok {
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
		value, _ = base64Conv.Encode(value)
		return

	case types.DECODE_GZIP:
		if gzipStr, ok := gzipConv.Encode(str); ok {
			value = gzipStr
		} else {
			err = errors.New("fail to build gzip")
		}
		return

	case types.DECODE_DEFLATE:
		if deflateStr, ok := deflateConv.Encode(str); ok {
			value = deflateStr
		} else {
			err = errors.New("fail to build deflate")
		}
		return

	case types.DECODE_ZSTD:
		if zstdStr, ok := zstdConv.Encode(str); ok {
			value = zstdStr
		} else {
			err = errors.New("fail to build zstd")
		}
		return

	case types.DECODE_BROTLI:
		if brotliStr, ok := brotliConv.Encode(str); ok {
			value = brotliStr
		} else {
			err = errors.New("fail to build brotli")
		}
		return

	case types.DECODE_MSGPACK:
		if msgpackStr, ok := msgpackConv.Encode(str); ok {
			value = msgpackStr
		} else {
			err = errors.New("fail to build msgpack")
		}
		return
	}
	return str, nil
}
