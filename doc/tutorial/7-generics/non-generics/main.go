package main

import "fmt"

func main() {
	mi := make(map[string]int64)
	mi["a"] = 1
	mi["b"] = 2
	fmt.Println(SumInts(mi))

	mf := map[string]float64{
		"a": 1.1,
		"b": 2.2,
		"c": 3.3,
	}

	fmt.Println(SumFloats(mf))
}

func SumInts(m map[string]int64) int64 {
	var sum int64
	for _, v := range m {
		sum += v
	}
	return sum
}

func SumFloats(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}
