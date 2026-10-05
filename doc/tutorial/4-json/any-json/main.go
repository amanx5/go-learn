package main

import (
	"encoding/json/v2"
	"fmt"
)

const book = `{
	"Id": 1,
	"Name": "Cosmos",
	"Keywords": ["science", "fiction"]
}`

func main() {
	var result any
	json.Unmarshal([]byte(book), &result)

	fmt.Printf("Result = %v\n", result)
	fmt.Printf("Result Type = %T\n", result)

	resultMap := result.(map[string]any)

	printValueAndType("KEY", "VALUE", "TYPE")

	for key, value := range resultMap {

		switch val := value.(type) {
		case float64:
			printValueAndType(key, val, "float64")
		case string:
			printValueAndType(key, val, "string")
		case []any:
			printValueAndType(key, val, "an array")
		}
	}
}

func printValueAndType(k string, v any, t string) {
	fmt.Printf("%-12s %-25s %-10s\n", k, fmt.Sprint(v), t)
}
