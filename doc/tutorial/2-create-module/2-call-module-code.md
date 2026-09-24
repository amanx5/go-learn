# [Call your code from another module](https://go.dev/doc/tutorial/call-module-code)

Now this is the time to use the `Greet` function (refer [create a module](./create-module.md)) 
from another module.


1. Create a `hello-people` directory as a sibling to  `greetings` directory. This will be the directory of consumer module.
    ```bash
    cd ..
    mkdir hello-people
    cd hello-people
    ```

2. Enable dependency tracking in the consumer module.
    ```go 
    go mod init example.com/hello-people
    ```

3. Create `hello-people.go` file 

4. Add code: 
    1. Make it part of `main` package:
        ```go
        package main
        ```

    2. Import the `greetings` package in the main package. You can import multiple packages in same import statement like below:
        ```go
        import ("fmt"; "example.com/greetings")
        ```
        or using block import (also called a factored import):
        ```go
        import (
            "fmt"
            "example.com/greetings"
        )
        ```

        Both are valid statements. 
        
        Note: If you're on VSC and Go extension installed, you should see [BrokenImport](https://pkg.go.dev/golang.org/x/tools/internal/typesinternal#BrokenImport) error on hovering the `example.com/greetings` package import since we haven't published the module that provides this package.

    3. Write code for the main function 
        ```go
        func main() {
            message := greetings.Greet("Go")
            fmt.Println(message)
        }
        ```

5. To use the local module `example.com/greetings` without publishing, we need to adapt the `example.com/hello-people` module so it can find the `example.com/greetings` code on local file system. 
    
    1. Run the following command in the `example.com/hello-people` module directory to achieve that:

        ```go
        go mod edit -replace example.com/greetings=../greetings
        ```

        This commands adds a [replace directive](https://go.dev/doc/modules/gomod-ref#replace).

        ```
        replace example.com/greetings => ../greetings
        ```
    
    2. Run tidy command to update the dependencies of the `example.com/hello-people` module
        ```go
        go mod tidy
        ```

        The tidy command will find out that the `example.com/hello-people` module depends on `example.com/greetings` and check that whether it is being tracked. 
        
        Since this was the first time, we ran the tidy command, the `go.mod` file wasn't tracking any dependency before this command. 
        
        So tidy command added a [require directive](https://go.dev/doc/modules/gomod-ref#require) to specify that `example.com/hello-people` module depends on `example.com/greetings` module since it imports the `greetings` package.

        ```
        require example.com/greetings v0.0.0-00010101000000-000000000000
        ```

        The number following the module path is a pseudo-version number -- a generated number used in place of a semantic version number (which the module doesn't have yet).

        To reference a published module, a `go.mod` file would typically omit the replace directive and use a require directive with a tagged [version number](https://go.dev/doc/modules/version-numbers) at the end.
        
        ```
        require example.com/greetings v1.1.0
        ```

6. Run the program
