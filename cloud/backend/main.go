package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type TestResult struct {
	ID        int64     `json:"id"`
	Scenario  string    `json:"scenario"`
	Passed    bool      `json:"passed"`
	Latency   int64     `json:"latency_ms"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./mockingjay.db"
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/results", createResult).Methods("POST")
	r.HandleFunc("/api/results", getResults).Methods("GET")
	r.HandleFunc("/api/health", healthCheck).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Printf("Database: %s", dbPath)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS test_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scenario TEXT NOT NULL,
		passed BOOLEAN NOT NULL,
		latency_ms INTEGER NOT NULL,
		error TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_created_at ON test_results(created_at DESC);
	`
	_, err := db.Exec(query)
	return err
}

func createResult(w http.ResponseWriter, r *http.Request) {
	var result TestResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO test_results (scenario, passed, latency_ms, error)
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, result.Scenario, result.Passed, result.Latency, result.Error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	result.ID = id
	result.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getResults(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}

	query := `
		SELECT id, scenario, passed, latency_ms, error, created_at
		FROM test_results
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []TestResult{}
	for rows.Next() {
		var r TestResult
		var errStr sql.NullString
		if err := rows.Scan(&r.ID, &r.Scenario, &r.Passed, &r.Latency, &errStr, &r.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if errStr.Valid {
			r.Error = errStr.String
		}
		results = append(results, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "OK")
}
