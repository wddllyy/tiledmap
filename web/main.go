package main

import (
	"log"
	"net/http"

	"mazemap/tiledmap"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/test/", testHandler)
	http.HandleFunc("/maze", mazeHandler)
	http.HandleFunc("/cellular", cellularHandler)
	http.HandleFunc("/perlin", tiledmap.PerlinHandler)
	http.HandleFunc("/perlingray", tiledmap.PerlinGrayHandler)
	http.HandleFunc("/dungeon", tiledmap.DungeonHandler)
	http.HandleFunc("/wfc", tiledmap.WFCHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
