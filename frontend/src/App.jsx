import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { ModuleProvider, useModules } from './contexts/ModuleContext';
import { authService } from './services/auth';
import syncManager from './db/syncManager';
import SyncIndicator from './components/common/SyncIndicator';
import Login from './pages/auth/Login';
import Dashboard from './pages/dashboard/Dashboard';
import PatientsList from './pages/patients/PatientsList';
import PatientDetails from './pages/patients/PatientDetails';
import SessionSchedule from './pages/sessions/SessionSchedule';
import LabResults from './pages/lab/LabResults';
import Invoices from './pages/billing/Invoices';
import StaffManagement from './pages/staff/StaffManagement';
import SettingsPage from './pages/settings/SettingsPage';
import Layout from './components/layout/Layout';

function ProtectedRoute({ children }) {
  return authService.isAuthenticated() ? (
    <Layout>{children}</Layout>
  ) : (
    <Navigate to="/login" replace />
  );
}

function AppContent() {
  const { isModuleEnabled } = useModules();

  useEffect(() => {
    // Start sync worker when user is authenticated
    if (authService.isAuthenticated()) {
      // Check if offline_sync module is enabled
      if (isModuleEnabled('offline_sync')) {
        syncManager.startBackgroundSync();

        return () => {
          syncManager.stopBackgroundSync();
        };
      }
    }
  }, [isModuleEnabled]);

  return (
    <>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/patients"
          element={
            <ProtectedRoute>
              <PatientsList />
            </ProtectedRoute>
          }
        />
        <Route
          path="/patients/:id"
          element={
            <ProtectedRoute>
              <PatientDetails />
            </ProtectedRoute>
          }
        />
        <Route
          path="/sessions"
          element={
            <ProtectedRoute>
              <SessionSchedule />
            </ProtectedRoute>
          }
        />
        <Route
          path="/lab"
          element={
            <ProtectedRoute>
              <LabResults />
            </ProtectedRoute>
          }
        />
        <Route
          path="/billing"
          element={
            <ProtectedRoute>
              <Invoices />
            </ProtectedRoute>
          }
        />
        <Route
          path="/staff"
          element={
            <ProtectedRoute>
              <StaffManagement />
            </ProtectedRoute>
          }
        />
        <Route
          path="/settings"
          element={
            <ProtectedRoute>
              <SettingsPage />
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Routes>

      {/* Sync Indicator - only show if offline_sync enabled and authenticated */}
      {authService.isAuthenticated() && isModuleEnabled('offline_sync') && <SyncIndicator />}
    </>
  );
}

function App() {
  return (
    <BrowserRouter>
      <ModuleProvider>
        <AppContent />
      </ModuleProvider>
    </BrowserRouter>
  );
}

export default App;
