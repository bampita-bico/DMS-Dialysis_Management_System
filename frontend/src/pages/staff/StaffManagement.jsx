import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import FormModal from '../../components/forms/FormModal';
import StaffForm from '../../components/forms/StaffForm';

export default function StaffManagement() {
  const [selectedCadre, setSelectedCadre] = useState('all');
  const [selectedDepartment, setSelectedDepartment] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState('cards'); // cards, list
  const [showStaffModal, setShowStaffModal] = useState(false);

  const { data: allStaff, loading } = useOfflineData('staff_profiles');
  const { data: shiftAssignments } = useOfflineData('shift_assignments');

  const handleCreateStaff = async (staffData) => {
    await offlineService.create('staff_profiles', staffData, 5);
  };

  const handleUpdateStaff = async (id, updates) => {
    await offlineService.update('staff_profiles', id, updates);
  };

  const handleAssignShift = async (shiftData) => {
    await offlineService.create('shift_assignments', shiftData, 6);
  };

  const staff = (allStaff || [])
    .filter(s => {
      if (selectedCadre !== 'all' && s.cadre !== selectedCadre) return false;
      if (selectedDepartment !== 'all' && s.department !== selectedDepartment) return false;
      if (searchQuery && !s.full_name?.toLowerCase().includes(searchQuery.toLowerCase())) return false;
      return true;
    })
    .map(user => ({
      ...user,
      name: user.full_name || `${user.first_name || ''} ${user.last_name || ''}`.trim(),
      status: user.is_active ? 'active' : 'inactive'
    }));

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Staff Management
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                Manage healthcare staff, schedules, and shifts
              </p>
            </div>
            <button
              onClick={() => setShowStaffModal(true)}
              className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm">
              Add Staff Member
            </button>
          </div>

          {/* Search and Filters */}
          <div className="mt-6 flex gap-4">
            <div className="flex-1">
              <input
                type="text"
                placeholder="Search by name or employee number..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent"
              />
            </div>

            <select
              value={selectedCadre}
              onChange={(e) => setSelectedCadre(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Cadres</option>
              <option value="doctor">Doctors</option>
              <option value="nurse">Nurses</option>
              <option value="technician">Technicians</option>
              <option value="administrator">Administrators</option>
            </select>

            <select
              value={selectedDepartment}
              onChange={(e) => setSelectedDepartment(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Departments</option>
              <option value="Clinical">Clinical</option>
              <option value="Laboratory">Laboratory</option>
              <option value="Administration">Administration</option>
            </select>

            <div className="flex gap-2">
              <button
                onClick={() => setViewMode('cards')}
                className={`px-4 py-2.5 rounded-lg font-medium transition-colors ${
                  viewMode === 'cards'
                    ? 'bg-sky-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-100'
                }`}
              >
                Cards
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`px-4 py-2.5 rounded-lg font-medium transition-colors ${
                  viewMode === 'list'
                    ? 'bg-sky-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-100'
                }`}
              >
                List
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Staff Grid/List */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
          </div>
        ) : (
          <div className={viewMode === 'cards' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'}>
            {staff.map((member) => (
              <div key={member.id} className="bg-white rounded-xl border border-gray-200 p-6 hover:shadow-md transition-shadow">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{member.name}</h3>
                    <p className="text-sm text-gray-500 mt-1">{member.email}</p>
                    <p className="text-xs text-gray-400 mt-1">{member.cadre || 'Staff'} • {member.department || 'General'}</p>
                  </div>
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                    member.status === 'active' ? 'bg-emerald-100 text-emerald-700' : 'bg-gray-100 text-gray-700'
                  }`}>
                    {member.status.toUpperCase()}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && staff.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No staff members found</p>
            <p className="text-gray-400 text-sm mt-2">Try adjusting your search or filters</p>
          </div>
        )}
      </div>

      <FormModal
        isOpen={showStaffModal}
        onClose={() => setShowStaffModal(false)}
        title="Add Staff Member"
        size="lg"
      >
        <StaffForm
          onSuccess={() => {
            setShowStaffModal(false);
          }}
          onCancel={() => setShowStaffModal(false)}
        />
      </FormModal>
    </div>
  );
}
