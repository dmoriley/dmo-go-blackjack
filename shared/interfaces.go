package shared

type Printable interface {
	Inspect() func() string
	Debug() func() string
}
