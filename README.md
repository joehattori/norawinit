# norawinit

`norawinit` is a tool to restrict `type` initialization to a specific function.

## Usage

There are some occations when you want to perform some specific tasks, for example validation, on initializing a struct.
`norawinit` can prevent the initialization by composite literal of such struct.
For example, when you want to make sure to run age validation of struct `User` on each initialization, you can write `onlyWrapper: NewUser()`
in the comment above the definition of struct `User` as described in the example below.
With this comment, `norawinit` responds with waring saying `"User should be initialized in NewUser."`. 

```go
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
```

## How to run
```
$ go get -u github.com/joehattori/norawinit
$ go vet -vettool=`which norawinit` YOUR_PACKAGE
```
