# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-based HTTP server implementing a user identity and authentication system following Domain-Driven Design (DDD) principles with Clean Architecture.

## Build and Run Commands

```bash
# Start PostgreSQL and Redis for development
docker-compose -f docker-compose.dev.yml up -d

# Run the server locally
go run ./cmd/server

# Run tests
go test ./...

# Run specific package tests
go test ./internal/identity/domain/user
go test ./internal/identity/application

# Run tests with verbose output
go test -v ./...

# Build binary
go build -o bin/server ./cmd/server

# Build Docker image
docker build -t test-server .
```

## Architecture Overview

### Domain-Driven Design Structure

The codebase follows DDD with Clean Architecture, organized into bounded contexts:

**Identity Bounded Context** (`internal/identity/`)
- **Domain Layer** (`domain/`): Core business logic and entities
  - `user/`: User aggregate root with value objects (Email, Password, Role, UserID)
  - `session/`: Session aggregate for authentication state
  - `verification/`: Email verification and password reset entities
  - Domain events are emitted from aggregates (UserRegistered, EmailVerified, PasswordChanged, etc.)

- **Application Layer** (`application/`): Use cases and application services
  - `AuthService`: Handles login, logout, session validation
  - `UserService`: User registration and management
  - `VerificationService`: Email verification and password reset flows

- **Infrastructure Layer** (`infrastructure/persistence/`): Data persistence
  - User, email verification, and password reset repositories use GORM with PostgreSQL
  - Session repository uses Redis for better performance and automatic expiration
  - Repositories implement domain repository interfaces
  - Maps between domain aggregates and persistence models
  - Publishes domain events after persistence

- **Handler Layer** (`handler/`): HTTP handlers and DTOs
  - Maps HTTP requests to application service calls
  - Returns JSON responses

### Shared Kernel

**Shared Domain** (`internal/shared/domain/`)
- Event bus infrastructure for domain events
- `SimpleEventBus`: In-memory event bus with subscribe/publish pattern

### Server Infrastructure

**Server** (`internal/server/`)
- HTTP middleware (authentication, rate limiting, logging)
- Health check handlers
- `RequireAuth`: Session-based authentication middleware
- `RequireAdmin`: Role-based authorization middleware

**Configuration** (`internal/config/`)
- Environment-based configuration loading
- Database, server, SMTP, session, and rate limit settings

## Key Architectural Patterns

### Aggregate Roots
User is the primary aggregate root. All modifications go through aggregate methods that enforce invariants and emit domain events:
- `NewUser()`: Factory method for new registrations
- `ReconstructUser()`: Factory for loading from persistence
- `Authenticate()`, `VerifyEmail()`, `ChangePassword()`, `ChangeRole()`, etc.

### Domain Events
Events are collected on aggregates and published after successful persistence:
1. Aggregate methods add events via `addEvent()`
2. Repository publishes events after `Save()`
3. Events are cleared from aggregate after publishing

### Repository Pattern
Domain repositories are interfaces in the domain layer, implemented in infrastructure layer:
- User, verification repositories: GORM with PostgreSQL
- Session repository: Redis with TTL-based expiration

### Value Objects
Email, Password, Role, UserID are immutable value objects with validation in constructors. Password uses bcrypt hashing.

### Session Management
- Cookie-based sessions stored in Redis
- Automatic expiration using Redis TTL (no cleanup goroutine needed)
- Session TTL configured via `SESSION_TTL` environment variable (default: 86400 seconds)
- Session data serialized as JSON with user ID, expiration, and creation timestamps

## Entry Point

`cmd/server/main.go` is the application entry point. It:
1. Loads configuration from environment variables
2. Connects to PostgreSQL and Redis
3. Runs auto-migrations for user, email verification, and password reset tables
4. Wires up dependencies (repositories, services, handlers)
5. Sets up HTTP routes with middleware chain
6. Implements graceful shutdown for both database and Redis connections

## Testing

Tests are located alongside source files with `_test.go` suffix:
- Domain layer: Value object validation and aggregate behavior
- Application layer: Service logic and use cases
- No integration tests currently present

## Environment Configuration

Required variables (see `.env.example`):
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`: PostgreSQL connection
- `REDIS_HOST`, `REDIS_PORT`: Redis connection

Optional variables:
- `REDIS_PASSWORD`, `REDIS_DB`: Redis authentication and database selection
- `SERVER_PORT` (default: 8080)
- `ENV` (development/production, affects logging format)
- `SESSION_TTL` (default: 86400 seconds = 24 hours)
- `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`: Rate limiting per IP
- `SMTP_*`: Email configuration for verification emails

## HTTP API Structure

Authentication endpoints:
- `POST /auth/signup`: Register new user
- `POST /auth/login`: Login and create session
- `POST /auth/logout`: Logout and destroy session
- `GET /me`: Get current user (requires auth)

Verification endpoints:
- `POST /verification/request`: Request email verification (requires auth)
- `GET /verification/verify`: Verify email via token
- `POST /password/reset/request`: Request password reset
- `POST /password/reset/confirm`: Confirm password reset

Admin endpoints (require admin role):
- `GET /admin/users`: List all users
- `GET /admin/users/{id}`: Get user by ID
- `PATCH /admin/users/{id}`: Update user
- `DELETE /admin/users/{id}`: Delete user

Health checks:
- `GET /health/live`: Liveness probe
- `GET /health/ready`: Readiness probe (checks DB)
- `GET /health`: Combined health check

## Important Implementation Details

- User passwords are hashed with bcrypt (cost factor 12) in `internal/identity/domain/user/password.go`
- Sessions stored in Redis with automatic TTL-based expiration
- Sessions are validated per request by checking cookie and Redis
- Rate limiting is per-IP with in-memory storage (resets on server restart)
- Domain events currently use synchronous in-memory bus
- GORM auto-generates user IDs (compromise noted in `user_repository.go:56-58`)
- Middleware chain order: Logging → RateLimit → Auth (per route)
- Redis session keys use pattern `session:{session_id}`
