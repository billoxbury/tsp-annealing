package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

// make polygon point set
func makePolygon(n int) tspProblem {

	var prob tspProblem

	for i := 0; i < n; i++ {
		x := math.Cos(2.0 * float64(i) * math.Pi / float64(n))
		y := math.Sin(2.0 * float64(i) * math.Pi / float64(n))
		pt := [2]float64{x, y}
		prob.labels = append(prob.labels, strconv.Itoa(i))
		prob.points = append(prob.points, pt)
	}
	prob.dist = distMatrix(prob.points)
	return prob
}

// read data file into points slice
func readCsv(dataFile string) tspProblem {

	var prob tspProblem
	df, _ := os.Open(dataFile)
	defer df.Close()

	// create scanner (bufio)
	scanner := bufio.NewScanner(df)
	scanner.Scan() // this skips the first line
	for scanner.Scan() {

		record := strings.Split(scanner.Text(), ",")
		prob.labels = append(prob.labels, record[0])
		x, _ := strconv.ParseFloat(record[1], 64)
		y, _ := strconv.ParseFloat(record[2], 64)
		pt := [2]float64{x, y}
		prob.points = append(prob.points, pt)
	}
	prob.dist = distMatrix(prob.points)
	return prob
}
