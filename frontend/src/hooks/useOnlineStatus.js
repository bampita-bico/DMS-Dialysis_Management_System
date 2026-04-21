import { useState, useEffect } from 'react';

/**
 * Hook to detect online/offline status
 * Returns true if online, false if offline
 */
export const useOnlineStatus = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => {
      console.log('📡 Connection restored - Back online');
      setIsOnline(true);
    };

    const handleOffline = () => {
      console.log('📡 Connection lost - Working offline');
      setIsOnline(false);
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // Periodic connectivity check (every 30 seconds)
    const intervalId = setInterval(() => {
      const currentStatus = navigator.onLine;
      if (currentStatus !== isOnline) {
        setIsOnline(currentStatus);
      }
    }, 30000);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      clearInterval(intervalId);
    };
  }, [isOnline]);

  return isOnline;
};

export default useOnlineStatus;
