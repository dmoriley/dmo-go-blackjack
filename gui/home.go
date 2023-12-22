package gui

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = "|= "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input := scanner.Text()
		io.WriteString(out, input+"\n")
	}
}
