package main

// specification of the TSP problem
type tspProblem struct {
	points [][2]float64
	labels []string
	dist   [][]float64
}

// parameters for explore/search
type annealParam struct {
	schedule    string
	temperature float64
	cooling     float64
	period      int
	burnin      int
	srate       int
	maxIter     int
	countdown   int
}

// structure of a walker
type tspWalker struct {
	id      int
	problem tspProblem
	param   annealParam
	state   []int
	move    func(int, int, []int)
	delta   func(int, int, []int, [][]float64) float64
	verbose bool
}

// each walker will send back data packets
type packet struct {
	id          int
	temperature float64
	energy      []float64
	best_e      float64
	best_s      []int
}

// DEPRECATED
//type walker interface {
//	testDelta(float64) int   // nr errors in iterLimit calls
//	timeMove() time.Duration // running time for iterLimit calls
//	timeEnergy() time.Duration
//	timeDelta() time.Duration
//	run() float64
//}
