package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"opt-nc-eligibilite/internal/api"
	"opt-nc-eligibilite/internal/cache"
	"opt-nc-eligibilite/internal/scraper"

	_ "opt-nc-eligibilite/docs" // Swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title OPT-NC Fiber Eligibility API
// @version 1.0
// @description API de vérification d'éligibilité à la fibre optique OPT Nouvelle-Calédonie
// @description
// @description Cette API permet de vérifier l'éligibilité d'un numéro de téléphone fixe à la fibre optique.
// @description
// @description ## Codes HTTP
// @description - **200 OK** : Numéro trouvé et vérifié (éligible ou non)
// @description - **400 Bad Request** : Paramètre manquant ou validation échouée
// @description - **404 Not Found** : Numéro inexistant dans la base OPT
// @description - **405 Method Not Allowed** : Seuls GET et POST acceptés
//
// @contact.name Support API
// @contact.url https://www.opt.nc
//
// @license.name MIT
//
// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	// Initialize scraper
	scraperInstance := scraper.NewScraper(60 * time.Second)

	// Initialize cache with 24h TTL
	cacheInstance := cache.NewCache(24 * time.Hour)

	// Check if running in API mode
	if len(os.Args) > 1 && os.Args[1] == "api" {
		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		// Initialize API server with cache
		server := api.NewServer(scraperInstance, cacheInstance)

		// Setup routes with middleware
		http.HandleFunc("/health", api.Logger(server.HealthHandler))
		http.HandleFunc("/api/v1/eligibility", api.Logger(api.CORS(server.CheckEligibilityHandler)))
		
		// Swagger UI
		http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

		fmt.Printf("🚀 OPT-NC Eligibility API Server started on port %s\n", port)
		fmt.Println("\nEndpoints:")
		fmt.Println("  GET  /health")
		fmt.Println("  GET  /api/v1/eligibility?phone=257364")
		fmt.Println("  POST /api/v1/eligibility {\"phone_number\":\"25.73.64\"}")
		fmt.Println("\n📚 Swagger UI:")
		fmt.Printf("  http://localhost:%s/swagger/\n", port)
		fmt.Println("\nPress Ctrl+C to stop")

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
		return
	}

	// CLI mode
	if len(os.Args) < 2 {
		log.Fatal("Usage:\n  CLI:  go run ./cmd/api <numero_telephone> [--json]\n  API:  go run ./cmd/api api [port]\n\nExemples:\n  go run ./cmd/api 257364\n  go run ./cmd/api 25.73.64 --json\n  go run ./cmd/api api 8080")
	}

	phoneNumber := os.Args[1]
	jsonOutput := len(os.Args) > 2 && os.Args[2] == "--json"

	result, err := scraperInstance.CheckEligibility(phoneNumber)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}

	if jsonOutput {
		result.RawHTML = "" // Remove HTML in JSON output
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("Erreur lors de la génération JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("Résultat d'éligibilité pour le numéro %s:\n\n", phoneNumber)
		if !result.Found {
			fmt.Println("❌ " + result.ErrorMessage)
		} else {
			if result.ADSL != nil {
				fmt.Printf("📡 ADSL: %s\n", result.ADSL.Status)
				if result.ADSL.Message != "" {
					fmt.Printf("   %s\n", result.ADSL.Message)
				}
			}
			if result.Fiber != nil {
				fmt.Printf("🌐 Fibre: %s (disponible: %v)\n", result.Fiber.Status, result.Fiber.Available)
				if result.Fiber.Message != "" {
					fmt.Printf("   %s\n", result.Fiber.Message)
				}
			}
			if result.ContactPhone != "" {
				fmt.Printf("\n📞 Contact: %s\n", result.ContactPhone)
			}
			if len(result.ISPProviders) > 0 {
				fmt.Printf("\n🏢 FAI disponibles: ")
				names := []string{}
				for _, isp := range result.ISPProviders {
					names = append(names, isp.Name)
				}
				fmt.Println(strings.Join(names, ", "))
			}
		}
	}
}
