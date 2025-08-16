package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

type JobPostgres struct {
	TransactionTime string `json:"transaction_time"`
	TransactionID   string `json:"transaction_id"`
	JobID           int64  `json:"job_id"`
	JobName         string `json:"job_name"`
	Duration        *int32 `json:"duration"`
}

type JobOracle struct {
	JobRef string `json:"job_ref"`
}

func main() {

	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/oracle", handlerOracle)
	http.HandleFunc("/postgres", handlerPostgres)

	port := ":8080"
	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {

	data := getPostgresData()

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "<!doctype html>")
	fmt.Fprintf(w, "<html lang=\"en\">")
	fmt.Fprintf(w, "<head>")
	fmt.Fprintf(w, "  <meta charset=\"utf-8\" />")
	fmt.Fprintf(w, "  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />")
	fmt.Fprintf(w, "  <title>Minimal HTML Page</title>")
	fmt.Fprintf(w, "</head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "  <h1>Hello, world!</h1>")
	fmt.Fprintf(w, "  <p>This is a minimal HTML page.</p>")
	fmt.Fprintf(w, "  <h1>Dropdown List Example</h1>")
	fmt.Fprintf(w, "  <label for=\"options\">Choose an option:</label>")
	fmt.Fprintf(w, "  <select id=\"options\" name=\"options\">")
	fmt.Fprintf(w, "    <option value=\"option1\">Option 1</option>")
	fmt.Fprintf(w, "    <option value=\"option2\">Option 2</option>")
	fmt.Fprintf(w, "    <option value=\"option3\">Option 3</option>")
	fmt.Fprintf(w, "    <option value=\"option4\">Option 4</option>")
	fmt.Fprintf(w, "  </select>")
	fmt.Fprintf(w, "	  <h2>Output</h2>")
	fmt.Fprintf(w, "  <pre id=\"output\">"+data+"</pre>")
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
}

func handlerPostgres(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	data := getPostgresData()
	fmt.Fprintf(w, "%s", data)
}

func handlerOracle(w http.ResponseWriter, r *http.Request) {

	connStr := "mcob/mcob@localhost:1521/MC19"
	db, err := sql.Open("godror", connStr)

	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	rows, err := db.Query("SELECT job_ref FROM SP111")
	if err != nil {
		log.Fatal(err)
	}

	var jobs []JobOracle
	for rows.Next() {
		var u JobOracle
		err := rows.Scan(&u.JobRef)
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

func getPostgresData() string {
	connStr := "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return "Error opening connection: " + err.Error()
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return "Cannot connect to database: " + err.Error()
	}

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return err.Error()
	}
	fmt.Println("PostgreSQL version:", version)

	rows, err := db.Query(`SELECT transaction_time, transaction_id, job_id, job_name, duration FROM jobs`)
	if err != nil {
		return err.Error()
	}

	var jobs []JobPostgres
	for rows.Next() {
		var u JobPostgres
		err := rows.Scan(&u.TransactionTime, &u.TransactionID, &u.JobID, &u.JobName, &u.Duration)
		if err != nil {
			log.Fatal(err)
		}
		jobs = append(jobs, u)
	}

	jsonData, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(jsonData)
}
