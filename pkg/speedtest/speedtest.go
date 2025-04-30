package speedtest

import (
	"fmt"
	"time"
)

// Result holds the result of a speedtest
type Result struct {
	DownMbps  float64
	UpMbps    float64
	PingMs    float64
	Timestamp time.Time
}

func RunSpeedtest() (Result, error) {
	fmt.Println("Running simulated speed test...")

	return Result{
		DownMbps:  100.0,
		UpMbps:    20.0,
		PingMs:    15.0,
		Timestamp: time.Now(),
	}, nil
}
