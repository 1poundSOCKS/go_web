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

	text, err := getPostgresData()

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
	if err != nil {
		fmt.Fprintf(w, "  <pre id=\"output\">"+text+": "+err.Error()+"</pre>")
	} else {
		fmt.Fprintf(w, "  <pre id=\"output\">"+text+"</pre>")
	}
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
}

func handlerPostgres(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	text, err := getPostgresData()
	if err != nil {
		fmt.Fprintf(w, "%s: %s", text, err.Error())
	} else {
		fmt.Fprintf(w, "%s", text)
	}
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

	rows, err := db.Query("SELECT * FROM SP111_JOBS")
	if err != nil {
		log.Fatal(err)
	}

	data := convertToJson(rows)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s", data)
}

func getPostgresData() (string, error) {

	connStr := "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return "Error opening connection", err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return "Cannot connect to database: ", err
	}

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return "error on SELECT", err
	}
	fmt.Println("PostgreSQL version:", version)

	rows, err := db.Query(`SELECT * FROM jobs`)
	if err != nil {
		return "error on SELECT", err
	}

	results := convertToJson(rows)

	return results, nil
}

func convertToJson(rows *sql.Rows) string {

	// Slice of maps to hold rows
	var results []map[string]interface{}

	cols, _ := rows.Columns()
	for rows.Next() {
		// Make a slice for the values
		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		// Scan the result into the pointers
		rows.Scan(valPtrs...)

		// Create a map for the row
		rowMap := make(map[string]interface{})
		for i, col := range cols {
			var v interface{}
			b, ok := vals[i].([]byte)
			if ok {
				v = string(b)
			} else {
				v = vals[i]
			}
			rowMap[col] = v
		}

		results = append(results, rowMap)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}
