package main

import (
	"fmt"
	"os"
)

func readFile(name string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from:", r)
		}
	}()
	file, err := os.Open(name)
	if err != nil {
		panic(err) // escalate unexpected error
	}
	defer file.Close()
	fmt.Println("File opened successfully:", name)
}

func main() {
	readFile("nonexistent.txt")
	fmt.Println("Program still running")
}
