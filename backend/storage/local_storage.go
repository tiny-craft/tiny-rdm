package storage

import (
	"os"
	"path"

	"github.com/vrischmann/userdir"
)

// localStorage provides reading and writing application data to the user's
// configuration directory.
type localStorage struct {
	ConfPath string
}

// NewLocalStore returns a localStore instance.
func NewLocalStore(filename string) *localStorage {
	return &localStorage{
		ConfPath: path.Join(userdir.GetConfigHome(), "TinyRDM", filename),
	}
}

// Load reads the given file in the user's configuration directory and returns
// its contents.
func (l *localStorage) Load() ([]byte, error) {
	d, err := os.ReadFile(l.ConfPath)
	if err != nil {
		return nil, err
	}
	return d, err
}

// Store writes data to the user's configuration directory at the given
// filename.
func (l *localStorage) Store(data []byte) error {
	dir := path.Dir(l.ConfPath)
	if err := ensureDirExists(dir); err != nil {
		return err
	}
	if err := os.WriteFile(l.ConfPath, data, 0600); err != nil {
		return err
	}
	return nil
}

// ensureDirExists checks for the existence of the directory at the given path,
// which is created if it does not exist.
func ensureDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0700); err != nil {
			return err
		}
	}
	return nil
}
