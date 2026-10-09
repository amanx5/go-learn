// this file is only compiled if the program is NOT running on windows
//go:build !windows

package main

import "fmt"

// no-op
func init() {
	// add a breakpoint manually on next line to test if this file is compiled
	fmt.Println("sample_other.go: inside init")
}
