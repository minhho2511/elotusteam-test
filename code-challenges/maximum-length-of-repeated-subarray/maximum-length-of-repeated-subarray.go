package main

import "fmt"

func findLength(nums1 []int, nums2 []int) int {
	m, n := len(nums1), len(nums2)
	
	if m < n {
		return findLength(nums2, nums1)
	}
	
	prev := make([]int, n+1)
	maxLen := 0
	
	for i := 1; i <= m; i++ {
		curr := make([]int, n+1)
		for j := 1; j <= n; j++ {
			if nums1[i-1] == nums2[j-1] {
				curr[j] = prev[j-1] + 1
				if curr[j] > maxLen {
					maxLen = curr[j]
				}
			}
		}
		prev = curr
	}
	
	return maxLen
}

func main() {
	nums1 := []int{1, 2, 3, 2, 1}
	nums2 := []int{3, 2, 1, 4, 7}
	result1 := findLength(nums1, nums2)
	fmt.Printf("Example: %d (expected: 3)\n", result1)
}
