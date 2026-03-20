package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type HealthResult struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

func checkHealth(url string) HealthResult {
	start := time.Now()
	res, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return HealthResult{URL: url, Duration: duration, Err: err}
	}

	return HealthResult{
		URL:      url,
		Status:   res.StatusCode,
		Duration: duration,
		Err:      err,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: 引数が足りません")
		return
	}
	url := os.Args[1]

	result := checkHealth(url)

	var mark string
	if result.Status >= 200 && result.Status < 300 {
		mark = "✓"
	} else {
		mark = "×"
	}

	if result.Err != nil {
		fmt.Printf("× %s --- %s\n", result.URL, result.Duration.Round(time.Millisecond))
	} else {
		fmt.Printf("%s %s %d %s\n", mark, result.URL, result.Status, result.Duration.Round(time.Millisecond))
	}
}
