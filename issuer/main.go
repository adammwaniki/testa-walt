package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adammwaniki/testa-walt/handlers"
)

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8082")
	waltIDURL := getEnv("WALTID_ISSUER_URL", "http://139.59.15.151:7002/openid4vc/sdjwt/issue")

	// Initialize handlers with configuration
	h := handlers.NewHandler(waltIDURL)

	// Routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/issue-credential", h.IssueCredential)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	addr := ":" + port
	log.Printf("Testa Gava server starting on %s", addr)
	log.Printf("Walt.id URL: %s", waltIDURL)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}