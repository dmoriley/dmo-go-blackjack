.PHONY: all run test testengine

all: run

run:
	go run main.go

# run all tests in project
test:
	go test ./...

testengine:
	go test ./engine/...
