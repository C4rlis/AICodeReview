package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/carlr/codereviewtool/internal/config"
	"github.com/carlr/codereviewtool/internal/queue"
	"github.com/carlr/codereviewtool/internal/webhook"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Code Review AI - Webhook Listener")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	log.Println("Connected to RabbitMQ")

	// Create webhook handler
	webhookHandler := webhook.NewHandler(cfg.GitHubWebhookSecret, rabbitMQ)

	// Setup HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/webhook/github", webhookHandler.HandleGitHub).Methods("POST")
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.WebhookPort)
	log.Printf("Webhook listener starting on %s", addr)
	log.Printf("GitHub webhook endpoint: http://localhost%s/webhook/github", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
