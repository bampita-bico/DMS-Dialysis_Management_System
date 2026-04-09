# Module 1 Foundation - Progress Tracker

**Date:** 2026-04-08
**Status:** IN PROGRESS

## Completed Tasks

### ✓ Migrations Created (12/12)
1. ✓ 002_create_hospitals.sql
2. ✓ 003_create_hospital_settings.sql
3. ✓ 004_create_departments.sql
4. ✓ 005_create_users.sql
5. ✓ 006_create_roles.sql
6. ✓ 007_create_user_roles.sql
7. ✓ 008_create_auth_sessions.sql
8. ✓ 009_create_audit_logs.sql
9. ✓ 010_create_sync_queue.sql
10. ✓ 011_create_sync_conflicts.sql
11. ✓ 012_create_notifications.sql
12. ✓ 013_create_file_attachments.sql

### ✓ sqlc Query Files Created (12/12)
1. ✓ hospitals.sql
2. ✓ hospital_settings.sql
3. ✓ departments.sql
4. ✓ users.sql
5. ✓ roles.sql
6. ✓ user_roles.sql
7. ✓ auth_sessions.sql
8. ✓ audit_logs.sql
9. ✓ sync_queue.sql
10. ✓ sync_conflicts.sql
11. ✓ notifications.sql
12. ✓ file_attachments.sql

## Pending Tasks

### ⏳ Database Setup
- [ ] Start PostgreSQL container (Docker credential issue in WSL)
- [ ] Run migrations: `goose -dir internal/db/migrations postgres "..." up`
- [ ] Verify all tables created

### ⏳ sqlc Code Generation
- [ ] Generate Go structs: `sqlc generate`
- [ ] Verify generated code in `internal/db/sqlc/`

### ✓ Phase 3 - Middleware (COMPLETED)
- [x] Write RLS middleware (`backend/internal/middleware/rls.go`)
- [x] Write JWT middleware (updated `backend/internal/http/middleware/jwt_auth.go`)
- [x] Write audit middleware (`backend/internal/middleware/audit.go`)
- [x] Complete JWT service (`backend/internal/security/jwt.go`)

### ✓ Seed Script (COMPLETED)
- [x] Write roles seed script (`backend/scripts/seed_roles.go`)

### ⏳ Phase 4 - API Endpoints
- [ ] Create API endpoints for each table

## Next Steps

1. **Fix Docker/PostgreSQL issue** - The Docker credential error needs to be resolved
2. **Run migrations** - Once DB is up, migrate all 12 tables
3. **Generate sqlc code** - Run `sqlc generate` after migrations succeed
4. **Test RLS** - Verify tenant isolation works
5. **Build middleware** - JWT + RLS + Audit
6. **Create endpoints** - REST API for all tables

## Notes

- All migrations follow the SCHEMA_PLAN.pdf spec exactly
- All migrations have proper DOWN sections for rollback
- RLS policies implemented on all tenant-scoped tables
- audit_logs has NO RLS (admins see all)
- hospitals table RLS uses `id = dms_current_hospital_id()` (tenant root)
- All other tables use `hospital_id = current_setting('app.current_hospital_id')::UUID`

## Blockers

1. **PostgreSQL not running** - Docker credential issue in WSL environment
   - Error: `docker: error getting credentials - err: exec: "docker-credential-desktop.exe"`
   - Solution needed: Fix Docker Desktop integration or use alternative postgres setup

## Files Created/Modified

### Migrations & Queries
- Deleted: `backend/internal/db/migrations/002_core_tables.sql` (wrong schema)
- Deleted: `backend/internal/db/query/patients.sql` (wrong schema)
- Created: 12 new migration files (002-013)
- Created: 12 new sqlc query files

### Middleware & Security
- Created: `backend/internal/middleware/rls.go` - RLS enforcement
- Created: `backend/internal/middleware/audit.go` - Audit logging
- Updated: `backend/internal/http/middleware/jwt_auth.go` - JWT validation
- Updated: `backend/internal/security/jwt.go` - Complete JWT service

### Scripts
- Created: `backend/scripts/seed_roles.go` - Default roles seeder
