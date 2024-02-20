package convutil

import (
	"bytes"
	"github.com/klauspost/compress/gzip"
	"io"
	"strings"
)

type GZipConvert struct{}

func (GZipConvert) Enable() bool {
	return true
}

func (GZipConvert) Encode(str string) (string, bool) {
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

func (GZipConvert) Decode(str string) (string, bool) {
	if reader, err := gzip.NewReader(strings.NewReader(str)); err == nil {
		defer reader.Close()
		var decompressed []byte
		if decompressed, err = io.ReadAll(reader); err == nil {
			return string(decompressed), true
		}
	}
	return str, false
}
