import db from '../db/schema';
import api from './api';
import { enqueueSync } from '../utils/queueManager';
import { v4 as uuidv4 } from 'uuid';

/**
 * Offline Service - Wrapper for CRUD operations with offline-first approach
 *
 * Pattern:
 * 1. Write operations → Save to IndexedDB immediately, queue for sync
 * 2. Read operations → Try IndexedDB first, fallback to API if online & stale
 */

class OfflineService {
  /**
   * Create entity (offline-first)
   * @param {string} entityType - Table name
   * @param {object} data - Entity data
   * @param {number} priority - Sync priority (1=highest)
   * @returns {object} Created entity with local ID
   */
  async create(entityType, data, priority = 5) {
    // Generate local ID
    const localId = data.id || `local_${uuidv4()}`;

    const entity = {
      ...data,
      id: localId,
      synced: false,
      updated_at: new Date().toISOString(),
      created_at: new Date().toISOString(),
    };

    // Save to IndexedDB immediately
    await db.table(entityType).put(entity);

    // Queue for sync to server
    await enqueueSync(entityType, localId, 'CREATE', entity, priority);

    console.log(`💾 Created ${entityType} locally: ${localId}`);

    return entity;
  }

  /**
   * Update entity (offline-first)
   * @param {string} entityType - Table name
   * @param {string} id - Entity ID
   * @param {object} updates - Fields to update
   * @param {number} priority - Sync priority
   * @returns {object} Updated entity
   */
  async update(entityType, id, updates, priority = 5) {
    // Get existing entity
    const existing = await db.table(entityType).get(id);
    if (!existing) {
      throw new Error(`${entityType} with id ${id} not found`);
    }

    const updated = {
      ...existing,
      ...updates,
      updated_at: new Date().toISOString(),
      synced: false,
    };

    // Update in IndexedDB
    await db.table(entityType).put(updated);

    // Queue for sync
    await enqueueSync(entityType, id, 'UPDATE', updated, priority);

    console.log(`💾 Updated ${entityType} locally: ${id}`);

    return updated;
  }

  /**
   * Delete entity (offline-first)
   * @param {string} entityType - Table name
   * @param {string} id - Entity ID
   * @param {number} priority - Sync priority
   */
  async delete(entityType, id, priority = 5) {
    const entity = await db.table(entityType).get(id);
    if (!entity) {
      throw new Error(`${entityType} with id ${id} not found`);
    }

    // Mark as deleted in local DB (soft delete)
    await db.table(entityType).update(id, {
      deleted_at: new Date().toISOString(),
      synced: false,
    });

    // Queue for sync
    await enqueueSync(entityType, id, 'DELETE', { id }, priority);

    console.log(`💾 Deleted ${entityType} locally: ${id}`);
  }

  /**
   * Get single entity by ID
   * @param {string} entityType - Table name
   * @param {string} id - Entity ID
   * @param {boolean} forceOnline - Force fetch from server
   * @returns {object|null} Entity or null
   */
  async getById(entityType, id, forceOnline = false) {
    // Try local first
    if (!forceOnline) {
      const local = await db.table(entityType).get(id);
      if (local && !local.deleted_at) {
        return local;
      }
    }

    // Fetch from server if online
    if (navigator.onLine) {
      try {
        const endpoint = this.getEndpoint(entityType);
        const response = await api.get(`${endpoint}/${id}`);
        const entity = response.data;

        // Cache in IndexedDB
        await db.table(entityType).put({
          ...entity,
          synced: true,
          updated_at: entity.updated_at || new Date().toISOString(),
        });

        return entity;
      } catch (error) {
        console.error(`Failed to fetch ${entityType}:${id}`, error);
        // Fall back to local copy
        return await db.table(entityType).get(id);
      }
    }

    return null;
  }

  /**
   * List entities with optional filtering
   * @param {string} entityType - Table name
   * @param {object} filters - Filter criteria
   * @param {boolean} forceOnline - Force fetch from server
   * @returns {Array} List of entities
   */
  async list(entityType, filters = {}, forceOnline = false) {
    // Try local first
    if (!forceOnline) {
      let query = db.table(entityType).filter(item => !item.deleted_at);

      // Apply filters
      if (filters.patient_id) {
        query = query.and(item => item.patient_id === filters.patient_id);
      }
      if (filters.status) {
        query = query.and(item => item.status === filters.status);
      }
      if (filters.date) {
        query = query.and(item => item.scheduled_date === filters.date || item.created_at?.startsWith(filters.date));
      }

      const results = await query.toArray();

      // If we have local data and not force online, return it
      if (results.length > 0 || !navigator.onLine) {
        return results;
      }
    }

    // Fetch from server if online
    if (navigator.onLine) {
      try {
        const endpoint = this.getEndpoint(entityType);
        const response = await api.get(endpoint, { params: filters });
        const entities = response.data.data || response.data;

        if (Array.isArray(entities)) {
          // Cache in IndexedDB
          await db.table(entityType).bulkPut(
            entities.map(e => ({
              ...e,
              synced: true,
              updated_at: e.updated_at || new Date().toISOString(),
            }))
          );

          return entities;
        }
      } catch (error) {
        console.error(`Failed to list ${entityType}`, error);
        // Fall back to local
      }
    }

    // Fallback to local data
    return await db.table(entityType)
      .filter(item => !item.deleted_at)
      .toArray();
  }

  /**
   * Search entities (local only for now)
   * @param {string} entityType - Table name
   * @param {string} searchTerm - Search query
   * @returns {Array} Matching entities
   */
  async search(entityType, searchTerm) {
    const lowerSearch = searchTerm.toLowerCase();

    return await db.table(entityType)
      .filter(item => {
        if (item.deleted_at) return false;

        // Search in search_terms array if exists
        if (item.search_terms && Array.isArray(item.search_terms)) {
          return item.search_terms.some(term =>
            term.toLowerCase().includes(lowerSearch)
          );
        }

        // Search in name/title fields
        const searchableFields = [
          item.name, item.title, item.first_name, item.last_name,
          item.mrn, item.generic_name, item.item_name
        ];

        return searchableFields.some(field =>
          field && field.toLowerCase().includes(lowerSearch)
        );
      })
      .toArray();
  }

  /**
   * Get API endpoint for entity type
   */
  getEndpoint(entityType) {
    const endpoints = {
      patients: '/patients',
      dialysis_sessions: '/dialysis-sessions',
      session_vitals: '/vitals',
      session_complications: '/session-complications',
      session_fluid_balance: '/session-fluid-balance',
      vascular_access: '/vascular-access',
      lab_orders: '/lab/orders',
      lab_results: '/lab-results',
      lab_critical_alerts: '/lab-critical-alerts',
      prescriptions: '/prescriptions',
      prescription_items: '/prescriptions',
      invoices: '/invoices',
      payments: '/payments',
      staff_profiles: '/staff-profiles',
      shift_assignments: '/shift-assignments',
    };

    return endpoints[entityType];
  }

  /**
   * Check if entity is synced
   * @param {string} entityType - Table name
   * @param {string} id - Entity ID
   * @returns {boolean} True if synced
   */
  async isSynced(entityType, id) {
    const entity = await db.table(entityType).get(id);
    return entity ? entity.synced === true : false;
  }

  /**
   * Get count of unsynced items for an entity type
   * @param {string} entityType - Table name
   * @returns {number} Count of unsynced items
   */
  async getUnsyncedCount(entityType) {
    return await db.table(entityType)
      .where('synced')
      .equals(false)
      .count();
  }
}

// Export singleton instance
const offlineService = new OfflineService();
export default offlineService;
