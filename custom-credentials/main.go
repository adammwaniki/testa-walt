package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Configuration
const (
	WaltIDBaseURL = "http://139.59.15.151:7002"
	IssuerDID     = "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5Iiwia2lkIjoicUZpVDJBeXVYNnVBZWY0OVE5Q19FdWxUT3VMNHZxTG1OZTYyR2NQNkZwbyIsIngiOiIzZVFIdHhMWURQSWtRT0s4MnRIcS1BWi1CVU1BX3U5XzFKMjdJVXo5TUdnIn0"
)

// FarmerCredentialRequest represents the API request structure
type FarmerCredentialRequest struct {
	FarmerType            string                 `json:"farmerType"`
	FirstName             string                 `json:"firstName"`
	FamilyName            string                 `json:"familyName"`
	PhoneNumber           string                 `json:"phoneNumber"`
	BirthDate             string                 `json:"birthDate,omitempty"`
	County                string                 `json:"county"`
	SubCounty             string                 `json:"subCounty"`
	FarmSize              *FarmSize              `json:"farmSize,omitempty"`
	DairySpecifics        *DairySpecifics        `json:"dairySpecifics,omitempty"`
	PoultrySpecifics      *PoultrySpecifics      `json:"poultrySpecifics,omitempty"`
	HorticultureSpecifics *HorticultureSpecifics `json:"horticultureSpecifics,omitempty"`
	AquacultureSpecifics  *AquacultureSpecifics  `json:"aquacultureSpecifics,omitempty"`
}

type FarmSize struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type DairySpecifics struct {
	CattleBreeds           []string          `json:"cattleBreeds"`
	NumberOfCattle         int               `json:"numberOfCattle"`
	MilkingCows            int               `json:"milkingCows"`
	AverageDailyProduction *ProductionMetric `json:"averageDailyProduction"`
	KDBNumber              string            `json:"kdbNumber,omitempty"`
}

type ProductionMetric struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type PoultrySpecifics struct {
	FarmingType            string             `json:"farmingType"`
	BirdPopulation         int                `json:"birdPopulation"`
	HousingType            string             `json:"housingType"`
	ProductionCapacity     *PoultryProduction `json:"productionCapacity,omitempty"`
	BiosecurityLevel       string             `json:"biosecurityLevel,omitempty"`
	VeterinaryRegistration string             `json:"veterinaryRegistration,omitempty"`
}

type PoultryProduction struct {
	EggsPerDay   int     `json:"eggsPerDay,omitempty"`
	MeatPerCycle float64 `json:"meatPerCycle,omitempty"`
}

type HorticultureSpecifics struct {
	Crops            []string `json:"crops"`
	FarmingMethod    string   `json:"farmingMethod"`
	GreenhouseCount  int      `json:"greenhouseCount,omitempty"`
	IrrigationSystem string   `json:"irrigationSystem"`
	Certifications   []string `json:"certifications,omitempty"`
	ExportMarket     bool     `json:"exportMarket,omitempty"`
	HCDNumber        string   `json:"hcdNumber,omitempty"`
}

type AquacultureSpecifics struct {
	Species                []string        `json:"species"`
	FarmingSystem          string          `json:"farmingSystem"`
	NumberOfPonds          int             `json:"numberOfPonds,omitempty"`
	WaterSource            string          `json:"waterSource"`
	ProductionCycle        *AquaProduction `json:"productionCycle,omitempty"`
	FeedingType            string          `json:"feedingType,omitempty"`
	FishDepartmentPermit   string          `json:"fishDepartmentPermit,omitempty"`
	WaterQualityManagement bool            `json:"waterQualityManagement,omitempty"`
}

type AquaProduction struct {
	CyclesPerYear int     `json:"cyclesPerYear"`
	FishPerCycle  int     `json:"fishPerCycle"`
	KgPerCycle    float64 `json:"kgPerCycle"`
}

// CredentialService handles credential operations
type CredentialService struct {
	client *http.Client
}

func NewCredentialService() *CredentialService {
	return &CredentialService{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// IssueCredentialHandler handles POST /credentials/issue
func (s *CredentialService) IssueCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req FarmerCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate required fields
	if err := s.validateRequest(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Build credential
	credential, err := s.buildCredential(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to build credential", err)
		return
	}

	// Issue via walt.id
	issuedCredential, err := s.issueToWaltID(credential)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to issue credential", err)
		return
	}

	// Return success response
	respondSuccess(w, http.StatusOK, issuedCredential)
}

// GetCredentialSchemaHandler handles GET /credentials/schemas/{type}
func (s *CredentialService) GetCredentialSchemaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	farmerType := vars["type"]

	schema, err := s.getSchema(farmerType)
	if err != nil {
		respondError(w, http.StatusNotFound, "Schema not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
}

// VerifyCredentialHandler handles POST /credentials/verify
func (s *CredentialService) VerifyCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// Read credential JWT from body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request", err)
		return
	}

	// Verify with walt.id
	verified, result, err := s.verifyWithWaltID(string(body))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Verification failed", err)
		return
	}

	response := map[string]interface{}{
		"verified": verified,
		"result":   result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListCredentialTypesHandler handles GET /credentials/types
func (s *CredentialService) ListCredentialTypesHandler(w http.ResponseWriter, r *http.Request) {
	types := map[string]interface{}{
		"types": []map[string]string{
			{
				"type":        "dairy",
				"name":        "Dairy Farmer",
				"icon":        "🐄",
				"description": "Cattle rearing and milk production",
			},
			{
				"type":        "poultry",
				"name":        "Poultry Farmer",
				"icon":        "🐔",
				"description": "Chicken farming for eggs and meat",
			},
			{
				"type":        "horticulture",
				"name":        "Horticulture Farmer",
				"icon":        "🥬",
				"description": "Vegetables, fruits, and flowers",
			},
			{
				"type":        "aquaculture",
				"name":        "Aquaculture Farmer",
				"icon":        "🐟",
				"description": "Fish and aquatic organism farming",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

// validateRequest validates the incoming request
func (s *CredentialService) validateRequest(req *FarmerCredentialRequest) error {
	if req.FarmerType == "" {
		return fmt.Errorf("farmerType is required")
	}
	if req.FirstName == "" {
		return fmt.Errorf("firstName is required")
	}
	if req.County == "" {
		return fmt.Errorf("county is required")
	}

	// Validate type-specific fields
	switch req.FarmerType {
	case "dairy":
		if req.DairySpecifics == nil {
			return fmt.Errorf("dairySpecifics is required for dairy farmers")
		}
	case "poultry":
		if req.PoultrySpecifics == nil {
			return fmt.Errorf("poultrySpecifics is required for poultry farmers")
		}
	case "horticulture":
		if req.HorticultureSpecifics == nil {
			return fmt.Errorf("horticultureSpecifics is required for horticulture farmers")
		}
	case "aquaculture":
		if req.AquacultureSpecifics == nil {
			return fmt.Errorf("aquacultureSpecifics is required for aquaculture farmers")
		}
	default:
		return fmt.Errorf("invalid farmerType: %s", req.FarmerType)
	}

	return nil
}

// buildCredential constructs the W3C credential
func (s *CredentialService) buildCredential(req *FarmerCredentialRequest) (map[string]interface{}, error) {
	// Generate holder DID (in production, this comes from farmer's wallet)
	holderDID := fmt.Sprintf("did:key:farmer_%d", time.Now().UnixNano())

	credential := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://w3id.org/security/suites/jws-2020/v1",
		},
		"type":         []string{"VerifiableCredential", "FarmerCredential", fmt.Sprintf("%sFarmerCredential", capitalize(req.FarmerType))},
		"issuer":       IssuerDID,
		"issuanceDate": time.Now().Format(time.RFC3339),
		"credentialSubject": map[string]interface{}{
			"id":               holderDID,
			"farmerType":       req.FarmerType,
			"firstName":        req.FirstName,
			"familyName":       req.FamilyName,
			"phoneNumber":      req.PhoneNumber,
			"birthDate":        req.BirthDate,
			"county":           req.County,
			"subCounty":        req.SubCounty,
			"farmSize":         req.FarmSize,
			"registrationDate": time.Now().Format(time.RFC3339),
		},
	}

	// Add type-specific data
	subject := credential["credentialSubject"].(map[string]interface{})
	switch req.FarmerType {
	case "dairy":
		subject["dairySpecifics"] = req.DairySpecifics
	case "poultry":
		subject["poultrySpecifics"] = req.PoultrySpecifics
	case "horticulture":
		subject["horticultureSpecifics"] = req.HorticultureSpecifics
	case "aquaculture":
		subject["aquacultureSpecifics"] = req.AquacultureSpecifics
	}

	return credential, nil
}

// issueToWaltID sends credential to walt.id for signing
func (s *CredentialService) issueToWaltID(credential map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(credential)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credential: %w", err)
	}

	url := fmt.Sprintf("%s/openid4vc/jwt/issue", WaltIDBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("walt.id returned status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// verifyWithWaltID verifies a credential with walt.id
func (s *CredentialService) verifyWithWaltID(credentialJWT string) (bool, map[string]interface{}, error) {
	url := fmt.Sprintf("%s/openid4vc/verify", WaltIDBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(credentialJWT))
	if err != nil {
		return false, nil, err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, nil, err
	}

	verified := resp.StatusCode == http.StatusOK
	return verified, result, nil
}

// getSchema returns the schema for a farmer type
func (s *CredentialService) getSchema(farmerType string) (map[string]interface{}, error) {
	// This would load from files in production
	schemas := map[string]map[string]interface{}{
		"dairy": {
			"type":        "dairy",
			"name":        "Dairy Farmer Credential",
			"description": "Credential for dairy farmers in Kenya",
			"fields": []map[string]interface{}{
				{"name": "cattleBreeds", "type": "array", "required": true},
				{"name": "numberOfCattle", "type": "integer", "required": true},
				{"name": "milkingCows", "type": "integer", "required": true},
				{"name": "averageDailyProduction", "type": "object", "required": true},
			},
		},
		"poultry": {
			"type":        "poultry",
			"name":        "Poultry Farmer Credential",
			"description": "Credential for poultry farmers in Kenya",
			"fields": []map[string]interface{}{
				{"name": "farmingType", "type": "string", "required": true},
				{"name": "birdPopulation", "type": "integer", "required": true},
				{"name": "housingType", "type": "string", "required": true},
			},
		},
		"horticulture": {
			"type":        "horticulture",
			"name":        "Horticulture Farmer Credential",
			"description": "Credential for horticulture farmers in Kenya",
			"fields": []map[string]interface{}{
				{"name": "crops", "type": "array", "required": true},
				{"name": "farmingMethod", "type": "string", "required": true},
				{"name": "irrigationSystem", "type": "string", "required": true},
			},
		},
		"aquaculture": {
			"type":        "aquaculture",
			"name":        "Aquaculture Farmer Credential",
			"description": "Credential for aquaculture farmers in Kenya",
			"fields": []map[string]interface{}{
				{"name": "species", "type": "array", "required": true},
				{"name": "farmingSystem", "type": "string", "required": true},
				{"name": "waterSource", "type": "string", "required": true},
			},
		},
	}

	schema, ok := schemas[farmerType]
	if !ok {
		return nil, fmt.Errorf("schema not found for type: %s", farmerType)
	}

	return schema, nil
}

// Helper functions
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func respondError(w http.ResponseWriter, code int, message string, err error) {
	log.Printf("Error: %s - %v", message, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
		"details": err.Error(),
	})
}

func respondSuccess(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// CORS middleware with enhanced configuration
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Allow all origins for now (restrict in production)
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7115"
	}

	service := NewCredentialService()
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "farmer-credential-service",
			"version": "1.0.0",
		})
	}).Methods("GET")

	// Credential endpoints
	r.HandleFunc("/credentials/issue", service.IssueCredentialHandler).Methods("POST")
	r.HandleFunc("/credentials/verify", service.VerifyCredentialHandler).Methods("POST")
	r.HandleFunc("/credentials/types", service.ListCredentialTypesHandler).Methods("GET")
	r.HandleFunc("/credentials/schemas/{type}", service.GetCredentialSchemaHandler).Methods("GET")

	// Apply middleware
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	// Get host interface
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // Bind to all interfaces by default
	}

	log.Printf("Farmer Credential Service")
	log.Printf("Starting on %s:%s", host, port)
	log.Printf("Health: http://localhost:%s/health", port)
	log.Printf("Issue: POST http://localhost:%s/credentials/issue", port)
	log.Printf("Verify: POST http://localhost:%s/credentials/verify", port)
	log.Printf("Types: GET http://localhost:%s/credentials/types", port)
	log.Printf("Schema: GET http://localhost:%s/credentials/schemas/{type}", port)
	log.Printf("Accessible from network on http://<your-ip>:%s", port)

	addr := host + ":" + port
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}