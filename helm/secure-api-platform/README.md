# Secure API Management Platform Helm Chart

This Helm chart deploys the Secure API Management Platform with its dependencies.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PV provisioner support in the underlying infrastructure

## Configuration

The following table lists the configurable parameters of the chart and their default values.


## Configuration

The following table lists the configurable parameters of the chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount.app` | Number of API application replicas | `2` |
| `image.app.repository` | API application image repository | `secure-api-platform` |
| `image.app.tag` | API application image tag | `latest` |
| `postgresql.database` | PostgreSQL database name | `apisecurity` |
| `postgresql.username` | PostgreSQL username | `postgres` |
| `postgresql.password` | PostgreSQL password | `random if not set` |
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.host` | Ingress hostname | `secure-api-platform.local` |

## Usage

1. Update values.yaml with your configuration
2. Install the chart:
   ```bash
   helm install secure-api-platform ./secure-api-platform
   ```
3. Port forward to access the services:
   ```bash
   kubectl port-forward svc/secure-api-platform 8080:8080 --address 0.0.0.0
   ```
4. Access the API application at http://localhost:8080/health to check the health of the API application