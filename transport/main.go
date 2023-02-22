package main

import "fmt"

func main() {
	fmt.Println(tableSizeFor(19))
}

func tableSizeFor(source int) int {
	maxCapacity := 1 << 30
	n := source - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return 1
	} else if n >= maxCapacity {
		return maxCapacity
	}
	return n + 1
}
