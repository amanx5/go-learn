// _windows suffix in the name of this file ensures that
// this file is only compiled if the program is running on windows
package main

import (
	"fmt"
	"runtime"
	"syscall"
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	isDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
)

func init() {
	fmt.Println("sample_windows.go: inside init")
	// check if program is called with a debugger
	// without this check, program will crash if ran without debugger
	if r, _, _ := isDebuggerPresent.Call(); r != 0 {
		// explicit breakpoint to pause and confirm this file is compiled
		runtime.Breakpoint()
	}
}
