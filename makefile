.PHONY: run build

run:
	go run . blegeafe as -n bla dasd

build:
	mkdir -p build/
	env GOOS=linux go build -o build/q.linux
	env GOOS=windows go build -o build/q.windows
	env GOOS=darwin go build -o build/q.macos
