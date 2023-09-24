package main

// 2-bond move class: reverse the subchain between 2 indices
func reverse(i int, j int, perm []int) {
	switch i < j {
	case true:
		reverseSlice(perm[i:(j + 1)])
	case false:
		reverseSlice(perm[j:(i + 1)])
	}
}
func reverseSlice(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// energy delta for 2-bond reverse
func reverseDelta(i int, j int, perm []int, dist [][]float64) float64 {
	np := len(perm)
	if i == j || (i == 0 && j == np-1) || (j == 0 && i == np-1) {
		return 0.0
	}
	dd := 0.0
	switch i < j {
	case true:
		// i before j
		dd -= dist[perm[j%np]][perm[(j+1)%np]]
		dd -= dist[perm[(np+i-1)%np]][perm[i%np]]
		dd += dist[perm[(np+i-1)%np]][perm[j%np]]
		dd += dist[perm[i%np]][perm[(j+1)%np]]
	case false:
		// j before i
		dd -= dist[perm[i%np]][perm[(i+1)%np]]
		dd -= dist[perm[(np+j-1)%np]][perm[j%np]]
		dd += dist[perm[(np+j-1)%np]][perm[i%np]]
		dd += dist[perm[j%np]][perm[(i+1)%np]]
	}
	return dd
}

// simplest move class: swap the permutation values at 2 indices
func swap(i int, j int, perm []int) {

	a, b := perm[i], perm[j]
	perm[i], perm[j] = b, a
}

// energy delta for swap
func swapDelta(i int, j int, perm []int, dist [][]float64) float64 {

	if i == j {
		return 0.0
	}
	dd := 0.0
	np := len(perm)
	if (i-j-1)%np == 0 {
		// j immediately before i
		dd -= dist[perm[i%np]][perm[(i+1)%np]]
		dd -= dist[perm[(np+j-1)%np]][perm[j%np]] // add np to first index to avoid -1%np
		dd += dist[perm[(np+j-1)%np]][perm[i%np]]
		dd += dist[perm[j%np]][perm[(i+1)%np]]
	} else if (i-j+1)%np == 0 {
		// i immediately before j
		dd -= dist[perm[j%np]][perm[(j+1)%np]]
		dd -= dist[perm[(np+i-1)%np]][perm[i%np]]
		dd += dist[perm[(np+i-1)%np]][perm[j%np]]
		dd += dist[perm[i%np]][perm[(j+1)%np]]
	} else {
		// i,j separated mod npoints
		dd -= dist[perm[(np+i-1)%np]][perm[i%np]]
		dd -= dist[perm[i%np]][perm[(i+1)%np]]
		dd -= dist[perm[(np+j-1)%np]][perm[j%np]]
		dd -= dist[perm[j%np]][perm[(j+1)%np]]
		dd += dist[perm[(np+i-1)%np]][perm[j%np]]
		dd += dist[perm[j%np]][perm[(i+1)%np]]
		dd += dist[perm[(np+j-1)%np]][perm[i%np]]
		dd += dist[perm[i%np]][perm[(j+1)%np]]
	}
	return dd
}
