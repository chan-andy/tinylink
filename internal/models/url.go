package models

import "time"

// URL represents a shortened URL in the system
type URL struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	LongURL   string    `json:"long_url" gorm:"not null"`
	ShortCode string    `json:"short_code" gorm:"unique;not null"`
	Clicks    uint64    `json:"clicks" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// CreateURLRequest represents the request to create a shortened URL
type CreateURLRequest struct {
	LongURL     string    `json:"long_url" binding:"required,url"`
	CustomCode  string    `json:"custom_code,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// URLResponse represents the response for URL-related operations
type URLResponse struct {
	ShortCode  string    `json:"short_code"`
	LongURL    string    `json:"long_url"`
	ShortURL   string    `json:"short_url"`
	Clicks     uint64    `json:"clicks"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
} 