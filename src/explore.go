/*

Explore landscape and return energies seen at a constant temperature.

make explore

TEMP=10.0
./bin/explore -f ./data/eire.csv  -o ./data/eire_$TEMP.csv -temp $TEMP -niters 500000

TEMP=0.1
./bin/explore -f ./data/gb_cities.csv  -o ./data/gb_$TEMP.csv -temp $TEMP -niters 500000


*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {

	// variables
	var dataFile, outFile string
	var moveclass string
	var temp float64
	var poly, niters int
	var period int
	var npoints int = 0

	// cmd line arguments
	flag.StringVar(&dataFile, "f", "", "cities file (CSV)")
	flag.StringVar(&outFile, "o", "data.csv", "output file")
	flag.IntVar(&poly, "poly", 0, "polygon size (option)")
	flag.IntVar(&niters, "niters", int(1e06), "nr iterations for search")
	flag.IntVar(&period, "per", int(1e04), "period between checks")
	flag.Float64Var(&temp, "temp", 1.0, "constant temperature")
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

	// initialise walker
	w := tspWalker{
		problem: prob,
		param:   par,
		state:   rand.Perm(npoints),
		move:    move,
		delta:   delta}

	// run explore
	w.explore(outFile)

	// report
	fmt.Printf("")
}

func (w tspWalker) explore(filename string) {

	prob := w.problem
	par := w.param
	npoints := len(prob.points)
	energy := travelDist(w.state, prob.dist)

	file, _ := os.Create(filename)
	defer file.Close()
	wrt := bufio.NewWriter(file)

	ct, accept := 0, 0
	mean := 0.0
	sd2 := 0.0
	for iter := 1; iter < par.maxIter; iter++ {

		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		delta_d := w.delta(i, j, w.state, prob.dist)
		if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
			// accept move
			w.move(i, j, w.state)
			energy += delta_d
			accept++
		}
		// update stats
		mean += energy
		sd2 += energy * energy

		fmt.Fprintf(wrt, "%v\n", energy)
		ct++

		// check statistics
		if iter%par.period == 0 {
			mean /= float64(par.period)
			sd2 /= float64(par.period)
			sd2 -= mean * mean
			fmt.Printf("Bin (%d): %v +/- %v\n", par.period, mean, math.Sqrt(sd2))
			mean = 0.0
			sd2 = 0.0
		}

	}
	wrt.Flush()
	fmt.Printf("Acceptance rate %g\n", float64(accept)/float64(ct))
	fmt.Printf("Written %d energy values to %s\n", ct, filename)
}
