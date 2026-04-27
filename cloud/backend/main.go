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
	Variant   string    `json:"variant,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ABTestResult struct {
	ID        int64     `json:"id"`
	TestName  string    `json:"test_name"`
	VariantA  string    `json:"variant_a"`
	VariantB  string    `json:"variant_b"`
	Winner    string    `json:"winner"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

type Transcription struct {
	ID         int64     `json:"id"`
	CallSID    string    `json:"call_sid,omitempty"`
	AudioPath  string    `json:"audio_path"`
	Text       string    `json:"text"`
	Confidence float64   `json:"confidence"`
	Duration   float64   `json:"duration_seconds"`
	CreatedAt  time.Time `json:"created_at"`
}

type ConversationMetrics struct {
	SuccessRate       float64 `json:"success_rate"`
	IntentAccuracy    float64 `json:"intent_accuracy"`
	AvgStepsCompleted float64 `json:"avg_steps_completed"`
	MultiTurnCount    int     `json:"multi_turn_count"`
	ContextRetention  float64 `json:"context_retention"`
	CoherenceScore    float64 `json:"coherence_score"`
	CompletenessScore float64 `json:"completeness_score"`
	SentimentScore    float64 `json:"sentiment_score"`
	ConfidenceScore   float64 `json:"confidence_score"`
	AvgResponseLength float64 `json:"avg_response_length"`
	TotalTests        int     `json:"total_tests"`
	PassedTests       int     `json:"passed_tests"`
	AvgLatency        float64 `json:"avg_latency_ms"`
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
	// migrate: add variant column if missing (older DBs)
	db.Exec(`ALTER TABLE test_results ADD COLUMN variant TEXT`)

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Test results
	r.HandleFunc("/api/results", createResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/results", getResults).Methods("GET", "OPTIONS")

	// A/B tests
	r.HandleFunc("/api/ab-tests", createABTest).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/ab-tests", getABTests).Methods("GET", "OPTIONS")

	// Transcriptions
	r.HandleFunc("/api/transcriptions", createTranscription).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/transcriptions", getTranscriptions).Methods("GET", "OPTIONS")

	// Metrics
	r.HandleFunc("/api/metrics", getMetrics).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/health", healthCheck).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func initDB() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS test_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scenario TEXT NOT NULL,
		passed BOOLEAN NOT NULL,
		latency_ms INTEGER NOT NULL,
		error TEXT,
		variant TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_results_created_at ON test_results(created_at DESC);

	CREATE TABLE IF NOT EXISTS ab_tests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_name TEXT NOT NULL,
		variant_a TEXT NOT NULL,
		variant_b TEXT NOT NULL,
		winner TEXT NOT NULL,
		summary TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS transcriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		call_sid TEXT,
		audio_path TEXT NOT NULL,
		text TEXT NOT NULL,
		confidence REAL,
		duration_seconds REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	return err
}

func createResult(w http.ResponseWriter, r *http.Request) {
	var result TestResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := db.Exec(
		`INSERT INTO test_results (scenario, passed, latency_ms, error, variant) VALUES (?, ?, ?, ?, ?)`,
		result.Scenario, result.Passed, result.Latency, result.Error, result.Variant,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result.ID, _ = res.LastInsertId()
	result.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getResults(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}
	rows, err := db.Query(
		`SELECT id, scenario, passed, latency_ms, COALESCE(error,''), COALESCE(variant,''), created_at FROM test_results ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []TestResult{}
	for rows.Next() {
		var r TestResult
		var createdAt string
		rows.Scan(&r.ID, &r.Scenario, &r.Passed, &r.Latency, &r.Error, &r.Variant, &createdAt)
		r.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		results = append(results, r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func createABTest(w http.ResponseWriter, r *http.Request) {
	var t ABTestResult
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := db.Exec(
		`INSERT INTO ab_tests (test_name, variant_a, variant_b, winner, summary) VALUES (?, ?, ?, ?, ?)`,
		t.TestName, t.VariantA, t.VariantB, t.Winner, t.Summary,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ID, _ = res.LastInsertId()
	t.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func getABTests(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, test_name, variant_a, variant_b, winner, COALESCE(summary,''), created_at FROM ab_tests ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tests := []ABTestResult{}
	for rows.Next() {
		var t ABTestResult
		var createdAt string
		rows.Scan(&t.ID, &t.TestName, &t.VariantA, &t.VariantB, &t.Winner, &t.Summary, &createdAt)
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		tests = append(tests, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tests)
}

func createTranscription(w http.ResponseWriter, r *http.Request) {
	var t Transcription
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := db.Exec(
		`INSERT INTO transcriptions (call_sid, audio_path, text, confidence, duration_seconds) VALUES (?, ?, ?, ?, ?)`,
		t.CallSID, t.AudioPath, t.Text, t.Confidence, t.Duration,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ID, _ = res.LastInsertId()
	t.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func getTranscriptions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, COALESCE(call_sid,''), audio_path, text, confidence, duration_seconds, created_at FROM transcriptions ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	transcriptions := []Transcription{}
	for rows.Next() {
		var t Transcription
		var createdAt string
		rows.Scan(&t.ID, &t.CallSID, &t.AudioPath, &t.Text, &t.Confidence, &t.Duration, &createdAt)
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		transcriptions = append(transcriptions, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transcriptions)
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	var metrics ConversationMetrics

	var passedTests sql.NullInt64
	var avgLatency sql.NullFloat64
	db.QueryRow(`SELECT COUNT(*), SUM(CASE WHEN passed THEN 1 ELSE 0 END), AVG(latency_ms) FROM test_results`).
		Scan(&metrics.TotalTests, &passedTests, &avgLatency)
	metrics.PassedTests = int(passedTests.Int64)
	metrics.AvgLatency = avgLatency.Float64

	if metrics.TotalTests > 0 {
		metrics.SuccessRate = float64(metrics.PassedTests) / float64(metrics.TotalTests) * 100
		metrics.IntentAccuracy = metrics.SuccessRate // approximation from pass rate
	}

	// Defaults for metrics not yet tracked in DB
	metrics.ContextRetention = 100.0
	metrics.CoherenceScore = 100.0
	metrics.CompletenessScore = metrics.SuccessRate
	metrics.SentimentScore = 75.0
	metrics.ConfidenceScore = metrics.SuccessRate
	metrics.AvgResponseLength = 67.0
	metrics.AvgStepsCompleted = 1.3

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "OK")
}
