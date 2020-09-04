// want package:`wrappers\(T: NewT, User: NewUser\)`
package b

import "a"

var _ = a.T{} // want "T should be initialized in NewT."

// User is ...
// initWrapper: NewUser()
type User struct {
	username string
	age      int
}

func NewUser(username string, age int) User {
	user := User{username, age}
	return user
}

var admin = &User{"admin", -1} // want "User should be initialized in NewUser."
