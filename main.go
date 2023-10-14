package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"runtime"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/services"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"

func main() {
	// Create an instance of the app structure
	sysSvc := services.System()
	connSvc := services.Connection()
	prefSvc := services.Preferences()
	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight := prefSvc.GetWindowSize()

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Tiny RDM",
		Width:     windowWidth,
		Height:    windowHeight,
		MinWidth:  consts.DEFAULT_WINDOW_WIDTH,
		MinHeight: consts.DEFAULT_WINDOW_HEIGHT,
		Frameless: runtime.GOOS != "darwin",
		Menu:      appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 0},
		OnStartup: func(ctx context.Context) {
			sysSvc.Start(ctx)
			connSvc.Start(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			connSvc.Stop(ctx)
		},
		Bind: []interface{}{
			sysSvc,
			connSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Tiny RDM " + version,
				Message: "A modern lightweight cross-platform Redis desktop client.\n\nCopyright © 2023",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		Linux: &linux.Options{
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
