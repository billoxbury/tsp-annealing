/*

Build with make.

Run with:

./bin/search -h

./bin/search -poly 10 -per 10 -pr
// ~0.5ms
./bin/search -poly 100 -per 50 -pr
// ~5ms
./bin/search -poly 1000 -niters 10000000
// ~1.4s

// 10,000-gon

// Best regime found:
./bin/search -poly 10000 -temp 1.0 -niters 1000000000 -v -sched sigmage
// Found distance 6.312087850677634 in time 3m21.880283693s

// I haven't been able to improve on this (i.e. cooling based on sigmage over bins)
// with standard schedule.

// GB 79 cities
./bin/search -dat ./data/gb_cities.csv -pr
Rscript ./R/drawRoute.R ./data/gb_cities.csv ./data/route.txt ./img/map.pdf

// Eire
// the Eire data set is much more challenging -  claimed optimal value = 206,171:
// https://www.math.uwaterloo.ca/tsp/world/eilog.html

./bin/search -dat ./data/eire.csv -v -niters 1000000000 -temp 10.0 -cool 0.9999 -per 100000
// Found distance 219027.2245719097 in time 5m7.335534685s

// with sigmage schedule

./bin/search -dat ./data/eire.csv  -temp 32.0 -cool 0.92 -niters 1000000000 -v -sched sigmage
// Found distance 229400.87286360632 in time 1m49.208398757s

// Initial temp 1000.0 suggested by landscape portrait, but this performs less well.

(NOTE that the move class _swap_ is dramatically worse then _reverse_ on all problems.)

*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	// variables
	var dataFile, outFile string
	var moveclass, schedule string
	var temp, cooling float64
	var period, countdown int
	var poly, nwalkers, niters int
	var npoints int = 0
	var verbose, pr bool

	// cmd line arguments
	flag.StringVar(&dataFile, "dat", "", "cities file (CSV)")
	flag.StringVar(&outFile, "out", "route.txt", "output file")
	flag.IntVar(&poly, "poly", 0, "polygon size (option)")
	flag.IntVar(&nwalkers, "nw", 1, "nr walkers")
	flag.IntVar(&period, "per", int(1e04), "period at each temparature")
	flag.IntVar(&countdown, "cd", 400, "countdown for acceptance condition")
	flag.IntVar(&niters, "niters", int(1e06), "max iterations for search")
	flag.Float64Var(&temp, "temp", 4.0, "initial temperature")
	flag.Float64Var(&cooling, "cool", 0.9, "cooling factor")
	flag.StringVar(&moveclass, "mc", "reverse", "move class (default: 2-bond chain reversal)")
	flag.StringVar(&schedule, "sched", "std", "cooling schedule (default: constant rate)")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&pr, "pr", false, "print route")
	flag.Parse()

	// initialise TSP problem
	var prob tspProblem
	if dataFile != "" {
		prob = readCsv(dataFile)
		npoints = len(prob.points)
	} else if poly > 0 {
		prob = makePolygon(poly)
		npoints = poly
	}
	if npoints == 0 {
		fmt.Println("No problem to process")
		return
	}

	// initialise Metropolis parameters
	par := annealParam{
		temperature: temp,
		cooling:     cooling,
		period:      period,
		maxIter:     niters,
		countdown:   countdown,
		schedule:    schedule}

	// set move class
	var move func(int, int, []int)
	var delta func(int, int, []int, [][]float64) float64
	switch moveclass {
	case "swap":
		move = swap
		delta = swapDelta
	default:
		move = reverse
		delta = reverseDelta
	}

	// channel for walkers to report on
	results := make(chan packet, nwalkers)

	// run walkers
	var wg sync.WaitGroup
	for i := 0; i < nwalkers; i++ {

		wg.Add(1)
		w := tspWalker{
			problem: prob,
			param:   par,
			state:   rand.Perm(npoints),
			move:    move,
			delta:   delta,
			verbose: verbose}

		go func() {
			defer wg.Done()
			w.search(results)
		}()
	}

	// collect results
	best_s := make([]int, npoints)
	best_e := float64(1 << 32)

	for i := 0; i < nwalkers; i++ {
		res := <-results
		if res.best_e < best_e {
			best_e = res.best_e
			copy(best_s, res.best_s)
		}
	}
	wg.Wait()

	// report results
	if pr {
		printRoute(best_s, prob.labels)
	}
	writePerm(best_s, "./data/"+outFile)
	fmt.Printf("Best distance found: %v\n", best_e)
	fmt.Printf("Best route written to %s\n", "./data/"+outFile)
}
