package main

import (
	"fmt"
	"sync"
)

type Result struct {
	ChunkIndex int
	ChunkSum   int
}

func processChunk(chunkIndex int, data []int, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	sum := 0

	for _, v := range data {
		sum += v
	}

	results <- Result{
		ChunkIndex: chunkIndex,
		ChunkSum: sum,
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	n := len(nums)
	chunkSize := 3

	results := make(chan Result)
	var wg sync.WaitGroup

	chunkIndex := 0

	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize

		if end > n {
			end = n
		}

		chunk := nums[i:end]

		wg.Add(1)
		go processChunk(chunkIndex, chunk, &wg, results)

		chunkIndex++
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	orderedResult := make([]int, chunkIndex)
	for output := range results {
		orderedResult[output.ChunkIndex] = output.ChunkSum
	}

	fmt.Println("Results of each chunk:")
	for index, sum := range orderedResult {
		fmt.Printf("Received result: Index = %d, Sum = %d\n", index, sum)
	}
}
