package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adammwaniki/testa-walt/models"
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

// Home renders the home page with credential selection
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

// ShowPDA1Form renders the PDA1 credential form
func (h *Handler) ShowPDA1Form(w http.ResponseWriter, r *http.Request) {
	err := h.Templates.ExecuteTemplate(w, "pda1-form.html", nil)
	if err != nil {
		log.Printf("Error rendering PDA1 form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ShowFarmerForm renders the Farmer credential form
func (h *Handler) ShowFarmerForm(w http.ResponseWriter, r *http.Request) {
	err := h.Templates.ExecuteTemplate(w, "farmer-form.html", nil)
	if err != nil {
		log.Printf("Error rendering Farmer form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// IssueCredential handles the PDA1 credential issuance request
func (h *Handler) IssueCredential(w http.ResponseWriter, r *http.Request) {
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

	// Extract farmer data from form
	farmer := h.extractFarmerData(r)

	// Build credential request
	credRequest := h.buildCredentialRequest(farmer)

	// Marshal request to JSON
	requestBody, err := json.Marshal(credRequest)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		h.renderError(w, "Failed to create credential request")
		return
	}

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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Walt.id: %v", err)
		h.renderError(w, "Failed to connect to credential service. Please try again.")
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
		h.renderError(w, fmt.Sprintf("Credential service error: %s", string(body)))
		return
	}

	// Response is the credential link
	credentialLink := string(body)

	// Render success response with HTMX
	h.renderSuccess(w, credentialLink, farmer.Forenames+" "+farmer.Surname, "PDA1")
}

// IssueFarmerCredential handles the Farmer credential issuance request
func (h *Handler) IssueFarmerCredential(w http.ResponseWriter, r *http.Request) {
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

	// Extract farmer credential data from form
	farmerCred := &models.SimpleFarmerCredential{
		GivenName:  r.FormValue("given_name"),
		FamilyName: r.FormValue("family_name"),
		FarmName:   r.FormValue("farm_name"),
		FarmType:   r.FormValue("farm_type"),
		LicenseNo:  r.FormValue("license_no"),
		Region:     r.FormValue("region"),
	}

	// Build credential request
	credRequest := h.buildFarmerCredentialRequest(farmerCred)

	// Marshal request to JSON
	requestBody, err := json.Marshal(credRequest)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		h.renderError(w, "Failed to create credential request")
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Use the farmer credential endpoint
	farmerURL := "http://139.59.15.151:7002/openid4vc/jwt/issue"
	req, err := http.NewRequest("POST", farmerURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		h.renderError(w, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Walt.id: %v", err)
		h.renderError(w, "Failed to connect to credential service. Please try again.")
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
		h.renderError(w, fmt.Sprintf("Credential service error: %s", string(body)))
		return
	}

	// Response is the credential link
	credentialLink := string(body)

	// Render success response
	h.renderSuccess(w, credentialLink, farmerCred.GivenName+" "+farmerCred.FamilyName, "Farmer")
}

// extractFarmerData extracts farmer information from the form
func (h *Handler) extractFarmerData(r *http.Request) *models.FarmerCredential {
	return &models.FarmerCredential{
		// Section 1: Personal Information
		PersonalIdentificationNumber: r.FormValue("personalId"),
		Sex:                          r.FormValue("sex"),
		Surname:                      r.FormValue("surname"),
		Forenames:                    r.FormValue("forenames"),
		DateOfBirth:                  r.FormValue("dateOfBirth"),
		Nationalities:                r.FormValue("nationalities"),

		// Residence Address
		ResidenceStreetNo: r.FormValue("residenceStreet"),
		ResidencePostCode: r.FormValue("residencePostCode"),
		ResidenceTown:     r.FormValue("residenceTown"),
		ResidenceCountry:  r.FormValue("residenceCountry"),

		// Stay Address
		StayStreetNo: r.FormValue("stayStreet"),
		StayPostCode: r.FormValue("stayPostCode"),
		StayTown:     r.FormValue("stayTown"),
		StayCountry:  r.FormValue("stayCountry"),

		// Section 2: Legislation
		MemberStateLegislation:   r.FormValue("memberState"),
		StartingDate:             r.FormValue("startingDate"),
		EndingDate:               r.FormValue("endingDate"),
		CertificateForDuration:   r.FormValue("certificateDuration") == "on",
		DeterminationProvisional: r.FormValue("provisional") == "on",
		TransitionRulesApply:     r.FormValue("transitionRules") == "on",

		// Section 3: Activity Type
		PostedEmployedPerson:         r.FormValue("postedEmployed") == "on",
		EmployedTwoOrMoreStates:      r.FormValue("employedMultiState") == "on",
		PostedSelfEmployedPerson:     r.FormValue("postedSelfEmployed") == "on",
		SelfEmployedTwoOrMoreStates:  r.FormValue("selfEmployedMultiState") == "on",
		CivilServant:                 r.FormValue("civilServant") == "on",
		ContractStaff:                r.FormValue("contractStaff") == "on",
		Mariner:                      r.FormValue("mariner") == "on",
		EmployedAndSelfEmployed:      r.FormValue("employedAndSelf") == "on",
		CivilAndEmployedSelfEmployed: r.FormValue("civilAndEmployed") == "on",
		FlightCrewMember:             r.FormValue("flightCrew") == "on",
		Exception:                    r.FormValue("exception") == "on",
		ExceptionDescription:         r.FormValue("exceptionDesc"),
		WorkingInStateUnder21:        r.FormValue("workingUnder21") == "on",

		// Section 4: Business
		Employee:            r.FormValue("employee") == "on",
		SelfEmployedActivity: r.FormValue("selfEmployed") == "on",
		BusinessName:        r.FormValue("businessName"),
		BusinessStreetNo:    r.FormValue("businessStreet"),
		BusinessPostCode:    r.FormValue("businessPostCode"),
		BusinessTown:        r.FormValue("businessTown"),
		BusinessCountry:     r.FormValue("businessCountry"),

		// Section 5: Work Location
		NoFixedAddress: r.FormValue("noFixedAddress") == "on",

		// Section 6: Institution
		InstitutionName:     r.FormValue("institutionName"),
		InstitutionStreetNo: r.FormValue("institutionStreet"),
		InstitutionPostCode: r.FormValue("institutionPostCode"),
		InstitutionTown:     r.FormValue("institutionTown"),
		InstitutionCountry:  r.FormValue("institutionCountry"),
		InstitutionID:       r.FormValue("institutionId"),
		OfficeFaxNo:         r.FormValue("officeFax"),
		OfficePhoneNo:       r.FormValue("officePhone"),
		Email:               r.FormValue("email"),
		IssueDate:           r.FormValue("issueDate"),
		Signature:           r.FormValue("signature"),
	}
}

// buildCredentialRequest builds the complete Walt.id credential request for PDA1
func (h *Handler) buildCredentialRequest(farmer *models.FarmerCredential) *models.CredentialRequest {
	// Parse nationalities (comma-separated)
	nationalities := []string{"BE"}
	if farmer.Nationalities != "" {
		nationalities = strings.Split(farmer.Nationalities, ",")
		for i := range nationalities {
			nationalities[i] = strings.TrimSpace(nationalities[i])
		}
	}

	return &models.CredentialRequest{
		IssuerKey: models.IssuerKey{
			Type: "jwk",
			JWK: models.JWK{
				Kty: "EC",
				X:   "SgfOvOk1TL5yiXhK5Nq7OwKfn_RUkDizlIhAf8qd2wE",
				Y:   "u_y5JZOsw3SrnNPydzJkoaiqb8raSdCNE_nPovt1fNI",
				Crv: "P-256",
				D:   "UqSi2MbJmPczfRmwRDeOJrdivoEy-qk4OEDjFwJYlUI",
			},
		},
		CredentialConfigurationID: "VerifiablePortableDocumentA1_jwt_vc",
		CredentialData: models.CredentialData{
			Context: []string{"https://www.w3.org/2018/credentials/v1"},
			ID:      "https://www.w3.org/2018/credentials/v1",
			Type: []string{
				"VerifiableCredential",
				"VerifiableAttestation",
				"VerifiablePortableDocumentA1",
			},
			Issuer:       "did:ebsi:zf39qHTXaLrr6iy3tQhT3UZ",
			IssuanceDate: "2020-03-10T04:24:12Z",
			CredentialSubject: models.CredentialSubject{
				ID: "did:key:z2dmzD81cgPx8Vki7JbuuMmFYrWPgYoytykUZ3eyqht1j9KbrvQgsKodq2xnfBMYGk99qtunHHQuvvi35kRvbH9SDnue2ZNJqcnaU7yAxeKqEqDX4qFzeKYCj6rdbFnTsf4c8QjFXcgGYS21Db9d2FhHxw9ZEnqt9KPgLsLbQHVAmNNZoz",
				Section1: models.Section1{
					PersonalIdentificationNumber: farmer.PersonalIdentificationNumber,
					Sex:                          farmer.Sex,
					Surname:                      farmer.Surname,
					Forenames:                    farmer.Forenames,
					DateBirth:                    farmer.DateOfBirth,
					Nationalities:                nationalities,
					StateOfResidenceAddress: models.Address{
						StreetNo:    farmer.ResidenceStreetNo,
						PostCode:    farmer.ResidencePostCode,
						Town:        farmer.ResidenceTown,
						CountryCode: farmer.ResidenceCountry,
					},
					StateOfStayAddress: models.Address{
						StreetNo:    farmer.StayStreetNo,
						PostCode:    farmer.StayPostCode,
						Town:        farmer.StayTown,
						CountryCode: farmer.StayCountry,
					},
				},
				Section2: models.Section2{
					MemberStateWhichLegislationApplies: farmer.MemberStateLegislation,
					StartingDate:                       farmer.StartingDate,
					EndingDate:                         farmer.EndingDate,
					CertificateForDurationActivity:     farmer.CertificateForDuration,
					DeterminationProvisional:           farmer.DeterminationProvisional,
					TransitionRulesApplyAsEC8832004:    farmer.TransitionRulesApply,
				},
				Section3: models.Section3{
					PostedEmployedPerson:         farmer.PostedEmployedPerson,
					EmployedTwoOrMoreStates:      farmer.EmployedTwoOrMoreStates,
					PostedSelfEmployedPerson:     farmer.PostedSelfEmployedPerson,
					SelfEmployedTwoOrMoreStates:  farmer.SelfEmployedTwoOrMoreStates,
					CivilServant:                 farmer.CivilServant,
					ContractStaff:                farmer.ContractStaff,
					Mariner:                      farmer.Mariner,
					EmployedAndSelfEmployed:      farmer.EmployedAndSelfEmployed,
					CivilAndEmployedSelfEmployed: farmer.CivilAndEmployedSelfEmployed,
					FlightCrewMember:             farmer.FlightCrewMember,
					Exception:                    farmer.Exception,
					ExceptionDescription:         farmer.ExceptionDescription,
					WorkingInStateUnder21:        farmer.WorkingInStateUnder21,
				},
				Section4: models.Section4{
					Employee:            farmer.Employee,
					SelfEmployedActivity: farmer.SelfEmployedActivity,
					NameBusinessName:    farmer.BusinessName,
					RegisteredAddress: models.Address{
						StreetNo:    farmer.BusinessStreetNo,
						PostCode:    farmer.BusinessPostCode,
						Town:        farmer.BusinessTown,
						CountryCode: farmer.BusinessCountry,
					},
				},
				Section5: models.Section5{
					NoFixedAddress: farmer.NoFixedAddress,
				},
				Section6: models.Section6{
					Name: farmer.InstitutionName,
					Address: models.Address{
						StreetNo:    farmer.InstitutionStreetNo,
						PostCode:    farmer.InstitutionPostCode,
						Town:        farmer.InstitutionTown,
						CountryCode: farmer.InstitutionCountry,
					},
					InstitutionID: farmer.InstitutionID,
					OfficeFaxNo:   farmer.OfficeFaxNo,
					OfficePhoneNo: farmer.OfficePhoneNo,
					Email:         farmer.Email,
					Date:          farmer.IssueDate,
					Signature:     farmer.Signature,
				},
			},
		},
		Mapping: models.Mapping{
			ID:     "<uuid>",
			Issuer: "<issuerDid>",
			CredentialSubject: map[string]interface{}{
				"id": "<subjectDid>",
			},
			IssuanceDate:   "<timestamp-ebsi>",
			Issued:         "<timestamp-ebsi>",
			ValidFrom:      "<timestamp-ebsi>",
			ExpirationDate: "<timestamp-ebsi-in:365d>",
			CredentialSchema: models.CredentialSchema{
				ID:   "https://api-conformance.ebsi.eu/trusted-schemas-registry/v3/schemas/z5qB8tydkn3Xk3VXb15SJ9dAWW6wky1YEoVdGzudWzhcW",
				Type: "FullJsonSchemaValidator2021",
			},
		},
		SelectiveDisclosure: h.buildSelectiveDisclosure(),
		IssuerDid:          "did:ebsi:zf39qHTXaLrr6iy3tQhT3UZ",
	}
}

// buildFarmerCredentialRequest builds the farmer credential request
func (h *Handler) buildFarmerCredentialRequest(farmer *models.SimpleFarmerCredential) *models.SimpleFarmerCredentialRequest {
	return &models.SimpleFarmerCredentialRequest{
		IssuerKey: models.FarmerIssuerKey{
			Type: "jwk",
			JWK: models.FarmerJWK{
				Kty: "OKP",
				D:   "uX8gZ8UPrWQhzOFaA5gmmkfOiIGCE7w1zbRUh9v6xb8",
				Crv: "Ed25519",
				Kid: "ynzK6u55SjO6hFEsW0kBKon_bpvpf5zrr-Q3FNHeAVE",
				X:   "e3CE1EOpYtE_6UyIN58UJwWmGGesV3kZHMVZABIQI3M",
			},
		},
		IssuerDid:                 "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5Iiwia2lkIjoieW56SzZ1NTVTak82aEZFc1cwa0JLb25fYnB2cGY1enJyLVEzRk5IZUFWRSIsIngiOiJlM0NFMUVPcFl0RV82VXlJTjU4VUp3V21HR2VzVjNrWkhNVlpBQklRSTNNIn0",
		CredentialConfigurationID: "FarmerCredential_jwt_vc_json",
		CredentialData: models.SimpleFarmerCredentialData{
			Context: []string{"https://www.w3.org/2018/credentials/v1"},
			ID:      "urn:uuid:{{$uuid}}",
			Type:    []string{"VerifiableCredential", "FarmerCredential"},
			Issuer: models.FarmerIssuer{
				ID:   "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5Iiwia2lkIjoieW56SzZ1NTVTak82aEZFc1cwa0JLb25fYnB2cGY1enJyLVEzRk5IZUFWRSIsIngiOiJlM0NFMUVPcFl0RV82VXlJTjU4VUp3V21HR2VzVjNrWkhNVlpBQklRSTNNIn0",
				Name: "Testa Gava",
			},
			CredentialSubject: models.SimpleFarmerCredentialSubject{
				GivenName:  farmer.GivenName,
				FamilyName: farmer.FamilyName,
				FarmName:   farmer.FarmName,
				FarmType:   farmer.FarmType,
				LicenseNo:  farmer.LicenseNo,
				Region:     farmer.Region,
			},
		},
		Mapping: models.SimpleFarmerMapping{
			ID:             "<uuid>",
			IssuanceDate:   "<timestamp>",
			ExpirationDate: "<timestamp-in:365d>",
		},
	}
}

// buildSelectiveDisclosure creates the selective disclosure configuration
func (h *Handler) buildSelectiveDisclosure() models.SelectiveDisclosureConfig {
	return models.SelectiveDisclosureConfig{
		Fields: models.SDFields{
			CredentialSubject: models.SDCredentialSubject{
				SD: false,
				Children: models.SDChildrenRoot{
					Fields: models.SDSections{
						Section1: h.createSDSection(map[string]bool{
							"personalIdentificationNumber": true,
							"sex":                          true,
							"surname":                      true,
							"forenames":                    true,
							"dateBirth":                    true,
							"nationalities":                true,
							"stateOfResidenceAddress":      true,
							"stateOfStayAddress":           true,
						}),
						Section3: h.createSDSection(map[string]bool{
							"postedEmployedPerson":         true,
							"employedTwoOrMoreStates":      true,
							"postedSelfEmployedPerson":     true,
							"selfEmployedTwoOrMoreStates":  true,
							"civilServant":                 true,
							"contractStaff":                true,
							"mariner":                      true,
							"employedAndSelfEmployed":      true,
							"civilAndEmployedSelfEmployed": true,
							"flightCrewMember":             true,
							"exception":                    true,
							"exceptionDescription":         true,
							"workingInStateUnder21":        true,
						}),
						Section4: h.createSDSection(map[string]bool{
							"employee":            true,
							"selfEmployedActivity": true,
							"nameBusinessName":    true,
							"registeredAddress":   true,
						}),
						Section5: h.createSDSection(map[string]bool{
							"noFixedAddress": true,
						}),
						Section6: h.createSDSection(map[string]bool{
							"name":          true,
							"address":       true,
							"institutionID": true,
							"officeFaxNo":   true,
							"officePhoneNo": true,
							"email":         true,
							"date":          true,
							"signature":     true,
						}),
					},
				},
			},
		},
	}
}

// createSDSection creates a selective disclosure section configuration
func (h *Handler) createSDSection(fields map[string]bool) models.SDSection {
	sdFields := make(map[string]models.SDField)
	for fieldName, sd := range fields {
		sdFields[fieldName] = models.SDField{SD: sd}
	}

	return models.SDSection{
		SD: false,
		Children: models.SDSectionChildren{
			Fields: sdFields,
		},
	}
}

// renderSuccess renders the success message with the credential link
func (h *Handler) renderSuccess(w http.ResponseWriter, credentialLink, name, credentialType string) {
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`
		<div id="result" class="success-message">
			<div class="success-icon">
				<i class="fa-solid fa-check" style="color: #28a745;"></i>
			</div>
			<h3>%s Credential Generated!</h3>
			<p>Credential issued for: <strong>%s</strong></p>
			
			<div class="credential-link-container">
				<label>Credential Link:</label>
				<div class="link-display">
					<textarea id="credentialLink" readonly>%s</textarea>
					<button onclick="copyToClipboard()" class="copy-btn">Copy Link</button>
				</div>
			</div>
			
			<div class="instructions">
				<h4>Next Steps:</h4>
				<ol>
					<li>Copy the credential link above</li>
					<li>Open your digital wallet app (e.g., Walt.id's Wallet)</li>
					<li>Paste the link to import your credential</li>
					<li>Your digital ID is now ready to use!</li>
				</ol>
			</div>
			
			<button onclick="location.href='/'" class="btn-secondary">Issue Another Credential</button>
		</div>

		<script>
		function copyToClipboard() {
			const linkInput = document.getElementById('credentialLink');
			linkInput.select();
			linkInput.setSelectionRange(0, 99999);
			document.execCommand('copy');
			
			const btn = event.target;
			const originalText = btn.textContent;
			btn.textContent = 'Copied!';
			btn.style.backgroundColor = '#27ae60';
			
			setTimeout(() => {
				btn.textContent = originalText;
				btn.style.backgroundColor = '';
			}, 2000);
		}
		</script>
	`, credentialType, name, credentialLink)

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
			<button onclick="location.href='/'" class="btn-secondary">Try Again</button>
		</div>
	`, message)

	w.Write([]byte(html))
}