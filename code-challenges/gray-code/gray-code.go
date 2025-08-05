package main

import "fmt"

func grayCode(n int) []int {
	result := make([]int, 1<<n)
	
	for i := 0; i < (1 << n); i++ {
		result[i] = i ^ (i >> 1)
	}
	
	return result
}

func main() {
	fmt.Println("Input: n = 2")
	fmt.Printf("Output: %v\n", grayCode(2))
}
