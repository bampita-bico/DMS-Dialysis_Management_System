# DMS Deployment Model

This document explains the recommended real-world deployment model for the Dialysis Management System in non-IT terms.

## Recommended Decision

Use two separate production environments:

1. **Government hospitals:** one Ministry-hosted national DMS for public hospitals.
2. **Private hospitals:** a separate private-sector DMS environment, not mixed with government hospital data.

For both environments, the default should be a **shared multi-hospital database**, with each hospital separated by `hospital_id` and PostgreSQL Row Level Security (RLS). Do not create one database per hospital by default; it is harder to support, back up, upgrade, and report from.

Use a dedicated database or dedicated deployment only when a large hospital group, legal agreement, or performance requirement genuinely needs it.

## What Gets Deployed

A production DMS installation has four main parts:

- **Frontend:** the web app staff use in the browser.
- **Backend API:** the secure server that receives requests and applies rules.
- **PostgreSQL database:** the clinical and operational records.
- **Operations layer:** backups, monitoring, logs, SSL/TLS certificates, and user access controls.

Hospitals should not run their own copy from a laptop for real production use. Production should run on Ministry servers, a trusted private cloud, or a properly managed hospital server.

## Government Deployment

For government hospitals, the most feasible model is:

- Ministry hosts one national government DMS.
- Every public dialysis unit is registered as a separate hospital/tenant.
- All public hospital data lives in one secured government PostgreSQL environment.
- RLS prevents one hospital from seeing another hospital's patient records.
- Ministry-approved platform users can access central reports across hospitals.
- Normal hospital staff only see their own hospital.

This model gives the Ministry national reporting without forcing each hospital to maintain its own database, server, backup process, and upgrade process.

## Private Hospital Deployment

Private hospitals should be kept separate from the government environment.

The practical default is:

- One private-sector DMS environment.
- One shared private-sector database separated by `hospital_id` and RLS.
- Each private hospital only sees its own patients, sessions, lab results, staff, billing, and reports.
- Platform admins manage hospital onboarding, support, and subscription status.

For a large private hospital chain or a hospital with strict contractual requirements, deploy a dedicated private instance and database for that group.

## Database Model

The system is designed for multi-tenancy:

- Most hospital-owned tables include `hospital_id`.
- Authentication identifies the user's hospital.
- The backend sets the tenant context for each request.
- PostgreSQL RLS enforces the hospital boundary inside the database.
- Platform-level actions are separate and must be protected by `super_admin` access.

This means a shared database can still keep hospitals separate when implemented and operated correctly.

## Admin Access Model

There are two different admin levels.

### Platform Admin

Platform admin is for the system owner, Ministry team, or DMS support team. It should live at:

```text
/platform
```

It should not be visible to ordinary hospital users in the top navigation. It is for:

- registering hospitals
- viewing hospitals across the platform
- managing platform-level settings
- central reporting
- support and audit work
- super admin users

Only `super_admin` users should access this page.

### Hospital Staff Management

Individual hospitals should use Staff Management, not the platform admin page.

Hospital admins should be able to:

- add and deactivate staff for their own hospital
- assign roles for their own hospital
- reset staff access for their own hospital
- manage schedules and shifts

They should not see other hospitals or platform-level controls.

## Reporting And Printing

The Ministry and hospital teams need printable and exportable reports. Required reporting surfaces include:

- **MMR:** Monthly Mortality Report for Ministry reporting.
- **Patient history:** printable longitudinal patient summary.
- **Lab results:** printable lab result reports.
- **Dialysis session records:** printable treatment/session summaries where clinically useful.
- **Operational reports:** patient volume, sessions, equipment use, missed treatments, and outcomes.
- **Financial reports:** billing and payment reports where billing is enabled.

Print buttons should appear where staff naturally need paper or PDF output: patient profile, lab results, MMR, clinical reports, billing, and selected session records.

## Feasible Rollout Plan

### Phase 1: Government Pilot

- Deploy on Ministry-approved secure servers.
- Onboard 2-3 public dialysis units first.
- Test patient registration, sessions, lab results, MMR, staff access, printing, backups, and role boundaries.
- Confirm Ministry reporting format before national rollout.

### Phase 2: Government Rollout

- Add all public dialysis units.
- Train each hospital's hospital admins.
- Keep platform admin limited to the Ministry/DMS team.
- Run national MMR and operational dashboards centrally.

### Phase 3: Private Sector Rollout

- Deploy a separate private-sector environment.
- Start with shared multi-tenant hosting for private hospitals.
- Offer dedicated deployments only for large chains or hospitals that require it.

### Phase 4: Integrations

- Add national registry exports, lab integrations, billing integrations, SMS reminders, and Ministry reporting APIs only after the core clinical workflow is stable.

## Production Checklist

Before real patient data is entered:

- Replace all demo passwords.
- Use HTTPS only.
- Use strong database passwords and restricted database access.
- Enable automated backups and test restore.
- Restrict platform admin accounts.
- Verify RLS isolation with test users from different hospitals.
- Enable audit logs for sensitive actions.
- Confirm Ministry MMR format and reporting schedule.
- Confirm data protection requirements for Uganda and each deployment environment.

## Current Local Demo Access

These are local/demo credentials only and must be rotated before production:

- Hospital user: `doctor@demo.com` / `password123`
- Platform super admin username: `bampita-bico`
- Platform super admin password: local/demo secret, rotate before production
- Platform route: `/platform`
