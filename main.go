package main

import (
	"fmt"
	"log"

	"github.com/deadpyxel/uptime-guard/internal/monitor"
	"github.com/deadpyxel/uptime-guard/internal/storage"
)

func main() {
	fmt.Println("Starting uptime-guard...")

	// Initialize storage (database)
	err := storage.InitializeDB("uptimeguard.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start monitoring
	monitor.StartMonitoring()

	// The application will exit here for now
}
