package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ashczar77/mockingjay/backend/handlers"
	"github.com/ashczar77/mockingjay/backend/repository"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./mockingjay.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err := repository.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	repo := repository.New(db)
	h := handlers.New(repo)

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/api/results", h.CreateResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/results", h.GetResults).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/ab-tests", h.CreateABTest).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/ab-tests", h.GetABTests).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/transcriptions", h.CreateTranscription).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/transcriptions", h.GetTranscriptions).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/metrics", h.GetMetrics).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/health", h.HealthCheck).Methods("GET")
	r.HandleFunc("/api/health/status", h.GetHealthStatus).Methods("GET", "OPTIONS")

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
