package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type CallRequest struct {
	Text string `json:"text"`
}

type CallResponse struct {
	Text    string `json:"text"`
	Intent  string `json:"intent,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/call", handleCall)
	http.HandleFunc("/call-bad", handleCallBad) // intentionally broken agent for demos
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/voice", handleVoice)

	port := "9000"
	log.Printf("🤖 Example Voice AI Server")
	log.Printf("Listening on http://localhost:%s", port)
	log.Printf("Endpoint: http://localhost:%s/call", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Simple intent detection
	text := strings.ToLower(req.Text)
	resp := CallResponse{Success: true}

	switch {
	case strings.Contains(text, "hello") || strings.Contains(text, "hi"):
		resp.Text = "Hello! How can I help you today?"
		resp.Intent = "greeting"

	case strings.Contains(text, "book") || strings.Contains(text, "appointment"):
		resp.Text = "I'd be happy to help you book an appointment. What date works for you?"
		resp.Intent = "booking_intent"

	case strings.Contains(text, "tomorrow") || strings.Contains(text, "7pm") || strings.Contains(text, "pm"):
		resp.Text = "Great! I've booked your appointment for tomorrow at 7pm. You'll receive a confirmation email shortly."
		resp.Intent = "confirmation"

	case strings.Contains(text, "cancel"):
		resp.Text = "I understand you'd like to cancel. Can you provide your booking reference?"
		resp.Intent = "cancellation"

	case strings.Contains(text, "hours") || strings.Contains(text, "open"):
		resp.Text = "We're open Monday to Friday, 9am to 6pm, and Saturday 10am to 4pm."
		resp.Intent = "hours_response"

	case strings.Contains(text, "help"):
		resp.Text = "I can help you book appointments, check our hours, or answer questions. What would you like to do?"
		resp.Intent = "help"

	default:
		resp.Text = "I'm not sure I understood that. Could you rephrase?"
		resp.Intent = "unknown"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// handleCallBad simulates a broken agent that misclassifies intents.
// Use endpoint http://localhost:9000/call-bad in mockingjay.yaml to demo failures.
func handleCallBad(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	time.Sleep(300 * time.Millisecond) // simulate slow response
	// Always returns wrong intent to demonstrate MockingJay catching failures
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CallResponse{
		Text:    "I'm sorry, I didn't understand that.",
		Intent:  "unknown",
		Success: true,
	})
}

func handleVoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	w.Header().Set("ngrok-skip-browser-warning", "true")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Say voice="Polly.Joanna">Hello! This is a MockingJay test call. Your voice AI agent is working correctly. Goodbye!</Say>
</Response>`)
}
