package main

import "fmt"

// EXPRESSION			DESCRIPTION									DECIMAL				EXPECTED TOOLTIP (on hovering the literal)			OBSERVED == EXPECTED
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
var a rune = 65      // decimal 									65					/* no tooltip */									YES
var b rune = 198     // decimal										198					/* no tooltip */									YES
var c rune = 0x41    // hex integer literal							65					65, 'A', U+0041, LATIN CAPITAL LETTER A				YES
var d rune = 0xc6    // hex integer literal							198					198, 'Æ', U+00C6, LATIN CAPITAL LETTER AE			YES
var e = 'A'          // single character rune literal				65					'A', U+0041, LATIN CAPITAL LETTER A					YES
var f = 'Æ'          // single character rune literal				198					'Æ', U+00C6, LATIN CAPITAL LETTER AE				YES
var g = '\x41'       // rune literal using byte escape 				65					'A', U+0041, LATIN CAPITAL LETTER A					YES
var h = '\xc6'       // rune literal using byte escape 				198					'Æ', U+00C6, LATIN CAPITAL LETTER AE				NO: No tooltip, it fails with "rune error"
var i = '\101'       // rune literal using octal escape 			65					'A', U+0041, LATIN CAPITAL LETTER A					YES
var j = '\306'       // rune literal using octal escape 			198					'Æ', U+00C6, LATIN CAPITAL LETTER AE				NO: No tooltip, it fails with "rune error"
var k = '\u0041'     // rune literal using code point escape 		65					'A', U+0041, LATIN CAPITAL LETTER A					YES
var l = '\u00c6'     // rune literal using code point escape 		198					'Æ', U+00C6, LATIN CAPITAL LETTER AE				YES
var m = '\U00000041' // rune literal using code point escape 		65					'A', U+0041, LATIN CAPITAL LETTER A					YES
var n = '\U000000c6' // rune literal using code point escape 		198					'Æ', U+00C6, LATIN CAPITAL LETTER AE				YES
var o = '\U0010FFFF' // rune literal using code point escape 		1,114,111			U+10FFFF,											YES
var p = '\a'         // control rune literal						7					U+0007, control										YES
// var x = ''		 // empty rune literal							- 					/* no tooltip */									YES
// var x = 'abc'	 // invalid rune literal						-					/* no tooltip */									YES
// var x = '\x0'	 // invalid rune literal						- 					illegal character U+0027 ''' in escape sequence		YES
// var x = '\U0020FFFF'	// invalid Unicode code point escape		-					escape sequence is invalid Unicode code point		YES

type Var struct {
	name  string
	value rune
}

var vars = []Var{
	{"a", a},
	{"b", b},
	{"c", c},
	{"d", d},
	{"e", e},
	{"f", f},
	{"g", g},
	{"h", h},
	{"i", i},
	{"j", j},
	{"k", k},
	{"l", l},
	{"m", m},
	{"n", n},
	{"o", o},
	{"p", p},
}

func main() {
	for _, v := range vars {
		fmt.Printf("%-5s %-10U %-6c %d\n", v.name, v.value, v.value, v.value)
	}
}
