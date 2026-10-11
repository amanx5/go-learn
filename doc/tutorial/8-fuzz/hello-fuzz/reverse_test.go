package main

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	inputs := []string{"Hello, world", " ", "!12345"}
	for _, i := range inputs {
		f.Add(i)
	}

	// Fuzz runs the fuzz target function for each input
	f.Fuzz(target)
}

func target(t *testing.T, str string) {
	rev, revErr := Reverse(str)
	if revErr != nil {
		t.Skip("[INFO] Reverse(str) returned an error. Skipping test criterias")
	}

	revrev, revrevErr := Reverse(rev)
	if revrevErr != nil {
		t.Skip("[INFO] Reverse(Reverse(str)) returned an error. Skipping test criterias")
	}

	debugLog := fmt.Sprintf("[DEBUG] Number of runes: str=%d, rev=%d, revrev=%d", utf8.RuneCountInString(str), utf8.RuneCountInString(rev), utf8.RuneCountInString(revrev))

	if utf8.ValidString(str) && !utf8.ValidString(rev) {
		t.Log(debugLog)
		t.Errorf("[FAILED] [CRITERIA 1]: rev=%q is an invalid UTF-8 string", rev)
	}
	if str != revrev {
		t.Log(debugLog)
		t.Errorf("[FAILED] [CRITERIA 2]: str=%q, revrev=%q", str, revrev)
	}
}

// type Testcase struct {
// 	in, want string
// }
// func TestReverse(t *testing.T) {
// 	testcases := []Testcase{
// 		{"Hello, world", "dlrow ,olleH"},
// 		{" ", " "},
// 		{"!12345", "54321!"},
// 	}

// 	for _, tcase := range testcases {
// 		in, want := tcase.in, tcase.want
// 		rev := Reverse(in)
// 		if rev != want {
// 			t.Errorf("Reverse(%q): %q, want %q", in, rev, want)
// 		}
// 	}
// }
