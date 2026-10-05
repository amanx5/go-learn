package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

type Version struct {
	Major, Minor, Patch int
}

type Semver struct {
	Major, Minor, Patch int
}

// MarshalJSONTo is bound to [Semver] type. This ensures [Semver] implements [json.MarshalerTo] interface.
//
// This ensures functions like [json.Marshal] and [json.MarshalWrite]
// invoke this custom logic whenever a [Semver] value is serialized
// (whether standalone instance, embedded in a struct, or inside slices/maps).
func (sv Semver) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch))
}

// Compile-time assertion ensuring [Semver] satisfies [json.MarshalerTo].
// This is a standard Go pattern that guarantees the compiler verifies interface compatibility.
var _ json.MarshalerTo = Semver{}

func main() {
	ver, _ := json.Marshal(Version{0, 0, 1})
	fmt.Println(string(ver)) // {"Major":0,"Minor":0,"Patch":1}

	semver, _ := json.Marshal(Semver{0, 0, 1})
	fmt.Println(string(semver)) // "0.0.1"
}
