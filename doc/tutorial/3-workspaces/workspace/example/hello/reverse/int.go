package reverse

import "strconv"

// Int returns the decimal reversal of the integer i.
func Int(i int) int {
	// convert int to ASCII
	a := strconv.Itoa(i)

	// reverse the input
	// note String() here is a function of `reverse` package, not the string() type convertor
	rev := String(a)

	// convert back to int
	i, _ = strconv.Atoi(rev)

	return i
}
