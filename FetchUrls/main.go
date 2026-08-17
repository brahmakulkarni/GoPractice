package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Job struct {
	Index int
	URL   string
}

type URLResult struct {
	Index int
	URL   string
	Body  string
}

func fetchResult(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error while creating req with context: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while executing GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(body), nil
}

func worker(ctx context.Context, jobs <-chan Job, results chan<- URLResult, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			body, err := fetchResult(ctx, job.URL)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			select {
			case <-ctx.Done():
				return
			case results <- URLResult{
				Index: job.Index,
				URL:   job.URL,
				Body:  body,
			}:
			}
		}
	}
}

func FetchAll(ctx context.Context, urls []string, workerCount int) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	n := len(urls)

	jobs := make(chan Job)
	results := make(chan URLResult)
	errChan := make(chan error, 1)

	// 1. Spawn workers
	for range workerCount {
		wg.Add(1)
		go worker(ctx, jobs, results, errChan, &wg)
	}

	// 2. Dispatch jobs
	go func() {
		defer close(jobs)

		for i, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- Job{
				Index: i,
				URL:   url,
			}:
			}
		}
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	answer := make([]string, n)
	for {
		select {
		case err := <-errChan:
			cancel()
			return nil, err
		case res, ok := <-results:
			if !ok {
				return answer, nil
			}
			answer[res.Index] = res.Body
		}
	}
}

func main() {
	ctx := context.Background()

	// all work
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/get?name=go",
		"https://httpbin.org/headers",
		"https://httpbin.org/user-agent",
		"https://httpbin.org/ip",
	}

	// one will fail
	// urls := []string{
	// 	"https://httpbin.org/get",
	// 	"https://httpbin.org/get?name=go",
	// 	"https://this-domain-definitely-does-not-exist-123456789.com",
	// 	"https://httpbin.org/headers",
	// 	"https://httpbin.org/ip",
	// }

	results, err := FetchAll(ctx, urls, 2)
	if err != nil {
		fmt.Println("FetchAll failed:", err)
		return
	}

	fmt.Printf("Successfully fetched %d URLs\n\n", len(results))

	for i, body := range results {
		fmt.Printf("URL %d: %d bytes\n", i, len(body))
	}
}
