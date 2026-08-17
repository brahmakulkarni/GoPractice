package main

import (
	"fmt"
	"sync"
)

func squareNumber(num int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()

	// out := num*num

	results <- num*num
}

func main() {
	nums := []int{1,2,3,4,5}

	var wg sync.WaitGroup
	results := make(chan int)

	for _, num := range nums {
		wg.Add(1)
		go squareNumber(num, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Processing numbers:")
	for res := range results {
		fmt.Printf("Squared number: %d\n", res)
	}
}