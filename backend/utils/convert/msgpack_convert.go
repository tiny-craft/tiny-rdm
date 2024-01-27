package convutil

import (
	"bytes"
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackConvert struct{}

func (MsgpackConvert) Encode(str string) (string, bool) {
	var err error
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	if err = enc.EncodeString(str); err == nil {
		return buf.String(), true
	}

	return str, false
}

func (MsgpackConvert) Decode(str string) (string, bool) {
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
