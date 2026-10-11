package main

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
)

func main() {
	// iterate over list of BCP 47 strings passed in command line arguments
	// https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag
	for _, arg := range os.Args[1:] {
		tag, err := language.Parse(arg)

		if err != nil {
			fmt.Printf("%s: error: %v\n", arg, err)
		} else {
			fmt.Printf("%s: tag %s\n", arg, tag)
		}
	}
}
