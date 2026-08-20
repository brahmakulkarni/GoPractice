package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	arr := []Person{
		{"Alice", 25},
		{"Bob", 20},
		{"Charlie", 23},
	}

	sort.Slice(arr, func(a, b int) bool {
		return arr[a].Age > arr[b].Age
	})

	fmt.Println(arr)
}
