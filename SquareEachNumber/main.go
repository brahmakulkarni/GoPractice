package main

import (
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

func worker(jobs <-chan Job, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	for job := range jobs {
		results <- Result{
			Index: job.Index,
			Value: job.Value * job.Value,
		}
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	n := len(nums)

	var wg sync.WaitGroup
	// 1. Create input and output channels
	jobs := make(chan Job)
	results := make(chan Result)

	// 2. Spawn the workers
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, results)
	}

	// 3. Feed jobs channel and close it once done sending all jobs
	go func() {
		defer close(jobs)
		for i, num := range nums {
			jobs <- Job{
				Index: i,
				Value: num,
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
	for res := range results {
		answer[res.Index] = res.Value
	}

	fmt.Println("Results in correct order:")
	for _, val := range answer {
		fmt.Printf("Squared number: %d\n", val)
	}
}
