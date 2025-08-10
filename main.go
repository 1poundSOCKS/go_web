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

type Job struct {
	TransactionTime string `json:"transaction_time"`
	TransactionID   string `json:"transaction_id"`
	JobID           int64  `json:"job_id"`
	JobName         string `json:"job_name"`
	Duration        *int32 `json:"duration"`
}

func handler(w http.ResponseWriter, r *http.Request) {

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

	rows, err := db.Query(`SELECT transaction_time, transaction_id, job_id, job_name, duration FROM jobs`)
	if err != nil {
		log.Fatal(err)
	}

	var jobs []Job
	for rows.Next() {
		var u Job
		err := rows.Scan(&u.TransactionTime, &u.TransactionID, &u.JobID, &u.JobName, &u.Duration)
		if err != nil {
			log.Fatal(err)
		}
		jobs = append(jobs, u)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s", string(jsonData))
}

func main() {

	http.HandleFunc("/", handler)

	port := ":8080"
	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
