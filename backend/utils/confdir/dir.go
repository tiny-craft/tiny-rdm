package confdir

import (
	"log"
	"os"
	"sync"
)

var confdirOnceValue = sync.OnceValue(func() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("get user config dir failed: %v, fallback to home dir", err)
		dir, _ = os.UserHomeDir()
	}
	return dir
})

func GetConfigDir() string {
	return confdirOnceValue()
}
