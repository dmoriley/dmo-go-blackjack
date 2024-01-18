.PHONY: all run test testgame testdeck

all: run

run:
	go run main.go

# run all tests in project
test:
	go test ./...

testgame:
	go test ./game/...

testdeck:
	go test ./decks/...
