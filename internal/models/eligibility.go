package models

import "time"

// EligibilityStatus représente le statut d'éligibilité (eligible, non-eligible, unknown)
type EligibilityStatus string

const (
	StatusEligible    EligibilityStatus = "eligible"
	StatusNotEligible EligibilityStatus = "non-eligible"
	StatusUnknown     EligibilityStatus = "unknown"
	StatusNotFound    EligibilityStatus = "not-found"
)

// ADSLEligibility contient les informations d'éligibilité ADSL
type ADSLEligibility struct {
	Status  EligibilityStatus `json:"status"`
	Message string            `json:"message,omitempty"`
}

// FiberEligibility contient les informations d'éligibilité fibre optique (THD)
type FiberEligibility struct {
	Status       EligibilityStatus `json:"status"`
	Available    bool              `json:"available"`
	Message      string            `json:"message,omitempty"`
	Installation string            `json:"installation,omitempty"`
}

// ISPProvider représente un fournisseur d'accès internet
type ISPProvider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Logo string `json:"logo,omitempty"`
}

// EligibilityResult représente le résultat complet du test d'éligibilité
type EligibilityResult struct {
	PhoneNumber  string            `json:"phone_number"`
	CheckedAt    time.Time         `json:"checked_at"`
	Found        bool              `json:"found"`
	ErrorMessage string            `json:"error_message,omitempty"`
	ADSL         *ADSLEligibility  `json:"adsl,omitempty"`
	Fiber        *FiberEligibility `json:"fiber,omitempty"`
	ContactPhone string            `json:"contact_phone,omitempty"`
	ISPProviders []ISPProvider     `json:"isp_providers,omitempty"`
	RawHTML      string            `json:"raw_html,omitempty"`
}

// NewEligibilityResult crée un nouveau résultat d'éligibilité
func NewEligibilityResult(phoneNumber string) *EligibilityResult {
	return &EligibilityResult{
		PhoneNumber: phoneNumber,
		CheckedAt:   time.Now(),
		Found:       false,
	}
}
