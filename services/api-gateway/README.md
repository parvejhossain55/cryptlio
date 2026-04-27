# Cryplio API Gateway

API Gateway for the Cryplio P2P cryptocurrency exchange platform. This service acts as the unified entry point for all client requests, routing them to appropriate backend microservices.

## Architecture

The API gateway sits at the edge of the microservices architecture, handling:

- **Authentication & Authorization**: JWT token validation
- **Request Routing**: Reverse proxy to backend services
- **Rate Limiting**: Per-IP and per-user rate limits
- **CORS**: Cross-origin request handling
- **Logging**: Structured request/response logging
- **Health Checks**: Liveness and readiness probes

## Services Proxied

| Service | Port | Responsibility |
|---------|------|----------------|
| Auth Service | TBD | Registration, login, session management |
| User Service | TBD | User profiles, settings |
| Wallet Service | 8081 | Deposits, withdrawals, balances |
| Trade Engine | 8082 | Ads, trades, escrow |
| KYC Service | 8083 | Identity verification |
| Notification | 8084 | Email, SMS, push |
| Dispute Service | 8085 | Dispute resolution |
| Merchant Service | TBD | Merchant dashboard, analytics |

## API Endpoints

All API endpoints are prefixed with `/api/v1`.

### Public Endpoints (no auth)

| Method | Endpoint | Service |
|--------|----------|---------|
| POST | `/api/v1/auth/register` | Auth |
| POST | `/api/v1/auth/login` | Auth |
| GET | `/api/v1/ads` | Trade Engine |
| GET | `/api/v1/market/rates` | Trade Engine |
| GET | `/api/v1/ads/:id` | Trade Engine (optional auth) |

### Authenticated Endpoints (requires Bearer JWT)

| Method | Endpoint | Service |
|--------|----------|---------|
| GET | `/api/v1/users/me` | User |
| PUT | `/api/v1/users/me` | User |
| POST | `/api/v1/kyc/submit` | KYC |
| GET | `/api/v1/kyc/status` | KYC |
| POST | `/api/v1/ads` | Trade Engine |
| PUT | `/api/v1/ads/:id` | Trade Engine |
| DELETE | `/api/v1/ads/:id` | Trade Engine |
| POST | `/api/v1/trades` | Trade Engine |
| POST | `/api/v1/trades/:id/paid` | Trade Engine |
| POST | `/api/v1/trades/:id/release` | Trade Engine |
| POST | `/api/v1/trades/:id/cancel` | Trade Engine |
| POST | `/api/v1/trades/:id/dispute` | Dispute |
| GET | `/api/v1/trades` | Trade Engine |
| GET | `/api/v1/trades/:id` | Trade Engine |
| GET | `/api/v1/wallet/balance` | Wallet |
| POST | `/api/v1/wallet/withdraw` | Wallet |
| POST | `/api/v1/wallet/deposit` | Wallet |
| GET | `/api/v1/wallet/history` | Wallet |
| GET | `/api/v1/notifications` | Notification |
| PATCH | `/api/v1/notifications/:id/read` | Notification |
| GET | `/api/v1/merchant/dashboard` | Merchant |

## Health & Monitoring

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Full health check with backend service status |
| `GET /live` | Liveness probe (K8s) |
| `GET /ready` | Readiness probe (K8s) |

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:

- `SERVER_PORT`: Port to listen on (default: 8080)
- `JWT_SECRET`: Secret for JWT token verification (required)
- `SERVICE_URLS`: Backend service URLs (wallet, trade-engine, kyc, etc.)
- `RATE_LIMIT_ENABLED`: Enable/disable rate limiting
- `RATE_LIMIT_REQUESTS`: Requests per window (default: 100)
- `RATE_LIMIT_WINDOW`: Window duration (default: 1m)

## Development

### Prerequisites

- Go 1.22+
- All backend services running

### Setup

1. Copy `.env.example` to `.env` and configure service URLs
2. Start backend services on configured ports
3. Run the gateway:

```bash
go run ./main.go
```

The gateway will start on `http://localhost:8080`.

### Testing

```bash
# Health check
curl http://localhost:8080/health

# Example API call (with JWT)
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/users/me
```

## Production

### Docker

Build and run with Docker:

```bash
docker build -t cryplio/api-gateway .
docker run -p 8080:8080 --env-file .env cryplio/api-gateway
```

### Kubernetes

Manifests should be created in `infrastructure/k8s/` for:

- Deployment (with replica count > 1)
- Service (ClusterIP or LoadBalancer)
- HorizontalPodAutoscaler
- PodDisruptionBudget
- NetworkPolicy

Standard K8s probes (`/live`, `/ready`, `/health`) are available.

### Rate Limiting

By default, rate limiting is enabled at 100 requests per second per IP. This can be tuned via environment variables. For authenticated endpoints, additional per-user rate limiting can be implemented in the future.

### Logging

Structured JSON logs are emitted using [zerolog](https://github.com/rs/zerolog). Log level can be controlled via `LOG_LEVEL` environment variable.

## Error Handling

Standardized error responses:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

Common HTTP status codes:

- `400` - Bad request
- `401` - Unauthorized (missing/invalid JWT)
- `403` - Forbidden (insufficient permissions)
- `404` - Not found
- `429` - Rate limited
- `500` - Internal server error
- `502` - Bad gateway (backend service unavailable)
- `503` - Service unavailable

## Security

- All endpoints require TLS in production
- JWT tokens should have appropriate expiration
- Backend service communication should be over mTLS in production
- Rate limiting prevents brute-force attacks
- CORS restricted to known origins in production

## License

Proprietary - Cryplio Platform
