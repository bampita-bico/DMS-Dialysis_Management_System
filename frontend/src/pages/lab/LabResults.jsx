import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import FormModal from '../../components/forms/FormModal';
import LabOrderForm from '../../components/forms/LabOrderForm';
import db from '../../db/schema';
import { authService } from '../../services/auth';
import { htmlTable, printHtml } from '../../utils/print';

export default function LabResults() {
  const [selectedPatient, setSelectedPatient] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [dateRange, setDateRange] = useState('7days');
  const [showNewOrderModal, setShowNewOrderModal] = useState(false);
  const [showQuickResultModal, setShowQuickResultModal] = useState(false);

  const { data: labResults, loading } = useOfflineData('lab_results');
  const { data: labOrders } = useOfflineData('lab_orders');
  const { data: labTests } = useOfflineData('lab_test_catalog');
  const { data: patients } = useOfflineData('patients');

  const results = (labResults || [])
    .filter(r => {
      if (selectedPatient !== 'all' && String(r.patient_id) !== String(selectedPatient)) return false;
      if (filterStatus !== 'all' && normalizeStatus(r.result_status || r.status) !== filterStatus) return false;
      if (!inDateRange(r.result_date || r.created_at, dateRange)) return false;
      return true;
    })
    .map(result => {
      const patient = patients?.find(p => String(p.id) === String(result.patient_id));
      const order = labOrders?.find(o => String(o.id) === String(result.lab_order_id || result.order_id));
      const test = labTests?.find(t => String(t.id) === String(result.test_id));
      return {
        ...result,
        patientName: patient?.full_name || 'Unknown',
        patientNumber: patient?.mrn || 'N/A',
        testName: result.test_name || test?.name || result.test_code || 'Lab result',
        status: normalizeStatus(result.result_status || result.status || 'pending'),
        orderDate: order?.created_at || result.result_date || result.created_at
      };
    });

  const handlePrintResults = () => {
    printHtml({
      title: 'Laboratory Results',
      subtitle: selectedPatient === 'all'
        ? `All patients | ${dateRangeLabel(dateRange)}`
        : `${results[0]?.patientName || 'Selected patient'} | ${dateRangeLabel(dateRange)}`,
      body: buildLabResultsPrintBody(results),
      footer: 'Printed from DMS laboratory records.',
    });
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Lab Results
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                View and manage laboratory test results
              </p>
            </div>
            <div className="flex flex-wrap gap-2">
              <button
                onClick={handlePrintResults}
                disabled={results.length === 0}
                className="px-5 py-3 border border-gray-300 bg-white text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
              >
                Print Results
              </button>
              <button
                onClick={() => setShowQuickResultModal(true)}
                className="px-5 py-3 border border-sky-200 bg-sky-50 text-sky-700 font-medium rounded-lg hover:bg-sky-100 transition-colors"
              >
                Record Result
              </button>
              <button
                onClick={() => setShowNewOrderModal(true)}
                className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm"
              >
                New Lab Order
              </button>
            </div>
          </div>

          {/* Filters */}
          <div className="mt-6 flex gap-4">
            <select
              value={selectedPatient}
              onChange={(e) => setSelectedPatient(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Patients</option>
              {patients?.map(p => (
                <option key={p.id} value={p.id}>
                  {p.full_name} ({p.mrn})
                </option>
              ))}
            </select>

            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Status</option>
              <option value="pending">Pending</option>
              <option value="collected">Collected</option>
              <option value="processing">Processing</option>
              <option value="completed">Completed</option>
            </select>

            <select
              value={dateRange}
              onChange={(e) => setDateRange(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="7days">Last 7 Days</option>
              <option value="30days">Last 30 Days</option>
              <option value="90days">Last 90 Days</option>
              <option value="all">All Time</option>
            </select>
          </div>
        </div>
      </div>

      {/* Results List */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
          </div>
        ) : (
          <div className="space-y-4">
            {results.map((result) => (
              <div key={result.id} className="bg-white rounded-xl border border-gray-200 p-6 hover:shadow-md transition-shadow">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{result.patientName}</h3>
                    <p className="text-sm text-gray-500 mt-1">MRN: {result.patientNumber}</p>
                  </div>
                  <span className="px-3 py-1 rounded-full text-xs font-medium bg-sky-100 text-sky-700">
                    {result.status.toUpperCase()}
                  </span>
                </div>
                <div className="mt-4 text-sm text-gray-600">
                  <p>Order Date: {result.orderDate ? new Date(result.orderDate).toLocaleDateString() : 'N/A'}</p>
                  <p className="mt-1">
                    {result.testName}: <span className="font-semibold text-gray-900">{formatResultValue(result)}</span>
                    <span className="ml-2 text-xs text-sky-700">{phaseLabel(result.result_phase)}</span>
                  </p>
                  {result.notes && <p className="mt-1 text-gray-500">{result.notes}</p>}
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && results.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No lab results found</p>
            <p className="text-gray-400 text-sm mt-2">Try adjusting your filters or create a new lab order</p>
          </div>
        )}
      </div>

      <FormModal
        isOpen={showNewOrderModal}
        onClose={() => setShowNewOrderModal(false)}
        title="Create Lab Order"
        size="lg"
      >
        <LabOrderForm
          onSuccess={() => {
            setShowNewOrderModal(false);
          }}
          onCancel={() => setShowNewOrderModal(false)}
        />
      </FormModal>

      <FormModal
        isOpen={showQuickResultModal}
        onClose={() => setShowQuickResultModal(false)}
        title="Record Lab Result"
        size="lg"
      >
        <QuickResultForm
          patients={patients || []}
          labTests={labTests || []}
          onSuccess={() => setShowQuickResultModal(false)}
          onCancel={() => setShowQuickResultModal(false)}
        />
      </FormModal>
    </div>
  );
}

function QuickResultForm({ patients, labTests, onSuccess, onCancel }) {
  const [formData, setFormData] = useState({
    patient_id: '',
    test_id: '',
    result_phase: 'pre_dialysis',
    value_numeric: '',
    value_text: '',
    unit: '',
    result_status: 'final',
    notes: '',
  });
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const selectedTest = labTests.find(test => String(test.id) === String(formData.test_id));

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
      ...(name === 'test_id' ? { unit: labTests.find(test => String(test.id) === String(value))?.unit || labTests.find(test => String(test.id) === String(value))?.default_unit || prev.unit } : {}),
    }));
    setError('');
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!formData.patient_id || !formData.test_id || (!formData.value_numeric && !formData.value_text)) {
      setError('Patient, test, and a numeric or text result are required.');
      return;
    }

    setSaving(true);
    try {
      const currentUser = authService.getCurrentUser();
      const numeric = formData.value_numeric === '' ? null : Number(formData.value_numeric);
      const result = {
        id: `local_lab_${Date.now()}`,
        patient_id: formData.patient_id,
        hospital_id: currentUser?.hospital_id || 'demo_hospital',
        test_id: formData.test_id,
        test_name: selectedTest?.name || 'Lab result',
        test_code: selectedTest?.code || '',
        result_phase: formData.result_phase,
        value_numeric: numeric,
        result_value: numeric,
        value_text: formData.value_text,
        unit: formData.unit || selectedTest?.unit || selectedTest?.default_unit || '',
        result_status: formData.result_status,
        status: formData.result_status,
        is_abnormal: isAbnormal(numeric, selectedTest),
        is_critical: isCritical(numeric, selectedTest),
        result_date: new Date().toISOString(),
        entered_by: currentUser?.id,
        notes: formData.notes,
        synced: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      await db.lab_results.put(result);
      window.dispatchEvent(new CustomEvent('dms-local-change', {
        detail: { entityType: 'lab_results', record: result },
      }));
      onSuccess?.();
    } catch (err) {
      setError(err.message || 'Failed to record result.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <SelectField label="Patient" name="patient_id" value={formData.patient_id} onChange={handleChange} required
          options={patients.filter(p => p.is_active).map(patient => ({ value: patient.id, label: `${patient.full_name} (${patient.mrn})` }))} />
        <SelectField label="Test" name="test_id" value={formData.test_id} onChange={handleChange} required
          options={labTests.map(test => ({ value: test.id, label: `${test.name} (${test.code})` }))} />
        <SelectField label="Clinical State" name="result_phase" value={formData.result_phase} onChange={handleChange}
          options={[
            { value: 'pre_dialysis', label: 'Pre-dialysis' },
            { value: 'routine', label: 'Routine/off-day' },
            { value: 'intra_dialysis', label: 'Intra-dialysis' },
            { value: 'intra_complication', label: 'Intra-dialysis complication' },
            { value: 'post_dialysis', label: 'Post-dialysis' },
            { value: 'off_session', label: 'Off-session urgent' },
          ]} />
        <SelectField label="Status" name="result_status" value={formData.result_status} onChange={handleChange}
          options={[
            { value: 'preliminary', label: 'Preliminary' },
            { value: 'final', label: 'Final' },
            { value: 'corrected', label: 'Corrected' },
          ]} />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <InputField label="Numeric Value" name="value_numeric" type="number" step="0.01" value={formData.value_numeric} onChange={handleChange} />
        <InputField label="Text Value" name="value_text" value={formData.value_text} onChange={handleChange} />
        <InputField label="Unit" name="unit" value={formData.unit} onChange={handleChange} />
      </div>
      <label className="block">
        <span className="block text-sm font-medium text-gray-700 mb-1">Notes</span>
        <textarea name="notes" value={formData.notes} onChange={handleChange} rows={3}
          className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500" />
      </label>
      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={saving} className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">Cancel</button>
        <button type="submit" disabled={saving} className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {saving ? 'Saving...' : 'Save Result'}
        </button>
      </div>
    </form>
  );
}

function SelectField({ label, name, value, onChange, options, required = false }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}{required && <span className="text-red-500 ml-1">*</span>}</span>
      <select name={name} value={value} onChange={onChange} required={required}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500">
        <option value="">Select {label}</option>
        {options.map(option => <option key={option.value} value={option.value}>{option.label}</option>)}
      </select>
    </label>
  );
}

function InputField({ label, name, value, onChange, type = 'text', step }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}</span>
      <input name={name} type={type} step={step} value={value} onChange={onChange}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500" />
    </label>
  );
}

function normalizeStatus(value) {
  return String(value || 'pending').replace('-', '_');
}

function phaseLabel(value) {
  const labels = {
    pre_dialysis: 'Pre-dialysis',
    routine: 'Routine',
    intra_dialysis: 'Intra-dialysis',
    intra_complication: 'Intra-complication',
    post_dialysis: 'Post-dialysis',
    off_session: 'Off-session',
  };
  return labels[value] || 'Routine';
}

function formatResultValue(result) {
  if (result.value_text) return result.value_text;
  const value = result.value_numeric ?? result.result_value ?? result.value;
  if (value === null || value === undefined || value === '') return 'Not recorded';
  return `${value} ${result.unit || ''}`.trim();
}

function buildLabResultsPrintBody(results) {
  if (results.length === 0) {
    return '<p>No lab results found for the selected filters.</p>';
  }

  return htmlTable(
    ['Date', 'Patient', 'MRN', 'Test', 'Value', 'Clinical State', 'Status', 'Notes'],
    results.map(result => [
      result.orderDate ? new Date(result.orderDate).toLocaleDateString() : 'N/A',
      result.patientName,
      result.patientNumber,
      result.testName,
      formatResultValue(result),
      phaseLabel(result.result_phase),
      result.status,
      result.notes || '',
    ])
  );
}

function dateRangeLabel(range) {
  const labels = {
    '7days': 'Last 7 days',
    '30days': 'Last 30 days',
    '90days': 'Last 90 days',
    all: 'All time',
  };
  return labels[range] || 'Selected dates';
}

function inDateRange(value, range) {
  if (range === 'all' || !value) return true;
  const time = new Date(value).getTime();
  if (Number.isNaN(time)) return true;
  const days = range === '7days' ? 7 : range === '30days' ? 30 : 90;
  return time >= Date.now() - days * 86400000;
}

function isAbnormal(value, test) {
  if (value === null || value === undefined || !test) return false;
  const low = Number(test.low ?? test.normal_low);
  const high = Number(test.high ?? test.normal_high);
  if (!Number.isNaN(low) && value < low) return true;
  if (!Number.isNaN(high) && value > high) return true;
  return false;
}

function isCritical(value, test) {
  if (value === null || value === undefined || !test) return false;
  const code = String(test.code || '').toUpperCase();
  if (code === 'K' && value >= 6) return true;
  if ((code === 'HB' || code === 'HGB') && value < 7) return true;
  if (code === 'GLU' && (value < 3 || value > 20)) return true;
  return false;
}
