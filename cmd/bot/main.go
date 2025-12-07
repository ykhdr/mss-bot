package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ykhdr/mss-bot/internal/app"
)

func main() {
	configPath := flag.String("config", "configs/config.kdl", "path to config file")
	flag.Parse()

	application, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Application error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")

	if err := application.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Application stopped")
}
