package main

import "fmt"

func main() {	
	fmt.Println(reverse(123))
}

func reverse(num int) int {
	res := 0
	copy := num

	for copy != 0 {
		qou, rem := divide(copy, 10)
		res = res*10 + rem
		copy = qou
	}

	return res
}

func divide(dividend int, divisor int) (int, int) {
	qou := dividend / divisor
	rem := dividend % divisor
	return qou, rem
}
