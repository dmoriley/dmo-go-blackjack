package utils

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

func TestBlankUserInput(t *testing.T) {
	expected := ""
	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)

	config := NewInputConfig(scanner)
	actual, err := GetUserInput(config)
	if err == nil {
		t.Fatalf("Expected error but got none. Actual value received is: |%s|", actual)
	}

	expected = " " // blank space added
	buf = strings.NewReader(expected)
	scanner = bufio.NewScanner(buf)

	config = NewInputConfig(scanner)
	actual, err = GetUserInput(config)
	if err == nil {
		t.Fatalf("Expected error but got none. Actual value received is: |%s|", actual)
	}

}

func TestUnrestrictedTextInput(t *testing.T) {
	expected := "testing"
	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)

	config := NewInputConfig(scanner)
	actual, err := GetUserInput(config)

	if err != nil {
		t.Fatalf("Got error %q. Expected none.", err.Error())
	}

	if actual != expected {
		t.Fatalf("Result wrong. Expected %q, but got %q", expected, actual)
	}
}

func TestTrimmingTextInput(t *testing.T) {
	expected := "testing"
	buf := strings.NewReader(" " + expected + " ")
	scanner := bufio.NewScanner(buf)

	config := NewInputConfig(scanner)
	actual, err := GetUserInput(config)

	if err != nil {
		t.Fatalf("Got error %q. Expected none.", err.Error())
	}

	if actual != expected {
		t.Fatalf("Result wrong. Expected %q, but got %q", expected, actual)
	}
}

func TestRestrictingTextInput(t *testing.T) {
	inputRestrictions := []string{
		"a", "b", "c",
	}
	expected := "a"

	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)
	config := NewInputConfig(scanner).SetExpectedValues(inputRestrictions...)
	actual, err := GetUserInput(config)

	if err != nil {
		t.Fatalf("Got error %q. Expected none.", err.Error())
	}

	if actual != expected {
		t.Fatalf("Result wrong. Expected %q, but got %q", expected, actual)
	}

	// now test that using value not in input restrictions throws and error

	expected = "d"

	buf = strings.NewReader(expected)
	scanner = bufio.NewScanner(buf)
	config = NewInputConfig(scanner).SetExpectedValues(inputRestrictions...)
	actual, err = GetUserInput(config)

	if err == nil {
		t.Fatalf("Expected error but got none.")
	}
}

func TestUnrestrictedIntegerInput(t *testing.T) {
	expected := "5"
	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)

	config := NewInputConfig(scanner)
	actual, err := GetUserInputInteger(config)

	if err != nil {
		t.Fatalf("Got error %q. Expected none.", err.Error())
	}

	if conv, _ := strconv.Atoi(expected); actual != conv {
		t.Fatalf("Result wrong. Expected %q, but got %q", expected, actual)
	}
}

func TestSupplyingNonIntegerForIntergerInput(t *testing.T) {
	expected := "abcd"
	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)

	config := NewInputConfig(scanner)
	actual, err := GetUserInputInteger(config)

	if err == nil {
		t.Fatalf("Expected error but got none. Actual value received is: %q", actual)
	}

}

func TestRestrictedIntegerInput(t *testing.T) {
	inputRestrictions := []string{
		"1", "2", "3",
	}
	expected := "1"

	buf := strings.NewReader(expected)
	scanner := bufio.NewScanner(buf)
	config := NewInputConfig(scanner).SetExpectedValues(inputRestrictions...)
	actual, err := GetUserInput(config)

	if err != nil {
		t.Fatalf("Got error %q. Expected none.", err.Error())
	}

	if actual != expected {
		t.Fatalf("Result wrong. Expected %q, but got %q", expected, actual)
	}

	// now test that using value not in input restrictions throws and error

	expected = "4"

	buf = strings.NewReader(expected)
	scanner = bufio.NewScanner(buf)
	config = NewInputConfig(scanner).SetExpectedValues(inputRestrictions...)
	actual, err = GetUserInput(config)

	if err == nil {
		t.Fatalf("Expected error but got none.")
	}
}
