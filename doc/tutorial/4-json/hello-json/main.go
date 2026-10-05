package main

import (
	"encoding/json/v2"
	"fmt"
	"time"
)

type Message struct {
	Name string
	Body string
	Time time.Time
}

func main() {
	// create Message instance
	msg := Message{"Aman", "Hello", time.Date(2026, 9, 28, 4, 0, 0, 0, time.UTC)}

	// encode Message to bytes
	encoded, err := json.Marshal(msg)

	fmt.Println(err)             // <nil>
	fmt.Println(string(encoded)) // {"Name":"Aman","Body":"Hello","Time":"2026-09-28T04:00:00Z"}

	var msg2 Message
	err2 := json.Unmarshal(encoded, &msg2)

	fmt.Println(err2) // <nil>
	fmt.Println(msg2) // {Aman Hello 2026-09-28 04:00:00 +0000 UTC}

}
