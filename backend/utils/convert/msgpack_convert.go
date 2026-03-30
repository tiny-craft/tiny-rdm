package convutil

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackConvert struct{}

func (MsgpackConvert) Enable() bool {
	return true
}

func (c MsgpackConvert) Encode(str string) (string, bool) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(str), &obj); err == nil {
		for k, v := range obj {
			obj[k] = c.TryFloatToInt(v)
		}
		if b, err := msgpack.Marshal(obj); err == nil {
			return string(b), true
		}
	}

	if b, err := msgpack.Marshal(str); err != nil {
		return string(b), true
	}

	return str, false
}

func (c MsgpackConvert) Decode(str string) (string, bool) {
	if !c.MaybeMsgpack(str) {
		return str, false
	}

	var decodedStr string
	if err := msgpack.Unmarshal([]byte(str), &decodedStr); err == nil {
		return decodedStr, true
	}

	var obj map[string]any
	if err := msgpack.Unmarshal([]byte(str), &obj); err == nil {
		if b, err := json.Marshal(obj); err == nil {
			return string(b), true
		}
	}

	var arr []any
	if err := msgpack.Unmarshal([]byte(str), &arr); err == nil {
		if b, err := json.Marshal(arr); err == nil {
			return string(b), true
		}
	}

	return str, false
}

func (c MsgpackConvert) MaybeMsgpack(input string) bool {
	byt := []byte(input)
	if len(byt) < 2 {
		return false
	}

	b := byt[0]

	// fixmap 0x80-0x8f, fixarray 0x90-0x9f
	// array16 0xdc, array32 0xdd, map16 0xde, map32 0xdf
	return (b >= 0x80 && b <= 0x9f) || (b >= 0xdc && b <= 0xdf)
}

func (c MsgpackConvert) TryFloatToInt(input any) any {
	switch val := input.(type) {
	case map[string]any:
		for k, v := range val {
			val[k] = c.TryFloatToInt(v)
		}
		return val
	case []any:
		for i, v := range val {
			val[i] = c.TryFloatToInt(v)
		}
		return val
	case float64:
		if val == float64(int(val)) {
			return int(val)
		}
		return val
	default:
		return val
	}
}
