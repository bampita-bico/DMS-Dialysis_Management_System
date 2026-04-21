import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import FormModal from '../../components/forms/FormModal';
import InvoiceForm from '../../components/forms/InvoiceForm';

export default function Invoices() {
  const [filterStatus, setFilterStatus] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [dateRange, setDateRange] = useState('30days');
  const [showInvoiceModal, setShowInvoiceModal] = useState(false);

  const { data: allInvoices, loading } = useOfflineData('invoices');
  const { data: insuranceSchemes } = useOfflineData('insurance_schemes');
  const { data: priceLists } = useOfflineData('price_lists');
  const { data: patients } = useOfflineData('patients');

  const handleCreateInvoice = async (invoiceData) => {
    await offlineService.create('invoices', invoiceData, 7);
  };

  const handleRecordPayment = async (invoiceId, paymentData) => {
    await offlineService.create('payments', { invoice_id: invoiceId, ...paymentData }, 8);
  };

  const invoices = (allInvoices || [])
    .filter(inv => {
      if (filterStatus !== 'all' && inv.invoice_status !== filterStatus) return false;
      const patient = patients?.find(p => p.id === inv.patient_id);
      const patientName = patient ? `${patient.first_name} ${patient.last_name}` : '';
      if (searchQuery && !patientName.toLowerCase().includes(searchQuery.toLowerCase())) return false;
      return true;
    })
    .map(invoice => {
      const patient = patients?.find(p => p.id === invoice.patient_id);
      return {
        ...invoice,
        patientName: patient ? `${patient.first_name} ${patient.last_name}` : 'Unknown',
        patientNumber: patient?.mrn || 'N/A',
        status: invoice.invoice_status || 'pending'
      };
    });

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-UG', {
      style: 'currency',
      currency: 'UGX',
      minimumFractionDigits: 0
    }).format(amount || 0);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Invoices & Billing
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                Manage patient invoices and payments
              </p>
            </div>
            <button
              onClick={() => setShowInvoiceModal(true)}
              className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm">
              New Invoice
            </button>
          </div>

          {/* Search and Filters */}
          <div className="mt-6 flex gap-4">
            <div className="flex-1">
              <input
                type="text"
                placeholder="Search by patient name..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent"
              />
            </div>

            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Status</option>
              <option value="paid">Paid</option>
              <option value="partial">Partial Payment</option>
              <option value="unpaid">Unpaid</option>
              <option value="overdue">Overdue</option>
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

      {/* Invoices List */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
          </div>
        ) : (
          <div className="space-y-4">
            {invoices.map((invoice) => (
              <div key={invoice.id} className="bg-white rounded-xl border border-gray-200 p-6 hover:shadow-md transition-shadow">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 font-mono">{invoice.invoice_number || 'INV-PENDING'}</h3>
                    <p className="text-sm text-gray-600 mt-1">{invoice.patientName}</p>
                    <p className="text-xs text-gray-500 mt-0.5">MRN: {invoice.patientNumber}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-bold text-gray-900">{formatCurrency(invoice.total_amount)}</p>
                    <span className={`inline-block mt-2 px-3 py-1 rounded-full text-xs font-medium ${
                      invoice.status === 'paid' ? 'bg-emerald-100 text-emerald-700' :
                      invoice.status === 'partial' ? 'bg-amber-100 text-amber-700' :
                      invoice.status === 'overdue' ? 'bg-red-100 text-red-700' :
                      'bg-gray-100 text-gray-700'
                    }`}>
                      {invoice.status.toUpperCase()}
                    </span>
                  </div>
                </div>
                <div className="mt-4 flex gap-4 text-sm text-gray-600">
                  <span>Issue: {invoice.created_at ? new Date(invoice.created_at).toLocaleDateString() : 'N/A'}</span>
                  <span>Due: {invoice.due_date ? new Date(invoice.due_date).toLocaleDateString() : 'N/A'}</span>
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && invoices.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No invoices found</p>
            <p className="text-gray-400 text-sm mt-2">Try adjusting your search or filters</p>
          </div>
        )}
      </div>

      <FormModal
        isOpen={showInvoiceModal}
        onClose={() => setShowInvoiceModal(false)}
        title="Create New Invoice"
        size="xl"
      >
        <InvoiceForm
          onSuccess={() => {
            setShowInvoiceModal(false);
          }}
          onCancel={() => setShowInvoiceModal(false)}
        />
      </FormModal>
    </div>
  );
}
