package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Buffered input/output
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// Read size
	var n int
	fmt.Println("Enter the number of elements in your array:")
	fmt.Fscan(in, &n)

	// Create slice of size n
	arr := make([]int, n)

	// Read elements
	for i := 0; i < n; i++ {
		fmt.Printf("Enter the number at index %d:", i)
		fmt.Fscan(in, &arr[i])
	}

	// Print the array
	for i := 0; i < n; i++ {
		fmt.Fprint(out, arr[i]*arr[i], " ")
	}
	fmt.Fprintln(out)
}
