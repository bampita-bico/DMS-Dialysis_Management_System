import { useState } from 'react';
import FormField from './FormField';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function InvoiceForm({ onSuccess, onCancel }) {
  const { data: patients } = useOfflineData('patients');
  const { data: priceLists } = useOfflineData('price_lists');
  const { data: sessions } = useOfflineData('dialysis_sessions');

  const [formData, setFormData] = useState({
    patient_id: '',
    session_id: '',
    invoice_date: new Date().toISOString().split('T')[0],
    due_date: '',
    notes: '',
    items: []
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const activePatients = (patients || []).filter(p => p.is_active);
  const patientSessions = formData.patient_id
    ? (sessions || []).filter(s => s.patient_id === formData.patient_id && s.status === 'completed')
    : [];

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    setErrors(prev => ({ ...prev, [name]: '' }));
  };

  const handleItemToggle = (priceItem) => {
    const exists = formData.items.find(i => i.service_code === priceItem.service_code);
    if (exists) {
      setFormData(prev => ({
        ...prev,
        items: prev.items.filter(i => i.service_code !== priceItem.service_code)
      }));
    } else {
      setFormData(prev => ({
        ...prev,
        items: [...prev.items, {
          service_code: priceItem.service_code,
          service_name: priceItem.service_name,
          unit_price: priceItem.unit_price,
          quantity: 1
        }]
      }));
    }
  };

  const handleQuantityChange = (serviceCode, quantity) => {
    setFormData(prev => ({
      ...prev,
      items: prev.items.map(item =>
        item.service_code === serviceCode ? { ...item, quantity: parseInt(quantity) || 1 } : item
      )
    }));
  };

  const calculateTotals = () => {
    const total = formData.items.reduce((sum, item) => sum + (item.unit_price * item.quantity), 0);
    return { total, net: total, balance: total };
  };

  const validate = () => {
    const newErrors = {};
    if (!formData.patient_id) newErrors.patient_id = 'Patient is required';
    if (!formData.invoice_date) newErrors.invoice_date = 'Invoice date is required';
    if (formData.items.length === 0) newErrors.items = 'Add at least one service';
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
      const totals = calculateTotals();

      // Generate invoice number
      const invoiceNumber = `INV-${Date.now()}`;

      const invoicePayload = {
        patient_id: formData.patient_id,
        session_id: formData.session_id || null,
        hospital_id: currentUser?.hospital_id,
        invoice_number: invoiceNumber,
        invoice_date: formData.invoice_date,
        due_date: formData.due_date || null,
        total_amount: totals.total,
        net_amount: totals.net,
        balance_due: totals.balance,
        paid_amount: 0,
        discount_amount: 0,
        tax_amount: 0,
        status: 'pending',
        issued_by: currentUser?.id,
        notes: formData.notes
      };

      const invoiceId = await offlineService.create('invoices', invoicePayload, 7);

      for (const item of formData.items) {
        await offlineService.create('invoice_items', {
          invoice_id: invoiceId,
          service_code: item.service_code,
          description: item.service_name,
          quantity: item.quantity,
          unit_price: item.unit_price,
          total_price: item.unit_price * item.quantity,
          hospital_id: currentUser?.hospital_id
        }, 7);
      }

      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to create invoice' });
    } finally {
      setLoading(false);
    }
  };

  const totals = calculateTotals();

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

      {patientSessions.length > 0 && (
        <FormField label="Related Session (Optional)" name="session_id" type="select"
          value={formData.session_id} onChange={handleChange}
          options={patientSessions.map(s => ({
            value: s.id,
            label: `Session ${new Date(s.scheduled_date).toLocaleDateString()}`
          }))} />
      )}

      <div className="grid grid-cols-2 gap-4">
        <FormField label="Invoice Date" name="invoice_date" type="date"
          value={formData.invoice_date} onChange={handleChange} error={errors.invoice_date} required />

        <FormField label="Due Date" name="due_date" type="date"
          value={formData.due_date} onChange={handleChange} />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Services & Items <span className="text-red-500">*</span>
        </label>
        <div className="max-h-48 overflow-y-auto border border-gray-300 rounded-lg p-3 space-y-2">
          {(priceLists || []).map(item => (
            <label key={item.id} className="flex items-center justify-between cursor-pointer hover:bg-gray-50 p-2 rounded">
              <div className="flex items-center flex-1">
                <input
                  type="checkbox"
                  checked={formData.items.some(i => i.service_code === item.service_code)}
                  onChange={() => handleItemToggle(item)}
                  className="h-4 w-4 text-sky-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">{item.service_name}</span>
              </div>
              <span className="text-sm font-medium text-gray-900">UGX {item.unit_price?.toLocaleString()}</span>
            </label>
          ))}
        </div>
        {errors.items && <p className="mt-1 text-sm text-red-600">{errors.items}</p>}
      </div>

      {formData.items.length > 0 && (
        <div className="border border-gray-200 rounded-lg p-4 bg-gray-50">
          <h4 className="font-medium text-gray-900 mb-3">Selected Items</h4>
          <div className="space-y-2">
            {formData.items.map(item => (
              <div key={item.service_code} className="flex items-center justify-between">
                <span className="text-sm text-gray-700">{item.service_name}</span>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    min="1"
                    value={item.quantity}
                    onChange={(e) => handleQuantityChange(item.service_code, e.target.value)}
                    className="w-16 px-2 py-1 border border-gray-300 rounded text-sm"
                  />
                  <span className="text-sm font-medium">UGX {(item.unit_price * item.quantity).toLocaleString()}</span>
                </div>
              </div>
            ))}
          </div>
          <div className="mt-3 pt-3 border-t border-gray-300">
            <div className="flex justify-between font-bold text-gray-900">
              <span>Total Amount:</span>
              <span>UGX {totals.total.toLocaleString()}</span>
            </div>
          </div>
        </div>
      )}

      <FormField label="Notes" name="notes" type="textarea"
        value={formData.notes} onChange={handleChange} rows={2} />

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Creating...' : 'Create Invoice'}
        </button>
      </div>
    </form>
  );
}
