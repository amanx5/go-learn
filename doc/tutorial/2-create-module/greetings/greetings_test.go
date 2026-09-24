package greetings

import (
	"regexp"
	"testing"
)

// TestGreetName calls greetings.Greet with a name, checking
// for a valid return value.
func TestGreetName(t *testing.T) {
	// valid name chosen for testing
	name := "Aman"

	// Build a regular expression that matches `name` as a complete word.
	// `\b` represents a word boundary, so the name won't match as part of
	// another word.
	// MustCompile is like Compile but panics if the expression cannot be parsed
	// "panic" means the program encounters a situation it cannot continue normally from and stops the current flow of execution (if not recovered).
	// We are using MustCompile instead of Compile here since the regex is known/fixed and an invalid pattern would indicate a programming erro
	want := regexp.MustCompile(`\b` + name + `\b`)

	// call the Greet function with name
	got, err := Greet(name)

	// test if there is any error or msg is not in expected pattern
	if !want.MatchString(got) || err != nil {
		t.Errorf(
			`Greet(%q) = (%q, %v), want (%#q, %v)`,
			name, // %q will print it as a quoted string
			got,
			err,  // %v will print using its default format
			want, // %#q will print it as a Go-syntax quoted string
			nil,
		)
	}
}

func TestGreetEmpty(t *testing.T) {
	empty := ""

	got, err := Greet(empty)

	if got != "" || err == nil {
		t.Errorf(`Greet(%q) = (%q, %v), want ("", error)`, empty, got, err)
	}
}
