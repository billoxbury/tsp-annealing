package main

import (
	"bufio"
	"fmt"
	"os"
)

// show a given route
func printRoute(perm []int, labels []string) {

	for _, v := range perm {
		fmt.Printf("%v --> ", labels[v])
	}
	fmt.Printf("%v\n", labels[perm[0]])
}

// output permutation to file
func writePerm(perm []int, fileName string) {

	file, _ := os.Create(fileName)
	defer file.Close()
	wrt := bufio.NewWriter(file)
	fmt.Fprintf(wrt, "route\n")
	for _, j := range perm {
		fmt.Fprintf(wrt, "%d\n", j)
	}
	wrt.Flush()
}
