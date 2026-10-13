# Field: Issue Title
x/tools/gopls: hover fails on rune literals with non-ASCII byte escapes (>= 0x80)

---

# Field: gopls version
```text
golang.org/x/tools/gopls v0.23.0
    golang.org/x/tools/gopls@v0.23.0 h1:Dn6mf9WXu9iLnTftDDMb9wV0c6Se7PjzEMqP0LEe08Y=
```
*(Also verified on master commit `31f699c742197c9eb5e5c2499e4cdf2b80a8f53f`)*

---

# Field: go env
```text
set GO111MODULE=
set GOAMD64=v1
set GOARCH=amd64
set GOBIN=C:\Users\<username>\go\bin
set GOCACHE=C:\Users\<username>\AppData\Local\go-build
set GOENV=C:\Users\<username>\AppData\Roaming\go\env
set GOEXE=.exe
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOMOD=NUL
set GOMODCACHE=C:\Users\<username>\go\pkg\mod
set GOOS=windows
set GOPATH=C:\Users\<username>\go
set GOPROXY=https://proxy.golang.org,direct
set GOROOT=C:\Program Files\Go
set GOSUMDB=sum.golang.org
set GOTOOLCHAIN=auto
set GOVERSION=go1.27.1
```

---

# Field: What did you do?
Hovered over a Go rune literal containing a hexadecimal or octal byte escape $\ge 128$ (`0x80`), such as `'\xc6'` or `'\306'`:

```go
package main

var a = '\xc6'   // byte escape for 198 ('Æ')
var b = '\306'   // octal escape for 198 ('Æ')
var c = '\u00c6' // Unicode escape for 198 ('Æ')
```

---

# Field: What did you see happen?
No tooltip is displayed. The LSP client logs a failure for `textDocument/hover`:

```text
[Error] Request textDocument/hover failed.
  Message: rune error
  Code: 0
```

### Behavior Comparison

| Literal | Constant Value | Expected `gopls` Hover | Observed `gopls` Hover |
| :--- | :--- | :--- | :--- |
| `'\xc6'` | `rune` (`198` = `'Æ'`) | `'Æ', U+00C6, LATIN CAPITAL LETTER AE` | **No tooltip** (logs `rune error`) |
| `'\306'` | `rune` (`198` = `'Æ'`) | `'Æ', U+00C6, LATIN CAPITAL LETTER AE` | **No tooltip** (logs `rune error`) |
| `'\u00c6'` | `rune` (`198` = `'Æ'`) | `'Æ', U+00C6, LATIN CAPITAL LETTER AE` | As expected |
| `0x00c6` | `untyped int` (`198`) | `198, 'Æ', U+00C6, LATIN CAPITAL LETTER AE` | As expected |

---

# Field: What did you expect to see?
Hovering over `'\xc6'` and `'\306'` should display the rune hover tooltip:
```text
'Æ', U+00C6, LATIN CAPITAL LETTER AE
```
(Matching the tooltip displayed when hovering over `'\u00c6'`).

---

# Field: Editor and settings
VS Code with the Go extension (`golang.go`). Default language server settings (`"go.useLanguageServer": true`).

---

# Field: Logs
```text
[Error] Request textDocument/hover failed.
  Message: rune error
  Code: 0
```

### Root Cause
In [`gopls/internal/golang/hover.go`](https://github.com/golang/tools/blob/31f699c742197c9eb5e5c2499e4cdf2b80a8f53f/gopls/internal/golang/hover.go#L1058-L1069), `hoverConstantExpr` decodes character literals (`token.CHAR`) using `strconv.Unquote` followed by `utf8.DecodeRuneInString`:

```go
case token.CHAR:
    s, err := strconv.Unquote(lit.Value)
    if err != nil {
        return protocol.Range{}, nil, nil
    }
    r, _ = utf8.DecodeRuneInString(s)
    if r == utf8.RuneError {
        return protocol.Range{}, nil, fmt.Errorf("rune error") // <--- Fails here
    }
    start, end = lit.Pos(), lit.End()
```

1. Per the [Go Language Specification on Rune Literals](https://go.dev/ref/spec#Rune_literals):
   > *"In rune literals, octal and hexadecimal escapes represent values in a range determined by the rune type: an octal or hexadecimal escape must represent a value between 0 and 255 inclusive."*
   Therefore, `'\xc6'` is a valid rune constant with integer value `198` (`'Æ'`).
2. When `strconv.Unquote("'\xc6'")` runs, its single-quote branch appends the raw byte (`buf = append(buf, byte(v))`) rather than UTF-8 encoding it, returning `"\xc6"` (`[]byte{0xC6}`).
3. Calling `utf8.DecodeRuneInString("\xc6")` on this non-UTF-8 single byte returns `utf8.RuneError`.
4. `gopls` then returns `fmt.Errorf("rune error")`, failing the hover request.

### Suggested Fix
Read the rune value directly from `tv.Value` (as evaluated by `go/types`), while retaining `strconv.Unquote(lit.Value)` to validate literal syntax:

```diff
 			case token.CHAR:
-				s, err := strconv.Unquote(lit.Value)
+				_, err := strconv.Unquote(lit.Value)
 				if err != nil {
 					// If the conversion fails, it's because of an invalid syntax, therefore
 					// there is no rune to be found.
 					return protocol.Range{}, nil, nil
 				}
-				r, _ = utf8.DecodeRuneInString(s)
-				if r == utf8.RuneError {
-					return protocol.Range{}, nil, fmt.Errorf("rune error")
-				}
-				start, end = lit.Pos(), lit.End()
+				if v, ok := constant.Int64Val(tv.Value); ok && v >= 0 && v <= utf8.MaxRune && utf8.ValidRune(rune(v)) {
+					r = rune(v)
+					start, end = lit.Pos(), lit.End()
+				}
```

#### Why `strconv.Unquote` syntax validation must be retained
`strconv.Unquote(lit.Value)` must be retained to guard against syntactically invalid character literals (such as `'abc'`, `''`, or malformed escape sequences).

During type checking of an invalid literal like `'abc'`, `go/types` still populates `tv.Value` with a partial integer constant (`97`, the ASCII code for `'a'`). This occurs because:
1. `types.Checker.basicLit` invokes `types.operand.setConst` with `constant.MakeFromLiteral`.
2. Inside [`constant.MakeFromLiteral`](https://pkg.go.dev/go/constant#MakeFromLiteral) (for `token.CHAR`), it calls `strconv.UnquoteChar` to decode the first character, but ignores the remaining `tail` rather than returning `constant.Unknown`:
   - Call chain: `types.Checker.basicLit` > `types.operand.setConst` > `constant.MakeFromLiteral` > `strconv.UnquoteChar`
3. As a result, `tv.Value` contains `constant.MakeInt64('a')` even though the literal failed syntax validation.

While `constant.MakeFromLiteral` is a candidate for an upstream fix in `golang/go`, `gopls` (belonging to `golang.org/x/tools`) must remain backward compatible across older Go toolchain versions. Thus, retaining the `strconv.Unquote` guard within `gopls` is necessary regardless of future standard library updates.

---

> **Note:** *I have a working fix and regression test ready for this in `x/tools`, and will be submitting a Gerrit CL shortly.*
