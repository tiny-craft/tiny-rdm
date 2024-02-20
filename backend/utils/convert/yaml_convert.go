package convutil

import (
	"gopkg.in/yaml.v3"
)

type YamlConvert struct{}

func (YamlConvert) Enable() bool {
	return true
}

func (YamlConvert) Encode(str string) (string, bool) {
	return str, true
}

func (YamlConvert) Decode(str string) (string, bool) {
	var obj map[string]any
	err := yaml.Unmarshal([]byte(str), &obj)
	return str, err == nil
}
