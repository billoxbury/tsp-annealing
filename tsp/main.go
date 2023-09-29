/*

Run polygon examples with randomised parameters, collect data set of results.

make makepolydata



*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// variables
	var outFile string
	var moveclass, schedule string
	var nruns, niters, minn, maxn int

	// cmd line arguments
	flag.StringVar(&outFile, "o", "data.csv", "output file")
	flag.IntVar(&minn, "min", int(100), "min polygon size")
	flag.IntVar(&maxn, "max", int(5000), "max polygon size")
	flag.IntVar(&nruns, "nrun", int(100), "nr experiments")
	flag.IntVar(&niters, "maxiter", int(1e08), "max iters per experiment")
	flag.StringVar(&moveclass, "mc", "reverse", "move class (default: 2-bond chain reversal)")
	flag.StringVar(&schedule, "sched", "std", "cooling schedule (default: constant rate)")
	flag.Parse()

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

	// run experiments
	for i := 0; i < nruns; i++ {

		// set randomised polygon
		npoints := minn + 100*rand.Intn((maxn-minn)/100)
		prob := makePolygon(npoints)

		// set randomised walker
		par := makeParam(niters, schedule)
		wkr := tspWalker{
			problem: prob,
			param:   par,
			state:   rand.Perm(npoints),
			move:    move,
			delta:   delta}

		start := time.Now()
		E, _ := wkr.run()
		t := time.Since(start).Seconds()

		// report
		fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n",
			npoints,
			E,
			t,
			par.temperature,
			par.cooling,
			par.period,
			par.schedule)
	}
}

func makeParam(niters int, schedule string) annealParam {

	var period int
	temp := 1.0 + 7.0*rand.Float64()
	cooling := 0.9 + rand.Float64()
	period := 100 * rand.Intn(100)

	par := annealParam{
		schedule:    schedule,
		maxIter:     niters,
		temperature: temp,
		cooling:     cooling,
		period:      period,
		countdown:   40}

	return par
}
