package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Colors []string

/*
// Every slice in Go at runtime is represented in memory by the 3-word struct header:

type slice struct {
    array unsafe.Pointer // 4 or 8 bytes (pointer to the first element in the backing array)
    len   int            // 4 or 8 bytes (Current length: number of elements visible in the slice)
    cap   int            // 4 or 8 bytes (Capacity: total elements the slice can reach before reallocating)
}

// Total: 24 bytes on 64-bit system (or 12 bytes on a 32-bit system).
*/

func main() {
	var c Colors

	// Strictly speaking, Go has NO pass-by-reference.
	// Everything in Go is pass-by-value (passed by copy).
	// When you pass a pointer (&c), you are passing
	// the pointer's value (an 8-byte hexadecimal memory address in x64 system) by copy.

	fmt.Println("=== BEFORE APPEND ===")
	printHeader()
	passValue(c)
	passValue(c)
	passPtrValue(&c)
	passPtrValue(&c)

	w.Flush()

	c = append(c, "red")

	fmt.Println("\n=== AFTER APPEND ===")
	printHeader()
	passValue(c)
	passValue(c)
	passPtrValue(&c)
	passPtrValue(&c)

	w.Flush()
}

// c is copy of 24-byte slice header
func passValue(c Colors) {
	fmt.Fprintf(w, "passValue(c)\t%v\t%p\t%p\tNew copy of header on stack\n", c, &c, c)
}

// c is copy of pointer (8-byte mem addr) to original argument
func passPtrValue(c *Colors) {
	fmt.Fprintf(w, "passPtrValue(&c)\t%v\t%p\t%p\tOriginal header in main (via pointer)\n", *c, c, *c)
}

var w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

func printHeader() {
	fmt.Fprintln(w, "CALL\tVALUE\tHEADER ADDR\t1ST ELEM ADDR\tEXPLANATION")
	fmt.Fprintln(w, "----\t-----\t-----------\t-------------\t-----------")
}
