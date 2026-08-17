package main

import (
	"fmt"
	"math"
	"sync"
)

type Input struct {
	Index int
	Value int
}

type Result struct {
	Index int
	Value int
}

func worker(nums []Input, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	for _, inp := range nums {
		results <- Result{
			Index: inp.Index,
			Value: inp.Value * inp.Value,
		}
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	n := len(nums)

	var modifiedNums = make([]Input, n)
	for i, val := range nums {
		modifiedNums[i] = Input{
			Index: i,
			Value: val,
		}
	}

	var wg sync.WaitGroup
	results := make(chan Result)

	numWorkers := 10
	chunk := int(math.Ceil(float64(n) / float64(numWorkers)))

	start := 0
	for i := 0; i < numWorkers && start < n; i++ {
		wg.Add(1)
		end := start + chunk
		if end > n {
			end = n
		}
		go worker(modifiedNums[start:end], &wg, results)
		start = end
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	answer := make([]int, n)
	for res := range results {
		answer[res.Index] = res.Value
	}

	fmt.Println("Results in correct order:")
	for _, val := range answer {
		fmt.Printf("Squared number: %d\n", val)
	}
}
