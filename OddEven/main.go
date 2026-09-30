package main

import (
	"fmt"
	"sync"
)

func printOdd(odd <-chan struct{}, even chan<- struct{}, done chan struct{}, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= n; i += 2 {
		select {
		case <-odd:
			fmt.Println("Odd:", i)
		case <-done:
			fmt.Println("Odd exiting")
			return
		}

		if i+1 <= n {
			select {
			case even <- struct{}{}:
			case <-done:
				fmt.Println("Odd exiting")
				return
			}
		}
	}
	close(done)
}

func printEven(odd chan<- struct{}, even <-chan struct{}, done chan struct{}, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= n; i += 2 {
		select {
		case <-even:
			fmt.Println("Even:", i)
		case <-done:
			fmt.Println("Even exiting")
			return
		}

		if i+1 <= n {
			select {
			case odd <- struct{}{}:
			case <-done:
				fmt.Println("Even exiting")
				return
			}
		}
	}
	close(done)
}

func main() {
	var wg sync.WaitGroup
	odd := make(chan struct{})
	even := make(chan struct{})
	done := make(chan struct{})

	wg.Add(2)
	go printOdd(odd, even, done, 10, &wg)
	go printEven(odd, even, done, 20, &wg)

	odd <- struct{}{}
	wg.Wait()
}
