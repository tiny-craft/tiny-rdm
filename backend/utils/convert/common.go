package convutil

import (
	"os"
	"path"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/utils/confdir"
)

func writeExecuteFile(content []byte, filename string) (string, error) {
	filepath := path.Join(confdir.GetConfigDir(), consts.APP_DATA_FOLDER, "decoder", filename)
	_ = os.Mkdir(path.Dir(filepath), 0777)
	err := os.WriteFile(filepath, content, 0777)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
