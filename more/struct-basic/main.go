package main

import "fmt"

// declaring struct type
type User struct {
	Name string
	Pass string
}

func main() {
	// creating instance without specifying keys 
	user1 := User{"Aman", "1234"}

	// creating instance
	user2 := User{Name: "John", Pass: "4321"}

	fmt.Println(user1, user2)
}
