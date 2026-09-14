package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	numStrings, _ := in.ReadString('\n')
	numStrings = strings.TrimSpace(numStrings)
	numStringInt, _ := strconv.Atoi(numStrings)
	fmt.Println("numStringInt:", numStringInt)

	for range numStringInt {
		greeting, _ := in.ReadString('\n')
		out.WriteString(greeting)
	}
}
