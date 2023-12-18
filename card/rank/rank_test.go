package rank

import (
	"testing"
)

func TestNewRank(t *testing.T) {

	// test a valid rank
	testRank, error := NewRank("Two", 2)
	if error != nil {
		t.Fatalf("Rank fail. Found error when expected none")
	}

	if testRank.Name != "Two" {
		t.Fatalf("Rank fail. Expected = %q, want = %q", testRank.Name, testRank.Name)
	}

	if testRank.Value != 2 {
		t.Fatalf("Rank fail. Expected = %d, want = %d", 2, testRank.Value)
	}

	testRank, error = NewRank(King, 13)

	if error != nil {
		t.Fatalf("Rank fail. Found error when expected none")
	}

	if testRank.Name != "King" {
		t.Fatalf("Rank fail. Expected = %q, want = %q", testRank.Name, testRank.Name)
	}

	if testRank.Value != 13 {
		t.Fatalf("Rank fail. Expected = %d, want = %d", 13, testRank.Value)
	}

	testRank, error = NewRank("test", 1)

	if error == nil {
		t.Fatalf("Rank fail. Expected error but got none")
	}

	if error.Error() != `"test" is not a valid card rank` {
		t.Fatalf(`Rank fail. Error message wrong. Got = "%s"`, error.Error())
	}
}
