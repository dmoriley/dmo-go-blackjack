package engine

const (
	One         Denomination = "One"
	Five                     = "Five"
	TwentyFive               = "Twenty Five"
	Fifty                    = "Fifty"
	Hundred                  = "One Hundred"
	FiveHundred              = "Five Hundred"
	Thousand                 = "One Thousand"
)

type Denomination string

// Single betting chip
type Chip struct {
	// Name of the chip
	Denomination Denomination
	// Associated value of the chip
	Value int
}
