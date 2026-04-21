import { useState } from 'react';
import FormField from './FormField';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function PatientForm({ patient = null, onSuccess, onCancel }) {
  const { data: staff } = useOfflineData('staff_profiles');

  const [formData, setFormData] = useState({
    // Basic Information
    full_name: patient?.full_name || '',
    preferred_name: patient?.preferred_name || '',
    mrn: patient?.mrn || '',
    national_id: patient?.national_id || '',
    date_of_birth: patient?.date_of_birth?.split('T')[0] || '',
    sex: patient?.sex || '',
    blood_type: patient?.blood_type || 'unknown',

    // Demographics
    marital_status: patient?.marital_status || 'unknown',
    nationality: patient?.nationality || 'Ugandan',
    religion: patient?.religion || '',
    occupation: patient?.occupation || '',
    education_level: patient?.education_level || '',

    // Language & Communication
    primary_language: patient?.primary_language || 'English',
    interpreter_needed: patient?.interpreter_needed || false,

    // Medical
    primary_doctor_id: patient?.primary_doctor_id || '',

    // Contact Info (will be saved separately)
    phone_number: '',
    email: '',
    physical_address: '',

    // Emergency Contact
    emergency_contact_name: '',
    emergency_contact_phone: '',
    emergency_contact_relationship: ''
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const doctors = (staff || []).filter(s =>
    (s.cadre === 'doctor' || s.cadre === 'nephrologist') && s.is_active
  );

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    setErrors(prev => ({ ...prev, [name]: '' }));
  };

  const validate = () => {
    const newErrors = {};
    if (!formData.full_name?.trim()) newErrors.full_name = 'Full name is required';
    if (!formData.mrn?.trim()) newErrors.mrn = 'MRN is required';
    if (!formData.date_of_birth) newErrors.date_of_birth = 'Date of birth is required';
    if (new Date(formData.date_of_birth) > new Date()) newErrors.date_of_birth = 'Cannot be in future';
    if (!formData.sex) newErrors.sex = 'Sex is required';
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const newErrors = validate();
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setLoading(true);
    try {
      const currentUser = authService.getCurrentUser();

      // Patient main data
      const patientPayload = {
        full_name: formData.full_name,
        preferred_name: formData.preferred_name || null,
        mrn: formData.mrn,
        national_id: formData.national_id || null,
        date_of_birth: formData.date_of_birth,
        sex: formData.sex,
        blood_type: formData.blood_type,
        marital_status: formData.marital_status,
        nationality: formData.nationality,
        religion: formData.religion || null,
        occupation: formData.occupation || null,
        education_level: formData.education_level || null,
        primary_language: formData.primary_language,
        interpreter_needed: formData.interpreter_needed,
        primary_doctor_id: formData.primary_doctor_id || null,
        hospital_id: currentUser?.hospital_id,
        registered_by: currentUser?.id,
        registration_date: new Date().toISOString().split('T')[0],
        is_active: true,
      };

      let patientId;
      if (patient) {
        await offlineService.update('patients', patient.id, patientPayload);
        patientId = patient.id;
      } else {
        patientId = await offlineService.create('patients', patientPayload, 10);
      }

      // Save contacts
      if (formData.phone_number) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'phone',
          value: formData.phone_number,
          label: 'Primary Phone',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      if (formData.email) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'email',
          value: formData.email,
          label: 'Primary Email',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      if (formData.physical_address) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'address',
          value: formData.physical_address,
          label: 'Home Address',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      // Save emergency contact
      if (formData.emergency_contact_name) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'emergency',
          value: `${formData.emergency_contact_name}|${formData.emergency_contact_phone}|${formData.emergency_contact_relationship}`,
          label: 'Emergency Contact',
          is_primary: false,
          is_verified: false
        }, 10);
      }

      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to save patient' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 max-h-[70vh] overflow-y-auto px-1">
      {errors.submit && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {errors.submit}
        </div>
      )}

      {/* Basic Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Basic Information</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Full Name" name="full_name" value={formData.full_name}
            onChange={handleChange} error={errors.full_name} required />

          <FormField label="Preferred Name" name="preferred_name" value={formData.preferred_name}
            onChange={handleChange} placeholder="Nickname or preferred name" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="MRN (Patient Number)" name="mrn" value={formData.mrn}
            onChange={handleChange} error={errors.mrn} required />

          <FormField label="National ID" name="national_id" value={formData.national_id}
            onChange={handleChange} placeholder="National ID or Passport" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Date of Birth" name="date_of_birth" type="date"
            value={formData.date_of_birth} onChange={handleChange} error={errors.date_of_birth} required />

          <FormField label="Sex" name="sex" type="select" value={formData.sex}
            onChange={handleChange} error={errors.sex} required
            options={[
              { value: 'male', label: 'Male' },
              { value: 'female', label: 'Female' },
              { value: 'intersex', label: 'Intersex' },
              { value: 'unknown', label: 'Unknown' }
            ]} />
        </div>

        <FormField label="Blood Type" name="blood_type" type="select"
          value={formData.blood_type} onChange={handleChange} required
          options={[
            { value: 'unknown', label: 'Unknown' },
            { value: 'A+', label: 'A+' }, { value: 'A-', label: 'A-' },
            { value: 'B+', label: 'B+' }, { value: 'B-', label: 'B-' },
            { value: 'AB+', label: 'AB+' }, { value: 'AB-', label: 'AB-' },
            { value: 'O+', label: 'O+' }, { value: 'O-', label: 'O-' }
          ]} />
      </div>

      {/* Demographics */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Demographics</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Marital Status" name="marital_status" type="select"
            value={formData.marital_status} onChange={handleChange}
            options={[
              { value: 'unknown', label: 'Unknown' },
              { value: 'single', label: 'Single' },
              { value: 'married', label: 'Married' },
              { value: 'divorced', label: 'Divorced' },
              { value: 'widowed', label: 'Widowed' }
            ]} />

          <FormField label="Nationality" name="nationality" value={formData.nationality}
            onChange={handleChange} />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Religion" name="religion" value={formData.religion}
            onChange={handleChange} placeholder="e.g., Catholic, Muslim, etc." />

          <FormField label="Occupation" name="occupation" value={formData.occupation}
            onChange={handleChange} />
        </div>

        <FormField label="Education Level" name="education_level" type="select"
          value={formData.education_level} onChange={handleChange}
          options={[
            { value: '', label: 'Not specified' },
            { value: 'none', label: 'None' },
            { value: 'primary', label: 'Primary' },
            { value: 'secondary', label: 'Secondary' },
            { value: 'tertiary', label: 'Tertiary/University' },
            { value: 'postgraduate', label: 'Postgraduate' }
          ]} />
      </div>

      {/* Contact Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Contact Information</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Phone Number" name="phone_number" value={formData.phone_number}
            onChange={handleChange} placeholder="+256..." />

          <FormField label="Email Address" name="email" type="email" value={formData.email}
            onChange={handleChange} />
        </div>

        <FormField label="Physical Address" name="physical_address" type="textarea"
          value={formData.physical_address} onChange={handleChange} rows={2}
          placeholder="Village, Parish, District..." />
      </div>

      {/* Language & Communication */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Language & Communication</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Primary Language" name="primary_language" value={formData.primary_language}
            onChange={handleChange} />

          <div className="flex items-center pt-8">
            <input type="checkbox" name="interpreter_needed" checked={formData.interpreter_needed}
              onChange={handleChange} className="h-4 w-4 text-sky-600 rounded mr-2" />
            <label className="text-sm text-gray-700">Interpreter Needed</label>
          </div>
        </div>
      </div>

      {/* Medical Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Medical Information</h3>

        <FormField label="Primary Doctor" name="primary_doctor_id" type="select"
          value={formData.primary_doctor_id} onChange={handleChange}
          options={doctors.map(d => ({
            value: d.id,
            label: d.full_name
          }))} />
      </div>

      {/* Emergency Contact */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Emergency Contact</h3>

        <FormField label="Emergency Contact Name" name="emergency_contact_name"
          value={formData.emergency_contact_name} onChange={handleChange} />

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Emergency Contact Phone" name="emergency_contact_phone"
            value={formData.emergency_contact_phone} onChange={handleChange} />

          <FormField label="Relationship" name="emergency_contact_relationship"
            value={formData.emergency_contact_relationship} onChange={handleChange}
            placeholder="e.g., Spouse, Parent, Sibling" />
        </div>
      </div>

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 sticky bottom-0 bg-white">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Saving...' : patient ? 'Update Patient' : 'Create Patient'}
        </button>
      </div>
    </form>
  );
}
