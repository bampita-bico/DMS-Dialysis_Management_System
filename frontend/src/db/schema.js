import Dexie from 'dexie';

// IndexedDB schema for offline storage
export const db = new Dexie('dms_local');

// Define database schema
db.version(2).stores({
  // Core entities
  patients: 'id, mrn, hospital_id, full_name, *search_terms, updated_at, synced',
  patient_contacts: 'id, patient_id, contact_type, hospital_id, updated_at, synced',

  // Dialysis sessions & machines
  dialysis_sessions: 'id, patient_id, scheduled_date, status, hospital_id, updated_at, synced',
  dialysis_machines: 'id, machine_number, model, operational_status, hospital_id, updated_at, synced',
  session_vitals: 'id, session_id, recorded_at, updated_at, synced',
  session_complications: 'id, session_id, complication_type, severity, updated_at, synced',
  session_fluid_balance: 'id, session_id, updated_at, synced',

  // Vascular access
  vascular_access: 'id, patient_id, access_type, access_status, updated_at, synced',

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
