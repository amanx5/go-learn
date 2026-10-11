# [Tutorial: Getting started with fuzzing](https://go.dev/doc/tutorial/fuzz)

Using fuzzing, you can test your app against random data to detect vulnerabilities in your code or inputs that can cause your app to crash.

Fuzzing helps you detect vulnerabilities like:
- SQL injection
- Buffer overflow
- Denial of service
- Cross-site scripting

Related docs:

[Go Fuzzing Glossary](https://go.dev/doc/security/fuzz/#glossary)

[Go Fuzzing Docs](https://go.dev/doc/security/fuzz/#requirements)

In this tutorial, we will write a fuzz test for a simple function and debug and fix issues in its code.

## Create a project folder for this tutorial

```shell
mkdir hello-fuzz
cd hello-fuzz
go mod init hello_fuzz
```

## Add a function

Create a `main.go` file with the following code:

- Package declaration and imports:

    ```go
    package main

    import "fmt"
    ```

- A function to reverse a string:

    ```go
    func Reverse(s string) string {
        rev := []byte(s)
        sLen := len(s)
        for i, j := 0, sLen-1; i < sLen/2; i, j = i+1, j-1 {
            rev[i], rev[j] = s[j], s[i]
        }
        return string(rev)
    }
    ```

- Main function:

    ```go
    func main() {
        str := "aman"
        rev := Reverse(str)
        revrev := Reverse(rev)

        fmt.Println("str:", str)
        fmt.Println("rev:", rev)
        fmt.Println("revrev:", revrev)
    }
    ```

## Add a unit test

Create a file `reverse_test.go`:

- Package declaration and imports:

    ```go
    package main

    import "testing"
    ```

- Define a type `Testcase` for the testcases:

    ```go
    type Testcase struct {
        in, want string
    }
    ```

- Declare the test function:
    - Signature:
        ```go
        func TestReverse(t *testing.T) {}
        ```
    - Body:
        - Define testcases:
            ```go
            testcases := []Testcase{
                {"Hello, world", "dlrow ,olleH"},
                {" ", " "},
                {"!12345", "54321!"},
            }
            ```
        - Call the function with testcases:
            ```go
            for _, tcase := range testcases {
                in, want := tcase.in, tcase.want
                rev := Reverse(in)
                if rev != want {
                    t.Errorf("Reverse(%q): %q, want %q", in, rev, want)
                }
            }
            ```

- Run the unit test:
    ```shell
    go test
    ```

## Add a fuzz test

The unit test has some limitations. One limitation is that the developer needs to manually add the test case inputs.

With fuzzing, the inputs are generated automatically and may cover edge cases that the developer might miss.

All test-related code (unit tests, benchmarks and fuzz tests) lives inside `*_test.go` files. So we can add a new fuzz test function in the same `reverse_test.go` file. Comment or remove the `TestReverse` function so that the program output stays clean.

Unlike unit tests, in fuzz tests we don't know the exact expected output of the function since we don't know the inputs as well (since inputs are auto-generated).

So, for our fuzz test of the `Reverse` function, we can't do the comparison of `Testcase.in` vs `Testcase.want` as we did in the Unit test section. But there are a few things that we can still verify:

#### Criterion 1

Reversing a string should remain as a valid UTF-8 string.

#### Criterion 2

Reversing a string twice should give us the original string.

#### Implementation

- Declare the [fuzz test](https://go.dev/doc/security/fuzz/#glos-fuzz-test):
    - Signature:

        ```go
        func FuzzReverse(f *testing.F) {}
        ```

        A fuzz test must start with the `Fuzz` prefix. A fuzz test receives a `*testing.F` pointer.

    - Body:
        - Declare some inputs for the [seed corpus](https://go.dev/doc/security/fuzz/#glos-seed-corpus):

            ```go
            inputs := []string{"Hello, world", " ", "!12345"}
            ```

        - Add [inputs/corpus entries](https://go.dev/doc/security/fuzz/#glos-corpus-entry) to the seed corpus:

            ```go
            for _, i := range inputs {
                f.Add(i)
            }
            ```

        - Pass the [fuzz target function](https://go.dev/doc/security/fuzz/#glos-fuzz-target) to the [(\*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz) function:
            ```go
            f.Fuzz(target)
            ```
            We'll define the fuzz target function separately outside the fuzz test function in the next step for explaining it separately. Generally, we pass the fuzz target function as an anonymous function directly to the `(*testing.F).Fuzz` function.

- Declare the [fuzz target function](https://go.dev/doc/security/fuzz/#glos-fuzz-target):
    - Signature:

        ```go
        func target(t *testing.T, input string) {}
        ```

        The fuzz target accepts `*testing.T` as the first parameter, followed by the [fuzzing arguments](https://go.dev/doc/security/fuzz/#glos-fuzzing-arguments) (we have a single argument `input` in this case).

    - Body:
        - Verify [Criterion 1](#criterion-1)

            ```go
            rev := Reverse(input)
            if utf8.ValidString(input) && !utf8.ValidString(rev) {
                t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
            }
            ```

            This requires importing `"unicode/utf8"`.

        - Verify [Criterion 2](#criterion-2)
            ```go
            revrev := Reverse(rev)
            if input != revrev {
                t.Errorf("Before: %q, after: %q", input, revrev)
            }
            ```

#### Run the code

- First run the fuzz test without fuzzing to check if seed inputs pass:

    ```shell
    go test -v
    ```

    Use the `-v` flag for verbosity.

- [Run the fuzz test](https://go.dev/doc/security/fuzz/#running-fuzz-tests) with fuzzing to check if the program fails against generated inputs:

    ```shell
    go test -v -fuzz=Fuzz
    ```

    - The `-fuzz` flag enables fuzzing. It takes a regular expression that selects fuzz tests whose names match it. By giving it the value `Fuzz`, all the fuzz tests that contain `Fuzz` will run. For running a single fuzz test named `FuzzReverse`, use `^FuzzReverse$`.

    - By default, fuzzing runs indefinitely if it doesn't find any errors. To restrict the time fuzzing takes, use the `-fuzztime` flag.
      For example, `-fuzztime 10s` would terminate the test when an error occurs or when 10 seconds have elapsed.

    - Learn more about testing flags [here](https://pkg.go.dev/cmd/go#hdr-Testing_flags).

- Check the output of the fuzz test:

    The command should fail with output like below:

    ```
    === RUN   FuzzReverse
    fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
    fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 16 workers
    fuzz: elapsed: 0s, execs: 2271 (4803/sec), new interesting: 11 (total: 14)
    --- FAIL: FuzzReverse (0.48s)
        --- FAIL: FuzzReverse (0.00s)
            reverse_test.go:41: Reverse produced invalid UTF-8 string "\x94\xca"

        Failing input written to testdata\fuzz\FuzzReverse\a8c8b98939d243c8
        To re-run:
        go test -run=FuzzReverse/a8c8b98939d243c8
    === NAME
    FAIL
    exit status 1
    FAIL    hello_fuzz      1.188s
    ```

    - The command fails and writes the failing test input in the file at path `<project>/testdata/fuzz/FuzzReverse/<contenthash>`.
        - The `<contenthash>` is a truncated SHA-256 hash.
          Go computes the hash of the corpus entry's content and uses the first 16 hexadecimal characters as the filename.
          Since it is content-derived, the same input always produces the exact same file name.
        - The `testdata` folder enables automatic regression testing.
          The next time `go test` (even without the `-fuzz` flag) is run, Go automatically includes the inputs from the `testdata` folder.
    - To check what input failed the fuzz test, open the file at path `<project>/testdata/fuzz/FuzzReverse/<contenthash>`.
      It should look like below:
        ```
        go test fuzz v1
        string("ʔ")
        ```

        - The corpus file has a strict structure:
            - **The Header (Line 1):** Specifies the encoding version.
            - **The Body (Subsequent Lines):** Every line directly maps to the parameters of the fuzz target function in the exact order they appear.  
              The values are written out in raw Go literal syntax.  
              There will be only 1 line in the body section since our fuzz target function has only one fuzz argument.
        - Your corpus file may contain a different string instead of `ʔ`.
    - To run a fuzz test against a specific corpus entry, use the previously generated `contenthash`/filename of the corpus file. For example:
        ```shell
        go test -run=FuzzReverse/a8c8b98939d243c8
        ```
        
## Handle multi-byte characters correctly
 
The current implementation reverses the string byte-by-byte, and that's the reason the fuzz test failed for an input string `"ʔ"` (which may be a different string in your case) that contains a multi-byte character. 

Not all Unicode characters take a single byte to be stored in memory. Multi-byte characters like `ʔ` require more than 1 byte. 

We must not change the order of bytes of a multi-byte character while reversing the string, otherwise the result would become an invalid string. That's why [utf8.ValidString](https://pkg.go.dev/unicode/utf8) returned `false` (as per docs: "ValidString reports whether s consists entirely of valid UTF-8-encoded runes") and ultimately caused the test to fail.

To fix this, we must reverse the string rune-by-rune so that the UTF-8 encoded bytes of multi-byte characters are not mistakenly reordered.

For better understanding, read the blog [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings).

 
- Fix the implementation by traversing the string rune-by-rune:

    ```go
    func Reverse(s string) string {
        r := []rune(s)
        for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
            r[i], r[j] = r[j], r[i]
        }
        return string(r)
    }
    ```

    > NOTE:
    > This implementation is still not completely correct since it doesn't handle [combining characters](https://en.wikipedia.org/wiki/Combining_character) correctly. But we will ignore that issue for the simplicity of this tutorial.

- Run the fuzz test again to check if there are any new issues.
    -   You will get output similar to below:
        ```text
        --- FAIL: FuzzReverse (0.02s)
            --- FAIL: FuzzReverse (0.00s)
                reverse_test.go:33: Before: "\xc6", after: "�"
            Failing input written to testdata/fuzz/FuzzReverse/e504fa449c34beb2
        ```
     
    -   This time the fuzz test failed due to Criterion 2 (the string is different from the result of reversing it twice). 
        This happened because the input string itself is invalid Unicode.

## Handle input strings containing invalid UTF-8 bytes

Let's first understand the problem. 

In Go, [a string is a read-only slice of bytes](https://go.dev/blog/strings), and can contain bytes that aren’t valid UTF-8. 


So, an input like `"\xc6"` is a valid string but that doesn't mean it is a valid UTF-8 string.  

**What is `\x`?**
-   The input `"\xc6"` represents a string containing only a single byte `0xc6` (198 in decimal).
-   The hexadecimal escape sequence `\x` is used to represent a single byte of data using 2 hex digits.
-   Go interprets the escape sequence `\x` followed by two [hex digits](https://go.dev/ref/spec#hex_digit) as a byte and stores it exactly the same in memory (without any encoding; encoding happens when we use `\u` or `\U`). 
-   As a result, the byte `0xc6` (198 in decimal) gets stored in memory. 
-   Learn more at [Rune literals](https://go.dev/ref/spec#Rune_literals) and [String literals](https://go.dev/ref/spec#String_literals).

**Why is `"\xc6"` an invalid UTF-8 string?**
-   The byte `0xc6` (11000110 in binary) is NOT a valid standalone UTF-8 character. 
-   UTF-8 uses [specific encoding schemes / byte patterns](https://lokalise.com/blog/what-is-character-encoding-exploring-unicode-utf-8-ascii-and-more/#byte-patterns-in-utf-8) to represent Unicode code points in memory. 
-   The byte `0xc6`/`11000110` starts with the bits `110`, which have a special meaning in UTF-8. The bits `110` signal the UTF-8 decoder that these bits are part of a leading byte of a 2-byte Unicode code point. 
-   So, the decoder expects the byte `11000110` to be followed by a continuation byte. 
-   But in the case of string `"\xc6"` there is only a single byte `11000110` without any continuation byte which causes the UTF-8 decoder to fail, hence the string is not a valid UTF-8 byte sequence.

**Why does it matter in our code?**
-   Now that we know why the input `"\xc6"` is invalid Unicode, let's understand the impact of it on our current logic of reversing a string. 
-   The logic in the current implementation of `Reverse` converts the input string to a `[]rune` for rune-by-rune reversal. 
-   Converting an invalid UTF-8 string to `[]rune` results in invalid UTF-8 bytes being converted to the [**Replacement Character** `U+FFFD`](https://pkg.go.dev/unicode/utf8#RuneError) (hex `0xFFFD`, decimal `65533`).
-   As a result, the input string is no longer equal to the resultant string produced after reversing it twice.

**Debug the test case:**
-   Add print statements to the `Reverse` function:
    ```diff
        func Reverse(s string) string {
    +       fmt.Printf("input: %q\n", s)
            r := []rune(s)
    +       fmt.Printf("runes: %q\n", r)
            // existing code 
        }
    ```

-   For our ease in debugging, we will run only the failing test this time using the `go test -run` command. For example:
    ```shell
    go test -run=FuzzReverse/e504fa449c34beb2
    ```

    The output would contain:
    ```text
    input: "\xc6"
    runes: ['�']
    input: "�"
    runes: ['�']
    ```

**Fix the issue:**

- Return an error if the input is not valid UTF-8:
    ```diff
    +   import (
    +       "errors"
    +       "unicode/utf8"
    +   )

    -   func Reverse(s string) string {
    +   func Reverse(s string) (string, error) {
    +        if !utf8.ValidString(s) {
    +            return s, errors.New("input is not valid UTF-8")
    +        }
            r := []rune(s)
            for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
                r[i], r[j] = r[j], r[i]
            }
    -        return string(r)
    +        return string(r), nil
        }
    ```
- Update all the call sites of `Reverse` to ignore or consume the returned error:
    -   Update the `main` function in `main.go`:
        ```diff
        -    rev := Reverse(str)
        -    revrev := Reverse(rev)
        +    rev, _ := Reverse(str)
        +    revrev, _ := Reverse(rev)
        ```

        You can log the returned errors if needed.

    - Update `target` in `reverse_test.go`:
        ```diff
        -    rev := Reverse(input)
        +    rev, err1 := Reverse(input)
        +    if err1 != nil {
        +        return
        +    }
             if utf8.ValidString(input) && !utf8.ValidString(rev) {
                 t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
             }
        -    revrev := Reverse(rev)
        +    revrev, err2 := Reverse(rev)
        +    if err2 != nil {
        +        return
        +    }
        ```

        In the above approach, we return from the target in case the `Reverse` function returns an error. This indirectly means that the test input will pass the fuzz test. 
        
        Instead of returning, you can call `t.Skip()` to skip further execution for that fuzz input. Skipped inputs are printed with the "SKIP" marker in the command output (when run with the `-v` flag) like below: 

        ```text
        ...

        --- PASS: FuzzReverse (0.00s)
            --- SKIP: FuzzReverse/e504fa449c34beb2 (0.00s)
            
        ...
        ```

**Run the code**:
-   Run `go test -fuzz=Fuzz`, wait a few seconds to verify that the command is running, and then press `Ctrl+C` to stop fuzzing.
-   If we neither pass the `-fuzztime` flag nor stop the fuzz test manually, then the command will continue to run indefinitely until it finds any input that fails the fuzz test.
-   Run the command again with the `-fuzztime` flag:
    ```shell
    go test -fuzz=Fuzz -fuzztime 30s
    ```
-   Either way, the fuzzing should pass. 