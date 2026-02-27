package convutil

import (
	"bytes"
	"io"
	"strings"

	"github.com/klauspost/compress/flate"
)

type DeflateConvert struct{}

func (d DeflateConvert) Enable() bool {
	return true
}

func (d DeflateConvert) Encode(str string) (string, bool) {
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

func (d DeflateConvert) Decode(str string) (string, bool) {
	reader := flate.NewReader(strings.NewReader(str))
	defer reader.Close()
	if decompressed, err := io.ReadAll(reader); err == nil {
		return string(decompressed), true
	}
	return str, false
}
