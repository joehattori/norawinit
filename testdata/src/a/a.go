package a

var a = T{}            // want "T should be initialized in NewT."
var b = []int{1, 2, 3} // OK

func f() T {
	if false {
		return T{} // want "T should be initialized in NewT."
	}
	t := &T{} // want "T should be initialized in NewT."
	return *t
}

type (
	s struct{}

	// T is struct
	// initWrapper: NewT()
	T struct{}
)

func NewT() *T {
	return &T{}
}
