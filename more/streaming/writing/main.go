package main

import (
	"fmt"
	"io"
)

type Terminal struct{}

// Write method writes the bytes to the terminal
func (t Terminal) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func main() {
	var w io.Writer = Terminal{}
	w.Write([]byte("Hello"))
}
