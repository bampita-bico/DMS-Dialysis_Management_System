# SYSTEM_MAP — Dialysis Management System (DMS)

This file is the project's **external memory**.

Update rules:
- Whenever I add a **table**, **relationship**, **Go service/module**, or **Flutter feature**, I update this file.
- Keep entries short, factual, and link to file paths.

---

## Architecture Snapshot

- Backend: Go + Gin, PostgreSQL (central), SQLite (edge/mobile cache)
- DB Access: `sqlc` (typed queries) + `goose` (SQL migrations)
- Multi-tenancy: PostgreSQL **RLS** enforced via **JWT claim** `hospital_id`
  - API middleware validates JWT
  - Each DB transaction runs `SET LOCAL app.hospital_id = '<uuid>'`
  - RLS policies use: `current_setting('app.hospital_id', true)::uuid`
- UUID primary keys for offline-first sync safety
- Soft delete: `deleted_at` on tenant tables

---

## Repo Layout

- `backend/` — Go API server
- `backend/internal/db/migrations/` — goose SQL migrations
- `backend/internal/db/query/` — sqlc query files
- `flutter/dms_app/` — Flutter multi-platform app
- `legacy/` — Access exports (CSV/SQL) when available

---

## Clinical + Infrastructure Tables (Planned / In Progress)

### Batch 01 — First 10 tables (foundation)

1) `hospitals`
2) `staff`
3) `patients`
4) `vascular_access`
5) `dialysis_sessions`
6) `lab_results`
7) `inventory_items`
8) `suppliers`
9) `payments`
10) `audit_logs`

Tenant enforcement:
- All tables except `hospitals` include `hospital_id uuid not null`.
- `hospitals` uses `id` as the tenant key.
