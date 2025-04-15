package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis() {
	// Get Redis URL from environment or use default
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Create Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0, // use default DB
	})

	// Test the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize Redis
	initRedis()

	// Set up Gin
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API Routes
	r.GET("/health", healthCheck)
	r.POST("/api/v1/shorten", createShortURL)
	r.GET("/api/v1/info/:code", getURLInfo)
	r.GET("/:code", redirectURL)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheck(c *gin.Context) {
	// Check Redis connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "Redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// validateAndNormalizeURL checks if the URL is valid and adds http:// if needed
func validateAndNormalizeURL(inputURL string) (string, error) {
	// Add http:// if no protocol is specified
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format")
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: no host found")
	}

	// Validate the host has at least one dot and no spaces
	if !strings.Contains(parsedURL.Host, ".") || strings.Contains(parsedURL.Host, " ") {
		return "", fmt.Errorf("invalid domain name")
	}

	// Check if the URL is reachable
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(inputURL)
	if err != nil {
		return "", fmt.Errorf("unable to reach the website")
	}
	defer resp.Body.Close()

	// Check if the response status code is successful
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("website returned error status: %d", resp.StatusCode)
	}

	return inputURL, nil
}

func createShortURL(c *gin.Context) {
	var request struct {
		LongURL string `json:"long_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and normalize the URL
	normalizedURL, err := validateAndNormalizeURL(request.LongURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique short code
	shortCode := generateShortCode()

	// Store in Redis with 1 year expiration
	err = rdb.Set(ctx, shortCode, normalizedURL, 365*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_code": shortCode,
		"short_url":  "http://" + c.Request.Host + "/" + shortCode,
		"long_url":   normalizedURL,
	})
}

func getURLInfo(c *gin.Context) {
	code := c.Param("code")

	// Get URL from Redis
	longURL, err := rdb.Get(ctx, code).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	// Get visit count
	visits, _ := rdb.Get(ctx, "visits:"+code).Int64()

	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"short_url":  "http://" + c.Request.Host + "/" + code,
		"long_url":   longURL,
		"visits":     visits,
	})
}

func redirectURL(c *gin.Context) {
	code := c.Param("code")

	// Get URL from Redis
	longURL, err := rdb.Get(ctx, code).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	// Increment visit count in a separate goroutine
	go func() {
		rdb.Incr(ctx, "visits:"+code)
	}()

	c.Redirect(http.StatusMovedPermanently, longURL)
}

// generateShortCode creates a random short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	for {
		code := make([]byte, length)
		charsetLength := big.NewInt(int64(len(charset)))

		for i := range code {
			n, err := rand.Int(rand.Reader, charsetLength)
			if err != nil {
				timestamp := time.Now().UnixNano()
				code[i] = charset[timestamp%int64(len(charset))]
				continue
			}
			code[i] = charset[n.Int64()]
		}

		// Check if code exists in Redis
		exists, err := rdb.Exists(ctx, string(code)).Result()
		if err != nil {
			log.Printf("Redis error checking code existence: %v", err)
			continue
		}
		if exists == 0 {
			return string(code)
		}
	}
}
