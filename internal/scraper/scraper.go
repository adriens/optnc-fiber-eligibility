package scraper

import (
	"context"
	"time"

	"opt-nc-eligibilite/internal/models"
	"opt-nc-eligibilite/pkg/validator"

	"github.com/chromedp/chromedp"
)

// Scraper gère le scraping du site OPT-NC
type Scraper struct {
	timeout time.Duration
}

// NewScraper crée un nouveau scraper avec un timeout
func NewScraper(timeout time.Duration) *Scraper {
	return &Scraper{
		timeout: timeout,
	}
}

// CheckEligibility vérifie l'éligibilité d'un numéro de téléphone
func (s *Scraper) CheckEligibility(phoneNumber string) (*models.EligibilityResult, error) {
	// Validate phone number
	cleanedPhone, err := validator.ValidatePhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Chrome options for Docker/container environment
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new browser context
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Navigate to the page and fill the form
	var htmlResult string
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://www.opt.nc/particuliers/telephonie-fixe/fibre-optique"),
		chromedp.WaitVisible(`#edit-phone-number`, chromedp.ByQuery),
		chromedp.SendKeys(`#edit-phone-number`, cleanedPhone, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#edit-gdpr`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#edit-submit`, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.InnerHTML(`#ajax-opt-eligibility-result`, &htmlResult, chromedp.ByQuery),
	)

	if err != nil {
		return nil, err
	}

	result := ParseEligibilityHTML(htmlResult, cleanedPhone)
	return result, nil
}
