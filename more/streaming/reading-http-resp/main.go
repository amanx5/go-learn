package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://dummyjson.com/users/?limit=1")

	buf := make([]byte, 4096)

	for {
		n, err := resp.Body.Read(buf)

		fmt.Println(string(buf[:n]))

		if err == io.EOF {
			break
		}
	}
}
