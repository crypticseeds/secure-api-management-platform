package handlers

import (
	"apisecurityplatform/pkg/database"
	"apisecurityplatform/pkg/models"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateAPIKeyInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// GenerateAPIKey generates a random 32-character API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 24) // 24 bytes will give us 32 characters in base64
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// @Summary Create API key
// @Description Generate a new API key for the authenticated user
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key body CreateAPIKeyInput true "API key details"
// @Success 201 {object} map[string]interface{} "API key created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/api-keys [post]
func CreateAPIKey(c *gin.Context) {
	var input CreateAPIKeyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("user_id")

	// Generate API key
	apiKey, err := GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Hash the API key before storing
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash API key"})
		return
	}

	// Create API key record
	apiKeyRecord := models.APIKey{
		UserID:      userID.(uint),
		Name:        input.Name,
		Key:         string(hashedKey),
		Description: input.Description,
	}

	if err := database.DB.Create(&apiKeyRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save API key"})
		return
	}

	// Return the unhashed key to the user (this is the only time they'll see it)
	c.JSON(http.StatusCreated, gin.H{
		"message": "API key created successfully",
		"api_key": apiKey,
		"id":      apiKeyRecord.ID,
		"name":    apiKeyRecord.Name,
	})
}

// @Summary List API keys
// @Description Get all API keys for the authenticated user
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of API keys"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/api-keys [get]
func ListAPIKeys(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var apiKeys []models.APIKey
	if err := database.DB.Where("user_id = ?", userID).Find(&apiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API keys"})
		return
	}

	// Don't return the hashed keys
	var response []gin.H
	for _, key := range apiKeys {
		response = append(response, gin.H{
			"id":          key.ID,
			"name":        key.Name,
			"description": key.Description,
			"created_at":  key.CreatedAt,
			"last_used":   key.LastUsedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": response})
}

// @Summary Delete API key
// @Description Delete a specific API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "API key ID"
// @Success 200 {object} map[string]interface{} "API key deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "API key not found"
// @Router /users/api-keys/{id} [delete]
func DeleteAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	keyID := c.Param("id")

	// Verify ownership and delete
	result := database.DB.Where("id = ? AND user_id = ?", keyID, userID).Delete(&models.APIKey{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}
