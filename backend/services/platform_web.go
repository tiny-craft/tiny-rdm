//go:build web

package services

import (
	"context"
)

// Callback functions - set by api package at startup to avoid import cycle
var EmitEventFunc func(event string, data any)
var RegisterHandlerFunc func(event string, handler func(data any))

// Stub types to replace Wails runtime types in web mode

type OpenDialogOptions struct {
	Title           string
	ShowHiddenFiles bool
	Filters         []FileFilter
}

type SaveDialogOptions struct {
	Title           string
	ShowHiddenFiles bool
	DefaultFilename string
	Filters         []FileFilter
}

type FileFilter struct {
	Pattern string
}

type Screen struct {
	IsCurrent bool
	Size      ScreenSize
}

type ScreenSize struct {
	Width  int
	Height int
}

// EventsEmit emits an event via WebSocket (web mode)
func EventsEmit(ctx context.Context, event string, data ...any) {
	if EmitEventFunc == nil {
		return
	}
	if len(data) == 1 {
		EmitEventFunc(event, data[0])
	} else {
		EmitEventFunc(event, data)
	}
}

// EventsOnce registers a one-time event listener via WebSocket (web mode)
func EventsOnce(ctx context.Context, event string, callback func(data ...any)) func() {
	if RegisterHandlerFunc == nil {
		return func() {}
	}
	RegisterHandlerFunc(event, func(data any) {
		RegisterHandlerFunc(event, func(data any) {})
		callback(data)
	})
	return func() {
		RegisterHandlerFunc(event, func(data any) {})
	}
}

// EventsOn registers an event listener via WebSocket (web mode)
func EventsOn(ctx context.Context, event string, callback func(data ...any)) {
	if RegisterHandlerFunc == nil {
		return
	}
	RegisterHandlerFunc(event, func(data any) {
		callback(data)
	})
}

// EventsOff removes an event listener (web mode)
func EventsOff(ctx context.Context, event string) {
	if RegisterHandlerFunc == nil {
		return
	}
	RegisterHandlerFunc(event, func(data any) {})
}

// OpenFileDialog is a no-op in web mode
func OpenFileDialog(ctx context.Context, opts OpenDialogOptions) (string, error) {
	return "", nil
}

// SaveFileDialog is a no-op in web mode
func SaveFileDialog(ctx context.Context, opts SaveDialogOptions) (string, error) {
	return "", nil
}

// ScreenGetAll returns empty in web mode
func ScreenGetAll(ctx context.Context) ([]Screen, error) {
	return nil, nil
}

// WindowMaximise is a no-op in web mode
func WindowMaximise(ctx context.Context) {}

// WindowIsFullscreen returns false in web mode
func WindowIsFullscreen(ctx context.Context) bool { return false }

// WindowGetSize returns defaults in web mode
func WindowGetSize(ctx context.Context) (int, int) { return 1024, 768 }

// WindowIsMaximised returns false in web mode
func WindowIsMaximised(ctx context.Context) bool { return false }

// WindowIsMinimised returns false in web mode
func WindowIsMinimised(ctx context.Context) bool { return false }

// WindowIsNormal returns true in web mode
func WindowIsNormal(ctx context.Context) bool { return true }
