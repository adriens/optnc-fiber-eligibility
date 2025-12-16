package cache

import (
	"sync"
	"time"

	"opt-nc-eligibilite/internal/models"
)

// CacheEntry représente une entrée dans le cache avec son expiration
type CacheEntry struct {
	Result    *models.EligibilityResult
	ExpiresAt time.Time
}

// Cache gère le cache des résultats d'éligibilité
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewCache crée un nouveau cache avec une durée de vie par défaut
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
	// Start cleanup goroutine
	go c.cleanup()
	return c
}

// Get récupère une entrée du cache si elle existe et n'est pas expirée
func (c *Cache) Get(phoneNumber string) (*models.EligibilityResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[phoneNumber]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Result, true
}

// Set ajoute une entrée dans le cache avec le TTL configuré
func (c *Cache) Set(phoneNumber string, result *models.EligibilityResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[phoneNumber] = &CacheEntry{
		Result:    result,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// cleanup supprime périodiquement les entrées expirées
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// Stats retourne les statistiques du cache
func (c *Cache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := len(c.entries)
	return map[string]interface{}{
		"entries": count,
		"ttl":     c.ttl.String(),
	}
}
