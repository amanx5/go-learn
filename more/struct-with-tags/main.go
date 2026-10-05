package main

import (
	"encoding/json/v2"
	"fmt"
)

// we can attach tags to fields in structs.
// Struct tags are metadata/annotations attached to fields.

type StructWithTags struct {
	Alias string `json:"Username"`
	Name  string
}

type StructWithConflict struct {
	Alias string `json:"Name"` // tag conflicting with a field
	Name  string
}

func main() {
	encode(StructWithTags{"amanx5", "aman"})
	encode(StructWithConflict{"amanx5", "aman"})
	decode(`{"Username":"amanx5","Name":"aman"}`, StructWithTags{})
	decode(`{"Name":"aman"}`, StructWithConflict{})
}

func encode(in any) {
	out, err := json.Marshal(in)

	var res any
	if err == nil {
		res = string(out)
	} else {
		res = err
	}

	fmt.Printf("\n%+v ===encoded===> %v\n", in, res)
}

func decode(in string, out any) {
	err := json.Unmarshal([]byte(in), &out)

	var res any
	if err == nil {
		res = out
	} else {
		res = err
	}

	fmt.Printf("\n%v ===decoded===> %+v\n", in, res)
}
