package convutil

import (
	"errors"
	"regexp"
	"tinyrdm/backend/types"
	strutil "tinyrdm/backend/utils/string"
)

type DataConvert interface {
	Enable() bool
	Encode(string) (string, bool)
	Decode(string) (string, bool)
}

var (
	jsonConv    JsonConvert
	uniJsonConv UnicodeJsonConvert
	yamlConv    YamlConvert
	xmlConv     XmlConvert
	base64Conv  Base64Convert
	binaryConv  BinaryConvert
	hexConv     HexConvert
	gzipConv    GZipConvert
	deflateConv DeflateConvert
	zstdConv    ZStdConvert
	lz4Conv     LZ4Convert
	brotliConv  BrotliConvert
	msgpackConv MsgpackConvert
	phpConv     = NewPhpConvert()
	pickleConv  = NewPickleConvert()
)

var BuildInFormatters = map[string]DataConvert{
	types.FORMAT_JSON:         jsonConv,
	types.FORMAT_UNICODE_JSON: uniJsonConv,
	types.FORMAT_YAML:         yamlConv,
	types.FORMAT_XML:          xmlConv,
	types.FORMAT_HEX:          hexConv,
	types.FORMAT_BINARY:       binaryConv,
}

var BuildInDecoders = map[string]DataConvert{
	types.DECODE_BASE64:  base64Conv,
	types.DECODE_GZIP:    gzipConv,
	types.DECODE_DEFLATE: deflateConv,
	types.DECODE_ZSTD:    zstdConv,
	types.DECODE_LZ4:     lz4Conv,
	types.DECODE_BROTLI:  brotliConv,
	types.DECODE_MSGPACK: msgpackConv,
	types.DECODE_PHP:     phpConv,
	types.DECODE_PICKLE:  pickleConv,
}

// ConvertTo convert string to specified type
// @param decodeType empty string indicates automatic detection
// @param formatType empty string indicates automatic detection
// @param custom decoder if any
func ConvertTo(str, decodeType, formatType string, customDecoder []CmdConvert) (value, resultDecode, resultFormat string) {
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
	value, resultDecode = decodeWith(str, decodeType, customDecoder)
	// then format content
	if len(formatType) <= 0 {
		value, resultFormat = autoViewAs(value)
	} else {
		value, resultFormat = viewAs(value, formatType)
	}
	return
}

func decodeWith(str, decodeType string, customDecoder []CmdConvert) (value, resultDecode string) {
	if len(decodeType) > 0 {
		value = str

		if buildinDecoder, ok := BuildInDecoders[decodeType]; ok {
			if decodedStr, ok := buildinDecoder.Decode(str); ok {
				value = decodedStr
			}
		} else if decodeType != types.DECODE_NONE {
			for _, decoder := range customDecoder {
				if decoder.Name == decodeType {
					if decodedStr, ok := decoder.Decode(str); ok {
						value = decodedStr
					}
					break
				}
			}
		}

		resultDecode = decodeType
		return
	}

	value, resultDecode = autoDecode(str, customDecoder)
	return
}

// attempt try possible decode method
// if no decode is possible, it will return the origin string value and "none" decode type
func autoDecode(str string, customDecoder []CmdConvert) (value, resultDecode string) {
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

			if value, ok = lz4Conv.Decode(str); ok {
				resultDecode = types.DECODE_LZ4
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

			if value, ok = phpConv.Decode(str); ok {
				resultDecode = types.DECODE_PHP
				return
			}

			if value, ok = pickleConv.Decode(str); ok {
				resultDecode = types.DECODE_PICKLE
				return
			}

			// try decode with custom decoder
			for _, decoder := range customDecoder {
				if decoder.Auto {
					if value, ok = decoder.Decode(str); ok {
						resultDecode = decoder.Name
						return
					}
				}
			}
		}
	}

	value = str
	resultDecode = types.DECODE_NONE
	return
}

func viewAs(str, formatType string) (value, resultFormat string) {
	if len(formatType) > 0 {
		value = str
		if buildinFormatter, ok := BuildInFormatters[formatType]; ok {
			if formattedStr, ok := buildinFormatter.Decode(str); ok {
				value = formattedStr
			}
		}
		resultFormat = formatType
		return
	}
	return
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

		if value, ok = yamlConv.Decode(str); ok {
			resultFormat = types.FORMAT_YAML
			return
		}

		if value, ok = xmlConv.Decode(str); ok {
			resultFormat = types.FORMAT_XML
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

func SaveAs(str, format, decode string, customDecoder []CmdConvert) (value string, err error) {
	value = str
	if buildingFormatter, ok := BuildInFormatters[format]; ok {
		if formattedStr, ok := buildingFormatter.Encode(str); ok {
			value = formattedStr
		} else {
			err = errors.New("invalid " + format + " data")
			return
		}
	}

	if buildinDecoder, ok := BuildInDecoders[decode]; ok {
		if encodedValue, ok := buildinDecoder.Encode(str); ok {
			value = encodedValue
		} else {
			err = errors.New("fail to build " + decode)
		}
		return
	} else if decode != types.DECODE_NONE {
		for _, decoder := range customDecoder {
			if decoder.Name == decode {
				if encodedStr, ok := decoder.Encode(str); ok {
					value = encodedStr
				} else {
					err = errors.New("fail to build " + decode)
				}
				return
			}
		}
	}
	return
}
