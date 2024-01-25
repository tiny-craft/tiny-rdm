package services

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	runtime2 "runtime"
	"sync"
	"time"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/types"
	sliceutil "tinyrdm/backend/utils/slice"
)

type systemService struct {
	ctx        context.Context
	appVersion string
}

var system *systemService
var onceSystem sync.Once

func System() *systemService {
	if system == nil {
		onceSystem.Do(func() {
			system = &systemService{
				appVersion: "0.0.0",
			}
			go system.loopWindowEvent()
		})
	}
	return system
}

func (s *systemService) Start(ctx context.Context, version string) {
	s.ctx = ctx
	s.appVersion = version

	// maximize the window if screen size is lower than the minimum window size
	if screen, err := runtime.ScreenGetAll(ctx); err == nil && len(screen) > 0 {
		for _, sc := range screen {
			if sc.IsCurrent {
				if sc.Size.Width < consts.MIN_WINDOW_WIDTH || sc.Size.Height < consts.MIN_WINDOW_HEIGHT {
					runtime.WindowMaximise(ctx)
					break
				}
			}
		}
	}
}

func (s *systemService) Info() (resp types.JSResp) {
	resp.Success = true
	resp.Data = struct {
		OS      string `json:"os"`
		Arch    string `json:"arch"`
		Version string `json:"version"`
	}{
		OS:      runtime2.GOOS,
		Arch:    runtime2.GOARCH,
		Version: s.appVersion,
	}
	return
}

// SelectFile open file dialog to select a file
func (s *systemService) SelectFile(title string, extensions []string) (resp types.JSResp) {
	filters := sliceutil.Map(extensions, func(i int) runtime.FileFilter {
		return runtime.FileFilter{
			Pattern: "*." + extensions[i],
		}
	})
	filepath, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title:           title,
		ShowHiddenFiles: true,
		Filters:         filters,
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"path": filepath,
	}
	return
}

// SaveFile open file dialog to save a file
func (s *systemService) SaveFile(title string, defaultName string, extensions []string) (resp types.JSResp) {
	filters := sliceutil.Map(extensions, func(i int) runtime.FileFilter {
		return runtime.FileFilter{
			Pattern: "*." + extensions[i],
		}
	})
	filepath, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		Title:           title,
		ShowHiddenFiles: true,
		DefaultFilename: defaultName,
		Filters:         filters,
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"path": filepath,
	}
	return
}

func (s *systemService) loopWindowEvent() {
	var fullscreen, maximised, minimised, normal bool
	var width, height int
	var dirty bool
	for {
		time.Sleep(300 * time.Millisecond)
		if s.ctx == nil {
			continue
		}

		dirty = false
		if f := runtime.WindowIsFullscreen(s.ctx); f != fullscreen {
			// full-screen switched
			fullscreen = f
			dirty = true
		}

		if w, h := runtime.WindowGetSize(s.ctx); w != width || h != height {
			// window size changed
			width, height = w, h
			dirty = true
		}

		if m := runtime.WindowIsMaximised(s.ctx); m != maximised {
			maximised = m
			dirty = true
		}

		if m := runtime.WindowIsMinimised(s.ctx); m != minimised {
			minimised = m
			dirty = true
		}

		if n := runtime.WindowIsNormal(s.ctx); n != normal {
			normal = n
			dirty = true
		}

		if dirty {
			runtime.EventsEmit(s.ctx, "window_changed", map[string]any{
				"fullscreen": fullscreen,
				"width":      width,
				"height":     height,
				"maximised":  maximised,
				"minimised":  minimised,
				"normal":     normal,
			})

			if !fullscreen && !minimised {
				// save window size and position
				Preferences().SaveWindowSize(width, height, maximised)
			}
		}
	}
}
