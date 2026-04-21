import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import useOfflineData from '../../hooks/useOfflineData';
import FormModal from '../../components/forms/FormModal';
import PatientForm from '../../components/forms/PatientForm';

export default function PatientDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { data: patients, loading } = useOfflineData('patients');
  const { data: sessions } = useOfflineData('dialysis_sessions');
  const { data: contacts } = useOfflineData('patient_contacts');
  const { data: vascularAccess } = useOfflineData('vascular_access');

  const [activeTab, setActiveTab] = useState('overview');
  const [showEditModal, setShowEditModal] = useState(false);

  // Support both string and numeric IDs from URL params
  const patient = patients?.find(p => p.id === id || p.id === parseInt(id));
  const patientSessions = sessions?.filter(s => s.patient_id === id || s.patient_id === parseInt(id)) || [];
  const patientContacts = contacts?.filter(c => c.patient_id === id || c.patient_id === parseInt(id)) || [];
  const patientAccess = vascularAccess?.filter(v => v.patient_id === id || v.patient_id === parseInt(id)) || [];

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
      </div>
    );
  }

  if (!patient) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Patient Not Found</h2>
          <p className="text-gray-600 mb-4">The patient you're looking for doesn't exist.</p>
          <Link to="/patients" className="text-sky-600 hover:underline">
            ← Back to Patients
          </Link>
        </div>
      </div>
    );
  }

  const age = patient.date_of_birth
    ? Math.floor((new Date() - new Date(patient.date_of_birth)) / 31557600000)
    : 'N/A';

  const tabs = [
    { id: 'overview', label: 'Overview' },
    { id: 'sessions', label: 'Sessions' },
    { id: 'labs', label: 'Lab Results' },
    { id: 'medications', label: 'Medications' },
    { id: 'access', label: 'Vascular Access' },
  ];

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <button
                onClick={() => navigate('/patients')}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                ← Back
              </button>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">
                  {patient.full_name || 'Unknown Patient'}
                </h1>
                <p className="text-sm text-gray-600 mt-1">
                  MRN: {patient.mrn} • Age: {age} • Sex: {patient.sex}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                patient.is_active
                  ? 'bg-green-100 text-green-800'
                  : 'bg-gray-100 text-gray-800'
              }`}>
                {patient.is_active ? 'Active' : 'Inactive'}
              </span>
              <button
                onClick={() => setShowEditModal(true)}
                className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700">
                Edit Patient
              </button>
            </div>
          </div>

          {/* Tabs */}
          <div className="mt-6 flex space-x-1 border-b border-gray-200">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`px-4 py-2 font-medium transition-colors ${
                  activeTab === tab.id
                    ? 'text-sky-600 border-b-2 border-sky-600'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {activeTab === 'overview' && <OverviewTab patient={patient} contacts={patientContacts} />}
        {activeTab === 'sessions' && <SessionsTab sessions={patientSessions} />}
        {activeTab === 'labs' && <LabsTab patientId={patient.id} />}
        {activeTab === 'medications' && <MedicationsTab patientId={patient.id} />}
        {activeTab === 'access' && <AccessTab patientId={patient.id} accessData={patientAccess} />}
      </div>

      {/* Edit Patient Modal */}
      <FormModal
        isOpen={showEditModal}
        onClose={() => setShowEditModal(false)}
        title="Edit Patient"
        size="lg"
      >
        <PatientForm
          patient={patient}
          onSuccess={() => {
            setShowEditModal(false);
            window.location.reload();
          }}
          onCancel={() => setShowEditModal(false)}
        />
      </FormModal>
    </div>
  );
}

function OverviewTab({ patient, contacts }) {
  const age = patient.date_of_birth
    ? Math.floor((new Date() - new Date(patient.date_of_birth)) / 31557600000)
    : 'N/A';

  const getContact = (type) => contacts.find(c => c.contact_type === type);
  const phoneContact = getContact('phone');
  const emailContact = getContact('email');
  const addressContact = getContact('address');
  const emergencyContact = getContact('emergency');

  // Parse emergency contact value (format: "name|phone|relationship")
  const emergencyParts = emergencyContact?.value?.split('|') || [];
  const emergencyName = emergencyParts[0] || 'Not provided';
  const emergencyPhone = emergencyParts[1] || 'Not provided';

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Demographics */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Demographics</h3>
        <dl className="space-y-3">
          <div>
            <dt className="text-sm text-gray-600">Full Name</dt>
            <dd className="text-base font-medium text-gray-900">{patient.full_name}</dd>
          </div>
          {patient.preferred_name && (
            <div>
              <dt className="text-sm text-gray-600">Preferred Name</dt>
              <dd className="text-base font-medium text-gray-900">{patient.preferred_name}</dd>
            </div>
          )}
          <div>
            <dt className="text-sm text-gray-600">MRN</dt>
            <dd className="text-base font-medium text-gray-900">{patient.mrn}</dd>
          </div>
          {patient.national_id && (
            <div>
              <dt className="text-sm text-gray-600">National ID</dt>
              <dd className="text-base font-medium text-gray-900">{patient.national_id}</dd>
            </div>
          )}
          <div>
            <dt className="text-sm text-gray-600">Date of Birth</dt>
            <dd className="text-base font-medium text-gray-900">
              {patient.date_of_birth
                ? new Date(patient.date_of_birth).toLocaleDateString()
                : 'N/A'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Age</dt>
            <dd className="text-base font-medium text-gray-900">{age} years</dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Sex</dt>
            <dd className="text-base font-medium text-gray-900 capitalize">{patient.sex}</dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Blood Type</dt>
            <dd className="text-base font-medium text-gray-900 uppercase">
              {patient.blood_type?.replace('_', ' ') || 'Unknown'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Marital Status</dt>
            <dd className="text-base font-medium text-gray-900 capitalize">
              {patient.marital_status || 'Unknown'}
            </dd>
          </div>
          {patient.nationality && (
            <div>
              <dt className="text-sm text-gray-600">Nationality</dt>
              <dd className="text-base font-medium text-gray-900">{patient.nationality}</dd>
            </div>
          )}
          {patient.religion && (
            <div>
              <dt className="text-sm text-gray-600">Religion</dt>
              <dd className="text-base font-medium text-gray-900">{patient.religion}</dd>
            </div>
          )}
          {patient.occupation && (
            <div>
              <dt className="text-sm text-gray-600">Occupation</dt>
              <dd className="text-base font-medium text-gray-900">{patient.occupation}</dd>
            </div>
          )}
          {patient.education_level && (
            <div>
              <dt className="text-sm text-gray-600">Education</dt>
              <dd className="text-base font-medium text-gray-900 capitalize">{patient.education_level}</dd>
            </div>
          )}
        </dl>
      </div>

      {/* Contact Information */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Contact Information</h3>
        <dl className="space-y-3">
          <div>
            <dt className="text-sm text-gray-600">Phone</dt>
            <dd className="text-base font-medium text-gray-900">
              {phoneContact?.value || 'Not provided'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Email</dt>
            <dd className="text-base font-medium text-gray-900">
              {emailContact?.value || 'Not provided'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Address</dt>
            <dd className="text-base font-medium text-gray-900">
              {addressContact?.value || 'Not provided'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Emergency Contact</dt>
            <dd className="text-base font-medium text-gray-900">
              {emergencyName}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Emergency Phone</dt>
            <dd className="text-base font-medium text-gray-900">
              {emergencyPhone}
            </dd>
          </div>
        </dl>
      </div>

      {/* Medical Information */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Medical Information</h3>
        <dl className="space-y-3">
          <div>
            <dt className="text-sm text-gray-600">Registration Date</dt>
            <dd className="text-base font-medium text-gray-900">
              {patient.registration_date
                ? new Date(patient.registration_date).toLocaleDateString()
                : 'N/A'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Status</dt>
            <dd className="text-base font-medium text-gray-900">
              {patient.is_active ? 'Active' : 'Inactive'}
            </dd>
          </div>
          <div>
            <dt className="text-sm text-gray-600">Nationality</dt>
            <dd className="text-base font-medium text-gray-900">
              {patient.nationality || 'Not specified'}
            </dd>
          </div>
        </dl>
      </div>
    </div>
  );
}

function SessionsTab({ sessions }) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">
        Dialysis Sessions ({sessions.length})
      </h3>
      {sessions.length === 0 ? (
        <p className="text-gray-600">No sessions recorded yet.</p>
      ) : (
        <div className="space-y-4">
          {sessions.map((session) => (
            <div key={session.id} className="border border-gray-200 rounded-lg p-4">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-medium text-gray-900">
                    {new Date(session.scheduled_date).toLocaleDateString()}
                  </p>
                  <p className="text-sm text-gray-600 mt-1">
                    Duration: {session.actual_duration_mins || session.prescribed_duration_mins} mins
                  </p>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                  session.status === 'completed'
                    ? 'bg-green-100 text-green-800'
                    : session.status === 'scheduled'
                    ? 'bg-blue-100 text-blue-800'
                    : 'bg-gray-100 text-gray-800'
                }`}>
                  {session.status}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function LabsTab({ patientId }) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Lab Results</h3>
      <p className="text-gray-600">Lab results will be displayed here.</p>
    </div>
  );
}

function MedicationsTab({ patientId }) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Medications</h3>
      <p className="text-gray-600">Medication history will be displayed here.</p>
    </div>
  );
}

function AccessTab({ patientId, accessData }) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Vascular Access</h3>
        <button className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 text-sm">
          + Add Access
        </button>
      </div>

      {accessData.length === 0 ? (
        <p className="text-gray-600">No vascular access records found.</p>
      ) : (
        <div className="space-y-4">
          {accessData.map(access => (
            <div key={access.id} className="border border-gray-200 rounded-lg p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="text-lg font-semibold text-gray-900 capitalize">
                      {access.access_type?.replace('_', ' ')}
                    </span>
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                      access.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : access.status === 'maturing'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}>
                      {access.status}
                    </span>
                    {access.is_primary_access && (
                      <span className="px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        Primary
                      </span>
                    )}
                  </div>

                  <dl className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                      <dt className="text-gray-600">Site</dt>
                      <dd className="font-medium text-gray-900 capitalize">
                        {access.access_site?.replace('_', ' ')} ({access.site_side})
                      </dd>
                    </div>
                    <div>
                      <dt className="text-gray-600">Insertion Date</dt>
                      <dd className="font-medium text-gray-900">
                        {access.insertion_date ? new Date(access.insertion_date).toLocaleDateString() : 'N/A'}
                      </dd>
                    </div>
                    {access.inserted_by && (
                      <div>
                        <dt className="text-gray-600">Inserted By</dt>
                        <dd className="font-medium text-gray-900">{access.inserted_by}</dd>
                      </div>
                    )}
                    {access.first_use_date && (
                      <div>
                        <dt className="text-gray-600">First Use</dt>
                        <dd className="font-medium text-gray-900">
                          {new Date(access.first_use_date).toLocaleDateString()}
                        </dd>
                      </div>
                    )}
                    {access.catheter_type && (
                      <div>
                        <dt className="text-gray-600">Catheter Type</dt>
                        <dd className="font-medium text-gray-900">{access.catheter_type}</dd>
                      </div>
                    )}
                    {access.fistula_vein && (
                      <div>
                        <dt className="text-gray-600">Fistula Vein</dt>
                        <dd className="font-medium text-gray-900">{access.fistula_vein}</dd>
                      </div>
                    )}
                  </dl>

                  {access.notes && (
                    <div className="mt-3 pt-3 border-t border-gray-200">
                      <dt className="text-sm text-gray-600 mb-1">Notes</dt>
                      <dd className="text-sm text-gray-900">{access.notes}</dd>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
