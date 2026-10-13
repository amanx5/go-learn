package main

import "fmt"

func main() {
	var a any

	a = "string"
	a = 1234
	a = 1.234

	// type assertion
	// Assigning back to 'a' re-boxes it into 'any', because static types never change in Go:
	a = a.(float64)

	// =========================================================================
	// 1. Single-value form: a2 := a.(float64)
	// - Narrows the type: a2 has static type float64.
	// - PANICS at runtime if the dynamic type does NOT match.
	// =========================================================================
	a2 := a.(float64) // Success: dynamic type is float64
	fmt.Printf("1. Single-value (success): a2 = %v (static type: %T)\n", a2, a2)

	// Demonstrating single-value failure (triggers runtime panic):
	// Uncommenting the line below causes:
	// panic: interface conversion: interface {} is float64, not string
	// _ = a.(string)

	func() {
		// defer schedules a function call to be executed right before the surrounding function returns
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("   Single-value (failure): panicked with -> %v\n\n", r)
			}
		}()
		_ = a.(string) // Mismatched type triggers panic
	}()

	// =========================================================================
	// 2. Comma-ok form: a2, ok := a.(float64)
	// - Safe type assertion: NEVER panics.
	// - If match:   ok is true,  a2 contains the underlying value.
	// - If mismatch: ok is false, a2 contains the zero value of the target type.
	// =========================================================================
	// Successful assertion:
	valFloat, ok := a.(float64)
	fmt.Printf("2. Comma-ok (success): val = %v, ok = %v\n", valFloat, ok)

	// Failed assertion (safe, no panic):
	valStr, ok := a.(string)
	fmt.Printf("   Comma-ok (failure): val = %q (zero-value), ok = %v\n", valStr, ok)
}
