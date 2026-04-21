import { useState } from 'react';
import FormField from './FormField';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function SessionForm({ onSuccess, onCancel }) {
  const { data: patients } = useOfflineData('patients');
  const { data: machines } = useOfflineData('dialysis_machines');
  const { data: staff } = useOfflineData('staff_profiles');

  const [formData, setFormData] = useState({
    patient_id: '',
    machine_id: '',
    scheduled_date: new Date().toISOString().split('T')[0],
    scheduled_start_time: '08:00',
    shift: 'morning',
    prescribed_duration_mins: 240,
    modality: 'hd',
    status: 'scheduled',
    primary_nurse_id: '',
    supervising_doctor_id: '',
    session_notes: '',
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const activePatients = (patients || []).filter(p => p.is_active);
  const availableMachines = (machines || []).filter(m => m.operational_status === 'operational');
  const nurses = (staff || []).filter(s => s.cadre === 'nurse' && s.is_active);
  const doctors = (staff || []).filter(s =>
    (s.cadre === 'doctor' || s.cadre === 'nephrologist') && s.is_active
  );

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    setErrors(prev => ({ ...prev, [name]: '' }));
  };

  const validate = () => {
    const newErrors = {};
    if (!formData.patient_id) newErrors.patient_id = 'Patient is required';
    if (!formData.machine_id) newErrors.machine_id = 'Machine is required';
    if (!formData.scheduled_date) newErrors.scheduled_date = 'Date is required';
    if (!formData.scheduled_start_time) newErrors.scheduled_start_time = 'Start time is required';
    if (!formData.primary_nurse_id) newErrors.primary_nurse_id = 'Primary nurse is required';
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
      const payload = {
        patient_id: formData.patient_id,
        machine_id: formData.machine_id,
        scheduled_date: formData.scheduled_date,
        scheduled_start_time: formData.scheduled_start_time,
        shift: formData.shift,
        prescribed_duration_mins: formData.prescribed_duration_mins,
        modality: formData.modality,
        status: 'scheduled',
        primary_nurse_id: formData.primary_nurse_id,
        supervising_doctor_id: formData.supervising_doctor_id || null,
        session_notes: formData.session_notes,
        hospital_id: currentUser?.hospital_id,
        was_patient_reviewed: false
      };

      await offlineService.create('dialysis_sessions', payload, 8);
      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to schedule session' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {errors.submit && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {errors.submit}
        </div>
      )}

      <FormField label="Patient" name="patient_id" type="select"
        value={formData.patient_id} onChange={handleChange} error={errors.patient_id} required
        options={activePatients.map(p => ({
          value: p.id,
          label: `${p.full_name} (${p.mrn})`
        }))} />

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Scheduled Date" name="scheduled_date" type="date"
          value={formData.scheduled_date} onChange={handleChange} error={errors.scheduled_date} required />

        <FormField label="Start Time" name="scheduled_start_time" type="time"
          value={formData.scheduled_start_time} onChange={handleChange}
          error={errors.scheduled_start_time} required />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Shift" name="shift" type="select"
          value={formData.shift} onChange={handleChange} required
          options={[
            { value: 'morning', label: 'Morning' },
            { value: 'afternoon', label: 'Afternoon' },
            { value: 'evening', label: 'Evening' }
          ]} />

        <FormField label="Duration (minutes)" name="prescribed_duration_mins" type="number"
          value={formData.prescribed_duration_mins} onChange={handleChange} required />
      </div>

      <FormField label="Machine" name="machine_id" type="select"
        value={formData.machine_id} onChange={handleChange} error={errors.machine_id} required
        options={availableMachines.map(m => ({
          value: m.id,
          label: `${m.model} - ${m.machine_number}`
        }))} />

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Primary Nurse" name="primary_nurse_id" type="select"
          value={formData.primary_nurse_id} onChange={handleChange} error={errors.primary_nurse_id} required
          options={nurses.map(n => ({
            value: n.id,
            label: n.full_name
          }))} />

        <FormField label="Supervising Doctor (Optional)" name="supervising_doctor_id" type="select"
          value={formData.supervising_doctor_id} onChange={handleChange}
          options={doctors.map(d => ({
            value: d.id,
            label: d.full_name
          }))} />
      </div>

      <FormField label="Modality" name="modality" type="select"
        value={formData.modality} onChange={handleChange} required
        options={[
          { value: 'hd', label: 'Hemodialysis (HD)' },
          { value: 'hdf', label: 'Hemodiafiltration (HDF)' }
        ]} />

      <FormField label="Session Notes" name="session_notes" type="textarea"
        value={formData.session_notes} onChange={handleChange} rows={2} />

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Scheduling...' : 'Schedule Session'}
        </button>
      </div>
    </form>
  );
}
