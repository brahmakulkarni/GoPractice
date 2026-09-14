package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Job struct {
	Index int
	Value int
}

type Result struct {
	Index int
	Value int
}

func process(ctx context.Context, nums []int, workerCount int) ([]int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	n := len(nums)

	// 1. Create input and output channels
	jobs := make(chan Job)
	results := make(chan Result)
	errChan := make(chan error, 1)

	// 2. Spawn the workers
	for range workerCount {
		wg.Add(1)
		go worker(ctx, jobs, results, errChan, &wg)
	}

	// 3. Feed jobs channel and close it once done sending all jobs
	go func() {
		defer close(jobs)
		for i, num := range nums {
			select {
			case jobs <- Job{
				Index: i,
				Value: num,
			}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// 4. Close the results channel once all workers have completed processing and shipping the results
	go func() {
		wg.Wait()
		close(results)
	}()


	// 5. Save results in the correct order
	answer := make([]int, n)
	for {
		select {
		case err := <-errChan:
			cancel()
			return nil, err
		case res, ok := <-results:
			if !ok {
				return answer, nil
			}
			answer[res.Index] = res.Value
		}
	}
}

func computeSquare(num int) (int, error) {
	if num == 6 {
		return 0, errors.New("cannot process 6")
	}
	return num * num, nil
}

func worker(ctx context.Context, jobs <-chan Job, results chan<- Result, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			squared, err := computeSquare(job.Value)
			if err != nil {
				select {
				case errChan <- err:
				default:
					fmt.Println("One error already received. Don't need to report again")
				}

				return
			}

			select {
			case <-ctx.Done():
				return
			case results <- Result{
				Index: job.Index,
				Value: squared,
			}:
			}
		}
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 7, 8, 9, 10}

	ctx := context.Background()

	answer, err := process(ctx, nums, 10)

	if err != nil {
		fmt.Println("Processing failed:", err)
		return
	}

	fmt.Println("Results in correct order:")
	for _, val := range answer {
		fmt.Printf("Squared number: %d\n", val)
	}
}
