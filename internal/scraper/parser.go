package scraper

import (
	"strings"

	"opt-nc-eligibilite/internal/models"
)

// ParseEligibilityHTML parse le HTML de réponse et extrait les informations d'éligibilité
func ParseEligibilityHTML(htmlResult string, phoneNumber string) *models.EligibilityResult {
	result := models.NewEligibilityResult(phoneNumber)
	result.RawHTML = htmlResult

	// Check if phone number not found
	if strings.Contains(htmlResult, "Oups, ce numéro est introuvable") {
		result.Found = false
		result.ErrorMessage = "Numéro introuvable. Contactez le 1000 si vous pensez qu'il s'agit d'une erreur."
		return result
	}

	result.Found = true

	// Parse ADSL eligibility
	if strings.Contains(htmlResult, "Eligibilité ADSL") {
		adsl := &models.ADSLEligibility{}
		if strings.Contains(htmlResult, "non éligible") {
			adsl.Status = models.StatusNotEligible
			adsl.Message = "L'offre souscrite sur votre ligne n'est pas compatible avec l'ADSL."
		} else if strings.Contains(htmlResult, "éligible") {
			adsl.Status = models.StatusEligible
		}
		result.ADSL = adsl
	}

	// Parse Fiber eligibility
	if strings.Contains(htmlResult, "Eligibilité THD") {
		fiber := &models.FiberEligibility{}
		if strings.Contains(htmlResult, "Fibre optique pas disponible") {
			fiber.Status = models.StatusNotEligible
			fiber.Available = false
			fiber.Message = "Votre ligne n'est pas encore éligible à la fibre optique. La fibre n'est pas encore disponible à votre adresse."
		} else if strings.Contains(htmlResult, "Fibre optique disponible") {
			fiber.Status = models.StatusEligible
			fiber.Available = true
		}
		result.Fiber = fiber
	}

	// Extract contact phone
	if strings.Contains(htmlResult, "1016") {
		result.ContactPhone = "1016"
	} else if strings.Contains(htmlResult, "1000") {
		result.ContactPhone = "1000"
	}

	// Extract ISP providers
	ispNames := []string{"can'l", "InternetNC", "Lagoon", "MLS", "Nautile"}
	ispURLs := map[string]string{
		"can'l":      "http://www.canl.nc/",
		"InternetNC": "http://www.internetnc.nc/",
		"Lagoon":     "http://www.lagoon.nc/",
		"MLS":        "http://www.mls.nc/",
		"Nautile":    "http://www.nautile.nc/",
	}

	for _, name := range ispNames {
		if strings.Contains(htmlResult, name) {
			result.ISPProviders = append(result.ISPProviders, models.ISPProvider{
				Name: name,
				URL:  ispURLs[name],
			})
		}
	}

	return result
}
