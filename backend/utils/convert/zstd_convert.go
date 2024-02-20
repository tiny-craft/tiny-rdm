package convutil

import (
	"bytes"
	"github.com/klauspost/compress/zstd"
	"io"
	"strings"
)

type ZStdConvert struct{}

func (ZStdConvert) Enable() bool {
	return true
}

func (ZStdConvert) Encode(str string) (string, bool) {
	var compress = func(b []byte) (string, error) {
		var buf bytes.Buffer
		writer, err := zstd.NewWriter(&buf)
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
	if zstdStr, err := compress([]byte(str)); err == nil {
		return zstdStr, true
	}
	return str, false
}

func (ZStdConvert) Decode(str string) (string, bool) {
	if reader, err := zstd.NewReader(strings.NewReader(str)); err == nil {
		defer reader.Close()
		if decompressed, err := io.ReadAll(reader); err == nil {
			return string(decompressed), true
		}
	}
	return str, false
}
