package utils

import (
	"bufio"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type InputConfig struct {
	scanner        *bufio.Scanner
	expectedValues []string
	anyKey         bool
}

func NewInputConfig(scanner *bufio.Scanner) *InputConfig {
	return &InputConfig{
		scanner:        scanner,
		anyKey:         false,
		expectedValues: []string{},
	}
}

func (ic *InputConfig) SetExpectedValues(values ...string) *InputConfig {
	ic.expectedValues = values
	return ic
}

// Set any key to continue
func (ic *InputConfig) SetAnyKey(anyKey bool) *InputConfig {
	ic.anyKey = anyKey
	return ic
}

func GetUserInput(config *InputConfig) (string, error) {

	scanned := config.scanner.Scan()
	if !scanned && config.scanner.Err() == nil {
		if config.anyKey {
			// user hit enter as the "any key"
			return "", nil
		}
		// returns nil Err when EOF, or nothing provided
		return "", fmt.Errorf("No input provided")
	} else if !scanned {
		// error for some other reason
		return "", config.scanner.Err()
	}
	trimmed := strings.TrimSpace(config.scanner.Text())

	if config.anyKey {
		return "", nil
	}

	if len(trimmed) == 0 {
		return "", fmt.Errorf("No input provided")
	}

	if len(config.expectedValues) > 0 {
		// expected values provided, check if input is included
		if !slices.Contains(config.expectedValues, trimmed) {
			return "", fmt.Errorf("Unexpected value entered: %q", trimmed)
		}
	}
	return trimmed, nil
}

func EnterToContinue(scanner *bufio.Scanner) {
	fmt.Print("\nPress ENTER to continue...")
	inputConfig := NewInputConfig(scanner).SetAnyKey(true)
	GetUserInput(inputConfig)
}

// Parse user intput for a number
func GetUserInputInteger(config *InputConfig) (int, error) {
	text, err := GetUserInput(config)
	if err != nil {
		return 0, err
	}

	integer, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return integer, nil
}
