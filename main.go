package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: 引数が足りません")
		return
	}
	url := os.Args[1]

	start := time.Now()
	resp, err := http.Get(url)
	end := time.Since(start)

	if err != nil {
		fmt.Println("Error")
	} else {
		fmt.Printf("%d %s\n", resp.StatusCode, end)
	}
}
