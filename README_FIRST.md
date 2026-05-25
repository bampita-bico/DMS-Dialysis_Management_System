# README First - Dialysis Management System

Start here when running the system locally or explaining the deployment model to a hospital, Ministry team, or private partner.

## What To Do Right Now

### 1. Start The Local Demo

```bash
cd /home/bampita/Projects/My-apps/DMS-Dialysis_Management_System
./START.sh
```

If the backend and frontend are already running, open the app directly:

```text
http://localhost:5173
```

### 2. Use The Demo Logins

Hospital user:

- Username/email: `doctor@demo.com`
- Password: `password123`
- Scope: Demo Dialysis Center only

Platform super admin:

- Username: `bampita-bico`
- Email: `msbico@gmail.com`
- Password: local/demo secret, rotate before production
- Route: `http://localhost:5173/platform`

These are local/demo credentials only. Rotate all credentials before any production deployment.

### 3. Read The Deployment Model

Read [DEPLOYMENT_MODEL.md](DEPLOYMENT_MODEL.md) before deciding where the real system will live.

The recommended production model is:

- Government hospitals: one Ministry-hosted national DMS for all public dialysis units.
- Private hospitals: a separate private-sector DMS environment.
- Database: one shared PostgreSQL database per environment by default, separated by `hospital_id` and Row Level Security.
- Dedicated database/deployment: only for large private groups or special legal/performance requirements.

## Documentation Files

| File | Purpose | When to Read |
|------|---------|--------------|
| [DEPLOYMENT_MODEL.md](DEPLOYMENT_MODEL.md) | Production deployment, database model, admin access, rollout plan | Before Ministry/private rollout decisions |
| [README.md](README.md) | Developer setup, commands, architecture | When changing or running the code |
| [frontend/README.md](frontend/README.md) | Frontend development notes | When working on the React app |
| [backend/seeds/README.md](backend/seeds/README.md) | Reference-data seed notes | When loading medications, labs, consumables, and pricing |

## Access URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:5173 | Main application |
| Backend API | http://localhost:8080/api/v1 | REST API |
| Health Check | http://localhost:8080/health | Server status |
| Platform Admin | http://localhost:5173/platform | Super admin console only |

## Admin Model

There are two admin levels.

### Platform Admin

The platform admin page should be hidden from normal hospital navigation and available only to `super_admin` users at:

```text
/platform
```

Use it for system-owner work such as hospital onboarding, platform support, cross-hospital reporting, and Ministry-level oversight.

### Hospital Staff Management

Individual hospitals should not use the platform admin page. They should use Staff Management for:

- adding and deactivating their own staff
- assigning roles inside their own hospital
- managing shifts and local access
- resetting staff access for their own hospital

Hospital users should not see other hospitals or platform-level controls.

## Core System Areas

- Patients and patient history
- Dialysis sessions
- Clinical tracking
- Laboratory orders and lab results
- Monthly Mortality Report (MMR)
- Billing and finance where enabled
- Staff Management
- Equipment and inventory where enabled
- Offline-first local storage and sync

## Printing And Reporting

Print actions should be available where staff naturally need paper or PDF output:

- patient history
- lab results
- Monthly Mortality Report (MMR)
- selected clinical reports
- selected dialysis session summaries
- billing and financial reports where enabled

## Production Reminder

Do not use demo passwords in production. Before real patient data is entered:

- deploy to Ministry-approved or professionally managed servers
- use HTTPS only
- enable and test backups
- restrict platform admin accounts
- verify hospital isolation with test users from different hospitals
- confirm the Ministry MMR format and reporting schedule
- follow the production checklist in [DEPLOYMENT_MODEL.md](DEPLOYMENT_MODEL.md)

## Local Troubleshooting

Backend will not start:

```bash
ls backend/.env
go version
lsof -i :8080
```

Frontend will not start:

```bash
node --version
cd frontend && npm install
lsof -i :5173
```

API health check:

```bash
curl http://localhost:8080/health
```
