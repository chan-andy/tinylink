package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"
	"url-shortener/internal/models"
)

type URLService struct {
	// TODO: Add repository interface
}

func NewURLService() *URLService {
	return &URLService{}
}

// GenerateShortCode generates a short code for a given URL
func (s *URLService) GenerateShortCode(url string) string {
	// Create SHA256 hash of the URL
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	// Encode the first 8 bytes of the hash to base64
	encoded := base64.URLEncoding.EncodeToString(hash[:8])

	// Remove any special characters and limit to 8 characters
	shortCode := strings.NewReplacer("+", "", "/", "", "=", "").Replace(encoded)
	return shortCode[:8]
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(req *models.CreateURLRequest) (*models.URL, error) {
	// Validate URL
	if req.LongURL == "" {
		return nil, errors.New("URL cannot be empty")
	}

	shortCode := req.CustomCode
	if shortCode == "" {
		shortCode = s.GenerateShortCode(req.LongURL)
	}

	// Create URL entity
	url := &models.URL{
		LongURL:   req.LongURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
		ExpiresAt: req.ExpiresAt,
	}

	// TODO: Save to repository

	return url, nil
}

// GetURL retrieves a URL by its short code
func (s *URLService) GetURL(shortCode string) (*models.URL, error) {
	// TODO: Implement repository lookup
	return nil, errors.New("not implemented")
}

// IncrementClicks increments the click count for a URL
func (s *URLService) IncrementClicks(shortCode string) error {
	// TODO: Implement click tracking
	return errors.New("not implemented")
} 