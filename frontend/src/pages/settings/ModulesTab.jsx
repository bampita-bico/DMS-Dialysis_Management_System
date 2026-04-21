import React, { useState, useEffect } from 'react';
import { useModules } from '../../contexts/ModuleContext';
import EmojiIcon from '../../components/ui/EmojiIcon';

const ModulesTab = () => {
  const { modules, updateModules, getModuleConfig, loading } = useModules();
  const [localModules, setLocalModules] = useState(modules);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState(null);

  useEffect(() => {
    setLocalModules(modules);
  }, [modules]);

  const moduleDefinitions = [
    {
      key: 'lab_management',
      name: 'Lab Management',
      description: 'Full lab orders, results tracking, critical alerts, and test catalog management',
      icon: '🧪',
    },
    {
      key: 'full_pharmacy',
      name: 'Full Pharmacy',
      description: 'Medication stock management, dispensing, drug interactions, and inventory tracking',
      icon: '💊',
    },
    {
      key: 'hr_management',
      name: 'HR Management',
      description: 'Staff schedules, shift assignments, leave management, and performance tracking',
      icon: '👨‍⚕️',
      darkEmoji: true,
    },
    {
      key: 'inventory_tracking',
      name: 'Inventory Tracking',
      description: 'Equipment maintenance, consumables usage, stock alerts, and reorder management',
      icon: '📦',
    },
    {
      key: 'advanced_billing',
      name: 'Advanced Billing',
      description: 'Insurance claims, payment plans, waivers, and detailed financial reporting',
      icon: '💰',
      darkEmoji: true,
    },
    {
      key: 'imaging_integration',
      name: 'Imaging Integration',
      description: 'X-ray, ultrasound, and fistulogram orders with results tracking',
      icon: '🩻',
    },
    {
      key: 'chw_program',
      name: 'Community Health Workers',
      description: 'CHW tracking, home visits, patient transport coordination',
      icon: '🚑',
    },
    {
      key: 'outcomes_reporting',
      name: 'Outcomes & Registry',
      description: 'Clinical outcomes tracking, mortality records, national registry sync',
      icon: '📊',
    },
    {
      key: 'offline_sync',
      name: 'Offline Mode',
      description: 'Work offline with local data storage and automatic cloud sync when online',
      icon: '🔄',
    },
  ];

  const handleToggle = (key) => {
    // All modules available in government edition
    setLocalModules({
      ...localModules,
      [key]: !localModules[key],
    });
  };

  const handleSave = async () => {
    setSaving(true);
    setMessage(null);

    const result = await updateModules(localModules);

    if (result.success) {
      setMessage({ type: 'success', text: 'Modules updated successfully!' });
    } else {
      setMessage({ type: 'error', text: result.error || 'Failed to update modules' });
      setLocalModules(modules); // Revert on error
    }

    setSaving(false);
    setTimeout(() => setMessage(null), 3000);
  };

  const hasChanges = JSON.stringify(localModules) !== JSON.stringify(modules);

  if (loading) {
    return (
      <div className="p-6 text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
        <p className="mt-4 text-gray-600">Loading modules...</p>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-800">Module Configuration</h2>
        <p className="text-sm text-gray-600 mt-1">
          All features are enabled for government deployment. Toggle modules to show/hide them from users.
        </p>
      </div>

      {/* Alert Message */}
      {message && (
        <div className={`mb-4 p-4 rounded-lg ${message.type === 'success' ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'}`}>
          {message.text}
        </div>
      )}

      {/* Module List */}
      <div className="space-y-4">
        {moduleDefinitions.map((module) => {
          const isEnabled = localModules[module.key];

          return (
            <div
              key={module.key}
              className="border rounded-lg p-4 transition-all bg-white border-gray-200"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center">
                    <span className="text-2xl mr-3">
                      <EmojiIcon dark={module.darkEmoji} size="lg">{module.icon}</EmojiIcon>
                    </span>
                    <div>
                      <h3 className="font-medium text-gray-900 flex items-center">
                        {module.name}
                        <span className="ml-2 text-xs bg-green-100 text-green-700 px-2 py-1 rounded">
                          Government Edition
                        </span>
                      </h3>
                      <p className="text-sm text-gray-600 mt-1">{module.description}</p>
                    </div>
                  </div>
                </div>

                {/* Toggle Switch */}
                <button
                  onClick={() => handleToggle(module.key)}
                  className={`
                    relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer
                    ${isEnabled ? 'bg-blue-600' : 'bg-gray-300'}
                  `}
                >
                  <span
                    className={`
                      inline-block h-4 w-4 transform rounded-full bg-white transition-transform
                      ${isEnabled ? 'translate-x-6' : 'translate-x-1'}
                    `}
                  />
                </button>
              </div>
            </div>
          );
        })}
      </div>

      {/* Save Button */}
      {hasChanges && (
        <div className="mt-6 flex items-center justify-end space-x-3">
          <button
            onClick={() => setLocalModules(modules)}
            className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-blue-400 disabled:cursor-not-allowed"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      )}
    </div>
  );
};

export default ModulesTab;
