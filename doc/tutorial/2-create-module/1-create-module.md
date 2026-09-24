# [Create a Go module](https://go.dev/doc/tutorial/create-module)

In this tutorial you'll create two modules. The first is a library which is intended to be imported by other libraries or applications. The second is a caller application which will use the first.

## Creating a Go module

1. Open terminal
2. Create a `greetings` directory 
3. Initialize the module 
    ```bash
    go mod init example.com/greetings
    ```
4. Create a file `greetings.go`. 
5. Add code:
    - Make the file part of `greetings` package.
        ```go
        package greetings
        ```
    - Import `fmt` package
        ```go
        import "fmt"
        ```
    - Add `Greet` function that returns the greeting.
    
        ```go
        func Greet(name string) string {
            message := fmt.Sprintf("नमस्ते %v", name)
            return message
        }
        ```
        - The name of the function starts with a capital letter, that makes it an [exported name](https://go.dev/tour/basics/3).
        - The function receives a name parameter of type `string`.
        - The function has the return type `string`.
        - The variable is declared and initialized together in same line using the `:=` operator. It is called as **short variable declaration operator**. The type of the variable is auto detected based on the provided value. 
            ```go
            // this is the equivalent longer version
            var message string
            message = fmt.Sprintf("नमस्ते %v", name);
            ``` 
