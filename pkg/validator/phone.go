package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidatePhoneNumber valide et nettoie un numéro de téléphone
// Retourne le numéro nettoyé (6 chiffres) ou une erreur
func ValidatePhoneNumber(phone string) (string, error) {
	// Remove spaces and dots
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Check if only digits
	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return "", fmt.Errorf("le numéro doit contenir uniquement des chiffres")
	}

	// Check length
	if len(cleaned) != 6 {
		return "", fmt.Errorf("le numéro doit contenir exactement 6 chiffres, trouvé: %d", len(cleaned))
	}

	return cleaned, nil
}
