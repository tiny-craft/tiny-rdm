package convutil

import (
	"encoding/base64"
	"strings"
	sliceutil "tinyrdm/backend/utils/slice"
)

type CmdConvert struct {
	Name       string
	Auto       bool
	DecodePath string
	DecodeArgs []string
	EncodePath string
	EncodeArgs []string
}

const replaceholder = "{VALUE}"

func (c CmdConvert) Enable() bool {
	return true
}

func (c CmdConvert) Encode(str string) (string, bool) {
	base64Content := base64.StdEncoding.EncodeToString([]byte(str))
	var containHolder bool
	args := sliceutil.Map(c.EncodeArgs, func(i int) string {
		arg := strings.TrimSpace(c.EncodeArgs[i])
		if strings.Contains(arg, replaceholder) {
			arg = strings.ReplaceAll(arg, replaceholder, base64Content)
			containHolder = true
		}
		return arg
	})
	if len(args) <= 0 || !containHolder {
		args = append(args, base64Content)
	}
	output, err := runCommand(c.EncodePath, args...)
	if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
		return str, false
	}

	outputContent := make([]byte, base64.StdEncoding.DecodedLen(len(output)))
	n, err := base64.StdEncoding.Decode(outputContent, output)
	if err != nil {
		return str, false
	}
	return string(outputContent[:n]), true
}

func (c CmdConvert) Decode(str string) (string, bool) {
	base64Content := base64.StdEncoding.EncodeToString([]byte(str))
	var containHolder bool
	args := sliceutil.Map(c.DecodeArgs, func(i int) string {
		arg := strings.TrimSpace(c.DecodeArgs[i])
		if strings.Contains(arg, replaceholder) {
			arg = strings.ReplaceAll(arg, replaceholder, base64Content)
			containHolder = true
		}
		return arg
	})
	if len(args) <= 0 || !containHolder {
		args = append(args, base64Content)
	}
	output, err := runCommand(c.DecodePath, args...)
	if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
		return str, false
	}

	outputContent := make([]byte, base64.StdEncoding.DecodedLen(len(output)))
	n, err := base64.StdEncoding.Decode(outputContent, output)
	if err != nil {
		return str, false
	}
	return string(outputContent[:n]), true
}
