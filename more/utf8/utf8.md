# UTF-8 Encoding Scheme

UTF-8 is a variable-length character encoding for Unicode, encoding each code point into 1 to 4 bytes.

---

## Memory Representation Table

| Unicode Code Point | Sequence | Byte 1 (Header) | Byte 2 (Continuation) | Byte 3 (Continuation) | Byte 4 (Continuation) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`U+0000` to `U+007F`**<br>(0 to 127) | 1 byte | `0xxxxxxx`<br>Dec: 0 to 127<br>Bin: `00000000`–`01111111` | — | — | — |
| **`U+0080` to `U+07FF`**<br>(128 to 2,047) | 2 bytes | `110xxxxx`<br>Dec: 194 to 223<br>Bin: `11000010`–`11011111` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` | — | — |
| **`U+0800` to `U+FFFF`**<br>(2,048 to 65,535) | 3 bytes | `1110xxxx`<br>Dec: 224 to 239<br>Bin: `11100000`–`11101111` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` | — |
| **`U+10000` to `U+10FFFF`**<br>(65,536 to 1,114,111) | 4 bytes | `11110xxx`<br>Dec: 240 to 244<br>Bin: `11110000`–`11110100` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` | `10xxxxxx`<br>Dec: 128 to 191<br>Bin: `10000000`–`10111111` |

---

## Key Rules

1. **Self-Synchronizing:**
   - Any byte matching `10xxxxxx` (`128` to `191`) is always a **continuation byte**.
   - Header bytes start with `0` (1-byte), `110` (2-byte), `1110` (3-byte), or `11110` (4-byte).
   - This allows decoders to identify character boundaries from any random position in a byte stream.

2. **No Overlong Encodings:**
   - 2-byte sequences start at `194` (`11000010`), skipping `192`–`193` (`0xC0`–`0xC1`) to prevent redundant ASCII representations.
   - 4-byte sequences stop at `244` (`11110100`) because Unicode's highest defined code point is `U+10FFFF`. Bytes `245`–`255` are invalid.
