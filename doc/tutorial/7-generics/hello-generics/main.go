package main

import "fmt"

func main() {
	mi := map[string]int64{"a": 1, "b": 2}
	mf := map[string]float64{"a": 1.1, "b": 2.2}
	fmt.Println(Sum(mi))
	fmt.Println(Sum(mf))
	fmt.Println(SumNum(mi))
	fmt.Println(SumNum(mf))
}

func Sum[K comparable, V int64 | float64](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

type Number interface {
	int64 | float64
}

func SumNum[K comparable, V Number](m map[K]V) V {
	return Sum(m)
}
