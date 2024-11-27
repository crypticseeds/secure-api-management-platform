# Build stage
FROM golang:1.23.3-alpine3.20 AS builder

# Add these lines near the top, after the FROM instruction
LABEL maintainer="Femi Akinlotan <femi.akinlotan@devopsfoundry.com>"
LABEL description="Secure API Management Platform"

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/main cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Set environment variables
ENV DB_HOST=postgres-db \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=apisecurity \
    DB_PORT=5432 \
    JWT_SECRET=

# Expose port 8080
EXPOSE 8080

# Set the binary as the entrypoint
ENTRYPOINT ["/app/main"]