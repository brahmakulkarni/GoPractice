package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func usingAtomics(init int) int {
	var counter atomic.Int64
	var wg sync.WaitGroup

	counter.Store(int64(init))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 1000 {
			counter.Add(1)
		}
	} ()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 1000 {
			counter.Add(1)
		}
	} ()

	wg.Wait()

	return int(counter.Load())
}

func usingChannels(init int) int {
	var wg sync.WaitGroup

	increments := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for range 1000 {
			increments <- 1
		}
	} ()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for range 1000 {
			increments <- 1
		}
	} ()

	go func() {
		wg.Wait()
		close(increments)
	} ()

	for increment := range increments {
		init += increment
	}

	return init
}

func main() {
	counter := 0

	resAtomic := usingAtomics(counter)
	resChannels := usingChannels(counter)

	fmt.Println("Using atomic:", resAtomic)
	fmt.Println("Using channels:", resChannels)
}