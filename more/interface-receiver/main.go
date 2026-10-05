package main

import "fmt"

type Username string
type User struct {
	name Username
}

// u is a pointer to an User instance
func (u *User) SetNameInOriginal(name Username) {
	u.name = name
	fmt.Println("Name of original user is now", u.name)
}

// u is a COPY of an User instance
func (u User) SetNameInCopy(name Username) {
	u.name = name
	fmt.Println("Name of copied user is now", u.name)
}

func main() {
	u := User{}

	print(u) // {}

	u.SetNameInOriginal("Aman")
	print(u) // {Aman}

	u.SetNameInCopy("John")
	print(u) // {Aman} // name remains same
}

func print(u User) {
	fmt.Println("Original User =", u)
}
