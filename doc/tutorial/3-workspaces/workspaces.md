# [Getting started with multi-module workspaces](https://go.dev/doc/tutorial/workspaces)

## Table of Contents

[Create a module for your code](#create-a-module-for-your-code)  
[Create the workspace](#create-the-workspace)  
[Download and modify the golang.org/x/example/hello module](#download-and-modify-the-golangorgxexamplehello-module)  
[Run the code](#run-the-code)  
[Typical next steps](#typical-next-steps)  
[Takeaway](#takeaway)

This tutorial introduces the basics of creating workspaces in Go. With Workspaces, you can work on multiple modules at same time and easily build and run code in those modules.

In this tutorial, we will create two modules in a shared workspace, make changes across those modules, and see the results of those changes in a build.

## Create a module for your code

1. Navigate to your desired path where you want to create a workspace directory.
2. Create a workspace directory at your desired path. In this tutorial, we have created it [here](./workspace) directory.

    ```bash
    mkdir workspace
    cd workspace
    ```

3. Initialize the module
    - Create a new module `example.com/hello-reverse` in the `hello-reverse` directory:
        ```bash
        mkdir hello-reverse
        cd hello-reverse
        go mod init example.com/hello-reverse
        ```

    <!-- NOTE: In the official docs, this step is ordered after "Add a dependency on the ..." step. We create hello-reverse.go first to avoid a minor confusion.
    If `go get <packageName>` is run before any code depends on the package, Go may add the dependency as an indirect dependency in go.mod. After hello-reverse.go is created and imports the package, the dependency is actually direct, so tools such as Go extension in VS Code may report that the `// indirect` annotation should be removed. Running `go mod tidy` will fix the dependency declaration. -->
    - Create `hello-reverse.go` in the `hello-reverse` directory:

        ```go
        package main

        import ("fmt"; "golang.org/x/example/hello/reverse")

        func main () {
            fmt.Println(reverse.String("Hello"))
        }
        ```

    - Add a dependency on the `golang.org/x/example/hello/reverse` package by using [go get](https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them):

        ```bash
        go get golang.org/x/example/hello/reverse
        ```

    - Run the program:

        ```bash
        go run .        # runs successfully
        ```

        This command executes the program succesfully since it is run in same directory where `go.mod` is present, so go is aware of the main module and runs the main package.

    - Move back to workspace directory and try running the program:
        ```bash
        cd ..       # move to workspace directory
        go run ./hello-reverse      # go: cannot find main module, ...
        ```
        The command fails this time, since it is run from a directory where there is no `go.mod` or a `go.work` file currently.
        So, even though "./hello-reverse" contains go files and a main package, architecturally go can't run the main package since go is not aware of the main module (which is needed for resolving the dependencies).  
         In next step, we will configure the workspace such that go is able to resolve the main module even outside `hello-reverse` directory.

## Create the workspace

We will create a special file called [`go.work`](https://go.dev/ref/mod#go-work-file).

### Initialize the workspace

In the workspace directory, initialize a go workspace using [go work init](https://pkg.go.dev/cmd/go#hdr-Initialize_workspace_file) command.  
Note: You can pass multiple module paths in the `go work init` command to add multiple modules in the workspace. For now, we are only passing one module path:

```bash
go work init ./hello-reverse
```

- This command creates a `go.work` file for a workspace containing the module rooted at the `hello-reverse` directory.
- The file would look like below:

    ```
    go 1.27.1

    use ./hello-reverse
    ```

    - The `use` directive tells Go that the module in the `hello-reverse` directory should be ["main modules" ](https://go.dev/ref/mod#modules-overview) when doing a build.
    - That means, the `example.com/hello-reverse` module will be treated as a main module for Go commands in any subdirectory of the [workspace](./workspace/).  
      For example: If you are in `workspace` or `workspace/some-sub-directory/`, and run a Go command, Go sees the `go.work` file and uses the `example.com/hello-reverse` module as a main module. Open the terminal in workspace directory and check:

        ```bash
        # currently in workspace directory
        go list -m  # example.com/hello-reverse
        go run ./hello-reverse # works now
        ```

## Download and modify the golang.org/x/example/hello module

In this step, we will clone the git repo containing the `golang.org/x/example/hello` module into the workspace directory, and then
add a new function to it which we will later use in the main function in the `hello-reverse.go` program

1. Clone the repo in the workspace directory:

    ```bash
    git clone https://go.googlesource.com/example
    ```

    You will see a new directory named "example" in the workspace directory. It contains multiple modules, but we will use the "hello" module only in this tutorial.

    <!-- NOTE: In this repository, we removed example/.git (and other unused subdirectories in example/) so that the modified hello module and int.go can be tracked directly in Git without embedded repository / submodule conflicts. -->

2. Add the "hello" module from cloned repo to the workspace:

    ```bash
    go work use ./example/hello
    ```

    The `go work use` command adds a new module to the `go.work` file.

    The workspace now includes two copies of the same module `golang.org/x/example/hello`:
    - First copy downloaded by `go get` command is present in the go module cache.  
      The module cache is generally outside the project directory and is managed automatically by Go as dependencies are needed.
      The path of module cache can be checked using following command:
        ```bash
        go env GOMODCACHE
        ```
    - Second copy we just created by cloning is present inside the workspace directory
      [./workspace/example/hello](./workspace/example/hello/).  
      We created this copy so that we can add some additional functionality to the `golang.org/x/example/hello/reverse` package.

3. add a new function to reverse a number:
    - Create a new file `int.go` in the [./workspace/example/hello/reverse/](./workspace/example/hello/reverse/) package directory:

        ```go
        package reverse

        import "strconv"

        // Int returns the decimal reversal of the integer i.
        func Int(i int) int {
            i, _ = strconv.Atoi(String(strconv.Itoa(i)))
            return i
        }
        ```

4. Modify the [`hello-reverse.go`](./workspace/hello-reverse/hello-reverse.go) file:

    ```diff
    -    fmt.Println(reverse.String("Hello"))
    +    fmt.Println(reverse.String("Hello"), reverse.Int(24601))
    ```

## Run the code

From the workspace directory, run:

```bash
go run ./hello-reverse
```

- Not only we are running the program outside its module directory,
  but also overriding its dependencies using workspace.
- `./hello-reverse/hello-reverse.go` uses `reverse.Int` function and
  remember the official `golang.org/x/example/hello/reverse` package
  downloaded in the module cache by `go get` doesn't ship any `Int` function.
- But program doesn't fail since we cloned and modified the
  `golang.org/x/example/hello/reverse` package and then added it into
  the workspace enclosing the `example.com/hello-reverse` module.
- The `use ./example/hello` directive present in the `go.work` file
  instructs go to override the one in the module cache.

Also, if we run the program from the module directory, the workspace still overrides its dependencies:

```bash
cd hello-reverse
go run .
```

This is because the `go` command traverses the parent directories from the directory where the command is run to search for a `go.work` file (which activates workspace mode):
    - It doesn't stop searching if it finds a `go.mod` (module boundary) file in the command directory.
    - It continues searching upward for a `go.work` file to see if the current module is part of a multi-module workspace.
    - If it finds `go.work` in a parent directory, it activates workspace mode.
    - It stops searching when it reaches the root of the filesystem (or `GOROOT`).

To temporarily turn off this behavior:
We can run commands with the environment variable `GOWORK=off`.

```bash
# In Bash / macOS / Linux:
GOWORK=off go run .

# In PowerShell:
$env:GOWORK="off"; go run .

# In Windows Command Prompt:
set GOWORK=off && go run .
```

Try running the program again with workspace turned off:

```bash
GOWORK=off go run .
```

Output:

```text
# example.com/hello-reverse
.\hello-reverse.go:11:22: undefined: reverse.Int
```

You see the command failed now since Go falls back to the version in the module cache which doesn't contain the `reverse.Int` function.

## Typical next steps

Typically when you add new changes in existing published packages, you do so to add new features/fix bugs and then republish them as new versions of their modules.

This is usually done by pushing the commit containing the change to the module's version control repository and creating a new tag with the updated version.

For example: Just saying, for this tutorial, the commit would need to be done in [this repo](https://cs.opensource.google/go/x/example) (don't do it, its just a tutorial) with tag like "v0.1.0".

Once the release is done, we could remove local clone of the module if needed, and just increase the dependency version for `golang.org/x/example/hello` module to its new version in our [hello-reverse/go.mod](./workspace/hello-reverse/go.mod) module by using following command:

```bash
cd hello-reverse
go get golang.org/x/example/hello@v0.1.0
```

[Learn more about module release workflow](https://go.dev/doc/modules/release-workflow).

## Takeaway

- With workspaces:
    - It is easier to make changes in one module and use it in another.
    - You don't need `replace` directives to work across multiple modules.
    
- The `go` command has couple of subcommands for working with workspaces:

    - [`go work init`](https://pkg.go.dev/cmd/go#hdr-Initialize_workspace_file)

    - [`go work use [-r] [moddirs]`](https://pkg.go.dev/cmd/go#hdr-Add_modules_to_workspace_file)
        <!-- NOTE: official docs contain wording issues:
            "removes the use directory" should be "removes the use directive".
            The command docs also contain an extra "go.work file" in
            "each argument directory listed on the command line go.work file".
            -->
        - For each directory provided in `moddirs`, the command adds a `use` directive
        to the `go.work` file if the directory exists; otherwise, it removes the
        corresponding `use` directive from the file.
        - The `-r` flag recursively searches the argument directories for modules.

    - [`go work edit`](https://pkg.go.dev/cmd/go#hdr-Edit_go_work_from_tools_or_scripts)
        - Edits the go.work file similarly to go mod edit
    - [`go work sync`](https://pkg.go.dev/cmd/go#hdr-Sync_workspace_build_list_to_modules)
        - Syncs dependencies from the workspace’s build list into each of the workspace modules.

Find more details about workspaces [here](https://go.dev/ref/mod#workspaces).