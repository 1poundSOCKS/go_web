package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	p := Person{Name: "Alice", Age: 30}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(jsonData))
}

func main() {

	http.HandleFunc("/", handler)

	port := ":8080"
	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
