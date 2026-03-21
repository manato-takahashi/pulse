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

func printResults(results []HealthResult) bool {
	// 1. 最長URLの長さを調べる
	// 2. 各結果を幅揃え + 色付き + ステータス名で表示
	// 3. 1つでも失敗があったら true を返す

	maxUrlLen := 0
	maxStatusTextLen := 0
	for _, r := range results {
		if len(r.URL) > maxUrlLen {
			maxUrlLen = len(r.URL)
		}
		if r.Err != nil {
			if len("Connection Error") > maxStatusTextLen {
				maxStatusTextLen = len("Connection Error")
			}
		} else {
			if len(http.StatusText(r.Status)) > maxStatusTextLen {
				maxStatusTextLen = len(http.StatusText(r.Status))
			}
		}
	}

	hasFailed := false

	for _, result := range results {
		var mark string
		if result.Status >= 200 && result.Status < 300 {
			mark = "\033[32m✓\033[0m"
		} else {
			mark = "\033[31m×\033[0m"
			hasFailed = true
		}

		if result.Err != nil {
			hasFailed = true
			fmt.Printf("%s %-*s  --- %-*s  %s\n", mark, maxUrlLen, result.URL, maxStatusTextLen, "Connection Error", result.Duration.Round(time.Millisecond))
		} else {
			fmt.Printf("%s %-*s  %d %-*s  %s\n", mark, maxUrlLen, result.URL, result.Status, maxStatusTextLen, http.StatusText(result.Status), result.Duration.Round(time.Millisecond))
		}
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
	results := []HealthResult{}
	for range config.Endpoints {
		result := <-ch
		results = append(results, result)
	}

	hasFailed = printResults(results)

	if *exitOnFail && hasFailed {
		os.Exit(1)
	}
}
