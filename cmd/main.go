package main

import (
	_ "apisecurityplatform/docs"
	"apisecurityplatform/pkg/database"
	"apisecurityplatform/pkg/handlers"
	"apisecurityplatform/pkg/middleware"
	"apisecurityplatform/pkg/observability"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/attribute"
)

// @title           Secure API Management Platform
// @version         1.0
// @description     A secure API management platform with authentication and API key management.

// @contact.name   Femi Akinlotan
// @contact.email  femi.akinlotan@devopsfoundry.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @tag.name auth
// @tag.description Authentication operations

// @tag.name users
// @tag.description User operations

func main() {
	// Initialize tracer
	cleanup, err := observability.InitTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer cleanup()

	// Initialize database connection first
	database.ConnectDatabase()

	router := gin.Default()

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Enable debug mode
	gin.SetMode(gin.DebugMode)

	// Health check endpoint with tracing
	router.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		tracer := observability.GetTracer()

		// Create a child span for simulated DB check
		_, dbSpan := tracer.Start(ctx, "health.check.database")
		// Simulate DB check
		time.Sleep(10 * time.Millisecond)
		dbSpan.SetAttributes(attribute.Bool("database.healthy", true))
		dbSpan.End()

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Add global middleware
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.TracingMiddleware())

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
	}

	// Protected routes
	api := router.Group("/users")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/me", handlers.GetUserProfile)
		api.DELETE("/:id", handlers.DeleteUser)

		// API Key routes
		api.POST("/api-keys", handlers.CreateAPIKey)
		api.GET("/api-keys", handlers.ListAPIKeys)
		api.DELETE("/api-keys/:id", handlers.DeleteAPIKey)
	}

	// Protected route with API key authentication
	apiKeyProtected := router.Group("/api")
	apiKeyProtected.Use(middleware.APIKeyAuth())
	{
		apiKeyProtected.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "API key is valid",
				"key_id":  c.GetUint("api_key_id"),
			})
		})
	}

	// Print out all registered routes for debugging
	fmt.Println("\nRegistered Routes:")
	for _, route := range router.Routes() {
		fmt.Printf("Method: %s, Path: %s, Handler: %v\n",
			route.Method,
			route.Path,
			route.Handler)
	}

	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
