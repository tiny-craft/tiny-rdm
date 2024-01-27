package convutil

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackConvert struct{}

func (MsgpackConvert) Encode(str string) (string, bool) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(str), &obj); err == nil {
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
			return string(b), true
		}
	}

	return str, false
}
