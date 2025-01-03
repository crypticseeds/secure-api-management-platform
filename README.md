# Secure API Management Platform

The Secure API Management Platform is a project demonstrating secure user authentication, API key management, Rate limiting, observability, and automated documentation. The application is built with Go, containerized using Docker, and deployed in Kubernetes. Monitoring and metrics are integrated using Prometheus and Jaeger.

## Features

### User Management:
- User registration, login, and profile management.
- JWT-based authentication with role-based access control (RBAC).

### API Key Management:
- Generate, list, and revoke scoped API keys for programmatic access.

### API Documentation:
- Auto-generated Swagger/OpenAPI documentation available at /docs.

### Monitoring and Observability:
- Request tracing using Jaeger.
- Metrics collection using Prometheus.

### Health and Metrics Endpoints:
- /health for application readiness.
- /metrics for Prometheus metrics.

### Rate Limiting:
- Rate limiting middleware to prevent abuse.

## Tech Stack

- Programming Language: Go
- Framework: Gin
- Database: PostgreSQL
- Observability: Prometheus, Jaeger
- Documentation: Swagger (Swaggo)
- Containerization: Docker
- Deployment: Docker Compose or Kubernetes

## Project Structure

- `/cmd`: Entry point for the application.
- `/pkg`: Reusable packages (e.g., auth, database, handlers).
- `/configs`: Configuration files (e.g., database, environment variables).
- `/docs`: Swagger/OpenAPI documentation files.

## Getting Started

### Prerequisites
- Go 1.20+
- Docker
- Kubernetes (Minikube/KIND for local development)
- PostgreSQL
- Git

## API Endpoints

### Authentication

| Endpoint | Method | Description | Request Body | Security |
|----------|---------|-------------|--------------|-----------|
| `/auth/register` | POST | Register a new user | `{ "username": "string", "email": "string", "password": "string" }` | Password hashing, validation |
| `/auth/login` | POST | Login and get a JWT token | `{ "email": "string", "password": "string" }` | Token expiration |
| `/auth/logout` | POST | Invalidate the user's token | None | JWT validation |

### User Management

| Endpoint | Method | Description | Security |
|----------|---------|-------------|-----------|
| `/users/me` | GET | Get logged-in user's profile | JWT Authentication |
| `/users/{id}` | DELETE | Delete a user (admin-only) | Role-based Access |

### API Key Management

| Endpoint | Method | Description | Request Body | Security |
|----------|---------|-------------|--------------|-----------|
| `/api-keys` | POST | Generate a new API key | None | JWT Authentication |
| `/api-keys` | GET | List all API keys | None | JWT Authentication |
| `/api-keys/{id}` | DELETE | Revoke an API key | None | Role-based Access |
| `/api/test` | GET | Get usage metrics for an API key | None | X-API-Key: API Token |

### Monitoring

| Endpoint | Method | Description |
|----------|---------|-------------|
| `/health` | GET | Health check endpoint |
| `/metrics` | GET | Prometheus metrics endpoint |


## API Documentation

The API documentation is available through Swagger UI:

1. Run the application:
   ```bash
   docker-compose up
   ```

2. Access Swagger UI:
   - Open [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)
   - Browse and test available endpoints
   - View request/response schemas and examples

3. Access the Jaeger UI:
   - Open [http://localhost:16686](http://localhost:16686)
   - Select "api-security-platform" from the Service dropdown
   - Click "Find Traces" to view traces


## Access Services:

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/docs/index.html#/
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000


## Reference Images

<img src="./docs/images/01-health.png" alt="Health Check" />

<img src="./docs/images/02-register-user.png" alt="Register" />

<img src="./docs/images/03-login.png" alt="Login" />

<img src="./docs/images/06-create-api-key.png" alt="Create API Key" />

<img src="./docs/images/07-list-api-keys.png" alt="List API Keys" />

<img src="./docs/images/08-api-key-test.png" alt="API Key Test" />

<img src="./docs/images/09-delete-api-key.png" alt="Delete API Key" />

<img src="./docs/images/12-swagger.png" alt="Swagger UI" />

<img src="./docs/images/13-jaeger.png" alt="Jaeger" />

<img src="./docs/images/15-dashboard.png" alt="Grafana" />

<img src="./docs/images/16-response-time.png" alt="Response Time" />


## Future Improvements

- Add Keycloak for authentication and authorization.
- Add refresh tokens for enhanced session management.
- Implement webhook support for user-defined event notifications.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
