# Go Environment Variables on Windows

## Overview

| Env Variable    | Purpose                                                                                                                                   | Value set by installer |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | ---------------------- |
| `GOROOT`        | Specifies the location of the Go installation (toolchain, standard library, tests)                        | -                         |
| `GOPATH` (user) | Specifies the location for Go user data (downloaded module cache in `pkg/mod`, installed binaries in `bin`) | `%USERPROFILE%\go`        |
| `GOBIN`         | Specifies the location where Go should store the user installed executables.                              | -                         |
| `PATH` (system) | For running `go.exe` without specifying its full path.                                                    | `C:\Program Files\Go\bin` |
| `PATH` (user)   | For running user installed executables without specifying their full path.                                | `%USERPROFILE%\go\bin`    |

---

### Resolved values

To check the values of `GOROOT`, `GOPATH`, `GOBIN`, run:

```shell
go env GOROOT GOPATH GOBIN
```

### GOROOT

`GOROOT` points to the Go installation itself.

It contains the Go toolchain:

```text
C:\Program Files\Go\
                  └── bin\
                       ├── go.exe
                       └── gofmt.exe
```

You normally don't need to set `GOROOT` manually. Go detects the installation folder automatically based on the directory layout relative to the `go.exe` executable itself.

Setting `GOROOT` manually is rarely needed (and generally discouraged) unless the Go binary is moved to an arbitrary directory outside the standard toolchain layout. Switching Go versions is natively managed via Go toolchains (`GOTOOLCHAIN` or versioned binaries like `go1.x.y download`), without needing to alter `GOROOT`.

### GOPATH

`GOPATH` is Go's user-level workspace/data location.

By default, installer sets it as `%USERPROFILE%\go` in user variables.

A typical structure is:

```text
C:\Users\<user>\go\
                 ├── bin\
                 └── pkg\
                      └── mod\
```

### GOBIN

`GOBIN` is where Go installs executable commands.

When `GOBIN` is not explicitly set, Go uses `<GOPATH>\bin` as default (running `go env GOBIN` displays this resolved path).

Suppose you run:

```shell
go install golang.org/x/tools/gopls@latest
```

Go compiles and stores the binary in `<GOPATH>\bin\gopls.exe` (or `<GOBIN>\gopls.exe` if `GOBIN` is explicitly set).

### PATH System variable

During installation, the installer appends the Go installation's `bin` directory to the `PATH` system variable.

```text
PATH
  └── C:\Program Files\Go\bin
                            └── go.exe
```

This allows users to run the go executable from any directory without specifying the full path:

```shell
go version
```

### PATH User variable

During installation, the installer appends `%USERPROFILE%\go\bin` (the default `GOPATH\bin` location) to the `PATH` user variable.

```text
PATH
  └── C:\Users\<user>\go\bin
                          └── dlv.exe
                          └── gopls.exe
```

This allows users to run their installed executables from any directory without specifying their full path:

```shell
dlv     # runs the Delve debugger (if installed)
gopls   # runs the official Go LSP server (if installed)
```
