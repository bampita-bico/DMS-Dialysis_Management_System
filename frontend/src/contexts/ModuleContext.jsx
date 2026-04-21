import React, { createContext, useContext, useState } from 'react';

const ModuleContext = createContext(null);

export const ModuleProvider = ({ children }) => {
  // All modules enabled by default (no subscription tiers for government deployment)
  const [modules] = useState({
    lab_management: true,
    full_pharmacy: true,
    hr_management: true,
    inventory_tracking: true,
    advanced_billing: true,
    offline_sync: true,
    chw_program: true,
    imaging_integration: true,
    outcomes_reporting: true,
  });

  const loading = false;
  const error = null;

  const updateModules = async (newModules) => {
    // No-op: all modules always enabled for government deployment
    console.log('All modules are enabled by default for government deployment');
    return { success: true };
  };

  const isModuleEnabled = (moduleName) => {
    // All modules always enabled
    return modules[moduleName] === true;
  };

  const getModuleConfig = (moduleName) => {
    return {
      enabled: true,
      available: true,
    };
  };

  const value = {
    modules,
    loading,
    error,
    isModuleEnabled,
    getModuleConfig,
    updateModules,
    refreshModules: () => {},
  };

  return (
    <ModuleContext.Provider value={value}>
      {children}
    </ModuleContext.Provider>
  );
};

export const useModules = () => {
  const context = useContext(ModuleContext);
  if (!context) {
    throw new Error('useModules must be used within ModuleProvider');
  }
  return context;
};
