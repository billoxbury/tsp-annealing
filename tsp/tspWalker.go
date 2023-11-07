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

	is_sigmage := (par.schedule == "sigmage")
	result_ct, bin_ct := 0, 0
	sigmage_wait := 2 // waiting time to cool under sigmage schedule

	// to track progress
	acceptance := 0
	energy := travelDist(w.state, prob.dist)
	best_e := travelDist(w.state, prob.dist)
	lastBest := 2 * best_e
	mean_e, previous_mean := 0.0, 0.0
	sd2_e, previous_sd2 := 0.0, 0.0
	// to track the best permutation, make a new slice and copy perm into it:
	best_p := make([]int, npoints)
	copy(best_p, w.state)

	start := time.Now()
	for iter := 1; iter < par.maxIter; iter++ {

		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		delta_d := w.delta(i, j, w.state, prob.dist)
		if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
			// accept move
			w.move(i, j, w.state)
			acceptance += 1
			energy += delta_d
			if energy < best_e {
				best_e = energy
				copy(best_p, w.state)
			}
		}
		// update stats
		if is_sigmage {
			mean_e += energy
			sd2_e += energy * energy
		}

		// report progress
		if iter%par.period == 0 {
			if w.verbose {
				fmt.Printf("%6d: temperature %v, acceptance %v best dist %v\n",
					iter,
					par.temperature,
					float64(acceptance)/float64(par.period),
					best_e)
			}
			// check countdown
			if best_e == lastBest {
				result_ct++
			} else {
				lastBest = best_e
				result_ct = 0
			}
			if result_ct >= par.countdown {
				break
			}
			// otherwise proceed to cooler temperature
			if is_sigmage {
				mean_e /= float64(par.period)
				sd2_e /= float64(par.period)
				sd2_e -= mean_e * mean_e
				if (mean_e-previous_mean)*(mean_e-previous_mean) < 2*previous_sd2 {
					bin_ct++
				} else {
					bin_ct = 0
				}
				if bin_ct >= sigmage_wait {
					par.temperature *= par.cooling
				}
				previous_mean = mean_e
				previous_sd2 = sd2_e
				mean_e = 0.0
				sd2_e = 0.0
			} else {
				par.temperature *= par.cooling
			}
			// reset variables
			acceptance = 0
		}
	}
	runtime := time.Since(start)
	distance := travelDist(best_p, prob.dist)
	fmt.Printf("Found distance %v in time %v\n", distance, runtime)

	return distance, best_p
}

/*************************************************************************************/

// walker 'explore' routine

// each walker will send back packets as follows
type packet struct {
	id          int
	temperature float64
	energy      []float64
}

func (w tspWalker) explore(numJobs int, srate int, results chan<- packet) {

	prob := w.problem
	par := w.param
	npoints := len(prob.points)
	energy := travelDist(w.state, prob.dist)

	// variables to report
	var res packet

	ct := 0
	for job := 0; job < numJobs; job++ {

		var energies []float64
		// burn-in
		for iter := 0; iter < par.period; iter++ {
			i := rand.Intn(npoints)
			j := rand.Intn(npoints)
			delta_d := w.delta(i, j, w.state, prob.dist)
			if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
				// accept move
				w.move(i, j, w.state)
				energy += delta_d
			}
		}

		// record energies
		for iter := 0; iter < par.period; iter++ {
			i := rand.Intn(npoints)
			j := rand.Intn(npoints)
			delta_d := w.delta(i, j, w.state, prob.dist)
			if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
				// accept move
				w.move(i, j, w.state)
				energy += delta_d
			}
			if iter%srate == 0 {
				energies = append(energies, energy)
			}
		}

		// return results
		res.id = w.id
		res.temperature = par.temperature
		res.energy = energies
		ct += len(energies)
		results <- res

		// cool
		par.temperature *= par.cooling

	}
	fmt.Printf("%d: written %d energy values to output channel\n", w.id, ct)
}
