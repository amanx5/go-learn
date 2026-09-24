# [Return a random greeting](https://go.dev/doc/tutorial/random-greeting)

Now lets change the implementation such that the `Greet` function returns the greeting in a random language.

1. Edit [greetings.go](greetings/greetings.go)
    - Import `math/rand`

    - Declare a unexported named function `randomGreeting`. Since its name starts with lowercase letter, it is only accessible to the code of its own package.

        ```go
        func randomGreeting () string {

            greetings:= []string {
                "greeting1",
                "greeting2",
            }

            const index := rand.Intn(len(formats))

            return greetings[index]
        }
        ```

        - `greetings` is a [slice](https://go.dev/blog/slices-intro) that contains greetings in different languages.

    - Call `randomGreeting` function in `Greet` to receive the random greeting. Create the message using the random greeting and the name.

    ```go 
    greeting := randomGreeting()
	message := fmt.Sprintf("%v %v", greeting, name)
    ```

2. Edit [hello-people.go](hello-people/hello-people.go)
    - Call the `Greet` function with your name. For example:
    ```go
    message, err := greetings.Greet("aman")
    ```

3. Run the `hello-people-people.go` multiple times to get randomized greetings.

