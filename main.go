package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/retawsolit/WeMeet-server/pkg/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// serve static từ client.path (ví dụ /app/client/dist)
	fs := http.FileServer(http.Dir(cfg.Client.Path))
	mux.Handle("/", fs)

	addr := fmt.Sprintf(":%d", cfg.Client.Port)
	log.Println("Serving static from:", cfg.Client.Path)
	log.Println("HTTP listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
