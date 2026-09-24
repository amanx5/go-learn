package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

// Greet returns a greeting for the named person.
// It checks if the name is empty and returns an error if so,
// otherwise it combines a random greeting with the provided name.
func Greet(name string) (string, error) {
	// If no name was provided, return an error.
	if name == "" {
		return "", errors.New("Name shouldn't be empty")
	}

	// Get a random greeting message and format it with the name.
	greeting := randomGreeting()
	message := fmt.Sprintf("%v %v", greeting, name)
	return message, nil
}

// GreetMany takes a slice of names and returns a map pairing each name
// with its own randomized greeting message or an error if a name is empty.
func GreetMany(names []string) (map[string]string, error) {
	// Create a messages map to associate each of the received names (as a key) 
	// with a generated message (as a value). In Go, you initialize a map with 
	// the following syntax: make(map[key-type]value-type).
	messages := make(map[string]string)

	// Loop through the received slice of names, calling Greet for each.
	for _, name := range names {
		message, err := Greet(name)

		// Return any error encountered immediately.
		if err != nil {
			return nil, err
		}

		// Save the message in the map.
		messages[name] = message
	}

	return messages, nil
}

// randomGreeting returns one of a set of predefined greeting messages.
// It uses pseudo-random number generation to select a greeting from the list.
func randomGreeting() string {

	greetings := []string{
		"Hello",
		"Hola",    // Spanish
		"Bonjour", // German
		"Ciao",    // Italian
		"你好",      // Mandarin Chinese: Nǐ hǎo
		"नमस्ते",  // Hindi: Namaste
		"こんにちは",   // Japanese: Kon'nichiwa
		"Olá",     // Portuguese
	}

	// Generate a random index to select a greeting.
	index := rand.Intn(len(greetings))

	return greetings[index]
}
