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
	"runtime"
	"tinyrdm/backend/services"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()
	connSvc := services.Connection()
	prefSvc := services.Preferences()

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Tiny RDM",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		Menu:      appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			connSvc.Start(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			connSvc.Stop(ctx)
		},
		Bind: []interface{}{
			app,
			connSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			About: &mac.AboutInfo{
				Title:   "Tiny RDM",
				Message: "Copyright © 2023",
				Icon:    icon,
			},
			//WebviewIsTransparent: true,
			//WindowIsTranslucent:  true,
		},
		Linux: &linux.Options{
			Icon:             icon,
			WebviewGpuPolicy: linux.WebviewGpuPolicyOnDemand,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
