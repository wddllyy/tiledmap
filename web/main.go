package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/test/", testHandler)
	http.HandleFunc("/maze", mazeHandler)
	http.HandleFunc("/cellular", cellularHandler)
	http.HandleFunc("/perlin", perlinHandler)
	http.HandleFunc("/perlingray", perlinGrayHandler)
	http.HandleFunc("/dungeon", dungeonHandler)
	http.HandleFunc("/wfc", wfcHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
