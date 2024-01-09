package utils

import (
	"bufio"
	"strconv"
	"strings"
)

func GetUserInput(scanner *bufio.Scanner) (string, error) {

	scanned := scanner.Scan()
	if !scanned {
		return "", scanner.Err()
	}

	return strings.TrimSpace(scanner.Text()), nil
}

// Parse user intput for a number
func GetUserInputInteger(scanner *bufio.Scanner) (int, error) {
	text, err := GetUserInput(scanner)
	if err != nil {
		return 0, err
	}

	integer, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return integer, nil
}
