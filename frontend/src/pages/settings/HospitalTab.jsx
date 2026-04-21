import React, { useState, useEffect } from 'react';
import { authService } from '../../services/auth';

const HospitalTab = () => {
  const [hospitalInfo, setHospitalInfo] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadHospitalInfo();
  }, []);

  const loadHospitalInfo = () => {
    try {
      const user = authService.getCurrentUser();
      // In a real implementation, fetch full hospital details
      setHospitalInfo({
        name: user?.hospital_name || 'Demo Dialysis Center',
        location: 'Nairobi, Kenya',
        phone: '+254 700 000 000',
        email: 'info@hospital.com',
        admin: user?.full_name || 'Dr. Admin',
        license: 'DL-2024-001',
        established: '2024',
        beds: 20,
        machines: 12,
      });
    } catch (error) {
      console.error('Failed to load hospital info:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6 text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
        <p className="mt-4 text-gray-600">Loading hospital info...</p>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-800">Hospital Information</h2>
        <p className="text-sm text-gray-600 mt-1">
          View your hospital details and contact information
        </p>
      </div>

      {/* Hospital Details */}
      <div className="bg-white border border-gray-200 rounded-lg p-6 mb-6">
        <div className="flex items-start">
          <div className="bg-blue-100 rounded-full p-4 mr-4">
            <span className="text-4xl">🏥</span>
          </div>
          <div className="flex-1">
            <h3 className="text-2xl font-bold text-gray-900">{hospitalInfo.name}</h3>
            <p className="text-gray-600 mt-1">{hospitalInfo.location}</p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
              {/* Contact Information */}
              <div>
                <h4 className="font-semibold text-gray-800 mb-3">Contact Information</h4>
                <div className="space-y-3 text-sm">
                  <div className="flex items-center">
                    <span className="text-gray-500 w-20">Phone:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.phone}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-gray-500 w-20">Email:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.email}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-gray-500 w-20">Admin:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.admin}</span>
                  </div>
                </div>
              </div>

              {/* Facility Details */}
              <div>
                <h4 className="font-semibold text-gray-800 mb-3">Facility Details</h4>
                <div className="space-y-3 text-sm">
                  <div className="flex items-center">
                    <span className="text-gray-500 w-32">License Number:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.license}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-gray-500 w-32">Established:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.established}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-gray-500 w-32">Dialysis Beds:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.beds}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-gray-500 w-32">Dialysis Machines:</span>
                    <span className="text-gray-900 font-medium">{hospitalInfo.machines}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="text-blue-600 text-sm font-medium mb-1">Total Patients</div>
          <div className="text-2xl font-bold text-blue-900">156</div>
        </div>
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <div className="text-green-600 text-sm font-medium mb-1">Active Sessions</div>
          <div className="text-2xl font-bold text-green-900">12</div>
        </div>
        <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
          <div className="text-purple-600 text-sm font-medium mb-1">Staff Members</div>
          <div className="text-2xl font-bold text-purple-900">24</div>
        </div>
        <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
          <div className="text-orange-600 text-sm font-medium mb-1">Uptime</div>
          <div className="text-2xl font-bold text-orange-900">99.8%</div>
        </div>
      </div>

      {/* System Information */}
      <div className="bg-gray-50 border border-gray-200 rounded-lg p-6">
        <h4 className="font-semibold text-gray-900 mb-4">System Information</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-600">DMS Version:</span>
            <span className="font-medium text-gray-900">v2.0.0</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Last Backup:</span>
            <span className="font-medium text-gray-900">2 hours ago</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Database Size:</span>
            <span className="font-medium text-gray-900">1.2 GB</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Storage Used:</span>
            <span className="font-medium text-gray-900">45%</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HospitalTab;
