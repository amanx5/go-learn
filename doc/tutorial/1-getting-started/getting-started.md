# [Getting started](https://go.dev/doc/tutorial/getting-started.html)


## Package vs Module

### Package

A package consists of all the Go source files in the same directory. Every Go file belongs to a package (declared with `package <name>` at the top of the file, such as `package main`). All files in that directory that share the same package declaration are part of that single package.

### Module

A module is a collection of related Go packages that are released together as a single unit. A module is defined by a go.mod file (usually located at the root of the project directory), which specifies the module path and its dependencies.

### Summary

Files in a directory make up a package.
Packages in a directory tree (managed by a `go.mod` file) make up a module.

## Dependency Tracking

To use packages that are contained in other modules, we need to add them as dependiences in our own module.

A module is defined by a `go.mod` file, it defines a registry of other modules that the current module is using and instructs go to track them.

To create a module, run following command:

```bash
go mod init example/hello-world
```

- The `example/hello-world` in the command should be the module's module path. In real world projects, it is the location of the repository which is hosting that module. For example: `github.com/mymodule`.  
If you want your module to be accessible to other users, you must specify the module path as a location from which Go tools can download your module. Learn more at [Managing dependencies](https://go.dev/doc/modules/managing-dependencies#naming_module).


- The above command creates a `go.mod` file with content similar to the following:
  ```text
  module example/hello-world

  go 1.27.1
  ```

  This `go.mod` file has two directives:
  - The `module` directive declares the module path of the Go module.
    The module path is the module's identity and is used when importing packages belonging to this module.
  - The `go` directive specifies the Go language and toolchain version requirements for the module.
    - It specifies the minimum version of Go required to use the module.
    - It specifies the Go language version whose language semantics apply to code in the module.

## Create a `hello-world.go` file

In the same directory where the `go.mod` file is present, create a `hello-world.go` file with code below:
```go
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
```


## Running Go module

Navigate to the directory where the `go.mod` file is initialized and then run:

```bash
go run .
```

`run` command compiles and runs the named `main` Go package. Typically the package is specified as a list of .go source files from a single directory, but it may also be an import path, file system path, or pattern matching a single known package, as in 'go run .' or 'go run my/cmd'. [Learn more](https://go.dev/cmd/go/#hdr-Compile_and_run_Go_program)

## Using external packages

Visit [packages](https://pkg.go.dev) and search for the package you want to use.

Example: "quote" package

1. Search for a [qoute package](https://pkg.go.dev/search?q=quote)
2. In search results, find "rsc.io/quote" package. Select the version you want to use.
3. Scroll down to [Documentation](https://pkg.go.dev/rsc.io/quote#section-documentation) > [Index](https://pkg.go.dev/rsc.io/quote#pkg-index). You will see the list of functions available in the package.
4. Now, use it in your code:

    ```go
    package main

    import "fmt"
    import "rsc.io/quote" // import the quote package

    func main() {
        fmt.Println(
            quote.Go(), // use the function available in the quote package
        )
    }
    ```
5. Add the quote module as a requirement:

    ```bash
    go mod tidy
    ```

    - it adds the quote module as a dependency
    - generates a `go.sum` file for use in [authenticating the module.](https://go.dev/ref/mod#authenticating)
