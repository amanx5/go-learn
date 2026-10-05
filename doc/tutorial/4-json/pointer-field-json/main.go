package main

import (
	"encoding/json/v2"
	"fmt"
)

// Reference Types in Go (slices, maps, and pointers):
// When unmarshaling into a struct containing reference types, if the field
// is uninitialized (nil), json.Unmarshal automatically allocates the underlying
// data structure (e.g., a backing array for a slice) behind the scenes.

type Foo struct {
	Bar *Bar
}

type Bar struct {
	Size int
}

const fooWithoutBar = `{}`

const fooWithEmptyBar = `{ "Bar": {} }`

const fooWithBar = `{ "Bar": { "Size": 5 } }`

func main() {
	row("Operation", "foo", "foo.Bar", "&foo.Bar")

	var foo Foo
	var log = func(operation string) {
		row(
			operation,
			fmt.Sprintf("%+v", foo),                  // value of foo
			fmt.Sprintf("%v (%p)", foo.Bar, foo.Bar), // value of Bar field
			fmt.Sprintf("%v", &foo.Bar),              // address of Bar field
		)
	}

	log("Declare Foo")

	json.Unmarshal([]byte(fooWithoutBar), &foo)
	log("Unmarshal " + fooWithoutBar)

	json.Unmarshal([]byte(fooWithEmptyBar), &foo)
	log("Unmarshal " + fooWithEmptyBar)

	json.Unmarshal([]byte(fooWithBar), &foo)
	log("Unmarshal " + fooWithBar)
}

func row(a string, b string, c string, d string) {
	fmt.Printf("%-40s %-25s %-30s %-20s\n", a, b, c, d)
}
