# [Tutorial: Find and fix vulnerable dependencies with govulncheck](https://go.dev/doc/tutorial/govulncheck)

In this tutorial, you will learn how to find vulnerabilities in your project. You will also get an overview of vulnerability management tools.

Govulncheck is a [low-noise](#low-noise-tools) tool that helps you find and fix vulnerable dependencies in your Go projects.

## Low-Noise vs High-Noise tool

In security and software tools, **"noise"** refers to irrelevant warnings or false alarms for issues that actually do not even exist in your projects. These warnings/alerts overwhelm developers (often called **alert fatigue**).

### Low-noise tools
A **low-noise** tool is designed to report only actionable and genuine issues which actually affect your application.

### High-noise tools

A **high-noise** tool lacks reachability analysis, it flags a vulnerability if a package is present anywhere in your dependency tree — even if your project never invokes it. This results in a high rate of false alarms.

Most package ecosystems rely on CVE descriptions that only specify package names and version ranges (e.g., `"package X < 1.4.0 is vulnerable"`). 

For example, `npm audit` is a purely structural package-level vulnerability scanner. It only checks the package names and exact versions in the `package-lock.json` file against a database like the GitHub Advisory Database.

## How `govulncheck` achieves "low noise":

`govulncheck` uses reachability / call-graph analysis:

### Symbol-level inspection
- Go doesn't rely on just CVE descriptions. Go also has its own [Vulnerability Database](https://vuln.go.dev/) which is manually reviewed and curated by the Go security team. 
- Every report in the database contains affected module, package, fixed version and **exact list of vulnerable symbols** (functions, methods, types, interfaces or variables).
- Using these reports, `govulncheck` gives warnings/alerts when you actually use the vulnerable symbols and not just when your project has a dependency containing vulnerable symbols which you don't even use. 
    
### Call trace verification
When you run `govulncheck`, it performs a 3-tier analysis:

```
[ Module Level ]: Is there a vulnerable module in go.mod?
                        ↓ (Yes)
[ Package Level ]: Is the vulnerable package actually imported?
                        ↓ (Yes)
[ Symbol Level ]: Does a call path exist from an entrypoint to the vulnerable symbol?
```

- **Builds a Call graph**: Using static analysis (via Go's compiler tools and AST), it starts from entry points (like `main` function or test functions). 
- **Follows Call chains**: It traverses each function call and the whole tree if nested including the calls to transitive(indirect) dependencies.
- **Determines Reachability**: It matches the symbols in the call graph against the vulnerable symbols in the database.

That's enough theory, let's get into a practical example:

## Create a sample Go module with a vulnerable dependency:

- Create module:
    ```shell
    mkdir hello-vuln
    cd hello-vuln
    go mod init hello_vuln
    ```

- Create `main.go` file with following code:
    ```go
    package main

    import (
        "fmt"
        "os"

        "golang.org/x/text/language"
    )

    func main() {
        for _, arg := range os.Args[1:] {
            tag, err := language.Parse(arg)
            if err != nil {
                fmt.Printf("%s: error: %v\n", arg, err)
            } else {
                fmt.Printf("%s: tag %s\n", arg, tag)
            }
        }
    }
    ```

    The above program iterates over a list of [BCP 47](https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag) strings passed in command line arguments and prints the parsed tag (or error if parsing failed) for each string.
- Add the `golang.org/x/text@v0.3.5` (deliberately an older version since it contains known vulnerabilities) package as a dependency.
    ```shell
    go get golang.org/x/text@v0.3.5
    ```

## Install and run govulncheck:
- Install `govulncheck`:
    ```shell
    go install golang.org/x/vuln/cmd/govulncheck@latest
    ```

- Run `govulncheck`:
    ```shell
    govulncheck ./...
    ```

    This should output something like below:

    ```
    === Symbol Results ===

    Vulnerability #1: GO-2021-0113
        Out-of-bounds read in golang.org/x/text/language
    More info: https://pkg.go.dev/vuln/GO-2021-0113
    Module: golang.org/x/text
        Found in: golang.org/x/text@v0.3.5
        Fixed in: golang.org/x/text@v0.3.7
        Example traces found:
        #1: main.go:14:29: hello_vuln.main calls language.Parse

    Your code is affected by 1 vulnerability from 1 module.
    This scan also found 2 vulnerabilities in packages you import and 14
    vulnerabilities in modules you require, but your code doesn't appear to call
    these vulnerabilities.
    ```

### Understanding the vulnerability report:

 - As stated in the output, the `main` function in our `main.go` file is calling the `language.Parse` function which is one of the symbols that were reported as vulnerable in the [GO-2021-0113](https://pkg.go.dev/vuln/GO-2021-0113) report.

- The output also gives us a short informational message that there are two more vulnerabilities in the imported package(s) and 14 vulnerabilities in added module(s) but our code doesn't call those vulnerabilities.

### Evaluate the vulnerability:
This section covers the thought process that a developer/team can consider while dealing with vulnerabilities.

- Understand the report:
    Based on the GO-2021-0113 report, `Parse` can result in a panic if an untrusted input is passed to it. 

- Evaluate applicability:
    Let's suppose we want our program to be bullet proof against untrusted inputs, and we are concerned about denial of service, then this vulnerability is clearly applicable in our project. 
    
    Also, our program uses command line arguments as a source of input for the vulnerable symbol `language.Parse` which makes it more vulnerable.

- Decide on an action:
    We have 2 options to mitigate GO-2021-0113:
        - **Option 1:** Upgrade to a version containing the fix. This requires updating the dependency and in some cases might need refactoring based on the new implementation of updated packages.
        - **Option 2: Stop using the vulnerable symbol(s).** This requires code refactoring like replacing all usages of the vulnerable symbol with an alternative or custom solution.

    In this case, since a fix is available in a newer version and the `language.Parse` function is an integral part of the program, we will upgrade the dependency to the "fixed in" version. 

##  Upgrade vulnerable dependencies

As per the GO-2021-0113 report, the reported vulnerabilities exist in versions prior to `v0.3.7`. So let's upgrade `golang.org/x/text` to `v0.3.8`: 

```shell
go get golang.org/x/text@v0.3.8
```
You can also upgrade to latest if needed.

Run `govulncheck` again:

```
=== Symbol Results ===

No vulnerabilities found.

Your code is affected by 0 vulnerabilities.
This scan also found 1 vulnerability in packages you import and 14
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
```

The output states that there are no vulnerabilities that affect our code.

Run below command to check the other vulnerabilities which currently don't affect our code, but are present in the command output:

```shell
govulncheck -show verbose ./...
```

The command will recheck the vulnerabilities and print more details, including the 14 vulnerabilities that currently do not affect our code.

```
Fetching vulnerabilities from the database...

Checking the code against the vulnerabilities...

The package pattern matched the following root package:
  hello_vuln
Govulncheck scanned the following 2 modules and the go1.27.1 standard library:
  hello_vuln
  golang.org/x/text@v0.3.8

=== Symbol Results ===

No vulnerabilities found.

=== Package Results ===

Vulnerability #1: GO-2026-6604
    Root.Mkdir(All) can follow junctions out of the root on Windows in os
  More info: https://pkg.go.dev/vuln/GO-2026-6604
  Standard library
    Found in: os@go1.27.1
    Fixed in: os@go1.27.2
    Platforms: windows

=== Module Results ===

Vulnerability #1: GO-2026-6629
    Panic parsing crafted input in x/text/secure/precis in golang.org/x/text
  More info: https://pkg.go.dev/vuln/GO-2026-6629
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.8
    Fixed in: golang.org/x/text@v0.41.0

...
...
...

Your code is affected by 0 vulnerabilities.
This scan also found 1 vulnerability in packages you import and 14
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
```

> [!NOTE]
> The output is trimmed and does not include details of the remaining 13 vulnerabilities.

Regularly scan your project for vulnerabilities using `govulncheck` and safeguard your codebase by identifying, prioritizing, and addressing them.

## Learn more
- [govulncheck command](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [Vulnerability Management for Go](https://go.dev/blog/vuln)