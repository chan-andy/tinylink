package handlers

import (
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL handles the creation of shortened URLs
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.urlService.CreateShortURL(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response with full short URL
	baseURL := "http://" + c.Request.Host
	response := &models.URLResponse{
		ShortCode: url.ShortCode,
		LongURL:   url.LongURL,
		ShortURL:  baseURL + "/" + url.ShortCode,
		Clicks:    url.Clicks,
		CreatedAt: url.CreatedAt,
		ExpiresAt: url.ExpiresAt,
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectURL handles the redirection of short URLs
func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("code")

	url, err := h.urlService.GetURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Increment click count asynchronously
	go h.urlService.IncrementClicks(shortCode)

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}

// GetURLInfo returns information about a shortened URL
func (h *URLHandler) GetURLInfo(c *gin.Context) {
	shortCode := c.Param("code")

	url, err := h.urlService.GetURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	baseURL := "http://" + c.Request.Host
	response := &models.URLResponse{
		ShortCode: url.ShortCode,
		LongURL:   url.LongURL,
		ShortURL:  baseURL + "/" + url.ShortCode,
		Clicks:    url.Clicks,
		CreatedAt: url.CreatedAt,
		ExpiresAt: url.ExpiresAt,
	}

	c.JSON(http.StatusOK, response)
}
