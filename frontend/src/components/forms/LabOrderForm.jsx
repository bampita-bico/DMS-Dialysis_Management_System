import { useState } from 'react';
import FormField from './FormField';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function LabOrderForm({ onSuccess, onCancel }) {
  const { data: patients } = useOfflineData('patients');
  const { data: labTests } = useOfflineData('lab_test_catalog');

  const [formData, setFormData] = useState({
    patient_id: '',
    priority: 'routine',
    clinical_notes: '',
    diagnosis_code: '',
    selected_tests: []
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const activePatients = (patients || []).filter(p => p.is_active);
  const testsGrouped = (labTests || []).reduce((acc, test) => {
    const category = test.category || 'Other';
    if (!acc[category]) acc[category] = [];
    acc[category].push(test);
    return acc;
  }, {});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    setErrors(prev => ({ ...prev, [name]: '' }));
  };

  const handleTestToggle = (testId) => {
    setFormData(prev => ({
      ...prev,
      selected_tests: prev.selected_tests.includes(testId)
        ? prev.selected_tests.filter(id => id !== testId)
        : [...prev.selected_tests, testId]
    }));
  };

  const validate = () => {
    const newErrors = {};
    if (!formData.patient_id) newErrors.patient_id = 'Patient is required';
    if (formData.selected_tests.length === 0) newErrors.selected_tests = 'Select at least one test';
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

      const orderPayload = {
        patient_id: formData.patient_id,
        ordered_by: currentUser?.id,
        hospital_id: currentUser?.hospital_id,
        priority: formData.priority,
        clinical_notes: formData.clinical_notes,
        diagnosis_code: formData.diagnosis_code,
        status: 'pending',
        order_date: new Date().toISOString().split('T')[0],
        order_time: new Date().toLocaleTimeString('en-GB', { hour12: false })
      };

      const orderId = await offlineService.create('lab_orders', orderPayload, 9);

      for (const testId of formData.selected_tests) {
        await offlineService.create('lab_order_items', {
          order_id: orderId,
          test_id: testId,
          hospital_id: currentUser?.hospital_id,
          specimen_status: 'pending'
        }, 9);
      }

      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to create lab order' });
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
        <FormField label="Priority" name="priority" type="select"
          value={formData.priority} onChange={handleChange} required
          options={[
            { value: 'routine', label: 'Routine' },
            { value: 'urgent', label: 'Urgent' },
            { value: 'stat', label: 'STAT' }
          ]} />

        <FormField label="Diagnosis Code" name="diagnosis_code"
          value={formData.diagnosis_code} onChange={handleChange}
          placeholder="e.g., N18.5" />
      </div>

      <FormField label="Clinical Notes" name="clinical_notes" type="textarea"
        value={formData.clinical_notes} onChange={handleChange} rows={2}
        placeholder="Clinical indication for tests..." />

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Select Lab Tests <span className="text-red-500">*</span>
        </label>
        <div className="max-h-64 overflow-y-auto border border-gray-300 rounded-lg p-4 space-y-3">
          {Object.entries(testsGrouped).map(([category, tests]) => (
            <div key={category}>
              <h4 className="font-semibold text-gray-900 mb-2">{category}</h4>
              <div className="space-y-1 ml-4">
                {tests.map(test => (
                  <label key={test.id} className="flex items-center cursor-pointer hover:bg-gray-50 p-1 rounded">
                    <input
                      type="checkbox"
                      checked={formData.selected_tests.includes(test.id)}
                      onChange={() => handleTestToggle(test.id)}
                      className="h-4 w-4 text-sky-600 rounded mr-2"
                    />
                    <span className="text-sm text-gray-700">{test.name}</span>
                  </label>
                ))}
              </div>
            </div>
          ))}
        </div>
        {errors.selected_tests && <p className="mt-1 text-sm text-red-600">{errors.selected_tests}</p>}
      </div>

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Creating...' : 'Create Lab Order'}
        </button>
      </div>
    </form>
  );
}
