.PHONY: all run test testgame testdeck build clean

all: run

clean:
	rm -rf bin

build: clean
	go build -o bin/blackjack

run:
	go run main.go

# run all tests in project
test:
	go test ./...

testgame:
	go test ./game/...

testdeck:
	go test ./decks/...
