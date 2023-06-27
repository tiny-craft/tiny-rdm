package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"tinyrdm/backend/services"
	"tinyrdm/backend/storage"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	preferences := storage.NewPreferences()
	//connections := storage.NewConnections()
	connSvc := services.Connection()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Tiny RDM",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
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
			preferences,
			//connections,
			connSvc,
		},
		Mac: &mac.Options{
			//TitleBar:             mac.TitleBarHiddenInset(),
			//WebviewIsTransparent: true,
			//WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
