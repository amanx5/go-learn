package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// EXPRESSION          	TYPE                  		SIZE / REPRESENTATION                  PURPOSE
// 'A'                  Untyped rune constant 		No runtime storage by itself           Represents a Unicode code point
// r := 'A'             rune (int32)          		4 bytes for the variable               Stores a Unicode code point
// "hello"              string                		5 bytes of string data                 Stores a sequence of bytes
// "\x??"               string                		1 byte of string data                  Specifies exactly one byte using 2 hex digits
// "\u????"             string                		Encoded as 1–3 UTF-8 bytes             Specifies a Unicode code point using 4 hex digits
// "\U????????"         string                		Encoded as 1–4 UTF-8 bytes             Specifies a Unicode code point using 8 hex digits

// 198 = 0x00c6 = '\xc6' = "Æ" = "\u00c6" = "\xc3\x86" != "\xc6"
var v1 rune = 198
var v2 rune = 0x00c6
var v3 rune = '\xc6'
var v4 = "Æ"
var v5 = "\u00c6"
var v6 = "\xc3\x86"

// "\xC6" is a string containing one byte (0xC6), which is NOT a valid UTF-8 byte sequence.
// In a Go string, any byte between 0x80 (128) and 0xFF (255) is an invalid standalone UTF-8 character.
// Maximum number that can be stored in a single byte UTF8 encoding is 0x7F (127)
// there is one BUG in intellisense(hover tooltip interprets below string as 'Æ' which is incorrect)
var v7 = "\xc6"

// ʔ (U+0294) = 0x294 = 0d660 = "\u0294"
var v8 = "ʔ"
var v9 = "\u0294"

var v10 = "é"
var v11 = "e\u0301" // same as "é"
var v12 = "泃"
var v13 = "\x91"

// put all values in dump so that go optimizer doesnt remove the unused values
var dump = []any{v1, v2, v3, v5, v5, v6, v7, v8, v9, v11}

var VALUE_TO_USE = v7
var IMPL_TO_USE = Impl4

func main() {
	log("val", VALUE_TO_USE)

	str := string(VALUE_TO_USE)
	log("str", str)

	rev, revErr := Reverse(str)
	log("rev", rev)
	log("revErr", revErr)

	revrev, revrevErr := Reverse(rev)
	log("revrev", revrev)
	log("revrevErr", revrevErr)

	log("Critera 1 (utf8.ValidString(rev))", utf8.ValidString(rev))
	log("Critera 2 (str == revrev)", str == revrev)
}

func Reverse(s string) (string, error) {
	return IMPL_TO_USE(s)
}

// this impl satisfies both criterias
// but does not handle combining characters like "e\u0301" correctly
func Impl4(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}

	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

// this impl fails the Criteria 2
func Impl3(s string) (string, error) {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

// this impl failed with str="\xc6" (invalid utf8 string), corpus file: e504fa449c34beb2
func Impl2(str string) (string, error) {
	temp := []byte(str) // "\xc6" = 0xc6 = [198]

	tempLastUsedIndex := len(str)

	// When we loop over a string using for range, Go acts as a UTF-8 decoder.
	// when it hits the invalid byte like 0xC6, the loop's decoder fails,
	// throws away the original byte, and forces the variable currentRune
	// to be 65533 (which is 0xFFFD, the Unicode replacement character ``)
	// which breaks the whole logic
	for _, currentRune := range str {
		currentRuneToBytes := []byte(string(currentRune))
		tempCurrentStartIndex := tempLastUsedIndex - len(currentRuneToBytes)
		for i, j := tempCurrentStartIndex, 0; i < tempLastUsedIndex; i, j = i+1, j+1 {
			temp[i] = currentRuneToBytes[j]
		}
		tempLastUsedIndex = tempCurrentStartIndex
	}

	return string(temp), nil
}

// this impl failed with s=ʔ , corpus file=a8c8b98939d243c8
func Impl1(s string) (string, error) {
	rev := []byte(s)
	len := len(s)
	for i, j := 0, len-1; i < len/2; i, j = i+1, j-1 {
		rev[i], rev[j] = s[j], s[i]
	}
	return string(rev), nil
}

func log(lhs string, rhs any) {
	rhsEnhanced := rhs

	switch r := rhs.(type) {
	case bool:
		rhsEnhanced = b2T(r, "passed", "failed")
	case string:
		isValid := utf8.ValidString(r)
		rhsEnhanced = fmt.Sprintf("%q (%v)", r, validity(isValid))
	case rune:
		isValid := utf8.ValidRune(r)
		rhsEnhanced = fmt.Sprintf("%v (%v)", r, validity(isValid))
	}

	fmt.Printf("%-34s = %v", lhs, rhsEnhanced)
	fmt.Println()
}

func validity(b bool) string {
	return b2T(b, "valid", "invalid")
}

// bool to text
func b2T(b bool, t string, f string) string {
	if b {
		return t
	}
	return f
}
