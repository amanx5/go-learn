package main

import "fmt"

type User struct {
	Name string
	Age  int
}

var general_verbs = []string{"%v", "%+v", "%#v", "%T"}

func main() {
	user := User{"john", 10}

	row("Verb", "Output")

	for _, verb := range general_verbs {
		// fmt.Printf("%"+verb+" = "+verb+"\n", user)
		row(verb, fmt.Sprintf(verb, user))
	}
}

func row(a string, b string) {
	fmt.Printf("%-20v %-20v\n", a, b)
}
