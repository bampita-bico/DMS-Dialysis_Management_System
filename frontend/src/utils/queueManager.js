import db from '../db/schema';

/**
 * Queue Manager - Manages local changes pending sync to server
 */

/**
 * Add operation to sync queue
 * @param {string} entityType - Table name (e.g., 'patients', 'dialysis_sessions')
 * @param {string} entityId - Local ID of the entity
 * @param {string} operation - 'CREATE', 'UPDATE', or 'DELETE'
 * @param {object} payload - The entity data to sync
 * @param {number} priority - Priority (1=highest, 10=lowest)
 */
export const enqueueSync = async (entityType, entityId, operation, payload, priority = 5) => {
  try {
    await db.sync_queue.add({
      entity_type: entityType,
      entity_id: entityId,
      operation: operation,
      payload: payload,
      synced: false,
      priority: priority,
      created_at: new Date().toISOString(),
      attempts: 0,
    });

    console.log(`📤 Queued ${operation} for ${entityType}:${entityId}`);
  } catch (error) {
    console.error('Failed to enqueue sync:', error);
    throw error;
  }
};

/**
 * Get pending sync queue items
 * @param {number} limit - Max items to retrieve
 * @returns {Array} Pending sync items ordered by priority
 */
export const getPendingSyncItems = async (limit = 100) => {
  return await db.sync_queue
    .where('synced')
    .equals(false)
    .and(item => item.attempts < 3) // Skip items that failed 3+ times
    .sortBy('priority');
};

/**
 * Mark sync item as completed
 * @param {number} queueId - Sync queue item ID
 */
export const markSyncCompleted = async (queueId) => {
  await db.sync_queue.update(queueId, {
    synced: true,
    synced_at: new Date().toISOString(),
  });
};

/**
 * Mark sync item as failed
 * @param {number} queueId - Sync queue item ID
 * @param {string} error - Error message
 */
export const markSyncFailed = async (queueId, error) => {
  const item = await db.sync_queue.get(queueId);
  await db.sync_queue.update(queueId, {
    attempts: (item.attempts || 0) + 1,
    last_error: error,
    last_attempt_at: new Date().toISOString(),
  });
};

/**
 * Get sync queue statistics
 * @returns {Object} Stats about sync queue
 */
export const getSyncStats = async () => {
  const pending = await db.sync_queue.where('synced').equals(false).count();
  const failed = await db.sync_queue.where('attempts').above(2).count();
  const total = await db.sync_queue.count();

  return {
    pending,
    failed,
    synced: total - pending,
    total,
  };
};

/**
 * Clear synced items older than N days
 * @param {number} daysOld - Age threshold in days
 */
export const cleanupSyncQueue = async (daysOld = 7) => {
  const cutoffDate = new Date();
  cutoffDate.setDate(cutoffDate.getDate() - daysOld);

  const oldItems = await db.sync_queue
    .where('synced')
    .equals(true)
    .and(item => new Date(item.synced_at) < cutoffDate)
    .toArray();

  const ids = oldItems.map(item => item.id);
  await db.sync_queue.bulkDelete(ids);

  console.log(`🧹 Cleaned up ${ids.length} old sync queue items`);
};

/**
 * Retry failed sync items (reset attempts to 0)
 * @returns {number} Number of items reset
 */
export const retryFailedSyncs = async () => {
  const failedItems = await db.sync_queue
    .where('attempts')
    .between(1, 2) // Only retry items with 1-2 attempts
    .and(item => !item.synced)
    .toArray();

  for (const item of failedItems) {
    await db.sync_queue.update(item.id, { attempts: 0 });
  }

  console.log(`🔄 Reset ${failedItems.length} failed sync items for retry`);
  return failedItems.length;
};

/**
 * Get sync queue items by entity type
 * @param {string} entityType - Table name
 * @returns {Array} Sync items for that entity type
 */
export const getSyncItemsByType = async (entityType) => {
  return await db.sync_queue
    .where('entity_type')
    .equals(entityType)
    .and(item => !item.synced)
    .toArray();
};

export default {
  enqueueSync,
  getPendingSyncItems,
  markSyncCompleted,
  markSyncFailed,
  getSyncStats,
  cleanupSyncQueue,
  retryFailedSyncs,
  getSyncItemsByType,
};
