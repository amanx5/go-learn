# Field: Issue Title
x/tools/gopls: hover displays no tooltip for null character rune literals ('\x00', '\u0000') and null byte escapes in strings ("...\x00...")

---

# Field: gopls version
```text
golang.org/x/tools/gopls v0.23.0
```
*(Also verified on master)*

---

# Field: go env
```text
set GO111MODULE=
set GOAMD64=v1
set GOARCH=amd64
set GOBIN=C:\Users\user\go\bin
set GOCACHE=C:\Users\user\AppData\Local\go-build
set GOENV=C:\Users\user\AppData\Roaming\go\env
set GOEXE=.exe
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOMOD=NUL
set GOMODCACHE=C:\Users\user\go\pkg\mod
set GOOS=windows
set GOPATH=C:\Users\user\go
set GOPROXY=https://proxy.golang.org,direct
set GOROOT=C:\Program Files\Go
set GOSUMDB=sum.golang.org
set GOTOOLCHAIN=auto
set GOVERSION=go1.27.1
```

---

# Field: What did you do?

Hovered over various rune and constant literals (both failing null-character cases and passing control-character/hex cases):

```go
package main

// Failing cases (null character U+0000):
var a = '\x00'        // null rune via hex escape
var b = '\000'        // null rune via octal escape
var c = '\u0000'      // null rune via Unicode escape
var d = '\U00000000'  // null rune via 8-digit Unicode escape
var s1 = "foo\x00bar" // hovering over \x00
var s2 = "foo\000bar" // hovering over \000

// Passing cases (other control characters / hex ints):
var e = '\a'          // control rune ('\a')
var s3 = "foo\a bar"  // hovering over \a
var x = 0x00          // untyped hex int (0x00)
```

---

# Field: What did you see happen?

No tooltip is displayed at all for any null character literals (`U+0000`) or null string escape sequences, whereas other control characters and literals display as expected:

| Variable Name | Observed Hover Tooltip |
| :--- | :--- |
| `a` (`'\x00'`) | *(No tooltip)* |
| `b` (`'\000'`) | *(No tooltip)* |
| `c` (`'\u0000'`) | *(No tooltip)* |
| `d` (`'\U00000000'`) | *(No tooltip)* |
| `s1` (`\x00` in `"foo\x00bar"`) | *(No tooltip)* |
| `s2` (`\000` in `"foo\000bar"`) | *(No tooltip)* |
| `e` (`'\a'`) | `U+0007, control` |
| `s3` (`\a` in `"foo\a bar"`) | `U+0007, control` |
| `x` (`0x00`) | `0` |

---

# Field: What did you expect to see?

Hovering over null character literals and null string escape sequences should display their Unicode code point and general category (`U+0000, control`), matching the behavior of other control characters:

| Variable Name | Expected Tooltip | Observed Tooltip |
| :--- | :--- | :--- |
| `a` (`'\x00'`) | `U+0000, control` | *(No tooltip)* |
| `b` (`'\000'`) | `U+0000, control` | *(No tooltip)* |
| `c` (`'\u0000'`) | `U+0000, control` | *(No tooltip)* |
| `d` (`'\U00000000'`) | `U+0000, control` | *(No tooltip)* |
| `s1` (`\x00` in `"foo\x00bar"`) | `U+0000, control` | *(No tooltip)* |
| `s2` (`\000` in `"foo\000bar"`) | `U+0000, control` | *(No tooltip)* |
| `e` (`'\a'`) | `U+0007, control` | *(As expected)* |
| `s3` (`\a` in `"foo\a bar"`) | `U+0007, control` | *(As expected)* |
| `x` (`0x00`) | `0` | *(As expected)* |

---

# Field: Editor and settings
VS Code with the Go extension (`golang.go`). Default language server settings (`"go.useLanguageServer": true`).

---

# Field: Logs
No error is logged; `textDocument/hover` simply returns `null` (empty result).

---

### Root Cause

In [`gopls/internal/golang/hover.go`](https://github.com/golang/tools/blob/b8f586f8588e28f0c9fc12bcafd850db114ba049/gopls/internal/golang/hover.go#L1040-L1044), `hoverConstantExpr` initializes:

```go
var (
    value      string    // if non-empty, a constant value to format in hover
    r          rune      // if non-zero, format a description of this rune in hover
    start, end token.Pos // hover span
)
```

1. **Under `case token.CHAR:`**
   The character literal is decoded into `r`. For null character literals (`'\x00'`, `'\000'`, `'\u0000'`), the decoded rune code point is `0`, so `r = 0`.
2. **Under `case constant.String:`**
   When hovering over an escape sequence:
   ```go
   if bytes.HasPrefix(rest, []byte(`\`)) {
       rr, _, tail, err := strconv.UnquoteChar(string(rest), lit.Value[0])
       ...
       r = rr
   }
   ```
   For `\x00` or `\000` in a string like `"foo\x00bar"`, `strconv.UnquoteChar` returns `rr = 0`, so `r = 0`.
3. **At the end of `hoverConstantExpr`:**
   ```go
   if value == "" && r == 0 { // nothing to format
       return protocol.Range{}, nil, nil
   }
   ...
   if r != 0 {
       runeName := runenames.Name(r)
       ...
       if strconv.IsPrint(r) {
           fmt.Fprintf(&b, "'%c', ", r)
       }
       fmt.Fprintf(&b, "U+%04X, %s", r, runeName)
   }
   ```
   Because `r == 0` and `value == ""`, the guard `if value == "" && r == 0` treats code point `0` as "nothing to format" and immediately returns `nil`. Furthermore, the subsequent formatting block is guarded by `if r != 0`, which would also suppress formatting.

---

### Suggested Fix

Use `r = -1` (or `r < 0`) as the sentinel for "no rune found" instead of `r == 0`. In Go, `rune` is a signed 32-bit integer (`int32`), while all valid Unicode code points are non-negative (`[0, 0x10FFFF]`). Using `-1` unambiguously distinguishes code point `0` (`U+0000`) from the absence of a rune without introducing extra state variables or touching the intermediate assignment branches:

```diff
 	var (
 		value      string    // if non-empty, a constant value to format in hover
-		r          rune      // if non-zero, format a description of this rune in hover
+		r          rune = -1 // if >= 0, format a description of this rune in hover
 		start, end token.Pos // hover span
 	)
 ...
-	if value == "" && r == 0 { // nothing to format
+	if value == "" && r < 0 { // nothing to format
 		return protocol.Range{}, nil, nil
 	}
 ...
-	if r != 0 {
+	if r >= 0 {
 		runeName := runenames.Name(r)
 		if len(runeName) > 0 && runeName[0] == '<' {
 			// Check if the rune looks like an HTML tag. If so, trim the surrounding <>
 			// characters to work around https://github.com/microsoft/vscode/issues/124042.
 			runeName = strings.TrimRight(runeName[1:], ">")
 		}
 		if b.Len() > 0 {
 			b.WriteString(", ")
 		}
 		if strconv.IsPrint(r) {
 			fmt.Fprintf(&b, "'%c', ", r)
 		}
 		fmt.Fprintf(&b, "U+%04X, %s", r, runeName)
 	}
```

### History & Regression Origin

1. **Original Support (April 2021):**
   In commit [`cd7c003c`](https://github.com/golang/tools/commit/cd7c003cd38984ebf592aa8e4755c832a8d632bf) (*"internal/lsp: add support for hovering runes"*, golang/go#38239), `findRune` used an explicit error (`ErrNoRuneFound`) to indicate when no rune was found. For null rune literals, `findRune` returned `r = 0, err = nil`, and `hoverRune` correctly displayed:
   ```text
   U+0000, control
   ```

2. **Subsequent Hover Refactorings (January 2023 & October 2025):**
   In commit [`b15a5bc4`](https://github.com/golang/tools/commit/b15a5bc45b8308c1b02168030361541ccd932107) (*"gopls/internal/lsp/source: rewrite hover"*, golang/go#57987), hover logic was restructured into `hoverLit`, replacing error returns with an `if r == 0` check.

   This check was later carried forward into `hoverConstantExpr` (commit [`1385a78d`](https://github.com/golang/tools/commit/1385a78d70e70be0e2624b441367353423c7ba22), golang/go#69058) with the comment:
   ```go
   r rune // if non-zero, format a description of this rune in hover
   ```
   Because `0` is the default zero-value for runes in Go, treating `r == 0` as the sentinel for "no rune found" unintentionally excluded the valid Unicode code point `0` (`U+0000`). Since existing tests in `basiclit.txt` only asserted non-zero characters (`'\a'`, `'∑'`, etc.), this edge case was not covered by the test suite.
   
---

> **Note:** Discovered while investigating #82196, but tracked separately as this is an independent regression where code point `0` (`U+0000`) is silently suppressed by the `r == 0` sentinel check, rather than returning an LSP error on non-ASCII bytes.