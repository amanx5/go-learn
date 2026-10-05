package main

import "fmt"

func main() {
	var a any

	a = "string"
	a = 1234
	a = 1.234

	// type assertion
	a = a.(float64)
	a2 := a.(float64)

	// a remains any
	// a2 is narrowed
	fmt.Println(a, a2)
}
