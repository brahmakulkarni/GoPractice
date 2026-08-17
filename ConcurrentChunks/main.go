package main

import (
	"fmt"
	"sync"
)

func processChunk(data []int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()

	sum := 0

	for _, v := range data {
		sum += v
	}

	results <- sum
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	n := len(nums)
	chunkSize := 3

	results := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize

		if end > n {
			end = n
		}

		chunk := nums[i:end]

		wg.Add(1)
		go processChunk(chunk, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Results of each chunk:")
	for output := range results {
		fmt.Printf("Received result: %d\n", output)
	}
}
