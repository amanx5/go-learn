// Declare a main package (a package is a way to group functions,
// and it's made up of all the files in the same directory)
package main

// Import the popular fmt package, which contains functions for 
// formatting text, including printing to the console. 
// This package is one of the standard library packages you got 
// when you installed Go.
import "fmt"

// Implement a main function to print a message to the console. 
// A main function executes by default when you run the main package.
func main () {
	fmt.Println("Hello World")
}
