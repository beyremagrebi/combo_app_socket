package main

import (
	"log"

	"github.com/beyremagrebi/combo_app_socket/internal/config"
	"github.com/beyremagrebi/combo_app_socket/internal/realtime"
)

func main() {
	cfg := config.Load()

	server := realtime.NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}