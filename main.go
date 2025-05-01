package main

import (
	"fmt"
	"log"
	"time"

	"github.com/deadpyxel/uptime-guard/internal/monitor"
	"github.com/deadpyxel/uptime-guard/internal/storage"
	"github.com/deadpyxel/uptime-guard/pkg/speedtest"
)

func main() {
	fmt.Println("Starting uptime-guard...")

	// Initialize storage (database)
	err := storage.InitializeDB("uptimeguard.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		err := storage.CloseDB()
		if err != nil {
			log.Printf("Error closing database: %v\n", err)
		}
	}()

	// Create a speed test provider using the showwin/speedtest-go library
	speedtestProvider, err := speedtest.NewSpeedtestGoProvider()
	if err != nil {
		log.Fatalf("Error creating speedtest provider: %v", err)
	}

	cfg := monitor.Config{
		SpeedtestInterval: 30 * time.Minute, // Run Speed test every 30 minutes
		ConnCheckInterval: 30 * time.Second, // Check connectivity every 30 seconds
		DBPath:            "uptime.db",      // Maybe this can be removed since we do not operate directly on the DB
	}
	// Start monitoring
	monitor.StartMonitoring(cfg, speedtestProvider)

	// The application will exit here for now
}
