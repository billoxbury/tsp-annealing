package main

import (
	"time"
)

type tspProblem struct {
	points [][2]float64
	labels []string
	dist   [][]float64
}

type annealParam struct {
	temperature float64
	cooling     float64
	period      int
	maxIter     int
	countdown   int
}

type walker interface {
	testDelta(float64) int   // nr errors in iterLimit calls
	timeMove() time.Duration // running time for iterLimit calls
	timeEnergy() time.Duration
	timeDelta() time.Duration
	run() float64
}

type tspWalker struct {
	problem tspProblem
	param   annealParam
	state   []int
	move    func(int, int, []int)
	delta   func(int, int, []int, [][]float64) float64
	verbose bool
}
