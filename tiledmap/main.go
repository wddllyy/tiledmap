package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/test/", testHandler)
	http.HandleFunc("/maze", mazeHandler)
	http.HandleFunc("/cellular", cellularHandler)
	http.HandleFunc("/perlin", perlinHandler)
	http.HandleFunc("/perlingray", perlinGrayHandler)
	http.HandleFunc("/dungeon", dungeonHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
