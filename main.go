package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/WeMeet-server/pkg/controllers"
	"github.com/retawsolit/WeMeet-server/pkg/factory"
	"github.com/retawsolit/WeMeet-server/pkg/services"
)

func main() {
	// Load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize connections
	if err := factory.NewDatabaseConnection(cfg); err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := factory.NewRedisConnection(cfg); err != nil {
		log.Fatal("Failed to connect redis:", err)
	}

	if err := factory.NewNatsConnection(cfg); err != nil {
		log.Fatal("Failed to connect NATS:", err)
	}

	// Initialize services
	roomService := services.NewRoomService(cfg)

	// Initialize controllers
	roomController := controllers.NewRoomController(roomService)
	healthController := controllers.NewHealthController(cfg)

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/health", healthController.Health)
	mux.HandleFunc("/api/room", roomController.CreateRoom)
	mux.HandleFunc("/api/room/", roomController.GetRoom)
	mux.HandleFunc("/api/rooms", roomController.ListRooms)

	// Static files
	if cfg.Client.Path != "" {
		fs := http.FileServer(http.Dir(cfg.Client.Path))
		mux.Handle("/", fs)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Client.Port),
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		// Close connections
		if cfg.DB != nil {
			if sqlDB, err := cfg.DB.DB(); err == nil {
				sqlDB.Close()
			}
		}
		if cfg.RDS != nil {
			cfg.RDS.Close()
		}
		if cfg.NatsConn != nil {
			cfg.NatsConn.Close()
		}
	}()

	log.Printf("Server starting on port %d", cfg.Client.Port)
	log.Printf("Serving static from: %s", cfg.Client.Path)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}
