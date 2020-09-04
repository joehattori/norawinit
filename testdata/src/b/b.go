// want package:`wrappers\(T: NewT\)`
package b

import "a"

var _ = a.T{} // want "T should be initialized in NewT."
