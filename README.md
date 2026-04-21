# Dialysis Management System (DMS)

A comprehensive, multi-tenant healthcare management system for dialysis centers built with Go, PostgreSQL, and React.

## Features

- **Multi-tenant Architecture**: Secure tenant isolation using PostgreSQL Row Level Security (RLS)
- **Comprehensive Patient Management**: Demographics, medical history, dialysis sessions, medications
- **Clinical Workflows**: Session management, vital signs tracking, assessments, complications
- **Staff & Scheduling**: User management, roles, shift scheduling
- **Equipment & Inventory**: Dialysis machines, consumables, maintenance tracking
- **Billing & Finance**: Treatment billing, insurance claims, financial reporting
- **Laboratory Integration**: Test orders, results tracking, trend analysis
- **Reporting**: Clinical, operational, and financial reports

## Tech Stack

- **Backend**: Go 1.26.1 + Gin web framework
- **Database**: PostgreSQL 16 with Row Level Security (RLS)
- **Frontend**: React 19 + Vite (responsive web app)
- **Database Tools**: sqlc (type-safe queries), goose (migrations)
- **Authentication**: JWT-based with multi-tenant context

## Prerequisites

- Go 1.26.1 or later
- PostgreSQL 16 (via Docker or remote VPS)
- Node.js 18+ and npm (for frontend)
- Docker & Docker Compose (for local PostgreSQL)

## Development Setup Options

Choose one of three setup modes based on your hardware resources:

### Option 1: FREE Remote Database (Supabase) - Recommended for 4GB RAM

**Best for:** Resource-constrained machines, FREE hosting, fast setup

**Benefits:**
- 100% FREE forever (no credit card)
- 15-30 minute setup
- Frees ~500MB local RAM
- Includes web dashboard
- Automatic backups

**Quick Start:**

```bash
# 1. Sign up at supabase.com (FREE, no credit card)
# 2. Create project, get connection string
# 3. Install goose locally:
go install github.com/pressly/goose/v3/cmd/goose@latest

# 4. Run migrations from local machine:
cd backend
export SUPABASE_CONN="your_connection_string_here"
goose -dir internal/db/migrations postgres "$SUPABASE_CONN" up

# 5. Configure backend:
cp backend/.env.example.supabase backend/.env.supabase
# Edit with Supabase details
cp backend/.env.supabase backend/.env

# 6. Start backend
go run cmd/api/main.go

# 7. Start frontend
cd ../frontend && npm install && npm run dev
```

**Complete guide:** See [`FREE_SETUP.md`](FREE_SETUP.md)

### Option 2: Local Development (Docker) - For 8GB+ RAM

**Best for:** Machines with sufficient resources, offline development

```bash
# Start PostgreSQL
docker-compose up -d

# Copy environment file
cp backend/.env.example backend/.env

# Install dependencies
cd backend && go mod download

# Run migrations
cd backend
go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" up

# Start backend
go run cmd/api/main.go

# In another terminal, start frontend
cd frontend
npm install
npm run dev
```

Access at: `http://localhost:5173`

### Option 3: Paid VPS (DigitalOcean/Linode) - For Full Control

**Best for:** Production-like setup, faster performance, full server control

**Benefits:**
- Frees ~500MB local RAM
- Full control over PostgreSQL
- Lower latency than free services
- $5-6/month for VPS

**Quick Start (after VPS setup):**

```bash
# Configure remote database
cp backend/.env.example.remote backend/.env.remote
# Edit .env.remote with your VPS IP and password

# Activate remote configuration
cp backend/.env.remote backend/.env

# Start backend (lightweight, ~50-100MB)
cd backend && go run cmd/api/main.go

# Start frontend
cd frontend && npm run dev
```

**Complete Remote Setup Guide:** See [`REMOTE_SETUP.md`](REMOTE_SETUP.md) for:
- VPS provisioning (DigitalOcean/Linode/Vultr)
- PostgreSQL installation and configuration
- Database migration
- Automated backups
- Troubleshooting

## Demo Credentials

After running migrations and loading demo data:

- **Email**: `doctor@demo.com`
- **Password**: `password123`
- **Hospital**: Demo Dialysis Center (DEMO)

## Project Structure

```
DMS/
├── backend/
│   ├── cmd/api/              # Application entry point
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── db/
│   │   │   ├── migrations/  # Database migrations (104 files)
│   │   │   ├── query/       # SQL query definitions (sqlc)
│   │   │   ├── sqlc/        # Generated Go code
│   │   │   ├── pool/        # Connection pool
│   │   │   └── tenant/      # Multi-tenant context
│   │   ├── http/
│   │   │   ├── handlers/    # HTTP handlers
│   │   │   ├── middleware/  # JWT auth, CORS, etc.
│   │   │   ├── routes/      # Route registration
│   │   │   └── server/      # HTTP server setup
│   │   └── security/        # JWT, password hashing
│   ├── scripts/             # Utility scripts, demo data
│   └── tools/               # Go tool dependencies
├── frontend/                # React web application
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/          # Page components
│   │   └── services/       # API service layer
│   └── public/             # Static assets
└── docs/                   # Documentation
```

## Common Commands

```bash
# Backend

# Build
cd backend && go build -o bin/api cmd/api/main.go

# Run tests
cd backend && go test ./...

# Create new migration
cd backend && go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations create migration_name sql

# Generate sqlc code (after query changes)
cd backend && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Frontend

# Install dependencies
cd frontend && npm install

# Run dev server
cd frontend && npm run dev

# Build for production
cd frontend && npm run build
```

## Multi-Tenancy Architecture

DMS uses **PostgreSQL Row Level Security (RLS)** for tenant isolation:

1. **JWT Middleware** extracts `hospital_id` from token claims
2. **Transaction Context** sets PostgreSQL session variable: `SET LOCAL app.hospital_id = '<uuid>'`
3. **RLS Policies** automatically filter all queries by tenant

**Example RLS Policy:**
```sql
CREATE POLICY patients_isolation ON patients
  USING (hospital_id = dms_current_hospital_id());
```

This ensures complete data isolation at the database level - no application-level filtering needed!

## Key Features

### Clinical Modules
- **Patient Management**: Demographics, contacts, insurance, medical history
- **Dialysis Sessions**: Treatment tracking, vital signs, complications
- **Vascular Access**: Catheter and fistula management
- **Laboratory**: Test orders, results, alerts
- **Medications**: Prescriptions, administration, tracking

### Administrative Modules
- **Staff Management**: Users, roles, permissions, scheduling
- **Equipment**: Dialysis machines, beds, maintenance schedules
- **Inventory**: Consumables, suppliers, stock tracking
- **Billing**: Treatment billing, insurance, payments

### Reporting
- Clinical reports (patient outcomes, treatment compliance)
- Operational reports (session volume, equipment utilization)
- Financial reports (revenue, outstanding payments)

## API Documentation

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: Bearer token (JWT)

**Key Endpoints:**
- `POST /auth/login` - User authentication
- `GET /patients` - List patients (tenant-scoped)
- `POST /patients` - Create patient
- `GET /dialysis-sessions` - List sessions
- `POST /dialysis-sessions` - Record session

See `backend/internal/http/routes/` for complete endpoint list.

## Environment Variables

**Backend** (`backend/.env`):

```bash
APP_ENV=dev                    # dev or prod
HTTP_ADDR=:8080                # Server listen address

# Database (local or remote)
DB_HOST=localhost              # Use VPS IP for remote
DB_PORT=5432
DB_NAME=dms
DB_USER=dms
DB_PASSWORD=dms_dev_password   # Use secure password for remote
DB_SSLMODE=disable             # Use require for production
DB_MAX_CONNS=4                 # Connection pool size

# Security
JWT_SECRET=CHANGE_ME_DEV_ONLY  # Must be secure in production
```

**Frontend** (`frontend/.env`):

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

## Security Notes

### Development
- JWT secret is hardcoded (acceptable for development)
- SSL disabled for database connections
- CORS allows localhost origins
- Demo credentials available

### Production (TODO)
- Use environment-specific JWT secrets (rotation recommended)
- Enforce SSL/TLS for all database connections
- Restrict CORS to production domains
- Remove demo data
- Implement rate limiting
- Add API key authentication for external integrations
- Setup monitoring and alerting

## Database Management

### Migrations

**Create new migration:**
```bash
cd backend
go run github.com/pressly/goose/v3/cmd/goose -dir internal/db/migrations create add_feature sql
```

**Apply migrations:**
```bash
goose -dir internal/db/migrations postgres "CONNECTION_STRING" up
```

**Rollback:**
```bash
goose -dir internal/db/migrations postgres "CONNECTION_STRING" down
```

### Backups

**Local (Docker):**
```bash
docker exec -t $(docker ps -qf "name=dms-postgres") pg_dump -U dms dms > backup.sql
```

**Remote (VPS):**
```bash
ssh VPS_IP "pg_dump -U dms dms | gzip" > backup.sql.gz
```

See `REMOTE_SETUP.md` for automated backup configuration.

## Development Workflow

### Switching Between Local and Remote Database

**Use Remote Database:**
```bash
cp backend/.env.remote backend/.env
```

**Use Local Database:**
```bash
docker-compose up -d
cp backend/.env.local backend/.env
```

### Adding a New Feature

1. **Database Changes:**
   - Create migration: `goose create feature_name sql`
   - Write SQL with RLS policy
   - Apply migration: `goose up`

2. **Add Queries:**
   - Create `backend/internal/db/query/feature.sql`
   - Write sqlc-annotated queries
   - Generate: `sqlc generate`

3. **Create Handler:**
   - Add handler in `backend/internal/http/handlers/`
   - Register route in `backend/internal/http/routes/`
   - Apply JWT middleware for tenant-scoped endpoints

4. **Frontend Integration:**
   - Add API service in `frontend/src/services/`
   - Create components in `frontend/src/components/`
   - Add page in `frontend/src/pages/`

5. **Test:**
   - Unit tests: `go test ./...`
   - Manual testing: Use frontend or API client
   - Verify multi-tenancy isolation

## Testing

### Backend Tests
```bash
cd backend
go test ./...                    # All tests
go test -v ./internal/http/...   # Specific package
go test -run TestName ./...      # Specific test
```

### Frontend Tests
```bash
cd frontend
npm test                         # Run tests
npm run test:watch              # Watch mode
```

### Integration Tests
```bash
# Start backend and frontend
cd tests
./api_test.sh                    # Basic API health check
```

## Troubleshooting

### Backend won't start

**Check:**
- Database is running: `docker ps` or `psql -h VPS_IP`
- Environment variables: `cat backend/.env`
- Connection string format
- Migrations applied: `goose status`

### Frontend can't reach API

**Check:**
- Backend is running: `curl http://localhost:8080/health`
- Frontend `.env` has correct API URL
- CORS configuration in `backend/internal/http/server/server.go`
- Browser console for errors

### Queries return no data (RLS issue)

**Verify:**
- JWT token contains `hospital_id` claim
- `tenant.SetLocalHospitalID()` called before queries
- RLS policies enabled: `SELECT tablename, rowsecurity FROM pg_tables;`

See `REMOTE_SETUP.md` for remote-specific troubleshooting.

## Documentation

- **`CLAUDE.md`**: Complete architecture and development guide
- **`REMOTE_SETUP.md`**: Remote PostgreSQL setup for resource-constrained machines
- **`SYSTEM_MAP.md`**: Project memory and architecture decisions
- **`backend/internal/db/migrations/`**: Database schema definitions

## Contributing

This is currently a solo project for a specific healthcare facility. If adapting for your use:

1. Review multi-tenancy implementation
2. Customize clinical workflows for your region
3. Update billing logic for your payment model
4. Ensure compliance with local healthcare regulations (HIPAA, GDPR, etc.)

## License

Proprietary - All rights reserved

## Support

For questions or issues:
- Review `CLAUDE.md` for architecture details
- Check `REMOTE_SETUP.md` for deployment issues
- Review migration files for database schema

---

**Status**: ✅ Backend production-ready | ⏳ Frontend in active development

**Last Updated**: 2026-04-11
