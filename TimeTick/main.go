package main

import (
	"fmt"
	"time"
)

func main() {
	limiter := time.Tick(500 * time.Millisecond)

	for i := 1; i <= 5; i++ {
		<-limiter
		fmt.Println("Request", i, "at", time.Now())
	}
}