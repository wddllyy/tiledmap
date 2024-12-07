package main

import (
	"fmt"
	"net/http"
)

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

// handler reverses the subpath string
func testHandler(w http.ResponseWriter, req *http.Request) {
	subPath := req.URL.Path[len("/test/"):]
	runes := []rune(subPath)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	fmt.Fprintf(w, "Reversed path: %s\n", string(runes))
}
