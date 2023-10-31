## TSP Annealing

The Travelling Salesman Problem (TSP) as a platform for studying the Metropolis-Hastings algorithm (simulated annealing).

In this directory:

    /tsp                library code for tsp set-up and simulated annealing  
    - distance.go
    - interface.go
    - main.go               (temp file that gets overwritten by make)
    - move.go
    - output.go
    - tspProblem.go
    - tspTests.go
    - tspWalker.go
    /src                source code for experiments
    - explore.go
    - makepolydata.go
    - onesearch.go
    - runtests.go
    /bin                binaries for experiments (one for each file in /src)
    /R                  R scripts
    - drawRoute.R
    - landscape.R
    - polyParams.R
    /data               data sets
    /img                images (output by R scripts)
        /movie              make animation as follows:
                            1. put sequence of .png or .pdf files here
                            2. run 
                            > convert -delay 75  -resize '100%' -density 288x288 ./*.pdf map.gif
    Makefile            run make to control building of Go source code
    README.md           this README
    go.mod
    .gitignore