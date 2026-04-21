import db from './schema';
import api from '../services/api';
import { enqueueSync, getPendingSyncItems, markSyncCompleted, markSyncFailed } from '../utils/queueManager';

/**
 * Sync Manager - Orchestrates data synchronization between IndexedDB and server
 */

class SyncManager {
  constructor() {
    this.isSyncing = false;
    this.syncInterval = null;
    this.lastSyncTime = null;
  }

  /**
   * Start background sync worker
   * Syncs every 10 seconds when online
   */
  startBackgroundSync() {
    if (this.syncInterval) {
      console.log('⚠️ Background sync already running');
      return;
    }

    console.log('🔄 Starting background sync worker');

    // Immediate sync on start
    this.sync();

    // Then sync every 10 seconds
    this.syncInterval = setInterval(() => {
      if (navigator.onLine && !this.isSyncing) {
        this.sync();
      }
    }, 10000);
  }

  /**
   * Stop background sync worker
   */
  stopBackgroundSync() {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
      this.syncInterval = null;
      console.log('🛑 Background sync stopped');
    }
  }

  /**
   * Main sync function - pushes local changes to server
   */
  async sync() {
    if (this.isSyncing) {
      return; // Already syncing, skip
    }

    if (!navigator.onLine) {
      console.log('📡 Offline - skipping sync');
      return;
    }

    this.isSyncing = true;

    try {
      const pendingItems = await getPendingSyncItems(50); // Process 50 items per batch

      if (pendingItems.length === 0) {
        this.lastSyncTime = new Date();
        this.isSyncing = false;
        return;
      }

      console.log(`📤 Syncing ${pendingItems.length} items...`);

      for (const item of pendingItems) {
        try {
          await this.syncItem(item);
          await markSyncCompleted(item.id);
        } catch (error) {
          console.error(`Failed to sync ${item.entity_type}:${item.entity_id}`, error);
          await markSyncFailed(item.id, error.message);
        }
      }

      this.lastSyncTime = new Date();
      console.log(`✅ Sync completed at ${this.lastSyncTime.toLocaleTimeString()}`);

    } catch (error) {
      console.error('Sync error:', error);
    } finally {
      this.isSyncing = false;
    }
  }

  /**
   * Sync individual item to server
   * @param {Object} item - Sync queue item
   */
  async syncItem(item) {
    const { entity_type, entity_id, operation, payload } = item;

    const endpoint = this.getEndpoint(entity_type);
    if (!endpoint) {
      throw new Error(`No endpoint configured for ${entity_type}`);
    }

    let response;

    switch (operation) {
      case 'CREATE':
        response = await api.post(endpoint, payload);
        // Update local record with server-assigned ID
        if (response.data.id && response.data.id !== entity_id) {
          await this.updateLocalId(entity_type, entity_id, response.data.id);
        }
        break;

      case 'UPDATE':
        response = await api.patch(`${endpoint}/${entity_id}`, payload);
        break;

      case 'DELETE':
        response = await api.delete(`${endpoint}/${entity_id}`);
        // Remove from local DB
        await db.table(entity_type).delete(entity_id);
        break;

      default:
        throw new Error(`Unknown operation: ${operation}`);
    }

    // Mark local record as synced
    if (operation !== 'DELETE') {
      await db.table(entity_type).update(entity_id, {
        synced: true,
        updated_at: new Date().toISOString(),
      });
    }

    return response;
  }

  /**
   * Update local ID to match server ID
   * @param {string} tableName - Entity table
   * @param {string} localId - Temporary local ID
   * @param {string} serverId - Server-assigned ID
   */
  async updateLocalId(tableName, localId, serverId) {
    const record = await db.table(tableName).get(localId);
    if (record) {
      await db.table(tableName).delete(localId);
      record.id = serverId;
      await db.table(tableName).put(record);
      console.log(`Updated local ID: ${localId} → ${serverId}`);
    }
  }

  /**
   * Get API endpoint for entity type
   * @param {string} entityType - Entity table name
   * @returns {string} API endpoint path
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
      prescription_items: '/prescriptions', // Nested under prescriptions
      invoices: '/invoices',
      payments: '/payments',
      staff_profiles: '/staff-profiles',
      shift_assignments: '/shift-assignments',
    };

    return endpoints[entityType];
  }

  /**
   * Pull fresh data from server (full sync down)
   * @param {string} entityType - Entity to sync
   * @param {Date} since - Only get records updated after this date
   */
  async pullFromServer(entityType, since = null) {
    const endpoint = this.getEndpoint(entityType);
    if (!endpoint) {
      console.warn(`No endpoint for ${entityType}, skipping pull`);
      return;
    }

    try {
      const params = since ? { updated_since: since.toISOString() } : {};
      const response = await api.get(endpoint, { params });

      const records = response.data.data || response.data;

      if (Array.isArray(records) && records.length > 0) {
        // Upsert records into local DB
        await db.table(entityType).bulkPut(
          records.map(record => ({
            ...record,
            synced: true, // Server data is already synced
            updated_at: record.updated_at || new Date().toISOString(),
          }))
        );

        // Update metadata
        await db._metadata.put({
          key: `last_sync_${entityType}`,
          value: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        });

        console.log(`⬇️ Pulled ${records.length} ${entityType} from server`);
      }
    } catch (error) {
      console.error(`Failed to pull ${entityType}:`, error);
      throw error;
    }
  }

  /**
   * Initial data load - pulls all reference data and recent entities
   */
  async initialSync() {
    console.log('🔄 Starting initial sync...');

    try {
      // Pull reference data (read-only)
      await this.pullFromServer('lab_test_catalog');
      await this.pullFromServer('lab_panels');
      await this.pullFromServer('medications');
      await this.pullFromServer('consumables');
      await this.pullFromServer('insurance_schemes');
      await this.pullFromServer('price_lists');

      // Pull recent operational data (last 30 days)
      const thirtyDaysAgo = new Date();
      thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

      await this.pullFromServer('patients', thirtyDaysAgo);
      await this.pullFromServer('dialysis_sessions', thirtyDaysAgo);
      await this.pullFromServer('lab_orders', thirtyDaysAgo);
      await this.pullFromServer('prescriptions', thirtyDaysAgo);
      await this.pullFromServer('invoices', thirtyDaysAgo);

      console.log('✅ Initial sync completed');
    } catch (error) {
      console.error('Initial sync failed:', error);
      throw error;
    }
  }

  /**
   * Get sync status
   * @returns {Object} Current sync status
   */
  getStatus() {
    return {
      isSyncing: this.isSyncing,
      isRunning: this.syncInterval !== null,
      lastSyncTime: this.lastSyncTime,
      isOnline: navigator.onLine,
    };
  }
}

// Export singleton instance
const syncManager = new SyncManager();
export default syncManager;
