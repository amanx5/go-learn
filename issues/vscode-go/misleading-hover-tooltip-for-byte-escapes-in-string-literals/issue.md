# Field: Issue Title
x/tools/gopls: misleading hover tooltip for non-ASCII byte escapes (>= 0x80) in string literals

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

Hovered over byte escape sequences (`\xNN`, `\nnn`) in Go **string literals**, comparing non-ASCII bytes ($\ge \text{0x80}$) against passing baseline escapes:

```go
package main

// Case 1: Standalone non-ASCII byte (invalid UTF-8 string)
var b = "\xc6"

// Case 2: Multi-byte UTF-8 sequence encoding 'Æ' (U+00C6)
var c = "\xc3\x86"

// Case 3: Incomplete/corrupted UTF-8 sequence
var h = "\xc3\x41"

// Passing baseline: ASCII byte escape and Unicode escape
var g = "\x41"     // ASCII 'A'
var e = "\u00c6"   // Unicode escape for 'Æ'
```

---

# Field: What did you see happen?

When hovering over a hexadecimal byte escape (`\xNN`) in a string literal, `gopls` directly casts the raw byte value to a Unicode `rune` (effectively treating all bytes as ISO-8859-1 / Latin-1 code points: `rune(byteValue)`).

While this happens to produce the expected character for ASCII bytes ($< \text{0x80}$) because ASCII, Latin-1, and UTF-8 overlap, it completely breaks down for non-ASCII bytes ($\ge \text{0x80}$), producing false and misleading tooltips:

### Case 1: `"\xc6"` (Single invalid UTF-8 byte)
Hovering over `\xc6` claims it is `'Æ', U+00C6, LATIN CAPITAL LETTER AE`:

<img width="370" alt="hover \xc6" src="https://github.com/user-attachments/assets/be59331a-9d1a-4468-bcde-b0947b591180" />

In Go runtime, this byte is not valid UTF-8 and never decodes to `'Æ'` (`[]rune("\xc6")` yields `[U+FFFD]`).

### Case 2: `"\xc3\x86"` (Valid 2-byte UTF-8 encoding of `'Æ'`)

* **Case 2a:** Hovering over `\xc3` claims it is `'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE`:

  <img width="380" alt="hover \xc3" src="https://github.com/user-attachments/assets/6e9b18b1-6615-456e-bd2e-2a71f4f17db0" />

* **Case 2b:** Hovering over `\x86` claims it is `U+0086, control`:

  <img width="300" alt="hover \x86" src="https://github.com/user-attachments/assets/959c1011-4661-4958-aac9-e3c37514c7aa" />

Neither `'Ã'` nor `U+0086` exists in this string—the two bytes together represent the single rune `'Æ'` (`U+00C6`). Claiming that `\xc3` alone is `'Ã'` is misleading, as `\xc3` is only the first byte of `'Æ'`.

### Case 3: `"\xc3\x41"` (invalid UTF-8 sequence)
Hovering over `\xc3` claims it is `'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE`, even though `0xC3` is an orphaned UTF-8 leading byte followed by `'A'`:

<img width="380" alt="hover \xc3 in \xc3\x41" src="https://github.com/user-attachments/assets/7200db9c-fb2e-4b0a-a99a-a7111f857d37" />

---

# Field: What did you expect to see?

Hovering over an individual byte escape in a string literal should accurately reflect that it is a **byte** (`uint8`), rather than misrepresenting it as an independent Unicode code point:

| Case | Expression | Bytes (Hex) | Hover Target | Expected Tooltip | Observed Tooltip |
| :---: | :--- | :--- | :--- | :--- | :--- |
| Case 1 | `"\xc6"` | `[C6]` | `\xc6` | `byte 0xC6 (198) · invalid UTF-8 sequence` | `'Æ', U+00C6, LATIN CAPITAL LETTER AE` |
| Case 2a | `"\xc3\x86"` | `[C3 86]` | `\xc3` | `byte 0xC3 (195) · byte 1 of 2 for 'Æ' (U+00C6)` | `'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE` |
| Case 2b | `"\xc3\x86"` | `[C3 86]` | `\x86` | `byte 0x86 (134) · byte 2 of 2 for 'Æ' (U+00C6)` | `U+0086, control` |
| Case 3 | `"\xc3\x41"` | `[C3 41]` | `\xc3` | `byte 0xC3 (195) · invalid UTF-8 sequence` | `'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE` |
| Baseline | `"\x41"` | `[41]` | `\x41` | `'A', U+0041, LATIN CAPITAL LETTER A` | (As expected) |
| Baseline | `"\u00c6"` | `[C3 86]` | `\u00c6` | `'Æ', U+00C6, LATIN CAPITAL LETTER AE` | (As expected) |

---

# Field: Editor and settings
VS Code with the Go extension (`golang.go`). Default language server settings (`"go.useLanguageServer": true`).

---

# Field: Logs
No error is logged; `textDocument/hover` returns the incorrect rune description.

---

### Root Cause in `gopls`

In [`gopls/internal/golang/hover.go`](https://github.com/golang/tools/blob/27170b2725e9bde63c1e6b4ede2a4f63fcc8e8dc/gopls/internal/golang/hover.go#L1116-L1128), `hoverConstantExpr` inspects string escapes as follows:

```go
// Is the selected character denoted by an escape sequence?
if bytes.HasPrefix(rest, []byte(`\`)) {
    rr, _, tail, err := strconv.UnquoteChar(string(rest), lit.Value[0])
    if err != nil {
        return protocol.Range{}, nil, err
    }

    rawSize := len(rest) - len(tail)

    start, end = pos, pos+token.Pos(rawSize)
    r = rr // <--- Direct assignment treats the raw byte as a Unicode rune
}
```

1. `strconv.UnquoteChar` returns `(value rune, multibyte bool, tail string, err error)`.
   Per its documentation, `value` is *"the decoded Unicode code point **or byte value**"*.
   When decoding `\xNN` or `\nnn`, `strconv.UnquoteChar` returns `rune(byteValue)`.
2. `gopls` assigns `r = rr`, assuming whatever `strconv.UnquoteChar` returned is a Unicode code point, and formats it via `runenames.Name(r)` and `fmt.Fprintf(&b, "U+%04X, %s", r, runeName)`.
3. Per the [Go Language Specification on String Literals](https://go.dev/ref/spec#String_literals):
   > *"The three-digit octal (`\nnn`) and two-digit hexadecimal (`\xnn`) escapes represent **individual bytes** of the resulting string; all other escapes represent the (possibly multi-byte) UTF-8 encoding of individual characters. Thus inside a string literal `\377` and `\xFF` represent a single byte of value `0xFF=255`, while `ÿ`, `\u00ff`, `\U000000ff` and `\xc3\xbf` represent the two-byte UTF-8 encoding of character `U+00FF`."*

---

### Context: The Visual Selection Mismatch

In LSP clients (such as VS Code), hovering returns a `protocol.Range` matching the length of the escape sequence (`start, end = pos, pos+token.Pos(rawSize)`). This causes the editor to visually highlight **only** the 4-character escape sequence (e.g. `\xc3`).

Because the visual highlight covers only the single byte:
* Claiming `\xc3` is `'Ã'` is wrong (Latin-1 code point assumption).
* Claiming `\xc3` is `'Æ'` is also misleading (claiming one byte is the entire 2-byte character).

---

### Suggested Solutions

1. **Option 1: Display Raw Byte Value (Simple & Universally Accurate)**
   For byte escapes (`\xNN`, `\nnn`) in string literals with value $\ge \text{0x80}$, format as a byte rather than a Unicode character:
   ```text
   byte 0xC3 (195)
   ```
   This is 100% accurate regardless of whether the string contains UTF-8 text or raw binary data (such as image headers or network protocols).

2. **Option 2: Raw Byte Value + UTF-8 Context**
   Provide the byte value and its role within the string's UTF-8 structure:
   * Multi-byte rune members: `byte 0xC3 (195) · byte 1 of 2 for 'Æ' (U+00C6)`
   * Multi-byte continuation: `byte 0x86 (134) · byte 2 of 2 for 'Æ' (U+00C6)`
   * Incomplete or invalid sequences: `byte 0xC3 (195) · invalid UTF-8 sequence`

3. **Option 3: Multi-Byte Hover Range Expansion**
   If `gopls` desires displaying the decoded rune `'Æ'` for multi-byte sequences, it should expand the hover `protocol.Range` to cover the entire sequence (`\xc3\x86`), so the editor visually highlights both bytes together. If the byte sequence is invalid or incomplete (`\xc6`), highlight only the single byte and show its raw byte value.

---

I would like to work on this.