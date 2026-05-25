import Dexie from 'dexie';

// IndexedDB schema for offline storage
export const db = new Dexie('dms_local');

// Define database schema
db.version(5).stores({
  // Core entities
  patients: 'id, mrn, hospital_id, full_name, *search_terms, updated_at, synced',
  patient_contacts: 'id, patient_id, contact_type, hospital_id, updated_at, synced',

  // Dialysis sessions & machines
  dialysis_sessions: 'id, patient_id, scheduled_date, status, hospital_id, updated_at, synced',
  dialysis_machines: 'id, machine_number, model, operational_status, hospital_id, updated_at, synced',
  session_vitals: 'id, session_id, recorded_at, updated_at, synced',
  session_complications: 'id, session_id, complication_type, severity, updated_at, synced',
  session_fluid_balance: 'id, session_id, updated_at, synced',
  dialysate_records: 'id, session_id, patient_id, recorded_at, updated_at, synced',

  // Vascular access
  vascular_access: 'id, patient_id, access_type, access_status, updated_at, synced',
  vascular_access_assessments: 'id, patient_id, access_id, session_id, assessed_at, updated_at, synced',

  // Consultant clinical tracking extensions
  patient_clinical_profiles: 'id, patient_id, renal_course, kidney_disease_cause, updated_at, synced',
  session_safety_checks: 'id, patient_id, session_id, check_status, override_required, checked_at, updated_at, synced',
  clinical_alerts: 'id, patient_id, session_id, alert_type, severity, status, created_at, synced',
  mortality_records: 'id, patient_id, date_of_death, death_setting, hospital_id, updated_at, synced',
  treatment_telemetry: 'id, patient_id, session_id, recorded_at, updated_at, synced',
  access_lifecycle_events: 'id, patient_id, access_id, event_type, event_date, updated_at, synced',
  infection_surveillance_events: 'id, patient_id, access_id, event_type, event_date, reported_to_registry, updated_at, synced',
  adequacy_reviews: 'id, patient_id, session_id, review_month, adequacy_status, doctor_review_required, updated_at, synced',
  medication_reconciliation_reviews: 'id, patient_id, review_date, review_type, status, updated_at, synced',
  patient_reported_events: 'id, patient_id, session_id, event_type, reported_at, updated_at, synced',
  staff_attendance_verifications: 'id, staff_id, user_id, verification_date, verification_result, updated_at, synced',
  unit_safety_events: 'id, event_type, severity, event_date, closed_at, updated_at, synced',
  interoperability_exports: 'id, patient_id, export_type, fhir_resource_type, export_status, updated_at, synced',
  ontology_relationships: 'id, patient_id, source_type, relation_type, target_type, is_active, updated_at, synced',

  // Medical history
  diagnoses: 'id, patient_id, diagnosis_type, icd10_code, diagnosed_at, updated_at, synced',
  comorbidities: 'id, patient_id, condition, status, diagnosed_at, updated_at, synced',
  consents: 'id, patient_id, consent_type, status, signed_at, updated_at, synced',

  // Lab management
  lab_orders: 'id, patient_id, order_status, ordered_at, hospital_id, updated_at, synced',
  lab_order_items: 'id, order_id, test_id, specimen_status, updated_at, synced',
  lab_results: 'id, order_item_id, patient_id, result_value, result_status, updated_at, synced',
  lab_critical_alerts: 'id, patient_id, test_id, acknowledged, created_at, synced',

  // Reference data (read-only, synced from server)
  lab_test_catalog: 'id, hospital_id, code, name, category',
  lab_panels: 'id, hospital_id, panel_code, panel_name',
  medications: 'id, hospital_id, generic_name, *search_terms',
  consumables: 'id, hospital_id, item_name, category',
  insurance_schemes: 'id, hospital_id, scheme_code, scheme_name',
  price_lists: 'id, hospital_id, service_code, service_name',

  // Prescriptions
  prescriptions: 'id, patient_id, status, prescribed_at, hospital_id, updated_at, synced',
  prescription_items: 'id, prescription_id, medication_id, updated_at, synced',

  // Billing
  invoices: 'id, patient_id, invoice_number, invoice_status, hospital_id, updated_at, synced',
  payments: 'id, invoice_id, payment_method, payment_date, hospital_id, updated_at, synced',

  // Staff
  staff_profiles: 'id, user_id, staff_cadre, hospital_id, updated_at, synced',
  shift_assignments: 'id, staff_id, shift_date, shift_type, hospital_id, updated_at, synced',

  // Sync queue - tracks pending changes to upload
  sync_queue: '++id, entity_type, entity_id, operation, payload, synced, priority, created_at, attempts',

  // Metadata - tracks last sync times per entity type
  _metadata: 'key, value, updated_at',
});

// Helper function to initialize metadata
export const initializeMetadata = async () => {
  const entityTypes = [
    'patients', 'dialysis_sessions', 'lab_orders', 'lab_results',
    'prescriptions', 'invoices', 'payments', 'vascular_access',
    'vascular_access_assessments', 'session_vitals', 'dialysate_records',
    'diagnoses', 'comorbidities', 'consents', 'patient_clinical_profiles',
    'session_safety_checks', 'clinical_alerts', 'mortality_records', 'treatment_telemetry',
    'access_lifecycle_events', 'infection_surveillance_events', 'adequacy_reviews',
    'medication_reconciliation_reviews', 'patient_reported_events',
    'staff_attendance_verifications', 'unit_safety_events', 'interoperability_exports',
    'ontology_relationships',
  ];

  for (const entityType of entityTypes) {
    const exists = await db._metadata.get(`last_sync_${entityType}`);
    if (!exists) {
      await db._metadata.put({
        key: `last_sync_${entityType}`,
        value: null,
        updated_at: new Date().toISOString(),
      });
    }
  }
};

// Helper function to mark entity as synced
export const markAsSynced = async (tableName, localId, serverId = null) => {
  const table = db.table(tableName);
  const updates = { synced: true, updated_at: new Date().toISOString() };

  if (serverId && serverId !== localId) {
    // Server assigned a different ID, update it
    updates.id = serverId;
  }

  await table.update(localId, updates);
};

// Helper function to get unsynced records
export const getUnsyncedRecords = async (tableName) => {
  return await db.table(tableName)
    .where('synced')
    .equals(false)
    .toArray();
};

// Helper function to clear all data (for logout)
export const clearAllData = async () => {
  await db.delete();
  await db.open();
  await initializeMetadata();
};

// Helper function to get database size
export const getDatabaseSize = async () => {
  if (!navigator.storage || !navigator.storage.estimate) {
    return { usage: 0, quota: 0 };
  }

  const estimate = await navigator.storage.estimate();
  return {
    usage: estimate.usage || 0,
    quota: estimate.quota || 0,
    usageInMB: ((estimate.usage || 0) / (1024 * 1024)).toFixed(2),
    quotaInMB: ((estimate.quota || 0) / (1024 * 1024)).toFixed(2),
  };
};

// Export database instance
export default db;
