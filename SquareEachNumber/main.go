package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Index int
	Value int
}

func squareNumber(num int, index int, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	// out := num*num

	results <- Result{
		Index: index,
		Value: num*num,
	}
}

func main() {
	nums := []int{1,2,3,4,5}

	var wg sync.WaitGroup
	results := make(chan Result)

	for i, num := range nums {
		wg.Add(1)
		go squareNumber(num, i, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	
	var answer [5]int
	// fmt.Println("Processing numbers:")
	for res := range results {
		answer[res.Index] = res.Value
		// fmt.Printf("Squared number: %d\n", res)
	}
	fmt.Println("Results in correct order:")
	for _, val := range answer {
		fmt.Printf("Squared number: %d\n", val)
	}
}