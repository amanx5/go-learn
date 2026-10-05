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

func (sv *Semver) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var str string

	err := json.UnmarshalDecode(dec, &str)
	if err != nil {
		return err
	}

	_, err2 := fmt.Sscanf(str, "%d.%d.%d", &sv.Major, &sv.Minor, &sv.Patch)
	return err2
}

func main() {
	var out1 Version
	err1 := json.Unmarshal([]byte(`"0.0.1"`), &out1)
	fmt.Println("out1: ", out1) // {0 0 0}
	fmt.Println("err1: ", err1) // json: cannot unmarshal JSON string into Go main.Version

	var out2 Semver
	err2 := json.Unmarshal([]byte(`"0.0.1"`), &out2)
	fmt.Println("out2: ", out2) // {0 0 1}
	fmt.Println("err2: ", err2) // <nil>
}
