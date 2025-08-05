package main

import "fmt"

func sumOfDistancesInTree(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}
	
	graph := make([][]int, n)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	
	subtreeSize := make([]int, n)
	subtreeSum := make([]int, n)
	answer := make([]int, n)
	
	// First DFS: Calculate subtree information from root 0
	var dfs1 func(node, parent int)
	dfs1 = func(node, parent int) {
		subtreeSize[node] = 1
		
		for _, child := range graph[node] {
			if child != parent {
				dfs1(child, node)
				subtreeSize[node] += subtreeSize[child]
				subtreeSum[node] += subtreeSum[child] + subtreeSize[child]
			}
		}
	}
	
	// Second DFS: Re-root to calculate answer for each node  
	var dfs2 func(node, parent int)
	dfs2 = func(node, parent int) {
		for _, child := range graph[node] {
			if child != parent {
				answer[child] = answer[node] - subtreeSize[child] + (n - subtreeSize[child])
				dfs2(child, node)
			}
		}
	}
	
	dfs1(0, -1)
	answer[0] = subtreeSum[0]
	dfs2(0, -1)
	
	return answer
}

func main() {
	n1 := 6
	edges1 := [][]int{{0, 1}, {0, 2}, {2, 3}, {2, 4}, {2, 5}}
	result1 := sumOfDistancesInTree(n1, edges1)
	fmt.Printf("Example: %v\n", result1)
}
