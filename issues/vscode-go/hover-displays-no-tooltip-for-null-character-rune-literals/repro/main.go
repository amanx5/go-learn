package main

import (
	"fmt"
	"strconv"
)

// EXPRESSION           DECIMAL		 EXPECTED HOVER TOOLTIP     OBSERVED == EXPECTED
// -------------------------------------------------------------------------------------------------
var a rune = 0x00       //  0			0                           YES
var b = '0'             //  48			'0', U+0030, DIGIT ZERO     YES
var c = '\x00'          //  0			U+0000, control             NO: No tooltip
var d = '\x01'          //  1			U+0001, control             YES
var e = '\000'          //  0			U+0000, control             NO: No tooltip
var f = '\001'          //  1			U+0001, control             YES
var g = '\u0000'        //  0			U+0000, control             NO: No tooltip
var h = '\u0001'        //  1			U+0001, control             YES
var i = '\U00000000'    //  0			U+0000, control             NO: No tooltip
var j = '\U00000001'    //  1			U+0001, control				YES
var k = unq("\x00test") //  0			U+0000, control  		 	NO: No tooltip
var l = unq("\000test") //  0			U+0000, control  		 	NO: No tooltip
var m = unq("\atest")   //  7			U+0007, control 		 	YES
var n = unq(" test")    //  32 			(No tooltip)				YES

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
	{"k", (k)},
	{"l", (l)},
	{"m", (m)},
	{"n", (n)},
}

func main() {
	for _, v := range vars {
		fmt.Printf("%-5s %-10U %-6c %d\n", v.name, v.value, v.value, v.value)
	}
}

func unq(s string) rune {
	r, _, _, _ := strconv.UnquoteChar(s, '"')
	return r
}
