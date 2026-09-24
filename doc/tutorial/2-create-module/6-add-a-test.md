# [Add a test](https://go.dev/doc/tutorial/add-a-test)

Go has built-in testing support for unit testing.

- `testing` package
- `go test` command
- file naming conventions

---

1. Create a file `greetings_test.go` in [greetings](./greetings/) directory.

    Ending a file's name with _test.go tells the `go test` command that this file contains test functions.

2. Edit [greetings_test.go](./greetings/greetings_test.go)

    - Import packages
        ```go
        import ("testing"; "regexp")
        ```
    
    - Add a function that tests the `greetings.Greet` method with a valid name

        - Signature: 

            ```go
            func TestGreetName (t *testing.T) {
                // logic
            }
            ```

            - Testing function names should be in format: **Test*Name*** where ***Name*** indicates about the test target and case.

            - Testing functions have one parameter `t` which is a pointer to the `testing` package's [`testing.T`](https://go.dev/pkg/testing/#T) type. 
            This pointer has methods for reporting and logging the test.

        - Body: 
            ```go
            name := "Aman"
            want := regexp.MustCompile(`\b` + name + `\b`)

            got, err := Greet(name)

            if !want.MatchString(got) || err != nil {
                t.Errorf(`Greet(%q) = (%q, %v), want (%#q, nil)`, name, got, err, want)
            }
            ```

            - `name` variable stores a valid name used for testing.
            - `want` variable stores a `*regexp.Regexp` (pointer to a regex) compiled by the `regexp.MustCompile` method. 
            - `Greet` function is called with `name`, if the result either doesn't match with the `want` regex or contains error, we log the details using [`t.Errorf`](https://go.dev/pkg/testing/#T.Errorf).


    - Add another function that tests the empty value case:
        ```go
        func TestGreetEmpty (t *testing.T) {
            empty := ""

            got, err := Greet(empty)

            if got != "" || err == nil {
                t.Errorf(`Greet(%q) = (%q, %v), want ("", error)`, empty, got, err)
            }
        }

3. Run tests

    - In [greetings](./greetings/) directory, run [`go test`](https://go.dev/cmd/go/#hdr-Test_packages) command: 

        ```bash
        go test
        ```

        - The `go test` command visits each test file (whose filename ends with _test.go) and runs each test function (whose names begin with Test) present in those files.

        - By default, `go test` command only logs details about the failed tests and a final summary. For verbose output (lists all the tests and their results (irrespective of pass/fail)), run:

            ```bash
            go test -v
            ```

    - For checking how the output would look like if the test fails, temporarily change the code of `Greet` function as below:

        ```diff
        - message := fmt.Sprintf("%v %v", greeting, name)
        + message := fmt.Sprintf("%v", greeting)
        ```

        We deliberately removed the `name` from the message, which will cause the `TestGreetName` test to fail. Run the test again with the changed code.
