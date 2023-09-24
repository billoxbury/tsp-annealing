/*

Explore landscape and return energies seen at a constant temperature.

go build -o walker ./tsp

./walker -f ./data/gb_cities.csv -temp 0.1 -niters 500000 -o ./data/gb_0.1.csv

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
	var npoints int = 0

	// cmd line arguments
	flag.StringVar(&dataFile, "f", "", "cities file (CSV)")
	flag.StringVar(&outFile, "o", "data.csv", "output file")
	flag.IntVar(&poly, "poly", 0, "polygon size (option)")
	flag.IntVar(&niters, "niters", int(1e06), "nr iterations for search")
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

	ct := 0
	for iter := 0; iter < par.maxIter; iter++ {

		i := rand.Intn(npoints)
		j := rand.Intn(npoints)
		delta_d := w.delta(i, j, w.state, prob.dist)
		if delta_d < 0 || rand.Float64() < math.Exp(-delta_d/par.temperature) {
			// accept move
			w.move(i, j, w.state)
			energy += delta_d
			fmt.Fprintf(wrt, "%v\n", energy)
			ct++
		}
	}
	wrt.Flush()
	fmt.Printf("Written %d energy values to %s", ct, filename)
}
