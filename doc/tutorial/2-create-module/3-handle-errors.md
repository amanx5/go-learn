# [Return and handle an error](https://go.dev/doc/tutorial/handle-errors)

## Handle empty name

1. Edit [greetings.go](./greetings/greetings.go)

    1. Import `errors` package:
        ```go 
        import "errors"
        ```

    2. Change the return type:
        ```go
        func Greet(name string) (string, error)
        ```

        Any Go function can return [multiple values](https://go.dev/doc/effective_go.html#multiple-returns).

    3. Handle empty name case:
        ```go
        if (name == "") {
            return "", errors.New("Name shouldn't be empty")
        }
        ```

        Return empty greeting and a [new error](https://pkg.go.dev/errors/#example-New) with a message.

    4. Return `nil` (meaning no error) in valid name case:
        ```go
        // ...
        // existing code
        // ...
        return message, nil 
        ```

2. Edit [hello-people.go](./hello-people/hello-people.go)
    1. Import `log` package 
    2. Set the properties of standard logger
        ```go
        log.SetPrefix("greetings: ") // adds a prefix in log entries
        log.SetFlags(0) // disables printing the time, source file, and line number
        ```
    3. Call `Greet` with empty string and assign variables:
        ```go
        message, err := greetings.Greet("")
        ```

    4. Handle empty name case
        ```go
        if err != nil {
            log.Fatal(err)
        }
        ```
3. Run the program