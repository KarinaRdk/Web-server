package main

import (
	"TestWebServer/internal/config"
	"TestWebServer/internal/storage"
	"TestWebServer/internal/subscriber"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	pg, err := storage.InitDatabase(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pg.Close()

	logger := setupLogger("development")
	logger.Info("Database initialized successfully")
	db := storage.NewDatabase(pg)
	sub := subscriber.NewSubscriber(db)
	sub.Receive()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
