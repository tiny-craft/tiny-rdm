//go:build web

package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tinyrdm/backend/api"
	"tinyrdm/backend/services"
)

//go:embed all:frontend/dist
var assets embed.FS

var version = "0.0.0"

func main() {
	// Wire up event bridge callbacks (breaks import cycle: services -> api)
	services.EmitEventFunc = func(event string, data any) {
		api.Hub().Emit(event, data)
	}
	services.RegisterHandlerFunc = api.RegisterHandler

	// Initialize all services with a background context
	ctx := context.Background()

	sysSvc := services.System()
	connSvc := services.Connection()
	browserSvc := services.Browser()
	cliSvc := services.Cli()
	monitorSvc := services.Monitor()
	pubsubSvc := services.Pubsub()
	prefSvc := services.Preferences()
	prefSvc.SetAppVersion(version)
	prefSvc.UpdateEnv()

	// Start services
	sysSvc.Start(ctx, version)
	connSvc.Start(ctx)
	browserSvc.Start(ctx)
	cliSvc.Start(ctx)
	monitorSvc.Start(ctx)
	pubsubSvc.Start(ctx)

	services.GA().SetSecretKey("", "")

	// Initialize auth
	api.InitAuth()

	// Get port from env or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// Setup HTTP server
	router := api.SetupRouter(assets)
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		browserSvc.Stop()
		cliSvc.CloseAll()
		monitorSvc.StopAll()
		pubsubSvc.StopAll()
		srv.Close()
	}()

	log.Printf("Tiny RDM Web starting on http://0.0.0.0:%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
