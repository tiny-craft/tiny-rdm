//go:build !web

package services

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Type aliases for Wails runtime types
type OpenDialogOptions = runtime.OpenDialogOptions
type SaveDialogOptions = runtime.SaveDialogOptions
type FileFilter = runtime.FileFilter
type Screen = runtime.Screen

// EventsEmit emits an event to the frontend (Wails desktop)
func EventsEmit(ctx context.Context, event string, data ...any) {
	runtime.EventsEmit(ctx, event, data...)
}

// EventsOnce registers a one-time event listener (Wails desktop)
func EventsOnce(ctx context.Context, event string, callback func(data ...any)) func() {
	return runtime.EventsOnce(ctx, event, callback)
}

// EventsOn registers an event listener (Wails desktop)
func EventsOn(ctx context.Context, event string, callback func(data ...any)) {
	runtime.EventsOn(ctx, event, callback)
}

// EventsOff removes an event listener (Wails desktop)
func EventsOff(ctx context.Context, event string) {
	runtime.EventsOff(ctx, event)
}

// OpenFileDialog opens a native file dialog (Wails desktop)
func OpenFileDialog(ctx context.Context, opts OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(ctx, opts)
}

// SaveFileDialog opens a native save dialog (Wails desktop)
func SaveFileDialog(ctx context.Context, opts SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(ctx, opts)
}

// ScreenGetAll gets all screens (Wails desktop)
func ScreenGetAll(ctx context.Context) ([]Screen, error) {
	return runtime.ScreenGetAll(ctx)
}

// WindowMaximise maximises the window (Wails desktop)
func WindowMaximise(ctx context.Context) {
	runtime.WindowMaximise(ctx)
}

// WindowIsFullscreen checks if window is fullscreen (Wails desktop)
func WindowIsFullscreen(ctx context.Context) bool {
	return runtime.WindowIsFullscreen(ctx)
}

// WindowGetSize gets window size (Wails desktop)
func WindowGetSize(ctx context.Context) (int, int) {
	return runtime.WindowGetSize(ctx)
}

// WindowIsMaximised checks if window is maximised (Wails desktop)
func WindowIsMaximised(ctx context.Context) bool {
	return runtime.WindowIsMaximised(ctx)
}

// WindowIsMinimised checks if window is minimised (Wails desktop)
func WindowIsMinimised(ctx context.Context) bool {
	return runtime.WindowIsMinimised(ctx)
}

// WindowIsNormal checks if window is normal (Wails desktop)
func WindowIsNormal(ctx context.Context) bool {
	return runtime.WindowIsNormal(ctx)
}
