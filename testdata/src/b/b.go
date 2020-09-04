package b

import "a"

var _ = a.T{} // want "NG"
