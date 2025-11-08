# Test Server

## Development Setup

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Quick Start

1. **Start PostgreSQL**

   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. **Set Environment Variables**

   ```bash
   cp .env.example .env
   # Edit .env and set JWT_SECRET
   ```

3. **Run Application**

   ```bash
   go run ./cmd/app
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
- `JWT_SECRET`: JWT signing key

Optional:

- `SERVER_PORT`: Server port (default: `8080`)
- `ENV`: Environment mode - `development` or `production`
