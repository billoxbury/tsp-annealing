/*

Explore landscape and return energies seen at a constant temperature.

make explore


TEMP=1.0
./bin/explore -f ./data/gb_cities.csv  -o ./data/gb_$TEMP.csv -temp $TEMP -niters 500000


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
	var dataFile, outFile string
	var moveclass string
	var temp, cooling float64
	var poly, niters, numWalkers, numJobs int
	var period, srate int
	var npoints int = 0

	// cmd line arguments
	flag.StringVar(&dataFile, "f", "", "cities file (CSV)")
	flag.StringVar(&outFile, "o", "data.csv", "output file")
	flag.IntVar(&poly, "poly", 0, "polygon size (option)")
	flag.IntVar(&niters, "niters", int(1e06), "nr iterations for search")
	flag.IntVar(&numWalkers, "nw", 2, "nr walkers")
	flag.IntVar(&numJobs, "nj", 1, "nr jobs per walker")
	flag.IntVar(&period, "per", int(1e04), "period before cooling")
	flag.IntVar(&srate, "srate", 100, "sampling rate")
	flag.Float64Var(&temp, "temp", 1.0, "initial temperature")
	flag.Float64Var(&cooling, "cool", 0.9, "cooling factor")
	flag.StringVar(&moveclass, "mc", "reverse", "move class (default: 2-bond chain reversal)")
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
		cooling:     cooling,
		temperature: temp,
		maxIter:     niters}

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

	// open file for writing results
	file, _ := os.Create(outFile)
	defer file.Close()
	wrt := bufio.NewWriter(file)

	// run walkers
	var wg sync.WaitGroup
	for i := 0; i < numWalkers; i++ {
		wg.Add(1)
		w := tspWalker{
			id:      i,
			problem: prob,
			param:   par,
			state:   rand.Perm(npoints),
			move:    move,
			delta:   delta}
		go func() {
			defer wg.Done()
			w.explore(numJobs, srate, results)
		}()
	}
	// collect and report results
	fmt.Fprintf(wrt, "walker,temperature,energy\n")
	ct := 0
	for i := 0; i < numWalkers*numJobs; i++ {

		res := <-results
		ct += len(res.energy)

		for _, e := range res.energy {
			fmt.Fprintf(wrt, "%d,%v,%v\n", res.id, res.temperature, e)
		}
	}
	wg.Wait()
	wrt.Flush()

	// report
	fmt.Printf("Written %d energy values to %s\n", ct, outFile)
}
