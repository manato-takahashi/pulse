package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	URL string `yaml:"url"`
}

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

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: yamlが読み込めません")
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error: yamlを構造体に変換できません")
		return
	}

	ch := make(chan HealthResult)
	for _, ep := range config.Endpoints {
		go func(url string) {
			ch <- checkHealth(url)
		}(ep.URL)
	}

	for range config.Endpoints {
		result := <-ch
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
}
