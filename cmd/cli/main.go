package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "loadtester",
		Short: "LoadTester is a simple CLI for performing load tests on a web service.",
		Long:  `LoadTester is a CLI tool to test the load capacity of a web service by sending multiple HTTP requests.`,
		Run:   runLoadTest,
	}

	rootCmd.Flags().StringVar(&url, "url", "", "URL of the web service to test (required)")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Total number of requests to send (required)")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of simultaneous requests (required)")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func runLoadTest(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	var wg sync.WaitGroup
	requestsChan := make(chan int, requests)
	statusCodes := make(map[int]int)
	requestRefused := 0
	statusCodesMutex := &sync.Mutex{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestsChan {
				resp, err := http.Get(url)
				if err != nil {
					statusCodesMutex.Lock()
					requestRefused++
					statusCodesMutex.Unlock()

				}

				if resp != nil {
					statusCodesMutex.Lock()
					statusCodes[resp.StatusCode]++
					statusCodesMutex.Unlock()
					resp.Body.Close()
				}
			}
		}()
	}

	for i := 0; i < requests; i++ {
		requestsChan <- i
	}
	close(requestsChan)

	wg.Wait()

	totalTime := time.Since(startTime)

	fmt.Println("Load Test Completed")
	fmt.Printf("Total time taken: %v\n", totalTime)
	fmt.Printf("Total requests made: %d\n", requests)
	fmt.Printf("Requests with status 200: %d\n", statusCodes[200])
	fmt.Println("Other status codes distribution:")
	for code, count := range statusCodes {
		if code != 200 {
			fmt.Printf("Status %d: %d requests\n", code, count)
		}
	}
	fmt.Println("Refused requests: ", requestRefused)
}
