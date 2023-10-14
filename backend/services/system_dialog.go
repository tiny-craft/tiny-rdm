package services

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"sync"
	"tinyrdm/backend/types"
)

type systemService struct {
	ctx context.Context
}

var system *systemService
var onceSystem sync.Once

func System() *systemService {
	if system == nil {
		onceSystem.Do(func() {
			system = &systemService{}
		})
	}
	return system
}

func (s *systemService) Start(ctx context.Context) {
	s.ctx = ctx
}

// SelectFile open file dialog to select a file
func (s *systemService) SelectFile(title string) (resp types.JSResp) {
	filepath, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title:           title,
		ShowHiddenFiles: true,
	})
	if err != nil {
		log.Println(err)
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"path": filepath,
	}
	return
}
