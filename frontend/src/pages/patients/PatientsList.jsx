import { useState } from 'react';
import { Link } from 'react-router-dom';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import FormModal from '../../components/forms/FormModal';
import PatientForm from '../../components/forms/PatientForm';

export default function PatientsList() {
  const { data: patients, loading, error } = useOfflineData('patients');
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [sortBy, setSortBy] = useState('name');
  const [showNewPatientModal, setShowNewPatientModal] = useState(false);

  const handleCreatePatient = async (patientData) => {
    await offlineService.create('patients', patientData, 10);
    setShowNewPatientModal(false);
  };

  const handleUpdatePatient = async (id, updates) => {
    await offlineService.update('patients', id, updates);
  };

  const filteredPatients = (patients || [])
    .filter(p => {
      const fullName = p.full_name || '';
      const matchesSearch = fullName.toLowerCase().includes(searchQuery.toLowerCase()) ||
                           (p.mrn || '').toLowerCase().includes(searchQuery.toLowerCase());
      const matchesStatus = filterStatus === 'all' ||
                           (filterStatus === 'active' && p.is_active) ||
                           (filterStatus === 'inactive' && !p.is_active);
      return matchesSearch && matchesStatus;
    })
    .sort((a, b) => {
      const nameA = a.full_name || '';
      const nameB = b.full_name || '';
      if (sortBy === 'name') return nameA.localeCompare(nameB);
      if (sortBy === 'number') return (a.mrn || '').localeCompare(b.mrn || '');
      if (sortBy === 'nextSession') return 0; // Can enhance with actual session lookup
      return 0;
    });

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Patient Directory
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                Manage patient records, sessions, and medical history
              </p>
            </div>
            <button
              onClick={() => setShowNewPatientModal(true)}
              className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm"
            >
              {/* Icon: UserPlus */}
              + New Patient
            </button>
          </div>

          {/* Search and Filters */}
          <div className="mt-8 flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <input
                  type="text"
                  placeholder="Search by name or patient number..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full px-4 py-3 pl-11 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent transition-shadow"
                />
                <div className="absolute left-3 top-3.5 text-gray-400">
                  {/* Icon: MagnifyingGlass */}
                  🔍
                </div>
              </div>
            </div>

            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Patients</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>

            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="name">Sort by Name</option>
              <option value="number">Sort by Number</option>
              <option value="nextSession">Sort by Next Session</option>
            </select>
          </div>
        </div>
      </div>

      {/* Patient Grid */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredPatients.map((patient) => (
              <PatientCard key={patient.id} patient={patient} />
            ))}
          </div>
        )}

        {!loading && filteredPatients.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No patients found</p>
            <p className="text-gray-400 text-sm mt-2">Try adjusting your search or filters</p>
          </div>
        )}
      </div>

      <FormModal
        isOpen={showNewPatientModal}
        onClose={() => setShowNewPatientModal(false)}
        title="Add New Patient"
        size="lg"
      >
        <PatientForm
          onSuccess={() => {
            setShowNewPatientModal(false);
          }}
          onCancel={() => setShowNewPatientModal(false)}
        />
      </FormModal>
    </div>
  );
}

function PatientCard({ patient }) {
  const statusColors = {
    active: 'bg-emerald-50 text-emerald-700 border-emerald-200',
    inactive: 'bg-gray-100 text-gray-600 border-gray-200'
  };

  const fullName = patient.full_name || 'Unknown Patient';
  const age = patient.date_of_birth ?
    Math.floor((new Date() - new Date(patient.date_of_birth)) / 31557600000) :
    'N/A';
  const status = patient.is_active ? 'active' : 'inactive';

  return (
    <Link to={`/patients/${patient.id}`}>
      <div className="bg-white rounded-xl border border-gray-200 p-6 hover:shadow-lg hover:border-sky-300 transition-all duration-200 cursor-pointer group">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 group-hover:text-sky-600 transition-colors">
              {fullName || 'Unknown Patient'}
            </h3>
            <p className="text-sm text-gray-500 font-mono mt-1">{patient.mrn}</p>
          </div>
        </div>

        {/* Patient Info */}
        <div className="space-y-2 mb-4 pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-sm">
            <span className="text-gray-500">Age:</span>
            <span className="text-gray-900 font-medium">{age} • {patient.sex || 'Unknown'}</span>
          </div>
          {patient.primary_diagnosis && (
            <div className="flex items-center gap-2 text-sm">
              <span className="text-gray-500">Diagnosis:</span>
              <span className="text-gray-900 font-medium text-xs">{patient.primary_diagnosis}</span>
            </div>
          )}
          {patient.blood_type && patient.blood_type !== 'unknown' && (
            <div className="flex items-center gap-2 text-sm">
              <span className="text-gray-500">Blood Type:</span>
              <span className="text-gray-900 font-medium text-xs">{patient.blood_type.toUpperCase().replace('_', '')}</span>
            </div>
          )}
        </div>

        {/* Contact */}
        {patient.primary_phone && (
          <div className="mb-4 text-sm">
            <span className="text-gray-500">Phone:</span>
            <span className="text-gray-900 font-medium ml-2">{patient.primary_phone}</span>
          </div>
        )}

        {/* Footer */}
        <div className="flex items-center justify-between">
          <span className={`px-3 py-1 rounded-full text-xs font-medium border ${statusColors[status]}`}>
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </span>
          <span className="text-sky-600 text-sm font-medium group-hover:translate-x-1 transition-transform inline-block">
            View Details →
          </span>
        </div>
      </div>
    </Link>
  );
}
