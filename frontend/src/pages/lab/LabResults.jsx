import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import FormModal from '../../components/forms/FormModal';
import LabOrderForm from '../../components/forms/LabOrderForm';

export default function LabResults() {
  const [selectedPatient, setSelectedPatient] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [dateRange, setDateRange] = useState('7days');
  const [showNewOrderModal, setShowNewOrderModal] = useState(false);

  const { data: labResults, loading } = useOfflineData('lab_results');
  const { data: labOrders } = useOfflineData('lab_orders');
  const { data: labTests } = useOfflineData('lab_test_catalog');
  const { data: patients } = useOfflineData('patients');

  const handleCreateLabOrder = async (orderData) => {
    await offlineService.create('lab_orders', orderData, 9);
    setShowNewOrderModal(false);
  };

  const handleRecordResult = async (orderId, resultData) => {
    await offlineService.update('lab_results', orderId, resultData);
  };

  const results = (labResults || [])
    .filter(r => {
      if (selectedPatient !== 'all' && r.patient_id !== selectedPatient) return false;
      if (filterStatus !== 'all' && r.result_status !== filterStatus) return false;
      return true;
    })
    .map(result => {
      const patient = patients?.find(p => p.id === result.patient_id);
      const order = labOrders?.find(o => o.id === result.lab_order_id);
      return {
        ...result,
        patientName: patient?.full_name || 'Unknown',
        patientNumber: patient?.mrn || 'N/A',
        status: result.result_status || 'pending',
        orderDate: order?.created_at || result.created_at
      };
    });

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
            <button
              onClick={() => setShowNewOrderModal(true)}
              className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm"
            >
              + New Lab Order
            </button>
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

      {/* New Lab Order Modal */}
      {showNewOrderModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-bold text-gray-900 mb-4">New Lab Order (Coming Soon)</h3>
            <p className="text-gray-600 mb-4">
              The lab order form with test selection will be available in the next update.
            </p>
            <p className="text-sm text-gray-500 mb-4">
              Features coming soon:
              • Select patient
              • Choose tests (CBC, U&E, Creatinine, etc.)
              • Set priority
              • Add clinical notes
            </p>
            <button
              onClick={() => setShowNewOrderModal(false)}
              className="w-full px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700"
            >
              Close
            </button>
          </div>
        </div>
      )}

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
    </div>
  );
}
