## TSP Annealing

The Travelling Salesman Problem (TSP) as a platform for studying the Metropolis-Hastings algorithm (simulated annealing).

In this directory:

    /tsp                library code for tsp set-up and simulated annealing  
    - interface.go          structs defined: tspProblem, annealParam, tspWalker
    - distance.go
    - move.go
    - output.go
    - tspProblem.go
    - tspTests.go
    - tspWalker.go
    /src                source code for experiments - see comments at top of each file
    - explore.go
    - makepolydata.go
    - search.go
    - runtests.go
    /bin                binaries for experiments (one for each file in /src)
    /R                  R scripts
    - drawRoute.R
    - landscape.R
    - assessConvergence.R
    - polyParams.R
    /data               data sets
    /img                images (output by R scripts)
        /movie              make animation as follows:
                            1. put sequence of .png or .pdf files here
                            2. run 
                            > convert -delay 75  -resize '100%' -density 288x288 ./*.pdf map.gif
    Makefile            run make to control building of binaries
    README.md           
    go.mod
    .gitignore