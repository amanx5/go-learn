# vscode-go

Notes on the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go).

## Configuring a project that requires an unreleased Go version

If your project requires an unreleased Go version such as 1.28 (unreleased as of Sept. 2026) and the path to the unreleased Go executable is not configured in your environment variables, the Go extension may throw the following error upon opening Go files:

```
stderr: go: downloading go1.28.0 (windows/amd64)
go: download go1.28.0 for windows/amd64: toolchain not available
```

### Cause

Starting in Go 1.21, Go includes an automatic toolchain management feature (`GOTOOLCHAIN=auto`). When a `go.mod` file specifies a Go version higher than the installed Go version, Go automatically attempts to download that newer toolchain from the Go proxy. Because the version is unreleased or only built locally, downloading fails.

### Fix

1. **Set the path of the executable in your environment variables (`PATH`).**  
   You might prefer not to do this if you want your default global `go` executable to remain on a stable version. In that case, use option 2.

2. **Set the path to the Go distribution in project settings.**
   ```json
   // .vscode/settings.json
   {
   	"go.goroot": "D:\\go"
   }
   ```
   *(Alternatively, configure `"go.alternateTools": { "go": "D:\\go\\bin\\go.exe" }`)*

## Configuring a project (opened in a Multi-root workspace) that requires an unreleased Go version

If you open a multi-root workspace (a workspace containing multiple independent Go projects) and one of the projects requires an unreleased version, the extension will throw a `"toolchain not available"` error when opening Go files, even if that specific project folder defines `go.goroot` in its `.vscode/settings.json`.

### Cause

In VS Code multi-root workspaces (`.code-workspace`), settings in an individual folder's `.vscode/settings.json` only apply to **Resource-scoped** settings (e.g., editor formatting, tab size).

For security and architectural reasons, `go.goroot` and `go.alternateTools` have **Machine / Window** scope in the Go extension. VS Code ignores them when placed inside an individual folder's `.vscode/settings.json` within a multi-root workspace:

- `gopls` (the Go language server) runs as a single process per VS Code window. It cannot run Go 1.28 for one folder and Go 1.27.1 for another folder simultaneously in the same window.

### Fix

Unlike the `typescript-language-features` extension, the `vscode-go` extension does not run separate toolchains per folder within the same multi-root workspace. To use completely independent toolchains, you must open each folder in a separate VS Code window.

However, because Go toolchains are backward-compatible, you can configure the highest required version at the **workspace** level so that all projects are analyzed by that executable:

For example:

```json
// example.code-workspace
{
	"folders": [
		{ "name": "A", "path": "A" }, // Go 1.28
		{ "name": "B", "path": "B" }, // Go 1.27
		{ "name": "C", "path": "C" }  // Go 1.26
	]
}
```

Go 1.28 (project A) can parse, build, and analyze modules targeting Go 1.27 (project B) and Go 1.26 (project C) without downloading anything. However, Go 1.27 cannot analyze Go 1.28.

Add the Go 1.28 `GOROOT` to your workspace settings file (`example.code-workspace`):

```json
// example.code-workspace
{
	"folders": [
		/* ... */
	],
	"settings": {
		"go.goroot": "D:\\go"
	}
}
```

> [!WARNING]
> With this configuration, the extension will use Go 1.28 for all projects in the workspace, even those targeting an older version. Automated tooling (such as `go mod tidy` or `gopls` automatic updates) may bump directives (e.g., adding `toolchain go1.28` or updating the `go` directive).
> 
> While this will not cause errors inside this workspace, these projects may fail to build if opened separately in an environment where Go 1.28 is not installed.

