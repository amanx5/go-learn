package main

import "fmt"

func main() {
	// The variable a has the static type any (which is an alias for interface{}).
	// The compiler enforces this at compile time (for example, you cannot call len(a) even after initializing it with a string value).
	var a any = ""

	a = 1234              // dynamictype is now int
	a = 1.234             // dynamictype is now float64
	a = "onetwothreefour" // dynamictype is now string

	// Go's v := a.(type) does not change the type of an existing variable.
	// Instead, the compiler creates a separate, new scoped variable for each branch:
	switch v := a.(type) { // or simply switch a.(type) {
	case string:
		// len(a) // invalid argument: a (variable of interface type any) for built-in len
		// static type of a is STILL any here as STATIC TYPES never change.
		// Go declares a new scoped var `v` with the narrowed type:
		// Go implicitly declared: var v string = ...
		// This v was born as a string; its type never changed.
		log(a, v)
	case int:
		// Go implicitly declared: var v int = ...
		// This is a completely different variable that also happens to be named 'v'.
		log(a, v)
	case float64:
		log(a, v)
	case byte, rune:
		// Go implicitly declared: var v any = a
		// when a case lists more than one type, the new variable is NOT narrowed.
		log(a, v)
	default:
		// Go also implicitly declared: var v any = a
		log(a, v)
	}

}

func log(a any, v any) {
	// %T always prints the DYNAMIC TYPE
	// %T verb in fmt.Printf uses reflection at runtime (reflect.TypeOf).
	// Reflection unpacks the interface and inspects the concrete dynamic type of the value contained within it
	// Runtime reflection cannot inspect compile-time variable types because static types exist only for the compiler and are boxed into interface values at runtime.
	fmt.Printf("type of a=%T\n", a)
	fmt.Printf("type of v=%T\n", v)
}
