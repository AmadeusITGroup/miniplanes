package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	//fmt.Println("Ciao")
	//body := []string{}
	b := bytes.NewReader([]byte{})
	r, err := http.NewRequest("GET", "localhost:8080/airports", b)
	if err != nil {
		fmt.Printf("Hey error... %v", err)
		return
	}

	fmt.Printf("Hey OK! %s\n", r)
}
