package main

import (
	"log"
	"net/http"

	"github.com/adammwaniki/testa-walt/handlers"
)

func main() {
	// Walt.id URL for PDA1 credentials
	waltIDURL := "http://139.59.15.151:7002/openid4vc/sdjwt/issue"
	
	// Create handler
	h := handlers.NewHandler(waltIDURL)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/form/pda1", h.ShowPDA1Form)
	http.HandleFunc("/form/farmer", h.ShowFarmerForm)
	http.HandleFunc("/issue-credential", h.IssueCredential)
	http.HandleFunc("/issue-farmer-credential", h.IssueFarmerCredential)

	// Start server
	port := ":8082"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}