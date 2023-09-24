package main

import (
	"math"
)

// compute L_2 distance
func distance(p1 [2]float64, p2 [2]float64) float64 {
	return math.Sqrt((p1[0]-p2[0])*(p1[0]-p2[0]) + (p1[1]-p2[1])*(p1[1]-p2[1]))
}

func distMatrix(points [][2]float64) [][]float64 {

	npoints := len(points)
	// initialise distance matrix
	dist := make([][]float64, npoints)
	for i := range dist {
		dist[i] = make([]float64, npoints)
	}
	// compute distance matrix
	for i := 0; i < npoints; i++ {
		for j := 0; j < npoints; j++ {
			dist[i][j] = distance(points[i], points[j])
		}
	}
	return dist
}

// total distance around a given route
func travelDist(state []int, dist [][]float64) float64 {

	td := 0.0
	np := len(state)
	for i := range state {
		td += dist[state[i%np]][state[(i+1)%np]]
	}
	return td
}
