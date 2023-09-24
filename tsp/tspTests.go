package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func (w tspWalker) testDelta(tolerance float64) int {

	prob := w.problem
	par := w.param

	errCount := 0
	npoints := len(prob.dist)
	for iter := 0; iter < par.maxIter; iter++ {

		old_d := travelDist(w.state, prob.dist)
		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		delta_d := w.delta(i, j, w.state, prob.dist)
		w.move(i, j, w.state)
		new_d := travelDist(w.state, prob.dist)
		// print errors
		if math.Abs(old_d+delta_d-new_d) > tolerance {
			fmt.Println(old_d, i, j, delta_d, new_d)
			errCount++
		}
	}

	fmt.Printf("Found %d errors out of %d at tolerance %v\n", errCount, par.maxIter, tolerance)
	return errCount
}

func (w tspWalker) timeMove() time.Duration {

	prob := w.problem
	par := w.param

	npoints := len(prob.dist)
	start := time.Now()
	for iter := 0; iter < par.maxIter; iter++ {
		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		w.move(i, j, w.state)
	}
	runtime := time.Since(start)
	fmt.Printf("%d moves in time %v\n", par.maxIter, runtime)
	return runtime
}

func (w tspWalker) timeDelta() time.Duration {

	prob := w.problem
	par := w.param

	npoints := len(prob.dist)
	start := time.Now()
	for iter := 0; iter < par.maxIter; iter++ {
		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		w.delta(i, j, w.state, prob.dist)
	}
	runtime := time.Since(start)
	fmt.Printf("%d delta comps in time %v\n", par.maxIter, runtime)
	return runtime
}

func (w tspWalker) timeEnergy() time.Duration {

	prob := w.problem
	par := w.param

	npoints := len(prob.dist)
	start := time.Now()
	for iter := 0; iter < par.maxIter; iter++ {
		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		w.move(i, j, w.state)
		travelDist(w.state, prob.dist)
	}
	runtime := time.Since(start)
	fmt.Printf("%d move+energy comps in time %v\n", par.maxIter, runtime)
	return runtime
}
