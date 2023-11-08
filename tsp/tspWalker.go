package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
Metropolis search - returns best energy and best state.

- pluggable cooling schedule
- stopping criterion by repetition countdown for best energy
- fast er than explore() with no data collection

TO DO:

- run parallel walkers as go routines
*/
func (w tspWalker) search() (float64, []int) {

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
			// accept proposal
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
/*
Explore routine:

- flat cooling schedule only
- specified number of constant-temperature periods
- burn-in before data collection in each period
- data collection and piping to client
- parallel walkers via go routines

*/
func (w tspWalker) explore(numJobs int, results chan<- packet) {

	// set-up
	prob := w.problem
	par := w.param
	npoints := len(prob.points)

	// for sigmage schedule
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

	// to track the best state, make a new slice and copy initial state into it:
	best_s := make([]int, npoints)
	copy(best_s, w.state)

	// data packet for reporting
	var res packet
	res.best_s = make([]int, npoints)

	ct := 0
	start := time.Now()

	// MAIN LOOP
	for job := 0; job < numJobs; job++ {

		var energies []float64
		// initial burn-in
		for iter := 0; iter < par.burnin; iter++ {

			{ // MOVE BLOCK
				i := rand.Intn(npoints)
				j := rand.Intn(npoints)
				delta_d := w.delta(i, j, w.state, prob.dist)
				if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
					// accept proposal
					w.move(i, j, w.state)
					energy += delta_d
				}
				// update best found
				if energy < best_e {
					best_e = energy
					copy(best_s, w.state)
				}
				// update stats
				if is_sigmage {
					mean_e += energy
					sd2_e += energy * energy
				}
			} // END OF MOVE BLOCK
		}

		// after burn-in record energies
		for iter := 0; iter < par.period; iter++ {

			{ // MOVE BLOCK
				i := rand.Intn(npoints)
				j := rand.Intn(npoints)
				delta_d := w.delta(i, j, w.state, prob.dist)
				if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
					// accept proposal
					w.move(i, j, w.state)
					energy += delta_d
				}
				// update best found
				if energy < best_e {
					best_e = energy
					copy(best_s, w.state)
				}
				// update stats
				if is_sigmage {
					mean_e += energy
					sd2_e += energy * energy
				}
			} // END OF MOVE BLOCK

			// sample energy
			if iter%par.srate == 0 {
				energies = append(energies, energy)
			}
		}

		// END OF THIS TEMPERATURE - pipe data packet back to client
		res.id = w.id
		res.temperature = par.temperature
		res.energy = energies
		res.best_e = best_e
		copy(res.best_s, best_s)
		ct += len(energies)
		results <- res

		// verbose output
		if w.verbose {
			fmt.Printf("%2d - %6d: temperature %v, acceptance %v best dist %v\n",
				w.id,
				job,
				par.temperature,
				float64(acceptance)/float64(par.period),
				best_e)
		}
		/********************************/
		// check countdown NEEDS EDITING
		if best_e == lastBest {
			result_ct++
		} else {
			lastBest = best_e
			result_ct = 0
		}
		//if result_ct >= par.countdown {
		//break
		//}
		/********************************/

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

	} // END OF MAIN LOOP
	runtime := time.Since(start)

	// report
	distance := travelDist(best_s, prob.dist)
	//fmt.Printf("%d: written %d records to output channel\n", w.id, ct)
	fmt.Printf("%d: found distance %v in time %v\n", w.id, distance, runtime)
}
