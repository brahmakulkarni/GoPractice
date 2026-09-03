package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	s, _ := in.ReadString('\n')
	s = strings.TrimSpace(s)

	fmt.Fprintln(out, s)
}
