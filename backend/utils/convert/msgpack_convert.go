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

func (MsgpackConvert) Decode(str string) (string, bool) {
	var decodedStr string
	if err := msgpack.Unmarshal([]byte(str), &decodedStr); err == nil {
		return decodedStr, true
	}

	var obj map[string]any
	if err := msgpack.Unmarshal([]byte(str), &obj); err == nil {
		if b, err := json.Marshal(obj); err == nil {
			if len(b) >= 10 {
				return string(b), true
			}
		}
	}

	return str, false
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
