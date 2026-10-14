package main

import (
	"fmt"
	"unicode/utf8"
)

// EXPRESSION           BYTES (DEC)    BYTES (HEX)    HOVER       EXPECTED HOVER TOOLTIP                                  OBSERVED HOVER TOOLTIP                                PASS    COMMENTS
// --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
var a = "Æ"          // [195, 134]     [C3 86]        Æ           (No tooltip)                                            (As expected)                                         YES     -
var b = "\xc6"       // [198]          [C6]           \xc6        byte 0xC6 (198)                                         'Æ', U+00C6, LATIN CAPITAL LETTER AE                  NO      misleading: string contains byte 0xC6, not rune 'Æ'
var c = "\xc3\x86"   // [195, 134]     [C3 86]        \xc3        byte 0xC3 (195) [or: byte 1 of 'Æ' (U+00C6)]            'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE        NO      misleading: highlighted \xc3 is not 'Ã'
var d = "\xc3\x86"   // [195, 134]     [C3 86]        \x86        byte 0x86 (134) [or: byte 2 of 'Æ' (U+00C6)]            U+0086, control                                       NO      misleading: highlighted \x86 is not U+0086
var e = "\u00c6"     // [195, 134]     [C3 86]        \u00c6      'Æ', U+00C6, LATIN CAPITAL LETTER AE                    (As expected)                                         YES     -
var f = "\U000000c6" // [195, 134]     [C3 86]        \U000000c6  'Æ', U+00C6, LATIN CAPITAL LETTER AE                    (As expected)                                         YES     -
var g = "\x41"       // [65]           [41]           \x41        'A', U+0041, LATIN CAPITAL LETTER A                     (As expected)                                         YES     -
var h = "\xc3\x41"   // [195, 65]      [C3 41]        \xc3        byte 0xC3 (195) [incomplete UTF-8]                      'Ã', U+00C3, LATIN CAPITAL LETTER A WITH TILDE        NO      misleading: highlighted \xc3 is not 'Ã'
var i = "\xc3\x41"   // [195, 65]      [C3 41]        \x41        'A', U+0041, LATIN CAPITAL LETTER A                     (As expected)                                         YES     -
var j = "\x80\x81\x82"

type item struct {
	name string
	val  string
}

var items = []item{
	{"a", a}, {"b", b}, {"c", c}, {"d", d}, {"e", e},
	{"f", f}, {"g", g}, {"h", h}, {"i", i}, {"j", j},
}

func main() {
	fmt.Printf("%-4s %-20s %-15s %-8s %s\n", "VAR", "BYTES(DEC)", "BYTES(HEX)", "VALID?", "RUNES")
	fmt.Println("---------------------------------------------------------")
	for _, it := range items {
		decBytes := fmt.Sprintf("%v", []byte(it.val))
		hexBytes := fmt.Sprintf("[% X]", []byte(it.val))
		fmt.Printf("%-4s %-20s %-15s %-8t %U\n",
			it.name,
			decBytes,
			hexBytes,
			utf8.ValidString(it.val),
			[]rune(it.val),
		)
	}
}
