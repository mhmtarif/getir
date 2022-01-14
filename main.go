package main

import (
	"log"
	"net/http"
)

func main() {
	initMongoDb()
	initMemoryDb()
	http.HandleFunc("/", hello)
	http.HandleFunc("/mongo", mongoApi)
	http.HandleFunc("/in-memory", memoryApi)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
