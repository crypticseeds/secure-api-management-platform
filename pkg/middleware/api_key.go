package middleware

import (
	"apisecurityplatform/pkg/database"
	"apisecurityplatform/pkg/models"
	"apisecurityplatform/pkg/observability"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		var storedKey models.APIKey
		// Get all API keys and compare with bcrypt
		var apiKeys []models.APIKey
		if err := database.DB.Find(&apiKeys).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate API key"})
			c.Abort()
			return
		}

		// Find matching API key
		keyFound := false
		for _, key := range apiKeys {
			if err := bcrypt.CompareHashAndPassword([]byte(key.Key), []byte(apiKey)); err == nil {
				storedKey = key
				keyFound = true
				break
			}
		}

		if !keyFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Update last used timestamp
		now := time.Now().Unix()
		storedKey.LastUsedAt = &now
		database.DB.Save(&storedKey)

		// Increment API key usage metric
		observability.APIKeyUsage.WithLabelValues(fmt.Sprintf("%d", storedKey.ID)).Inc()

		c.Set("api_key_id", storedKey.ID)
		c.Set("user_id", storedKey.UserID)
		c.Next()
	}
}
