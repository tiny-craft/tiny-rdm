package convutil

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"
)

type LZ4Convert struct{}

func (LZ4Convert) Enable() bool {
	return true
}

func (LZ4Convert) Encode(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer := lz4.NewWriter(&buf)
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

func (LZ4Convert) Decode(str string) (string, bool) {
	reader := lz4.NewReader(bytes.NewReader([]byte(str)))
	if decompressed, err := io.ReadAll(reader); err == nil {
		return string(decompressed), true
	}
	return str, false
}
