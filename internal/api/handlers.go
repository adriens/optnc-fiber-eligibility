package api

import (
	"encoding/json"
	"net/http"
	"time"

	"opt-nc-eligibilite/internal/cache"
	"opt-nc-eligibilite/internal/models"
	"opt-nc-eligibilite/internal/scraper"
)

// APIErrorResponse représente une réponse d'erreur de l'API
type APIErrorResponse struct {
	Error   string `json:"error" example:"not_found"`
	Message string `json:"message" example:"Numéro introuvable"`
}

// APISuccessResponse représente une réponse de succès de l'API
type APISuccessResponse struct {
	Success   bool                      `json:"success" example:"true"`
	Data      *models.EligibilityResult `json:"data"`
	FromCache bool                      `json:"from_cache,omitempty" example:"false"`
}

// Server représente le serveur API
type Server struct {
	scraper *scraper.Scraper
	cache   *cache.Cache
}

// NewServer crée un nouveau serveur API avec cache
func NewServer(s *scraper.Scraper, c *cache.Cache) *Server {
	return &Server{
		scraper: s,
		cache:   c,
	}
}

// HealthHandler godoc
// @Summary Health check endpoint
// @Description Vérifie que l'API est en ligne et fonctionnelle
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Router /health [get]
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"service":   "opt-nc-eligibility-api",
		"timestamp": time.Now(),
	})
}

// CheckEligibilityHandler godoc
// @Summary Vérifier l'éligibilité à la fibre
// @Description Vérifie si un numéro de téléphone fixe est éligible à la fibre optique OPT-NC
// @Description
// @Description ### Codes de retour
// @Description - **200** : Numéro trouvé et vérifié (éligible ou non à la fibre)
// @Description - **400** : Paramètre manquant ou format invalide
// @Description - **404** : Numéro introuvable dans la base OPT
// @Description - **405** : Méthode HTTP non supportée
// @Tags eligibility
// @Accept json
// @Produce json
// @Param phone query string false "Numéro de téléphone (GET)" example(257364)
// @Param request body object false "Numéro de téléphone (POST)" SchemaExample({"phone_number": "257364"})
// @Success 200 {object} APISuccessResponse "Numéro trouvé et vérifié"
// @Failure 400 {object} APIErrorResponse "Erreur de validation"
// @Failure 404 {object} APIErrorResponse "Numéro introuvable"
// @Failure 405 {object} APIErrorResponse "Méthode non autorisée"
// @Router /api/v1/eligibility [get]
// @Router /api/v1/eligibility [post]
func (s *Server) CheckEligibilityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var phoneNumber string

	// Support both GET and POST
	if r.Method == "GET" {
		phoneNumber = r.URL.Query().Get("phone")
	} else if r.Method == "POST" {
		var req struct {
			PhoneNumber string `json:"phone_number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIErrorResponse{
				Error:   "invalid_request",
				Message: "Invalid JSON body",
			})
			return
		}
		phoneNumber = req.PhoneNumber
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIErrorResponse{
			Error:   "method_not_allowed",
			Message: "Only GET and POST methods are allowed",
		})
		return
	}

	if phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrorResponse{
			Error:   "missing_parameter",
			Message: "Parameter 'phone' or 'phone_number' is required",
		})
		return
	}

	// Check cache first
	if cachedResult, found := s.cache.Get(phoneNumber); found {
		// Remove raw HTML from cached response
		cachedResult.RawHTML = ""
		
		// If cached result is not found (404), return 404
		if !cachedResult.Found {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIErrorResponse{
				Error:   "not_found",
				Message: cachedResult.ErrorMessage,
			})
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(APISuccessResponse{
			Success:   true,
			Data:      cachedResult,
			FromCache: true,
		})
		return
	}

	// Check eligibility via scraping
	result, err := s.scraper.CheckEligibility(phoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Store in cache (even if not found - 404)
	s.cache.Set(phoneNumber, result)

	// If phone number not found in the database, return 404
	if !result.Found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIErrorResponse{
			Error:   "not_found",
			Message: result.ErrorMessage,
		})
		return
	}

	// Remove raw HTML from API response for cleaner output
	result.RawHTML = ""

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APISuccessResponse{
		Success:   true,
		Data:      result,
		FromCache: false,
	})
}
