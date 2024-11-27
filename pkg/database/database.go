package database

import (
	"fmt"
	"log"
	"os"

	"apisecurityplatform/pkg/models"
	"apisecurityplatform/pkg/observability"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	ctx := context.Background()
	tracer := observability.GetTracer()
	ctx, span := tracer.Start(ctx, "database.connect")
	defer span.End()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Set default values if environment variables are not set
	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}
	if dbname == "" {
		dbname = "apisecurity"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	// Add database info to span
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.name", dbname),
		attribute.String("db.host", host),
		attribute.String("db.port", port),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to connect to database")
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the database schemas
	err = database.AutoMigrate(&models.User{}, &models.APIKey{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to migrate database")
		log.Fatal("Failed to migrate database:", err)
	}

	span.SetStatus(codes.Ok, "Connected to database")
	DB = database
}
