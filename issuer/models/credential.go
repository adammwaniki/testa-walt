package models

// FarmerCredential represents a farmer's information for credential issuance
type FarmerCredential struct {
	// Section 1: Personal Information
	PersonalIdentificationNumber string `json:"personalIdentificationNumber"`
	Sex                          string `json:"sex"`
	Surname                      string `json:"surname"`
	Forenames                    string `json:"forenames"`
	DateOfBirth                  string `json:"dateBirth"`
	Nationalities                string `json:"nationalities"`

	// Residence Address
	ResidenceStreetNo   string `json:"residenceStreetNo"`
	ResidencePostCode   string `json:"residencePostCode"`
	ResidenceTown       string `json:"residenceTown"`
	ResidenceCountry    string `json:"residenceCountry"`

	// Stay Address
	StayStreetNo   string `json:"stayStreetNo"`
	StayPostCode   string `json:"stayPostCode"`
	StayTown       string `json:"stayTown"`
	StayCountry    string `json:"stayCountry"`

	// Section 2: Legislation Information
	MemberStateLegislation       string `json:"memberStateLegislation"`
	StartingDate                 string `json:"startingDate"`
	EndingDate                   string `json:"endingDate"`
	CertificateForDuration       bool   `json:"certificateForDuration"`
	DeterminationProvisional     bool   `json:"determinationProvisional"`
	TransitionRulesApply         bool   `json:"transitionRulesApply"`

	// Section 3: Activity Type
	PostedEmployedPerson         bool   `json:"postedEmployedPerson"`
	EmployedTwoOrMoreStates      bool   `json:"employedTwoOrMoreStates"`
	PostedSelfEmployedPerson     bool   `json:"postedSelfEmployedPerson"`
	SelfEmployedTwoOrMoreStates  bool   `json:"selfEmployedTwoOrMoreStates"`
	CivilServant                 bool   `json:"civilServant"`
	ContractStaff                bool   `json:"contractStaff"`
	Mariner                      bool   `json:"mariner"`
	EmployedAndSelfEmployed      bool   `json:"employedAndSelfEmployed"`
	CivilAndEmployedSelfEmployed bool   `json:"civilAndEmployedSelfEmployed"`
	FlightCrewMember             bool   `json:"flightCrewMember"`
	Exception                    bool   `json:"exception"`
	ExceptionDescription         string `json:"exceptionDescription"`
	WorkingInStateUnder21        bool   `json:"workingInStateUnder21"`

	// Section 4: Employer/Business Information
	Employee            bool   `json:"employee"`
	SelfEmployedActivity bool  `json:"selfEmployedActivity"`
	BusinessName        string `json:"businessName"`
	BusinessStreetNo    string `json:"businessStreetNo"`
	BusinessPostCode    string `json:"businessPostCode"`
	BusinessTown        string `json:"businessTown"`
	BusinessCountry     string `json:"businessCountry"`

	// Section 5: Work Location
	NoFixedAddress bool `json:"noFixedAddress"`

	// Section 6: Issuing Institution
	InstitutionName     string `json:"institutionName"`
	InstitutionStreetNo string `json:"institutionStreetNo"`
	InstitutionPostCode string `json:"institutionPostCode"`
	InstitutionTown     string `json:"institutionTown"`
	InstitutionCountry  string `json:"institutionCountry"`
	InstitutionID       string `json:"institutionID"`
	OfficeFaxNo         string `json:"officeFaxNo"`
	OfficePhoneNo       string `json:"officePhoneNo"`
	Email               string `json:"email"`
	IssueDate           string `json:"issueDate"`
	Signature           string `json:"signature"`
}

// CredentialRequest represents the request to Walt.id
type CredentialRequest struct {
	IssuerKey                   IssuerKey                 `json:"issuerKey"`
	CredentialConfigurationID   string                    `json:"credentialConfigurationId"`
	CredentialData              CredentialData            `json:"credentialData"`
	Mapping                     Mapping                   `json:"mapping"`
	SelectiveDisclosure         SelectiveDisclosureConfig `json:"selectiveDisclosure"`
	IssuerDid                   string                    `json:"issuerDid"`
}

// IssuerKey represents the issuer's cryptographic key
type IssuerKey struct {
	Type string `json:"type"`
	JWK  JWK    `json:"jwk"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Crv string `json:"crv"`
	D   string `json:"d"`
}

// CredentialData holds the complete credential structure
type CredentialData struct {
	Context         []string          `json:"@context"`
	ID              string            `json:"id"`
	Type            []string          `json:"type"`
	Issuer          string            `json:"issuer"`
	IssuanceDate    string            `json:"issuanceDate"`
	CredentialSubject CredentialSubject `json:"credentialSubject"`
}

// CredentialSubject contains all the farmer's credential data
type CredentialSubject struct {
	ID       string   `json:"id"`
	Section1 Section1 `json:"section1"`
	Section2 Section2 `json:"section2"`
	Section3 Section3 `json:"section3"`
	Section4 Section4 `json:"section4"`
	Section5 Section5 `json:"section5"`
	Section6 Section6 `json:"section6"`
}

// Section1 - Personal Information
type Section1 struct {
	PersonalIdentificationNumber string   `json:"personalIdentificationNumber"`
	Sex                          string   `json:"sex"`
	Surname                      string   `json:"surname"`
	Forenames                    string   `json:"forenames"`
	DateBirth                    string   `json:"dateBirth"`
	Nationalities                []string `json:"nationalities"`
	StateOfResidenceAddress      Address  `json:"stateOfResidenceAddress"`
	StateOfStayAddress           Address  `json:"stateOfStayAddress"`
}

// Section2 - Legislation Information
type Section2 struct {
	MemberStateWhichLegislationApplies string `json:"memberStateWhichLegislationApplies"`
	StartingDate                       string `json:"startingDate"`
	EndingDate                         string `json:"endingDate"`
	CertificateForDurationActivity     bool   `json:"certificateForDurationActivity"`
	DeterminationProvisional           bool   `json:"determinationProvisional"`
	TransitionRulesApplyAsEC8832004    bool   `json:"transitionRulesApplyAsEC8832004"`
}

// Section3 - Activity Type
type Section3 struct {
	PostedEmployedPerson         bool   `json:"postedEmployedPerson"`
	EmployedTwoOrMoreStates      bool   `json:"employedTwoOrMoreStates"`
	PostedSelfEmployedPerson     bool   `json:"postedSelfEmployedPerson"`
	SelfEmployedTwoOrMoreStates  bool   `json:"selfEmployedTwoOrMoreStates"`
	CivilServant                 bool   `json:"civilServant"`
	ContractStaff                bool   `json:"contractStaff"`
	Mariner                      bool   `json:"mariner"`
	EmployedAndSelfEmployed      bool   `json:"employedAndSelfEmployed"`
	CivilAndEmployedSelfEmployed bool   `json:"civilAndEmployedSelfEmployed"`
	FlightCrewMember             bool   `json:"flightCrewMember"`
	Exception                    bool   `json:"exception"`
	ExceptionDescription         string `json:"exceptionDescription"`
	WorkingInStateUnder21        bool   `json:"workingInStateUnder21"`
}

// Section4 - Employer/Business Information
type Section4 struct {
	Employee            bool    `json:"employee"`
	SelfEmployedActivity bool   `json:"selfEmployedActivity"`
	NameBusinessName    string  `json:"nameBusinessName"`
	RegisteredAddress   Address `json:"registeredAddress"`
}

// Section5 - Work Location
type Section5 struct {
	NoFixedAddress bool `json:"noFixedAddress"`
}

// Section6 - Issuing Institution
type Section6 struct {
	Name          string  `json:"name"`
	Address       Address `json:"address"`
	InstitutionID string  `json:"institutionID"`
	OfficeFaxNo   string  `json:"officeFaxNo"`
	OfficePhoneNo string  `json:"officePhoneNo"`
	Email         string  `json:"email"`
	Date          string  `json:"date"`
	Signature     string  `json:"signature"`
}

// Address represents a physical address
type Address struct {
	StreetNo    string `json:"streetNo"`
	PostCode    string `json:"postCode"`
	Town        string `json:"town"`
	CountryCode string `json:"countryCode"`
}

// Mapping for dynamic field mapping
type Mapping struct {
	ID                string                 `json:"id"`
	Issuer            string                 `json:"issuer"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	IssuanceDate      string                 `json:"issuanceDate"`
	Issued            string                 `json:"issued"`
	ValidFrom         string                 `json:"validFrom"`
	ExpirationDate    string                 `json:"expirationDate"`
	CredentialSchema  CredentialSchema       `json:"credentialSchema"`
}

// CredentialSchema defines the schema reference
type CredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// SelectiveDisclosureConfig defines which fields can be selectively disclosed
type SelectiveDisclosureConfig struct {
	Fields SDFields `json:"fields"`
}

// SDFields represents the selective disclosure field structure
type SDFields struct {
	CredentialSubject SDCredentialSubject `json:"credentialSubject"`
}

// SDCredentialSubject defines selective disclosure for credential subject
type SDCredentialSubject struct {
	SD       bool           `json:"sd"`
	Children SDChildrenRoot `json:"children"`
}

// SDChildrenRoot contains all section configurations
type SDChildrenRoot struct {
	Fields SDSections `json:"fields"`
}

// SDSections contains all sections
type SDSections struct {
	Section1 SDSection `json:"section1"`
	Section3 SDSection `json:"section3"`
	Section4 SDSection `json:"section4"`
	Section5 SDSection `json:"section5"`
	Section6 SDSection `json:"section6"`
}

// SDSection represents a section's selective disclosure config
type SDSection struct {
	SD       bool               `json:"sd"`
	Children SDSectionChildren  `json:"children"`
}

// SDSectionChildren contains field configurations
type SDSectionChildren struct {
	Fields map[string]SDField `json:"fields"`
}

// SDField represents a single field's selective disclosure setting
type SDField struct {
	SD bool `json:"sd"`
}