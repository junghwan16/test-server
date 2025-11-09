# Test Server

## Development Setup

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Quick Start

1. **Start PostgreSQL and Redis**

   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. **Set Environment Variables**

   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Run Application**

   ```bash
   go run ./cmd/server
   ```

4. **Test**

   ```bash
   # Register
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"test123"}'

   # Login
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"test123"}'
   ```

### Environment Variables

Required:

- `DB_HOST`: PostgreSQL host (default: `localhost`)
- `DB_USER`: PostgreSQL user (default: `postgres`)
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: Database name (default: `postgres`)
- `DB_PORT`: PostgreSQL port (default: `5432`)
- `REDIS_HOST`: Redis host (default: `localhost`)
- `REDIS_PORT`: Redis port (default: `6379`)

Optional:

- `REDIS_PASSWORD`: Redis password (default: empty)
- `REDIS_DB`: Redis database number (default: `0`)
- `SERVER_PORT`: Server port (default: `8080`)
- `ENV`: Environment mode - `development` or `production`
- `SESSION_TTL`: Session TTL in seconds (default: `86400`)
- `RATE_LIMIT_RPS`: Rate limit requests per second (default: `10`)
- `RATE_LIMIT_BURST`: Rate limit burst size (default: `20`)
