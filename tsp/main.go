/*

Code to run standard test on two move classes, using variable-sized polygon problems

Build with make, run with

./bin/runtests -h
./bin/runtests -n 100

etc

*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {

	var n int
	flag.IntVar(&n, "n", 100, "nr points to test on")
	flag.Parse()

	prob := makePolygon(n)
	par := annealParam{maxIter: 1e06}

	w_swap := tspWalker{
		problem: prob,
		param:   par,
		state:   rand.Perm(n),
		move:    swap,
		delta:   swapDelta}
	w_rev := tspWalker{
		problem: prob,
		param:   par,
		state:   rand.Perm(n),
		move:    reverse,
		delta:   reverseDelta}

	fmt.Printf("Testing for problem on %d points\n", n)
	fmt.Println("Swapping moves:")
	w_swap.testDelta(1e-10)
	w_swap.timeMove()
	w_swap.timeDelta()
	w_swap.timeEnergy()
	fmt.Println("Reversing moves:")
	w_rev.testDelta(1e-10)
	w_rev.timeMove()
	w_rev.timeDelta()
	w_rev.timeEnergy()
}
