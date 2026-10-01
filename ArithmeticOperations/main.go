package main

import (
	"fmt"
	"sync"
)

// 0 - add, 1 - subtract, 2- multiply, 3 - divide 4- modulo

type Result struct {
	Index int
	Value int
}

func add(a int, b int, res chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	res <- Result{
		Index: 0,
		Value: a + b,
	}
}

func subtract(a int, b int, res chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	res <- Result{
		Index: 1,
		Value: a - b,
	}
}

func multiply(a int, b int, res chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	res <- Result{
		Index: 2,
		Value: a * b,
	}
}

func divide(a int, b int, res chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	if b == 0 {
		res <- Result{
			Index: 3,
			Value: 0,
		}
		return
	}
	res <- Result{
		Index: 3,
		Value: a / b,
	}
}

func modulo(a int, b int, res chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	if b == 0 {
		res <- Result{
			Index: 4,
			Value: 0,
		}
		return
	}
	res <- Result{
		Index: 4,
		Value: a % b,
	}
}

func main() {
	inp := [][]int{{10, 2}, {7, 3}, {9, 0}, {5, 5}, {20, 4}, {8, 3}}
	var wg sync.WaitGroup
	// res := make(chan Result)

	answer := [][]int{}

	for _, pair := range inp {
		res := make(chan Result)
		wg.Add(5)
		go add(pair[0], pair[1], res, &wg)
		go subtract(pair[0], pair[1], res, &wg)
		go multiply(pair[0], pair[1], res, &wg)
		go divide(pair[0], pair[1], res, &wg)
		go modulo(pair[0], pair[1], res, &wg)

		go func() {
			wg.Wait()
			close(res)
		} ()

		currAns := make([]int, 5)
		for currRes := range res {
			currAns[currRes.Index] = currRes.Value
		}
		answer = append(answer, currAns)
		// close(res)
	}

	fmt.Println("Final Answer:", answer)
}
