package main

import (
	"TestWebServer/internal/cache"
	"TestWebServer/internal/config"
	"TestWebServer/internal/handler"
	"TestWebServer/internal/logger"
	"TestWebServer/internal/service"
	"TestWebServer/internal/storage"
	"TestWebServer/internal/subscriber"
	"TestWebServer/server"
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	// Ensure functioning of logger.
	l := logger.SetupLogger("development")
	l.Info("Database initialized successfully")

	// Initializing database
	pg, err := storage.InitDatabase(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pg.Close()
	db := storage.NewDatabase(pg)

	// Initializing cache.
	c := cache.New()

	// Creating a subscriber to receive messages via a broker.
	sub := subscriber.New(db, c)
	go sub.Receive()

	s := service.New(db, c)
	h := handler.New(s)

	// Create a ServeMux and register routes.
	mux := http.NewServeMux()
	handler.RegisterRoutes(h, mux)

	// Create and start the server with the ServeMux.
	serv := server.New(cfg.HTTPServer.Address, mux)
	serv.Start()

	// Wait for an interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shut down the server.
	serv.Stop()
}
