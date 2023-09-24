# directories
bin = ./bin
src = ./src
tsp = ./tsp

# build all
all:	runtests onesearch
	@ echo 'make complete'

# build targets
runtests: $(src)/runtests.go $(tsp)/*.go
	mv $(tsp)/main.go .
	cp $(src)/runtests.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	$(bin)/runtests -n 10
	mv ./main.go $(tsp)
onesearch: $(src)/onesearch.go $(tsp)/*.go
	mv $(tsp)/main.go .
	cp $(src)/onesearch.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	mv ./main.go $(tsp)
explore: $(src)/explore.go $(tsp)/*.go
	mv $(tsp)/main.go .
	cp $(src)/explore.go $(tsp)/main.go
	go build -o $(bin)/$@ $(tsp)
	mv ./main.go $(tsp)