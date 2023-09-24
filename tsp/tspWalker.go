package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Metropolis algorithm - returns best energy and best state
func (w tspWalker) run() (float64, []int) {

	prob := w.problem
	par := w.param
	npoints := len(prob.points)

	ct := par.maxIter

	// track progress
	acceptance := 0
	this_d := travelDist(w.state, prob.dist)
	best_d := travelDist(w.state, prob.dist)
	lastBest := 2 * best_d
	// to track the best permutation, make a new slice and copy perm into it:
	best_p := make([]int, npoints)
	copy(best_p, w.state)

	start := time.Now()
	for iter := 0; iter < par.maxIter; iter++ {

		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		delta_d := w.delta(i, j, w.state, prob.dist)
		if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
			// accept move
			w.move(i, j, w.state)
			acceptance += 1
			this_d += delta_d
			if this_d < best_d {
				best_d = this_d
				copy(best_p, w.state)
			}
		}

		// report progress
		if iter%par.period == 0 {
			if w.verbose {
				fmt.Printf("%6d: temperature %v, acceptance %v best dist %v\n",
					iter,
					par.temperature,
					float64(acceptance)/float64(par.period),
					best_d)
			}
			// check countdown
			if best_d == lastBest {
				ct++
			} else {
				lastBest = best_d
				ct = 0
			}
			if ct >= par.countdown {
				break
			}
			// otherwise proceed to cooler temperature
			par.temperature *= par.cooling
			acceptance = 0
		}
	}
	runtime := time.Since(start)
	distance := travelDist(best_p, prob.dist)
	fmt.Printf("Found distance %v in time %v\n", distance, runtime)

	return distance, best_p
}
