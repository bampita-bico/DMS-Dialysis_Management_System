# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Dialysis Management System (DMS)** - A multi-tenant healthcare management system for dialysis centers.

- **Backend**: Go 1.26.1 + Gin framework + PostgreSQL 16
- **Frontend**: Flutter (multi-platform, in early development)
- **Multi-tenancy**: PostgreSQL Row Level Security (RLS) enforced via JWT claims
- **Database**: Type-safe queries via sqlc, migrations via goose
- **Architecture**: Offline-first design with UUID primary keys

## Development Setup

### Prerequisites
- Go 1.26.1 (specified in `.tool-versions`)
- Docker & Docker Compose
- sqlc and goose (managed via `backend/tools/tools.go`)

### Initial Setup

```bash
# Start PostgreSQL
docker-compose up -d

# Copy environment file
cp backend/.env.example backend/.env

# Install Go tools
cd backend
go mod download

# Run migrations
cd backend
go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" up

# Generate sqlc code (after modifying queries)
cd backend
go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Run the API server
cd backend
go run cmd/api/main.go
```

### Common Commands

```bash
# Build the backend
cd backend && go build -o bin/api cmd/api/main.go

# Run tests
cd backend && go test ./...

# Run a single test
cd backend && go test -v -run TestName ./internal/package

# Create a new migration
cd backend && go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations create migration_name sql

# Check migration status
cd backend && go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" status

# Rollback last migration
cd backend && go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" down

# Regenerate sqlc after query changes
cd backend && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
```

## Architecture

### Multi-Tenancy via PostgreSQL RLS

**Critical security pattern** - All tenant data isolation is enforced at the database level:

1. **JWT Middleware** (`internal/http/middleware/jwt_auth.go`):
   - Extracts `hospital_id` and `staff_id` from JWT claims
   - Stores in Gin context for request lifecycle

2. **Transaction-Level Context** (`internal/db/tenant/tenant.go`):
   - Every transaction must call `tenant.SetLocalHospitalID(ctx, tx, hospitalID)`
   - Sets PostgreSQL session variable: `SET LOCAL app.hospital_id = '<uuid>'`
   - RLS policies use this via `dms_current_hospital_id()` function

3. **Database Schema**:
   - `hospitals` table uses `id` as tenant root
   - All other tables have `hospital_id uuid NOT NULL` foreign key
   - Every tenant table has RLS enabled with isolation policy
   - Soft deletes via `deleted_at timestamptz` column

**Example RLS Policy**:
```sql
ALTER TABLE patients ENABLE ROW LEVEL SECURITY;
CREATE POLICY patients_isolation ON patients
  USING (hospital_id = dms_current_hospital_id());
```

### Project Structure

```
backend/
├── cmd/api/              # Application entry point (main.go)
├── internal/
│   ├── config/          # Environment-based configuration
│   ├── db/
│   │   ├── migrations/  # Goose SQL migrations (001_init.sql, 002_core_tables.sql)
│   │   ├── query/       # sqlc query definitions (*.sql)
│   │   ├── sqlc/        # Generated Go code from sqlc (do not edit manually)
│   │   ├── pool/        # Database connection pool setup
│   │   └── tenant/      # Tenant context helper (SetLocalHospitalID)
│   ├── http/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # JWT auth, request ID, etc.
│   │   ├── routes/      # Route registration
│   │   └── server/      # HTTP server initialization
│   └── security/        # JWT service, password hashing
├── tools/               # Go tool dependencies (sqlc, goose)
├── go.mod
├── sqlc.yaml           # sqlc configuration
└── .env.example        # Environment variable template

flutter/                # Flutter mobile/desktop app (in development)
legacy/                 # Access database exports/references
tests/                  # Integration test scripts
```

### Database Layer

- **Migrations**: Located in `backend/internal/db/migrations/`, managed by goose
  - `001_init.sql`: Extensions (pgcrypto, citext) and `dms_current_hospital_id()` function
  - `002_core_tables.sql`: Core entities (hospitals, staff, patients, vascular_access, etc.)

- **Queries**: sqlc query files in `backend/internal/db/query/`
  - Write SQL with annotations like `-- name: GetPatient :one`
  - Run `sqlc generate` to produce type-safe Go functions in `internal/db/sqlc/`

- **Connection Pool**: Uses pgx/v5 driver, configured via environment variables
  - `DB_MAX_CONNS`: Connection pool size (default: 4 for development)

### HTTP Layer

- **Framework**: Gin web framework
- **Middleware Chain**:
  1. `middleware.RequestID()` - Adds request ID for tracing
  2. `middleware.JWTAuth()` - Validates JWT, extracts hospital_id and staff_id

- **Handler Pattern**: Handlers receive dependencies via struct fields, access tenant context from Gin context

### Configuration

Environment variables loaded via `internal/config/config.go`:
- `APP_ENV`: dev/prod
- `HTTP_ADDR`: Server listen address (default: :8080)
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`
- `DB_MAX_CONNS`: PostgreSQL connection pool size
- `JWT_SECRET`: JWT signing secret (must be secure in production)

## Development Patterns

### Adding a New Database Table

1. **Create Migration**:
```bash
cd backend
go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations create add_table_name sql
```

2. **Write Migration** in `internal/db/migrations/XXX_add_table_name.sql`:
```sql
-- +goose Up
CREATE TABLE table_name (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id uuid NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
  -- other columns
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE INDEX ix_table_name_hospital ON table_name (hospital_id);

CREATE TRIGGER trg_table_name_updated_at
BEFORE UPDATE ON table_name
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE table_name ENABLE ROW LEVEL SECURITY;
CREATE POLICY table_name_isolation ON table_name
  USING (hospital_id = dms_current_hospital_id());

-- +goose Down
DROP TABLE IF EXISTS table_name CASCADE;
```

3. **Run Migration**:
```bash
go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" up
```

4. **Add sqlc Queries** in `internal/db/query/table_name.sql`:
```sql
-- name: CreateTableName :one
INSERT INTO table_name (hospital_id, ...) VALUES ($1, ...) RETURNING *;

-- name: GetTableName :one
SELECT * FROM table_name WHERE id = $1 LIMIT 1;
```

5. **Generate Code**:
```bash
go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
```

### Adding a New HTTP Endpoint

1. Create handler in `internal/http/handlers/`
2. Register route in `internal/http/routes/`
3. Apply `middleware.JWTAuth()` for tenant-scoped endpoints
4. In handler, extract tenant context and use transactions with `tenant.SetLocalHospitalID()`

### Tenant-Scoped Query Pattern

```go
// In handler
hospitalID := c.GetString(middleware.CtxHospitalID)

tx, err := pool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// CRITICAL: Set tenant context
if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
    return err
}

// Now all queries respect RLS policies
queries := sqlc.New(tx)
patient, err := queries.GetPatient(ctx, patientID)
// ...

return tx.Commit(ctx)
```

## Testing

- Unit tests: Standard Go tests alongside source files
- Integration tests: Scripts in `tests/` directory
- Manual API testing: `tests/api_test.sh` (basic health check example)

## Important Notes

- **Never skip tenant context**: All tenant-scoped database operations MUST call `tenant.SetLocalHospitalID()` after beginning a transaction
- **UUID Primary Keys**: All tables use UUIDs for offline-first mobile sync
- **Soft Deletes**: Tenant tables use `deleted_at` instead of hard deletes. Queries must filter `WHERE deleted_at IS NULL`
- **Database-First Security**: RLS policies are the source of truth for tenant isolation. Do not rely solely on application-level filtering
- **Generated Code**: Files in `internal/db/sqlc/` are generated by sqlc. Edit the SQL files in `query/` directory instead
- **Go Workspace**: This project uses Go workspaces (`go.work`). All Go commands should be run from the `backend/` directory

## Reference Documents

- `SYSTEM_MAP.md`: Project memory file tracking architecture decisions, tables, and features
- `SCHEMA_PLAN.pdf`: Database schema planning document
- `TO_DO.pdf`: Project task list
- `accdb_tables.txt`: List of legacy Microsoft Access tables for reference
