# [Compile and install the application](https://go.dev/doc/tutorial/compile-install)

The `go run` command compiles and runs the source code but doesn't generate a binary executable. It is designed exclusively for rapid local development and prototyping.

For production deployments, the standard approach is to compile your application into a standalone binary using `go build` and execute that binary directly.

This tutorial introduces two additional commands for building code:

- The [go build](https://go.dev/cmd/go/#hdr-Compile_packages_and_dependencies) command compiles the packages, along with their dependencies, but it doesn't install the results.

- The [go install](https://go.dev/ref/mod#go-install) command compiles and installs the packages.

## Steps
1. Run build command in the [hello-people](./hello-people/) directory:

    ```bash
    go build
    ```

2. Run the generated executable `hello-people.exe`:

    On Bash / Git Bash / Powershell:

    ```bash
    ./hello-people.exe
    ```

    On Windows Command Prompt:

    ```cmd
    hello-people.exe
    ```

    Note: Right now it will only run if you are in the same directory where the `hello-people.exe` file exists or if you have provided full path of the executable.
    Next, we will install the executable so that it can be run from any directory.

3. First, check the Go install path (where the go command will install the current package) using the [`go list`](https://go.dev/cmd/go/#hdr-List_packages_or_modules) command with `-f` option:

    ```bash
    go list -f '{{.Target}}'
    ```

    - In windows, it might output `C:\Users\<username>\go\bin\hello-people.exe`. This indicates that `C:\Users\meisa\go\bin` is the install directory. This can also be checked using:
        ```bash
        go env GOBIN
        ```

4. Depending on your preference, you might want to do either of the following:
    - If you prefer the current install directory, you just need to ensure that the install directory path is present in the system's shell path ($PATH variable). Check the $PATH variable:

        ```bash
        echo $PATH
        ```

    - If you prefer some other install directory, you can change the install directory:

        ```bash
        go env -w GOBIN=C:\path\to\your\bin
        ```

    - In case, your desired install directory is not present in $PATH variable, set it:  
       On Unix-like systems:

        ```bash
        export PATH="$PATH:/path/to/your/install/directory"
        ```

        On Windows command prompt:

        ```cmd
        set PATH=%PATH%;C:\path\to\your\install\directory
        ```

        On Powershell:

        ```powershell
        $env:Path += ";C:\path\to\your\install\directory"
        ```

5. Now run the install command to compile and install the package:

    ```bash
    go install
    ```

6. Run the app easily from any directory just by name of the executable:
    ```bash
    hello-world
    ```

## Further Reading

[Managing dependencies
](https://go.dev/doc/modules/managing-dependencies)

[Developing and publishing modules
](https://go.dev/doc/modules/developing)
