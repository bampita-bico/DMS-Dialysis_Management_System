# DMS Tiered Subscription Plans with Feature Flags

## Overview

DMS implements a **hybrid licensing model** combining:
1. **Tiered Plans** (Basic, Standard, Enterprise) - Hard limits on available modules
2. **Feature Flags** - Soft toggles to enable/disable specific modules within a plan

This approach provides maximum flexibility while maintaining clear pricing tiers.

---

## Subscription Tiers

### Basic Tier (48 Essential Tables) — $199/month

**Target**: Small dialysis units (<10 patients), rural clinics  
**Focus**: Core clinical operations only

#### Included Modules:
- ✅ **Foundation** (6 tables): hospitals, users, roles, user_roles, auth_sessions, file_attachments
- ✅ **Patient Core** (5 tables): patients, patient_contacts, allergies, comorbidities, patient_flags
- ✅ **Dialysis Clinical** (17 tables): All session management, vitals, vascular access, water treatment
- ✅ **Pharmacy Essentials** (1 table): medication_administrations (intra-dialysis meds only)
- ✅ **Basic Billing** (2 tables): invoices, payments
- ✅ **Outcomes** (4 tables): clinical_outcomes, mortality_records, quality_indicators, registry_sync

#### NOT Included:
- ❌ Lab management (use hospital LIS)
- ❌ Full pharmacy module (prescriptions, stock, interactions)
- ❌ HR management (shift scheduling, performance, training)
- ❌ Equipment/consumables inventory
- ❌ Advanced billing (insurance claims, payment plans, waivers)
- ❌ Imaging integration
- ❌ CHW program
- ❌ Offline sync

---

### Standard Tier (68 Tables) — $499/month ⭐ RECOMMENDED

**Target**: Medium dialysis units (10-50 patients), district hospitals  
**Focus**: Essential + Recommended features for quality care

#### Everything in Basic PLUS:
- ✅ **Departments** - Organizational structure
- ✅ **Audit Logs** - Compliance and security tracking
- ✅ **Notifications** - Clinical alerts (K+ panic values, machine faults)
- ✅ **Next of Kin** - Emergency contacts
- ✅ **Patient Identifiers** - NHIF, national ID for billing
- ✅ **Lab Alerts** - Critical lab result alerts (limited lab module)
- ✅ **Lab Orders & Results** - Track pending labs and results
- ✅ **Medication Prescriptions** - Active medication lists
- ✅ **Drug Interactions** - Safety checking
- ✅ **Equipment Maintenance** - Preventive maintenance scheduling
- ✅ **Equipment Faults** - Fault reporting and tracking
- ✅ **Billing Accounts** - Patient account balances
- ✅ **Insurance Claims** - NHIF claim submission
- ✅ **Waivers** - Charity care and fee waivers
- ✅ **Staff Schedules** - Shift scheduling
- ✅ **Staff Qualifications** - License expiry tracking
- ✅ **Hospitalizations** - Hospitalization tracking
- ✅ **Power Outage Logs** - Medico-legal documentation

#### NOT Included:
- ❌ Full lab module (test catalog, panels, microbiology)
- ❌ Full pharmacy (stock management, refills, ADR reporting)
- ❌ HR module (performance reviews, training, leave management)
- ❌ Full inventory tracking
- ❌ CHW program
- ❌ Imaging integration (use hospital PACS)
- ❌ Offline sync

---

### Enterprise Tier (All 93 Tables) — $999/month

**Target**: Large dialysis centers (50+ patients), teaching hospitals  
**Focus**: Full feature set for comprehensive management

#### Everything in Standard PLUS:
- ✅ **Full Lab Module** - Test catalog, panels, reference ranges, microbiology, imaging
- ✅ **Full Pharmacy Module** - Medication catalog, stock management, prescription refills, ADR reporting
- ✅ **Full HR Module** - Staff profiles, performance reviews, training records, leave management, incident reports
- ✅ **Full Inventory** - Consumables catalog, stock tracking, usage per session
- ✅ **Advanced Billing** - Service catalog, price lists, invoice items, payment plans, insurance schemes
- ✅ **CHW Program** - Community health worker assignments (rural outreach)
- ✅ **Offline Sync** - Sync queue and conflict resolution for offline operations
- ✅ **Imaging Integration** - Imaging orders and results

#### Feature Flag Options:
- 🔘 Enable/disable offline sync (mobile/field clinics)
- 🔘 Enable/disable CHW program (rural only)
- 🔘 Enable/disable imaging integration (if no PACS)

---

## Feature Flag Configuration

Each hospital has an `enabled_modules` JSONB field in the `hospitals` table:

```json
{
  "lab_management": true,          // Full lab module (Enterprise only)
  "full_pharmacy": true,            // Full pharmacy (Enterprise only)
  "hr_management": true,            // HR features (Enterprise only)
  "inventory_tracking": true,       // Consumables inventory (Enterprise only)
  "advanced_billing": true,         // Complex billing (Enterprise only)
  "offline_sync": false,            // Offline-first sync (Enterprise optional)
  "chw_program": false,             // Community health workers (Enterprise optional)
  "imaging_integration": false,     // Imaging orders (Enterprise optional)
  "outcomes_reporting": true        // Outcomes tracking (Standard+)
}
```

### How Feature Flags Work

1. **Tier Check**: API middleware first checks if the module is allowed in the subscription plan
2. **Flag Check**: If allowed in plan, checks if hospital has enabled the module via feature flag
3. **Access Granted**: Only if BOTH checks pass

Example:
- Basic plan tries to access lab_management → ❌ Blocked by tier (not in Basic)
- Enterprise plan with `lab_management: false` → ❌ Blocked by feature flag
- Enterprise plan with `lab_management: true` → ✅ Allowed

---

## Module-to-Table Mapping

### Core Modules (Always Available)

| Module            | Tables                                                                 | Min Tier |
|-------------------|------------------------------------------------------------------------|----------|
| patients          | patients, patient_contacts, allergies, comorbidities, patient_flags    | Basic    |
| sessions          | dialysis_sessions, session_vitals, session_fluid_balance, etc.        | Basic    |
| vascular_access   | vascular_access, vascular_access_assessments                           | Basic    |
| medications       | medication_administrations                                             | Basic    |
| outcomes          | clinical_outcomes, mortality_records, quality_indicators, registry     | Basic    |
| billing           | invoices, payments                                                     | Basic    |

### Optional Modules (Tier-Restricted)

| Module                | Tables                                           | Min Tier   | Feature Flag          |
|-----------------------|--------------------------------------------------|------------|-----------------------|
| lab_management        | lab_test_catalog, lab_panels, lab_orders, etc.   | Enterprise | `lab_management`      |
| full_pharmacy         | medication_catalog, pharmacy_stock, etc.         | Enterprise | `full_pharmacy`       |
| hr_management         | staff_profiles, leave_records, performance, etc. | Enterprise | `hr_management`       |
| inventory_tracking    | consumables_catalog, consumables_stock, etc.     | Enterprise | `inventory_tracking`  |
| advanced_billing      | service_catalog, price_lists, payment_plans, etc.| Enterprise | `advanced_billing`    |
| offline_sync          | sync_queue, sync_conflicts                       | Enterprise | `offline_sync`        |
| chw_program           | community_health_workers, chw_assignments        | Enterprise | `chw_program`         |
| imaging_integration   | imaging_orders, imaging_results                  | Enterprise | `imaging_integration` |
| lab_alerts            | lab_critical_alerts, lab_orders (limited)        | Standard   | Always enabled        |
| equipment_maintenance | equipment_maintenance, equipment_faults          | Standard   | Always enabled        |
| staff_schedules       | staff_schedules, shift_assignments               | Standard   | Always enabled        |

---

## API Middleware Usage

### Protecting Routes with Module Access

```go
// In routes/routes.go
auth := r.Group("/api/v1")
auth.Use(middleware.JWTAuth(jwtSvc))

// Lab endpoints - Enterprise only
lab := auth.Group("/lab")
lab.Use(middleware.RequireModule(pool, "lab_management"))
{
    lab.GET("/tests", labCatalogHandler.ListTests)
    lab.POST("/orders", labOrdersHandler.CreateOrder)
}

// Pharmacy endpoints - Enterprise only
pharmacy := auth.Group("/pharmacy")
pharmacy.Use(middleware.RequireModule(pool, "full_pharmacy"))
{
    pharmacy.GET("/stock", pharmacyHandler.GetStockLevels)
    pharmacy.POST("/stock/movements", pharmacyHandler.RecordMovement)
}

// Equipment maintenance - Standard+
equipment := auth.Group("/equipment")
equipment.Use(middleware.RequireModule(pool, "equipment_maintenance"))
{
    equipment.POST("/maintenance", equipmentHandler.ScheduleMaintenance)
}
```

### Error Responses

When module access is denied:

```json
{
  "error": "Module not available in your subscription plan",
  "module": "lab_management",
  "plan": "standard",
  "upgrade": "Please upgrade to Enterprise to access this feature"
}
```

Or when module is disabled:

```json
{
  "error": "Module is disabled for your hospital",
  "module": "offline_sync"
}
```

---

## Configuration Management

### Updating Feature Flags (SQL)

```sql
-- Enable offline sync for a specific hospital
UPDATE hospitals
SET enabled_modules = jsonb_set(enabled_modules, '{offline_sync}', 'true')
WHERE id = 'hospital_uuid';

-- Disable CHW program
UPDATE hospitals
SET enabled_modules = jsonb_set(enabled_modules, '{chw_program}', 'false')
WHERE id = 'hospital_uuid';
```

### Updating Feature Flags (API)

```go
// PUT /api/v1/hospitals/:id/modules
{
  "offline_sync": true,
  "chw_program": false
}
```

---

## Migration Strategy

### Running the Migration

```bash
cd backend
go run github.com/pressly/goose/v3/cmd/goose \
  -dir internal/db/migrations \
  postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" \
  up
```

### Default Settings

All existing hospitals will default to:
- **Plan**: `standard`
- **Modules**: All enabled except `offline_sync`, `chw_program`, `imaging_integration`

### Rollback

```bash
go run github.com/pressly/goose/v3/cmd/goose \
  -dir internal/db/migrations \
  postgres "postgres://dms:dms_dev_password@localhost:5432/dms?sslmode=disable" \
  down
```

---

## Benefits of This Approach

1. **Flexibility**: Hospitals can enable/disable features without schema changes
2. **No Data Loss**: All tables remain in place, just access-controlled
3. **Easy Upgrades**: Change subscription_plan, features unlock instantly
4. **Customization**: Enterprise customers can disable unused modules
5. **No Migration Risks**: No table deletions, no data archival needed
6. **Gradual Rollout**: Enable features incrementally as hospital is ready

---

## Next Steps

1. ✅ Run migration to add `subscription_plan` and `enabled_modules` to hospitals
2. ✅ Add `RequireModule` middleware to protected routes
3. ⏳ Create admin UI to manage hospital plans and feature flags
4. ⏳ Update API documentation with tier requirements
5. ⏳ Add billing integration for plan upgrades/downgrades
6. ⏳ Create usage analytics to recommend plan upgrades

---

## Support & Upgrades

- **Basic → Standard**: $300 upgrade (billed annually)
- **Standard → Enterprise**: $500 upgrade (billed annually)
- **Custom Enterprise**: Contact sales for hospitals with >100 patients

For plan changes, contact: support@dms.africa
