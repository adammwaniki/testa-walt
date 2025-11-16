package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/adammwaniki/testa-walt/verifier/models"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	WaltIDURL string
	Templates *template.Template
}

// NewHandler creates a new handler with dependencies
func NewHandler(waltIDURL string) *Handler {
	// Parse templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}

	return &Handler{
		WaltIDURL: waltIDURL,
		Templates: templates,
	}
}

// Home renders the home page with verification options
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := h.Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// VerifyCredential handles the credential verification request
func (h *Handler) VerifyCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		h.renderError(w, "Failed to parse form data")
		return
	}

	// Extract verification options from form
	options := h.extractVerificationOptions(r)

	// Build verification request
	verifyRequest := h.buildVerificationRequest(options)

	// Marshal request to JSON
	requestBody, err := json.Marshal(verifyRequest)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		h.renderError(w, "Failed to create verification request")
		return
	}

	log.Printf("Sending verification request: %s", string(requestBody))

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Walt.id
	req, err := http.NewRequest("POST", h.WaltIDURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		h.renderError(w, "Failed to create request")
		return
	}

	// Set headers as per curl example
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("authorizeBaseUrl", "openid4vp://authorize")
	req.Header.Set("responseMode", "direct_post")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Walt.id: %v", err)
		h.renderError(w, "Failed to connect to verification service. Please try again.")
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		h.renderError(w, "Failed to read response")
		return
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Walt.id returned status %d: %s", resp.StatusCode, string(body))
		h.renderError(w, fmt.Sprintf("Verification service error (status %d): %s", resp.StatusCode, string(body)))
		return
	}

	// Response is the verification link
	verificationLink := string(body)
	log.Printf("Verification link received: %s", verificationLink)

	// Render success response with HTMX
	h.renderSuccess(w, verificationLink, options)
}

// extractVerificationOptions extracts verification settings from the form
func (h *Handler) extractVerificationOptions(r *http.Request) *models.VerificationOptions {
	return &models.VerificationOptions{
		CredentialType:     r.FormValue("credentialType"),
		CheckSignature:     r.FormValue("checkSignature") == "on",
		CheckExpiration:    r.FormValue("checkExpiration") == "on",
		CheckNotBefore:     r.FormValue("checkNotBefore") == "on",
		CheckRevokedStatus: r.FormValue("checkRevokedStatus") == "on",
	}
}

// buildVerificationRequest builds the complete Walt.id verification request
func (h *Handler) buildVerificationRequest(options *models.VerificationOptions) *models.VerificationRequest {
	// Build policies array based on user selection
	policies := []string{}
	
	if options.CheckSignature {
		policies = append(policies, "signature")
	}
	if options.CheckExpiration {
		policies = append(policies, "expired")
	}
	if options.CheckNotBefore {
		policies = append(policies, "not-before")
	}
	if options.CheckRevokedStatus {
		policies = append(policies, "revoked-status-list")
	}

	// If no policies selected, use all by default
	if len(policies) == 0 {
		policies = []string{"signature", "expired", "not-before", "revoked-status-list"}
	}

	// Determine credential type (default to VerifiablePortableDocumentA1)
	credentialType := options.CredentialType
	if credentialType == "" {
		credentialType = "VerifiablePortableDocumentA1"
	}

	return &models.VerificationRequest{
		VcPolicies: policies,
		RequestCredentials: []models.RequestCredential{
			{
				Format: "jwt_vc",
				InputDescriptor: models.InputDescriptor{
					ID: "e3d700aa-0988-4eb6-b9c9-e00f4b27f1d8",
					Constraints: models.Constraints{
						Fields: []models.Field{
							{
								Path: []string{"$.vc.type"},
								Filter: models.Filter{
									Contains: models.FilterContains{
										Const: credentialType,
									},
									Type: "array",
								},
							},
						},
					},
				},
			},
		},
	}
}

// renderSuccess renders the success message with the verification link
func (h *Handler) renderSuccess(w http.ResponseWriter, verificationLink string, options *models.VerificationOptions) {
	w.Header().Set("Content-Type", "text/html")
	
	// Build policies list for display
	policiesChecked := []string{}
	if options.CheckSignature {
		policiesChecked = append(policiesChecked, "Signature")
	}
	if options.CheckExpiration {
		policiesChecked = append(policiesChecked, "Expiration")
	}
	if options.CheckNotBefore {
		policiesChecked = append(policiesChecked, "Not-Before")
	}
	if options.CheckRevokedStatus {
		policiesChecked = append(policiesChecked, "Revocation Status")
	}

	policiesHTML := ""
	if len(policiesChecked) > 0 {
		policiesHTML = "<ul>"
		for _, policy := range policiesChecked {
			policiesHTML += fmt.Sprintf("<li>%s</li>", policy)
		}
		policiesHTML += "</ul>"
	} else {
		policiesHTML = "<p>All verification policies (default)</p>"
	}

	credType := options.CredentialType
	if credType == "" {
		credType = "VerifiablePortableDocumentA1 (default)"
	}

	html := fmt.Sprintf(`
		<div id="result" class="success-message">
			<div class="success-icon">✓</div>
			<h3>Verification Request Generated!</h3>
			<p>A credential verification request has been successfully created.</p>
			
			<div class="verification-details">
				<h4>Verification Configuration:</h4>
				<div class="detail-item">
					<strong>Credential Type:</strong> %s
				</div>
				<div class="detail-item">
					<strong>Policies Checked:</strong>
					%s
				</div>
			</div>
			
			<div class="verification-link-container">
				<label>Verification Link:</label>
				<div class="link-display">
					<textarea id="verificationLink" readonly>%s</textarea>
					<button onclick="copyToClipboard()" class="copy-btn">Copy Link</button>
				</div>
			</div>
			
			<div class="instructions">
				<h4>Next Steps:</h4>
				<ol>
					<li>Copy the verification link above</li>
					<li>Share this link with the credential holder</li>
					<li>Open the link in your digital wallet (e.g., Walt.id's Wallet)</li>
					<li>The wallet presents the credential for verification</li>
					<li>You'll receive the verification response</li>
				</ol>
			</div>
			
			<div class="info-box">
				<p><strong>Note:</strong> This link initiates an OpenID4VP (Verifiable Presentation) flow. 
				The credential holder's wallet will use this to present their credential securely.</p>
			</div>
			
			<button onclick="location.reload()" class="btn-secondary">Verify Another Credential</button>
		</div>

		<script>
		function copyToClipboard() {
			const linkInput = document.getElementById('verificationLink');
			linkInput.select();
			linkInput.setSelectionRange(0, 99999);
			document.execCommand('copy');
			
			const btn = event.target;
			const originalText = btn.textContent;
			btn.textContent = 'Copied!';
			btn.style.backgroundColor = '#2563eb';
			
			setTimeout(() => {
				btn.textContent = originalText;
				btn.style.backgroundColor = '';
			}, 2000);
		}
		</script>
	`, credType, policiesHTML, verificationLink)

	w.Write([]byte(html))
}

// renderError renders an error message
func (h *Handler) renderError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`
		<div id="result" class="error-message">
			<div class="error-icon">✗</div>
			<h3>Error</h3>
			<p>%s</p>
			<button onclick="location.reload()" class="btn-secondary">Try Again</button>
		</div>
	`, message)

	w.Write([]byte(html))
}