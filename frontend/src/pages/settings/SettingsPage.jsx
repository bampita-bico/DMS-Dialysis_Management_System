import React, { useState } from 'react';
import ModulesTab from './ModulesTab';
import HospitalTab from './HospitalTab';

const SettingsPage = () => {
  const [activeTab, setActiveTab] = useState('hospital');

  const tabs = [
    { id: 'hospital', label: 'Hospital Info', icon: '🏥' },
    { id: 'modules', label: 'Modules', icon: '⚙️' },
  ];

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Settings</h1>
        <p className="text-gray-600 mt-1">Manage your hospital configuration and features</p>
      </div>

      {/* Tab Navigation */}
      <div className="border-b border-gray-200 mb-6">
        <div className="flex space-x-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`
                pb-4 px-2 font-medium text-sm border-b-2 transition-colors
                ${activeTab === tab.id
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }
              `}
            >
              <span className="mr-2">{tab.icon}</span>
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div className="bg-white rounded-lg shadow">
        {activeTab === 'hospital' && <HospitalTab />}
        {activeTab === 'modules' && <ModulesTab />}
      </div>
    </div>
  );
};

export default SettingsPage;
