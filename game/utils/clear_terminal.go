package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// a map of clear functions per os
var clear map[string]func()

func init() {
	clear = make(map[string]func()) // initialize

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func ClearTerminal() {
	// check if os type included in map
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic(fmt.Sprintf("Can't clear terminal window. Your platform (%s) is not supported.", runtime.GOOS))
	}
}
