/*

Run polygon examples with randomised parameters, collect data set of results.

make makepolydata

./bin/makepolydata -nrun 990 -o data/polydata_sigmage.csv -sched sigmage -maxiter 1000000000
./bin/makepolydata -nrun 1000 -o data/polydata_std.csv -sched std -maxiter 1000000000


*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	// variables
	var outFile string
	var moveclass, schedule string
	var nruns, niters, minn, maxn int

	// cmd line arguments
	flag.StringVar(&outFile, "o", "data/polydata.csv", "output file")
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

	// open output file
	file, _ := os.Create(outFile)
	defer file.Close()
	wrt := bufio.NewWriter(file)
	fmt.Fprintf(wrt, "npoints,energy,time,temperature,cooling,period,schedule\n")

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
		E, _ := wkr.search() // <---- CORRECTION NEEDED: EDIT FOR USE AS GO ROUTINE
		t := time.Since(start).Seconds()

		// report
		fmt.Fprintf(wrt, "%v,%v,%v,%v,%v,%v,%v\n",
			npoints,
			E,
			t,
			par.temperature,
			par.cooling,
			par.period,
			par.schedule)
		wrt.Flush()
	}
}

func makeParam(niters int, schedule string) annealParam {

	temp := 1.0 + 3.0*rand.Float64()
	cooling := 0.8 + 0.2*rand.Float64()
	period := 20 * (50 + rand.Intn(949))

	par := annealParam{
		schedule:    schedule,
		maxIter:     niters,
		temperature: temp,
		cooling:     cooling,
		period:      period,
		countdown:   40}

	return par
}
