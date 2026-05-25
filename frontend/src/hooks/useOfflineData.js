import { useCallback, useEffect, useState } from 'react';
import offlineService from '../services/offlineService';

/**
 * Hook for offline-first data fetching
 * @param {string} entityType - Entity table name
 * @param {object} filters - Optional filters
 * @param {boolean} autoRefresh - Auto-refresh when online
 * @returns {object} { data, loading, error, refresh }
 */
export const useOfflineData = (entityType, filters = {}, autoRefresh = true) => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const filterKey = JSON.stringify(filters);

  const loadData = useCallback(async (forceOnline = false) => {
    try {
      setLoading(true);
      setError(null);
      const result = await offlineService.list(entityType, JSON.parse(filterKey || '{}'), forceOnline);
      setData(result);
    } catch (err) {
      console.error(`Failed to load ${entityType}:`, err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [entityType, filterKey]);

  useEffect(() => {
    queueMicrotask(() => loadData());

    const handleLocalChange = (event) => {
      if (!event.detail?.entityType || event.detail.entityType === entityType) {
        loadData(false);
      }
    };
    window.addEventListener('dms-local-change', handleLocalChange);

    // Auto-refresh when coming back online
    if (autoRefresh) {
      const handleOnline = () => {
        console.log(`🔄 Back online, refreshing ${entityType}...`);
        loadData(true); // Force online fetch
      };

      window.addEventListener('online', handleOnline);
      return () => {
        window.removeEventListener('dms-local-change', handleLocalChange);
        window.removeEventListener('online', handleOnline);
      };
    }

    return () => window.removeEventListener('dms-local-change', handleLocalChange);
  }, [autoRefresh, entityType, loadData]);

  return {
    data,
    loading,
    error,
    refresh: () => loadData(true),
  };
};

/**
 * Hook for offline-first single entity fetching
 * @param {string} entityType - Entity table name
 * @param {string} id - Entity ID
 * @returns {object} { data, loading, error, refresh }
 */
export const useOfflineEntity = (entityType, id) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadData = useCallback(async (forceOnline = false) => {
    if (!id) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const result = await offlineService.getById(entityType, id, forceOnline);
      setData(result);
    } catch (err) {
      console.error(`Failed to load ${entityType}:${id}:`, err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [entityType, id]);

  useEffect(() => {
    queueMicrotask(() => loadData());
  }, [loadData]);

  return {
    data,
    loading,
    error,
    refresh: () => loadData(true),
  };
};

export default useOfflineData;
