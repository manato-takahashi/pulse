package main

import (
	"flag"
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

func printResult(result HealthResult) bool {
	var mark string
	hasFailed := false

	if result.Status >= 200 && result.Status < 300 {
		mark = "✓"
	} else {
		mark = "×"
		hasFailed = true
	}

	if result.Err != nil {
		hasFailed = true
		fmt.Printf("× %s --- %s\n", result.URL, result.Duration.Round(time.Millisecond))
	} else {
		fmt.Printf("%s %s %d %s\n", mark, result.URL, result.Status, result.Duration.Round(time.Millisecond))
	}

	return hasFailed
}

func main() {
	configFile := flag.String("config", "", "設定ファイルのパス")
	exitOnFail := flag.Bool("exit-on-fail", false, "1つでも×があればエラーとするか(CI/CD用の設定)")
	flag.Parse()

	if *configFile == "" {
		fmt.Println("Error: --config フラグが必要です")
		return
	}

	data, err := os.ReadFile(*configFile)
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

	hasFailed := false
	for range config.Endpoints {
		result := <-ch
		if printResult(result) {
			hasFailed = true
		}
	}

	if *exitOnFail && hasFailed {
		os.Exit(1)
	}
}
