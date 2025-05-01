package main

import (
	"fmt"
	"log"

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

	// WIP: Run a speedtest once to verify execution
	fmt.Println("\nRunning a test speed test...")
	res, err := speedtestProvider.RunTest()
	if err != nil {
		log.Fatalf("Error running speedtest: %v", err)
	}
	fmt.Printf("Results: %v\n", res)

	// TODO: Save the result to database
	err = storage.SaveSpeedtestResult(res.DownMbps, res.UpMbps, res.PingMs, res.Timestamp)
	if err != nil {
		log.Printf("Error saving test speed test result: %v\n", err)
	}

	// Start monitoring
	monitor.StartMonitoring()

	// The application will exit here for now
}
