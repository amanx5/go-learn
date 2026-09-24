package main

import (
	"fmt"
	"log"
	
	"example.com/greetings"
)


func main() {
    // Set properties of the standard logger
    log.SetPrefix("greetings: ") // adds a prefix in log entries
    log.SetFlags(0) // disables printing the time, source file, and line number

	        names := []string {"Aman" , "John", "Adel"}

        messages, err := greetings.GreetMany(names)

	if err != nil {
		log.Fatal(err) // prints the error in the console and exits the program
	}


	fmt.Println(messages)
}
