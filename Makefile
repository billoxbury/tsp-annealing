# directories
bin = ./bin
src = ./src
tsp = ./tsp

# build all
all:	runtests onesearch explore makepolydata
	@ echo 'make complete'

# build targets
runtests: $(src)/runtests.go $(tsp)/*.go
	cp $(src)/runtests.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	$(bin)/runtests -n 10
	rm $(tsp)/main.go
onesearch: $(src)/onesearch.go $(tsp)/*.go
	cp $(src)/onesearch.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	rm $(tsp)/main.go
explore: $(src)/explore.go $(tsp)/*.go
	cp $(src)/explore.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	rm $(tsp)/main.go
makepolydata: $(src)/makepolydata.go $(tsp)/*.go
	cp $(src)/makepolydata.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	rm $(tsp)/main.go
