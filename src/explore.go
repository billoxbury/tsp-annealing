/*

Reads TSP problem, sets up annealing parameters, dispatches parallel walkers as go routines.
Each walker send periodic data packets back to the search client, which writes a diagnostic file,
tracks the best solution found and writes that to a best-route file.

./bin/search is faster but without diagnostic functionality.

To build:

make explore

To run e.g.:

TEMP=1.0
./bin/explore -f ./data/gb_cities.csv  \
	-d ./data/gb_$TEMP.csv \
	-temp $TEMP \
	-per 10000 \
	-nw 10 \
	-cool 0.97 \
	-srate 100 \
	-nj 100 \
	-pr

Rscript ./R/drawRoute.R ./data/gb_cities.csv ./data/route.txt ./img/map.pdf

EXAMPLES:
./bin/explore -poly 10 -per 10 -nj 20 -nw 1 -pr
// 0.06 ms
./bin/explore -poly 100 -nj 100 -per 500 -pr
// 13ms
./bin/explore -poly 1000 -nj 500
// 1.8s

// 10,000-gon
./bin/explore -poly 10000 -temp 1.0 -nw 8 -nj 1000 -v


*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

func main() {

	// variables
	var dataFile, diagFile, routeFile string
	var moveclass string
	var temp, cooling float64
	var poly, numWalkers, numJobs int
	var period, srate int
	var npoints int = 0
	var verbose, pr bool

	// cmd line arguments
	flag.StringVar(&dataFile, "f", "", "cities file (CSV)")
	flag.StringVar(&diagFile, "d", "./data/data.csv", "diagnostics file")
	flag.StringVar(&routeFile, "r", "./data/route.txt", "output route file")
	flag.IntVar(&poly, "poly", 0, "polygon size (option)")
	flag.IntVar(&numWalkers, "nw", 2, "nr walkers")
	flag.IntVar(&numJobs, "nj", 1, "nr jobs per walker")
	flag.IntVar(&period, "per", int(1e04), "period before cooling")
	flag.IntVar(&srate, "srate", 100, "sampling rate")
	flag.Float64Var(&temp, "temp", 1.0, "initial temperature")
	flag.Float64Var(&cooling, "cool", 0.9, "cooling factor")
	flag.StringVar(&moveclass, "mc", "reverse", "move class (default: 2-bond chain reversal)")
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
		period:      period,
		srate:       srate,
		cooling:     cooling,
		temperature: temp}

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
	results := make(chan packet, numWalkers*numJobs)

	// open diagnostics file for writing
	dfile, _ := os.Create(diagFile)
	defer dfile.Close()
	wrt := bufio.NewWriter(dfile)

	// run walkers
	var wg sync.WaitGroup
	best_s := make([]int, npoints)
	var best_e float64
	best_e = float64(1 << 32)

	for i := 0; i < numWalkers; i++ {

		wg.Add(1)
		w := tspWalker{
			id:      i,
			problem: prob,
			param:   par,
			state:   rand.Perm(npoints),
			move:    move,
			delta:   delta,
			verbose: verbose}

		go func() {
			defer wg.Done()
			w.explore(numJobs, results)
		}()
	}
	// collect and report results
	fmt.Fprintf(wrt, "walker,temperature,iteration,energy\n")
	ct := 0
	for i := 0; i < numWalkers*numJobs; i++ {

		res := <-results
		ct += len(res.energy)

		// check for global winner so far
		if res.best_e < best_e {
			best_e = res.best_e
			copy(best_s, res.best_s)
		}

		// write diagnostics
		for iter, e := range res.energy {
			fmt.Fprintf(wrt, "%d,%v,%d,%v\n", res.id, res.temperature, iter, e)
		}
	}
	wg.Wait()
	wrt.Flush()

	// write winning state
	writePerm(best_s, routeFile)
	if pr {
		printRoute(best_s, prob.labels)
	}

	// report
	fmt.Printf("Best distance found: %v\n", best_e)
	fmt.Printf("Best route written to %s\n", routeFile)
	fmt.Printf("Written %d diagnostic records to %s\n", ct, diagFile)
}
