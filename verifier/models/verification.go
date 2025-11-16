package models

// VerificationRequest represents the request to Walt.id verifier
type VerificationRequest struct {
	VcPolicies         []string            `json:"vc_policies"`
	RequestCredentials []RequestCredential `json:"request_credentials"`
}

// RequestCredential defines what credentials to verify
type RequestCredential struct {
	Format          string          `json:"format"`
	InputDescriptor InputDescriptor `json:"input_descriptor"`
}

// InputDescriptor defines constraints for credential verification
type InputDescriptor struct {
	ID          string      `json:"id"`
	Constraints Constraints `json:"constraints"`
}

// Constraints defines field constraints
type Constraints struct {
	Fields []Field `json:"fields"`
}

// Field defines a field constraint
type Field struct {
	Path   []string `json:"path"`
	Filter Filter   `json:"filter"`
}

// Filter defines the filter criteria
type Filter struct {
	Contains FilterContains `json:"contains"`
	Type     string         `json:"type"`
}

// FilterContains defines what the field should contain
type FilterContains struct {
	Const string `json:"const"`
}

// VerificationResponse represents the response from Walt.id
type VerificationResponse struct {
	VerificationID string `json:"verification_id"`
	URL            string `json:"url"`
	Status         string `json:"status"`
}

// VerificationOptions represents user-selected verification options
type VerificationOptions struct {
	CredentialType     string
	CheckSignature     bool
	CheckExpiration    bool
	CheckNotBefore     bool
	CheckRevokedStatus bool
}