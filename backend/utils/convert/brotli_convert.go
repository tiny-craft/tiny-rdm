package convutil

import (
	"bytes"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
)

type BrotliConvert struct{}

func (BrotliConvert) Enable() bool {
	return true
}

func (BrotliConvert) Encode(str string) (string, bool) {
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

func (BrotliConvert) Decode(str string) (string, bool) {
	reader := brotli.NewReader(strings.NewReader(str))
	if decompressed, err := io.ReadAll(reader); err == nil {
		return string(decompressed), true
	}
	return str, false
}
