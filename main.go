package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/retawsolit/WeMeet-server/pkg/config"
)

func main() {
	// Trong dev: copy config.yaml đã có sẵn ở root khi chạy app
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	addr := fmt.Sprintf(":%d", cfg.Client.Port)
	log.Println("HTTP listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
