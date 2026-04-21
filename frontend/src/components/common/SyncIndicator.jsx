import React, { useState, useEffect } from 'react';
import useOnlineStatus from '../../hooks/useOnlineStatus';
import syncManager from '../../db/syncManager';
import { getSyncStats } from '../../utils/queueManager';

const SyncIndicator = () => {
  const isOnline = useOnlineStatus();
  const [syncStats, setSyncStats] = useState({ pending: 0, failed: 0, synced: 0 });
  const [isSyncing, setIsSyncing] = useState(false);
  const [lastSyncTime, setLastSyncTime] = useState(null);

  useEffect(() => {
    // Update stats every 5 seconds
    const updateStats = async () => {
      const stats = await getSyncStats();
      setSyncStats(stats);

      const status = syncManager.getStatus();
      setIsSyncing(status.isSyncing);
      setLastSyncTime(status.lastSyncTime);
    };

    updateStats();
    const interval = setInterval(updateStats, 5000);

    return () => clearInterval(interval);
  }, []);

  const getStatusColor = () => {
    if (!isOnline) return 'bg-red-500';
    if (isSyncing) return 'bg-yellow-500 animate-pulse';
    if (syncStats.pending > 0) return 'bg-orange-500';
    return 'bg-green-500';
  };

  const getStatusText = () => {
    if (!isOnline) return 'Offline';
    if (isSyncing) return 'Syncing...';
    if (syncStats.pending > 0) return `${syncStats.pending} pending`;
    return 'Synced';
  };

  const getStatusIcon = () => {
    if (!isOnline) return '📡';
    if (isSyncing) return '🔄';
    if (syncStats.pending > 0) return '⏳';
    return '✅';
  };

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <div className="bg-white rounded-lg shadow-lg border border-gray-200 p-3 min-w-[200px]">
        {/* Status Badge */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center">
            <div className={`w-2 h-2 rounded-full ${getStatusColor()} mr-2`}></div>
            <span className="text-sm font-medium text-gray-700">{getStatusText()}</span>
          </div>
          <span className="text-lg">{getStatusIcon()}</span>
        </div>

        {/* Stats */}
        {syncStats.pending > 0 && (
          <div className="text-xs text-gray-600 space-y-1">
            <div className="flex justify-between">
              <span>Pending:</span>
              <span className="font-medium">{syncStats.pending}</span>
            </div>
            {syncStats.failed > 0 && (
              <div className="flex justify-between text-red-600">
                <span>Failed:</span>
                <span className="font-medium">{syncStats.failed}</span>
              </div>
            )}
          </div>
        )}

        {/* Last Sync Time */}
        {lastSyncTime && isOnline && syncStats.pending === 0 && (
          <div className="text-xs text-gray-500 mt-2">
            Last synced: {lastSyncTime.toLocaleTimeString()}
          </div>
        )}

        {/* Offline Warning */}
        {!isOnline && (
          <div className="mt-2 text-xs text-red-600 bg-red-50 rounded p-2">
            Working offline. Changes will sync when online.
          </div>
        )}
      </div>
    </div>
  );
};

export default SyncIndicator;
