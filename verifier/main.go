package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adammwaniki/testa-walt/verifier/handlers"
)

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8081")
	waltIDURL := getEnv("WALTID_VERIFIER_URL", "http://139.59.15.151:7003/openid4vc/verify")

	// Initialize handlers with configuration
	h := handlers.NewHandler(waltIDURL)

	// Routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/verify-credential", h.VerifyCredential)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	addr := ":" + port
	log.Printf("Testa SACCO verifier starting on %s", addr)
	log.Printf("Walt.id Verifier URL: %s", waltIDURL)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}