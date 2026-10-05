# Go Formatting Verbs (`fmt`)

A quick reference guide for formatting verbs and modifiers in Go.

---

## 1. General Verbs

These universal verbs work across multiple data types.

| Verb | Description | Example Output |
| :--- | :--- | :--- |
| `%v` | The value in its default format | `123`, `true`, `map[a:1]` |
| `%+v` | Structs with field names included | `{Name:Alice Age:30}` |
| `%#v` | The value represented in Go syntax | `main.User{Name:"Alice"}` |
| `%T` | The type of the value | `int`, `string`, `main.User` |
| `%%` | A literal percent sign | `%` |

---

## 2. Integers

| Verb | Description | Example (for `15`) |
| :--- | :--- | :--- |
| `%d` | Base 10 (decimal) | `15` |
| `%b` | Base 2 (binary) | `1111` |
| `%o` | Base 8 (octal) | `17` |
| `%O` | Octal with a `0o` prefix | `0o17` |
| `%x` | Base 16 (hexadecimal), lowercase | `f` |
| `%X` | Base 16 (hexadecimal), uppercase | `F` |
| `%c` | Character represented by corresponding Unicode code point | `\x0f` (or equivalent character) |

---

## 3. Floating-Point Numbers

| Verb | Description | Example (for `123.456`) |
| :--- | :--- | :--- |
| `%f` | Standard decimal notation (default 6 decimal places) | `123.456000` |
| `%e` | Scientific notation, lowercase | `1.234560e+02` |
| `%E` | Scientific notation, uppercase | `1.234560E+02` |
| `%g` | Uses `%e` or `%f`, whichever is more compact | `123.456` |

---

## 4. Strings and Byte Slices

| Verb | Description | Example (for `"hello"`) |
| :--- | :--- | :--- |
| `%s` | Raw string or uninterpreted bytes | `hello` |
| `%q` | Double-quoted string safely escaped with Go syntax | `"hello"` |
| `%x` | Hex dump of the string/bytes, lowercase | `68656c6c6f` |

---

## 5. Booleans and Pointers

| Verb | Description | Example |
| :--- | :--- | :--- |
| `%t` | Boolean (`true` or `false`) | `true` |
| `%p` | Pointer address in hexadecimal (with leading `0x`) | `0xc0000b2008` |

---

## 6. Width and Precision Modifiers

You can place optional flags, width parameters, and precision limits between the `%` symbol and the verb character.

| Modifier | Target Type | Description | Example | Result |
| :--- | :--- | :--- | :--- | :--- |
| `%5d` | Integer | Minimum width of 5, right-justified (padded with spaces) | `fmt.Printf("%5d", 12)` | `   12` |
| `%-5d` | Integer | Minimum width of 5, left-justified (padded with spaces) | `fmt.Printf("%-5d", 12)` | `12   ` |
| `%05d` | Integer | Minimum width of 5, padded with leading zeros | `fmt.Printf("%05d", 12)` | `00012` |
| `%+d` | Integer | Always print a sign (`+` or `-`) | `fmt.Printf("%+d", 12)` | `+12` |
| `%.2f` | Float | Custom precision (rounds to 2 decimal places) | `fmt.Printf("%.2f", 1.234)` | `1.23` |
| `%6.2f` | Float | Minimum width of 6, precision of 2 decimal places | `fmt.Printf("%6.2f", 1.234)` | `  1.23` |
| `%#x` | Hex/Octal | Forces a leading `0x` or `0o` prefix | `fmt.Printf("%#x", 15)` | `0xf` |