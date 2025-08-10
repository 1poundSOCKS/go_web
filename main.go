package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	transaction_time string
}

func handler(w http.ResponseWriter, r *http.Request) {

	// "db_connection": "host=localhost port=5432 dbname=mydatabase user=myuser password=mypassword",

	connStr := "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PostgreSQL version:", version)

	var u User
	err = db.QueryRow(`SELECT transaction_time FROM jobs`).Scan(&u.transaction_time)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("transaction time:", u.transaction_time)

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
