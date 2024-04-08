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
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
	"runtime"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/services"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"
var gaMeasurementID, gaSecretKey string

func main() {
	// Create an instance of the app structure
	sysSvc := services.System()
	connSvc := services.Connection()
	browserSvc := services.Browser()
	cliSvc := services.Cli()
	monitorSvc := services.Monitor()
	pubsubSvc := services.Pubsub()
	prefSvc := services.Preferences()
	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()
	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	app := &options.App{
		Title:                    "Tiny RDM",
		Width:                    windowWidth,
		Height:                   windowHeight,
		MinWidth:                 consts.MIN_WINDOW_WIDTH,
		MinHeight:                consts.MIN_WINDOW_HEIGHT,
		WindowStartState:         windowStartState,
		Frameless:                runtime.GOOS != "darwin",
		Menu:                     appMenu,
		EnableDefaultContextMenu: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: options.NewRGBA(27, 38, 54, 0),
		StartHidden:      true,
		OnStartup: func(ctx context.Context) {
			sysSvc.Start(ctx, version)
			connSvc.Start(ctx)
			browserSvc.Start(ctx)
			cliSvc.Start(ctx)
			monitorSvc.Start(ctx)
			pubsubSvc.Start(ctx)

			services.GA().SetSecretKey(gaMeasurementID, gaSecretKey)
			services.GA().Startup(version)
		},
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			runtime2.WindowSetPosition(ctx, x, y)
			runtime2.WindowShow(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := runtime2.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		OnShutdown: func(ctx context.Context) {
			browserSvc.Stop()
			cliSvc.CloseAll()
			monitorSvc.StopAll()
			pubsubSvc.StopAll()
		},
		Bind: []interface{}{
			sysSvc,
			connSvc,
			browserSvc,
			cliSvc,
			monitorSvc,
			pubsubSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Tiny RDM " + version,
				Message: "A modern lightweight cross-platform Redis desktop client.\n\nCopyright © 2024",
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
			ProgramName:         "Tiny RDM",
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	}

	// run application
	if err := wails.Run(app); err != nil {
		println("Error:", err.Error())
	}
}
