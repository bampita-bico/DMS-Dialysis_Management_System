import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';

export default function LabResults() {
  const [selectedPatient, setSelectedPatient] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [dateRange, setDateRange] = useState('7days');

  const { data: labResults, loading } = useOfflineData('lab_results');
  const { data: labOrders } = useOfflineData('lab_orders');
  const { data: labTests } = useOfflineData('lab_test_catalog');
  const { data: patients } = useOfflineData('patients');

  const handleCreateLabOrder = async (orderData) => {
    await offlineService.create('lab_orders', orderData, 9);
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
        patientName: patient ? `${patient.first_name} ${patient.last_name}` : 'Unknown',
        patientNumber: patient?.mrn || 'N/A',
        status: result.result_status || 'pending'
      };
    });

  // Old mock data below (keeping for reference)
