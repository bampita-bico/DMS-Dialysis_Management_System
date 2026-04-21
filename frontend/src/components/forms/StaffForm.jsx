import { useState } from 'react';
import FormField from './FormField';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function StaffForm({ staff = null, onSuccess, onCancel }) {
  const [formData, setFormData] = useState({
    full_name: staff?.full_name || '',
    email: staff?.email || '',
    phone_number: staff?.phone_number || '',
    cadre: staff?.cadre || '',
    specialty: staff?.specialty || '',
    license_number: staff?.license_number || '',
    is_active: staff?.is_active ?? true,
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

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
    if (!formData.email?.trim()) newErrors.email = 'Email is required';
    if (formData.email && !/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Invalid email format';
    }
    if (!formData.cadre) newErrors.cadre = 'Cadre is required';
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
        ...formData,
        hospital_id: currentUser?.hospital_id,
      };

      if (staff) {
        await offlineService.update('staff_profiles', staff.id, payload);
      } else {
        await offlineService.create('staff_profiles', payload, 5);
      }

      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to save staff member' });
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

      <FormField label="Full Name" name="full_name" value={formData.full_name}
        onChange={handleChange} error={errors.full_name} required />

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Email" name="email" type="email" value={formData.email}
          onChange={handleChange} error={errors.email} required />

        <FormField label="Phone Number" name="phone_number" value={formData.phone_number}
          onChange={handleChange} />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Cadre" name="cadre" type="select" value={formData.cadre}
          onChange={handleChange} error={errors.cadre} required
          options={[
            { value: 'doctor', label: 'Doctor' },
            { value: 'nephrologist', label: 'Nephrologist' },
            { value: 'nurse', label: 'Nurse' },
            { value: 'dialysis_technician', label: 'Dialysis Technician' },
            { value: 'pharmacist', label: 'Pharmacist' },
            { value: 'lab_technician', label: 'Lab Technician' },
            { value: 'receptionist', label: 'Receptionist' },
            { value: 'administrator', label: 'Administrator' }
          ]} />

        <FormField label="Specialty" name="specialty" value={formData.specialty}
          onChange={handleChange} placeholder="e.g., Nephrology, Critical Care" />
      </div>

      <FormField label="License Number" name="license_number" value={formData.license_number}
        onChange={handleChange} />

      <div className="flex items-center">
        <input type="checkbox" name="is_active" checked={formData.is_active}
          onChange={handleChange} className="h-4 w-4 text-sky-600 rounded" />
        <label className="ml-2 text-sm text-gray-700">Active Staff Member</label>
      </div>

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Saving...' : staff ? 'Update Staff' : 'Add Staff'}
        </button>
      </div>
    </form>
  );
}
