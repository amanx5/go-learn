package main

import (
	"encoding/json/v2"
	"fmt"
)

// Reference Types in Go (slices, maps, and pointers):
// When unmarshaling into a struct containing reference types, if the field
// is uninitialized (nil), json.Unmarshal automatically allocates the underlying
// data structure (e.g., a backing array for a slice) behind the scenes.

type Book struct {
	Id       int
	Name     string
	Keywords []string // Slice is a reference type; its zero-value is nil
}

const book = `{
	"Id": 1,
	"Name": "Cosmos",
	"Keywords": ["science"]
}`

func main() {
	var b Book

	// In Go:
	// - &b.Keywords (%p): Address of the slice header { Data unsafe.Pointer, Len int, Cap int } inside struct 'b'.
	// - b.Keywords  (%p): Address of the first element in the backing array on the heap.
	fmt.Printf("%-20s %-30s %-20s\n", "Keywords", "Keywords header addr", "Keywords addr (1st elem)")

	var details = func() {
		fmt.Printf("%-20v %-30p %-20p\n", fmt.Sprint(b.Keywords), &b.Keywords, b.Keywords)
	}

	// 1. Before Unmarshal:
	// 'var b Book' allocated the struct, but b.Keywords is currently a nil slice.
	// Its header has a null data pointer (0x0), meaning no backing array exists yet.
	details() // Keywords addr = 0x0 (nil)

	// 2. Unmarshal:
	// Unmarshal detects that b.Keywords is nil.
	// To populate the field, Unmarshal allocates a new backing array on the heap
	// behind the scenes and updates b.Keywords's slice header (Data pointer, len, and cap).
	json.Unmarshal([]byte(book), &b)

	// 3. After Unmarshal:
	// - "Keywords header addr" remains the same: the field slot in 'b' never moved.
	// - "Keywords addr" is now a valid heap address: proves Unmarshal allocated
	//   a new backing array for the slice behind the scenes.
	details() // Keywords addr = 0x... (newly allocated backing array)
}
