package main

import (
	"blackjack/game"
	"fmt"
	"os"
	"os/user"
)

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(`
  ____  _            _     _            _    
 | __ )| | __ _  ___| | __(_) __ _  ___| | __
 |  _ \| |/ _' |/ __| |/ /| |/ _' |/ __| |/ /
 | |_) | | (_| | (__|   < | | (_| | (__|   < 
 |____/|_|\__,_|\___|_|\_\/ |\__,_|\___|_|\_\
                        |__/ 
	`)
	fmt.Printf("Welcome, don't be caught counting cards %s...\n", user.Username)
	game.Start(os.Stdin, os.Stdout)
}
