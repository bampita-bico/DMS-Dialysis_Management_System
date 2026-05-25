import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import FormModal from '../../components/forms/FormModal';
import SessionForm from '../../components/forms/SessionForm';
import api from '../../services/api';
import db from '../../db/schema';

export default function SessionSchedule() {
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [viewMode, setViewMode] = useState('day'); // day, week
  const [filterShift, setFilterShift] = useState('all');
  const [showScheduleModal, setShowScheduleModal] = useState(false);
  const [workflow, setWorkflow] = useState({ mode: null, session: null });
  const [actionError, setActionError] = useState('');

  const { data: allSessions, loading } = useOfflineData('dialysis_sessions');
  const { data: patients } = useOfflineData('patients');
  const { data: dialysateRecords } = useOfflineData('dialysate_records');

  // Map database sessions to UI format
  const sessions = (allSessions || [])
    .filter(s => dateOnly(s.scheduled_date)?.startsWith(selectedDate))
    .map(session => {
      const patient = patients?.find(p => String(p.id) === String(session.patient_id));
      const fullName = patient?.full_name || (patient ? `${patient.first_name || ''} ${patient.last_name || ''}`.trim() : 'Unknown');
      const dialysate = dialysateRecords?.find(d => String(d.session_id) === String(session.id));
      const schedTime = timeOnly(session.scheduled_start_time || session.scheduled_date) || '00:00';
      const hour = parseInt(schedTime.split(':')[0]);
      const shift = hour < 12 ? 'morning' : hour < 17 ? 'afternoon' : 'evening';
      const duration = session.planned_duration_hours || Math.round((session.prescribed_duration_mins || 240) / 60);
      const status = normalizeSessionStatus(session.session_status || session.status || 'scheduled');

      return {
        raw: session,
        id: session.id,
        patientId: session.patient_id,
        machineId: session.machine_id,
        patientName: fullName,
        patientNumber: patient?.mrn || 'N/A',
        time: schedTime,
        duration,
        machine: session.machine_id || 'Unassigned',
        shift: shift,
        status,
        nurse: session.nurse_assigned || 'Unassigned',
        startTime: timeOnly(session.actual_start_time),
        endTimeExpected: calculateExpectedEnd(schedTime, duration),
        dialysateMagnesium: dialysate?.magnesium_meq_l,
        dialysateSodium: dialysate?.sodium_meq_l,
        dialysatePotassium: dialysate?.potassium_meq_l,
        dialysateVerified: dialysate?.composition_verified
      };
    });

  const shifts = {
    morning: { name: 'Morning Shift', time: '07:00 - 12:00', color: 'text-amber-600 bg-amber-50' },
    afternoon: { name: 'Afternoon Shift', time: '12:00 - 17:00', color: 'text-sky-600 bg-sky-50' },
    evening: { name: 'Evening Shift', time: '17:00 - 22:00', color: 'text-indigo-600 bg-indigo-50' }
  };

  const filteredSessions = filterShift === 'all'
    ? sessions
    : sessions.filter(s => s.shift === filterShift);

  const sessionsByShift = {
    morning: filteredSessions.filter(s => s.shift === 'morning'),
    afternoon: filteredSessions.filter(s => s.shift === 'afternoon'),
    evening: filteredSessions.filter(s => s.shift === 'evening')
  };

  const stats = {
    total: sessions.length,
    inProgress: sessions.filter(s => s.status === 'in_progress').length,
    scheduled: sessions.filter(s => s.status === 'scheduled').length,
    completed: sessions.filter(s => s.status === 'completed').length
  };

  const openWorkflow = (mode, session) => {
    setActionError('');
    setWorkflow({ mode, session });
  };

  const closeWorkflow = () => {
    setActionError('');
    setWorkflow({ mode: null, session: null });
  };

  const handleWorkflowSubmit = async (mode, payload) => {
    setActionError('');
    try {
      await applySessionWorkflow(workflow.session, mode, payload);
      closeWorkflow();
    } catch (error) {
      const message = error.response?.data?.message || error.response?.data?.error || error.message || 'Session action failed';
      setActionError(message);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Session Schedule
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                Dialysis session scheduling and real-time status monitoring
              </p>
            </div>
            <button
              onClick={() => setShowScheduleModal(true)}
              className="px-6 py-3 bg-sky-600 text-white font-medium rounded-lg hover:bg-sky-700 transition-colors shadow-sm">
              {/* Icon: Plus */}
              Schedule Session
            </button>
          </div>

          {/* Stats */}
          <div className="mt-6 grid grid-cols-4 gap-4">
            <div className="bg-gray-100 rounded-lg p-4 border border-gray-200">
              <div className="text-2xl font-bold text-gray-900">{stats.total}</div>
              <div className="text-sm text-gray-600 mt-1">Total Sessions</div>
            </div>
            <div className="bg-emerald-50 rounded-lg p-4 border border-emerald-200">
              <div className="text-2xl font-bold text-emerald-700">{stats.inProgress}</div>
              <div className="text-sm text-emerald-700 mt-1">In Progress</div>
            </div>
            <div className="bg-sky-50 rounded-lg p-4 border border-sky-200">
              <div className="text-2xl font-bold text-sky-700">{stats.scheduled}</div>
              <div className="text-sm text-sky-700 mt-1">Scheduled</div>
            </div>
            <div className="bg-gray-100 rounded-lg p-4 border border-gray-300">
              <div className="text-2xl font-bold text-gray-700">{stats.completed}</div>
              <div className="text-sm text-gray-600 mt-1">Completed</div>
            </div>
          </div>

          {/* Date Picker and Filters */}
          <div className="mt-6 flex items-center gap-4">
            <input
              type="date"
              value={selectedDate}
              onChange={(e) => setSelectedDate(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent"
            />

            <select
              value={filterShift}
              onChange={(e) => setFilterShift(e.target.value)}
              className="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-transparent bg-white"
            >
              <option value="all">All Shifts</option>
              <option value="morning">Morning</option>
              <option value="afternoon">Afternoon</option>
              <option value="evening">Evening</option>
            </select>

            <div className="flex-1"></div>

            <div className="flex gap-2">
              <button
                onClick={() => setViewMode('day')}
                className={`px-4 py-2.5 rounded-lg font-medium transition-colors ${
                  viewMode === 'day'
                    ? 'bg-sky-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-100'
                }`}
              >
                Day View
              </button>
              <button
                onClick={() => setViewMode('week')}
                className={`px-4 py-2.5 rounded-lg font-medium transition-colors ${
                  viewMode === 'week'
                    ? 'bg-sky-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-100'
                }`}
              >
                Week View
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Sessions by Shift */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
          </div>
        ) : (
          <div className="space-y-8">
            {Object.entries(sessionsByShift).map(([shiftKey, shiftSessions]) => {
              if (shiftSessions.length === 0 && filterShift !== 'all' && filterShift !== shiftKey) return null;

              return (
                <div key={shiftKey}>
                  <div className={`inline-flex items-center gap-3 px-4 py-2 rounded-lg mb-4 ${shifts[shiftKey].color}`}>
                    <div className="text-lg font-semibold">{shifts[shiftKey].name}</div>
                    <div className="text-sm opacity-75">{shifts[shiftKey].time}</div>
                    <div className="ml-2 px-2 py-0.5 bg-white/50 rounded-full text-sm font-medium">
                      {shiftSessions.length} sessions
                    </div>
                  </div>

                  {shiftSessions.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      No sessions scheduled for this shift
                    </div>
                  ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {shiftSessions.map((session) => (
                        <SessionCard key={session.id} session={session} onAction={openWorkflow} />
                      ))}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Schedule Modal */}
      <FormModal
        isOpen={showScheduleModal}
        onClose={() => setShowScheduleModal(false)}
        title="Schedule Dialysis Session"
        size="xl"
      >
        <SessionForm
          onSuccess={() => {
            setShowScheduleModal(false);
          }}
          onCancel={() => setShowScheduleModal(false)}
        />
      </FormModal>

      <FormModal
        isOpen={Boolean(workflow.mode)}
        onClose={closeWorkflow}
        title={workflowTitle(workflow.mode)}
        size="lg"
      >
        <SessionWorkflowModal
          mode={workflow.mode}
          session={workflow.session}
          error={actionError}
          onSubmit={handleWorkflowSubmit}
          onCancel={closeWorkflow}
        />
      </FormModal>
    </div>
  );
}

function SessionCard({ session, onAction }) {
  const statusConfig = {
    scheduled: {
      bg: 'bg-sky-50',
      border: 'border-sky-200',
      text: 'text-sky-700',
      badge: 'bg-sky-100 text-sky-700'
    },
    in_progress: {
      bg: 'bg-emerald-50',
      border: 'border-emerald-300',
      text: 'text-emerald-700',
      badge: 'bg-emerald-100 text-emerald-700'
    },
    completed: {
      bg: 'bg-gray-100',
      border: 'border-gray-200',
      text: 'text-gray-600',
      badge: 'bg-gray-100 text-gray-700'
    }
  };

  const config = statusConfig[session.status] || {
    bg: 'bg-amber-50',
    border: 'border-amber-200',
    text: 'text-amber-700',
    badge: 'bg-amber-100 text-amber-700'
  };

  return (
    <div className={`${config.bg} rounded-xl border ${config.border} p-5 hover:shadow-md transition-shadow`}>
      {/* Time and Status */}
      <div className="flex items-start justify-between mb-3">
        <div>
          <div className="text-2xl font-bold text-gray-900">{session.time}</div>
          <div className="text-sm text-gray-600 mt-0.5">{session.duration}h session</div>
        </div>
        <span className={`px-3 py-1 rounded-full text-xs font-medium ${config.badge}`}>
          {humanizeStatus(session.status).toUpperCase()}
        </span>
      </div>

      {/* Patient Info */}
      <div className="mb-3 pb-3 border-b border-gray-200">
        <div className="font-semibold text-gray-900">{session.patientName}</div>
        <div className="text-xs text-gray-500 font-mono mt-1">{session.patientNumber}</div>
      </div>

      {/* Session Details */}
      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-gray-600">Machine:</span>
          <span className="font-semibold text-gray-900">{session.machine}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Nurse:</span>
          <span className="font-medium text-gray-900">{session.nurse}</span>
        </div>
        {session.startTime && (
          <div className="flex justify-between">
            <span className="text-gray-600">Started:</span>
            <span className="font-medium text-emerald-600">{session.startTime}</span>
          </div>
        )}
        <div className="flex justify-between">
          <span className="text-gray-600">Expected End:</span>
          <span className="font-medium text-gray-900">{session.endTimeExpected}</span>
        </div>
        <div className="grid grid-cols-3 gap-2 pt-2">
          <div className="rounded-lg bg-white/70 border border-gray-200 px-2 py-1">
            <div className="text-[10px] uppercase text-gray-500">Na</div>
            <div className="font-semibold text-gray-900">{session.dialysateSodium || '--'}</div>
          </div>
          <div className="rounded-lg bg-white/70 border border-gray-200 px-2 py-1">
            <div className="text-[10px] uppercase text-gray-500">K</div>
            <div className="font-semibold text-gray-900">{session.dialysatePotassium || '--'}</div>
          </div>
          <div className="rounded-lg bg-white/70 border border-gray-200 px-2 py-1">
            <div className="text-[10px] uppercase text-gray-500">Mg</div>
            <div className="font-semibold text-gray-900">{session.dialysateMagnesium || '--'}</div>
          </div>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Dialysate check:</span>
          <span className={`font-medium ${session.dialysateVerified ? 'text-emerald-700' : 'text-amber-700'}`}>
            {session.dialysateVerified ? 'Verified' : 'Pending'}
          </span>
        </div>
      </div>

      {/* Actions */}
      <div className="mt-4 pt-3 border-t border-gray-200 space-y-2">
        <button
          type="button"
          onClick={() => onAction(session.status === 'scheduled' ? 'start' : session.status === 'in_progress' ? 'complete' : 'details', session)}
          className={`w-full py-2 rounded-lg font-medium transition-colors ${
          session.status === 'scheduled'
            ? 'bg-sky-600 text-white hover:bg-sky-700'
            : session.status === 'in_progress'
              ? 'bg-emerald-600 text-white hover:bg-emerald-700'
            : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-100'
        }`}>
          {session.status === 'scheduled' && 'Start Session'}
          {session.status === 'in_progress' && 'Complete Session'}
          {session.status === 'completed' && 'View Report'}
          {!['scheduled', 'in_progress', 'completed'].includes(session.status) && 'View Details'}
        </button>
        {session.status === 'in_progress' && (
          <button
            type="button"
            onClick={() => onAction('abort', session)}
            className="w-full py-2 rounded-lg border border-red-200 bg-red-50 font-medium text-red-700 hover:bg-red-100 transition-colors"
          >
            Abort / Pause Session
          </button>
        )}
      </div>
    </div>
  );
}

function SessionWorkflowModal({ mode, session, error, onSubmit, onCancel }) {
  const [formData, setFormData] = useState(() => ({
    pre_weight_kg: session?.raw?.pre_weight_kg || '',
    pre_bp_systolic: session?.raw?.pre_bp_systolic || '',
    pre_bp_diastolic: session?.raw?.pre_bp_diastolic || '',
    pre_hr: session?.raw?.pre_hr || '',
    pre_temp: session?.raw?.pre_temp || '36.7',
    actual_duration_mins: session?.raw?.actual_duration_mins || (session?.duration ? session.duration * 60 : 240),
    post_weight_kg: session?.raw?.post_weight_kg || '',
    post_bp_systolic: session?.raw?.post_bp_systolic || '',
    post_bp_diastolic: session?.raw?.post_bp_diastolic || '',
    post_hr: session?.raw?.post_hr || '',
    was_patient_reviewed: true,
    session_notes: session?.raw?.session_notes || '',
    aborted_reason: '',
  }));
  const [saving, setSaving] = useState(false);

  if (!session) return null;

  const handleChange = (event) => {
    const { name, value, type, checked } = event.target;
    setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setSaving(true);
    try {
      await onSubmit(mode, formData);
    } finally {
      setSaving(false);
    }
  };

  if (mode === 'details') {
    return (
      <div className="space-y-4">
        <div className="rounded-lg border border-gray-200 bg-gray-50 p-4">
          <p className="text-lg font-semibold text-gray-900">{session.patientName}</p>
          <p className="text-sm text-gray-600 mt-1">
            {session.patientNumber} | {session.time} | {humanizeStatus(session.status)}
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
          <Readonly label="Machine" value={session.machine} />
          <Readonly label="Nurse" value={session.nurse} />
          <Readonly label="Dialysate Na / K / Mg" value={`${session.dialysateSodium || '--'} / ${session.dialysatePotassium || '--'} / ${session.dialysateMagnesium || '--'}`} />
          <Readonly label="Dialysate check" value={session.dialysateVerified ? 'Verified' : 'Pending'} />
        </div>
        <div className="flex justify-end pt-4 border-t border-gray-200">
          <button type="button" onClick={onCancel} className="px-4 py-2 rounded-lg bg-gray-900 text-white hover:bg-gray-800">
            Close
          </button>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {error}
        </div>
      )}

      <div className="rounded-lg border border-sky-200 bg-sky-50 p-4">
        <p className="font-semibold text-gray-900">{session.patientName}</p>
        <p className="text-sm text-sky-800 mt-1">
          MRN {session.patientNumber} | {session.time} | machine {session.machine}
        </p>
      </div>

      {mode === 'start' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <WorkflowInput label="Pre weight (kg)" name="pre_weight_kg" type="number" step="0.1" value={formData.pre_weight_kg} onChange={handleChange} required />
          <WorkflowInput label="Heart rate" name="pre_hr" type="number" value={formData.pre_hr} onChange={handleChange} required />
          <WorkflowInput label="Systolic BP" name="pre_bp_systolic" type="number" value={formData.pre_bp_systolic} onChange={handleChange} required />
          <WorkflowInput label="Diastolic BP" name="pre_bp_diastolic" type="number" value={formData.pre_bp_diastolic} onChange={handleChange} required />
          <WorkflowInput label="Temperature" name="pre_temp" type="number" step="0.1" value={formData.pre_temp} onChange={handleChange} />
        </div>
      )}

      {mode === 'complete' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <WorkflowInput label="Delivered minutes" name="actual_duration_mins" type="number" value={formData.actual_duration_mins} onChange={handleChange} required />
            <WorkflowInput label="Post weight (kg)" name="post_weight_kg" type="number" step="0.1" value={formData.post_weight_kg} onChange={handleChange} required />
            <WorkflowInput label="Post systolic BP" name="post_bp_systolic" type="number" value={formData.post_bp_systolic} onChange={handleChange} required />
            <WorkflowInput label="Post diastolic BP" name="post_bp_diastolic" type="number" value={formData.post_bp_diastolic} onChange={handleChange} required />
            <WorkflowInput label="Post heart rate" name="post_hr" type="number" value={formData.post_hr} onChange={handleChange} required />
          </div>
          <label className="flex items-center gap-2 text-sm text-gray-700">
            <input type="checkbox" name="was_patient_reviewed" checked={formData.was_patient_reviewed} onChange={handleChange} className="h-4 w-4 rounded text-sky-600" />
            Doctor/nurse review completed before closure
          </label>
          <WorkflowTextArea label="Session notes" name="session_notes" value={formData.session_notes} onChange={handleChange} />
        </div>
      )}

      {mode === 'abort' && (
        <WorkflowTextArea
          label="Abort / pause reason"
          name="aborted_reason"
          value={formData.aborted_reason}
          onChange={handleChange}
          required
        />
      )}

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button type="button" onClick={onCancel} disabled={saving} className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={saving} className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {saving ? 'Saving...' : workflowButtonLabel(mode)}
        </button>
      </div>
    </form>
  );
}

function WorkflowInput({ label, name, value, onChange, type = 'text', required = false, step }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}{required && <span className="text-red-500 ml-1">*</span>}</span>
      <input
        name={name}
        type={type}
        step={step}
        value={value}
        onChange={onChange}
        required={required}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500"
      />
    </label>
  );
}

function WorkflowTextArea({ label, name, value, onChange, required = false }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}{required && <span className="text-red-500 ml-1">*</span>}</span>
      <textarea
        name={name}
        value={value}
        onChange={onChange}
        required={required}
        rows={3}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500"
      />
    </label>
  );
}

function Readonly({ label, value }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-gray-50 p-3">
      <p className="text-xs font-semibold uppercase tracking-wide text-gray-500">{label}</p>
      <p className="mt-1 font-medium text-gray-900">{value || 'Not recorded'}</p>
    </div>
  );
}

async function applySessionWorkflow(session, mode, formData) {
  if (!session?.id) return;

  const now = new Date().toISOString();
  const apiPayloads = {
    start: {
      pre_weight_kg: Number(formData.pre_weight_kg),
      pre_bp_systolic: Number(formData.pre_bp_systolic),
      pre_bp_diastolic: Number(formData.pre_bp_diastolic),
      pre_hr: Number(formData.pre_hr),
      pre_temp: Number(formData.pre_temp || 0),
    },
    complete: {
      actual_duration_mins: Number(formData.actual_duration_mins),
      post_weight_kg: Number(formData.post_weight_kg),
      post_bp_systolic: Number(formData.post_bp_systolic),
      post_bp_diastolic: Number(formData.post_bp_diastolic),
      post_hr: Number(formData.post_hr),
      was_patient_reviewed: Boolean(formData.was_patient_reviewed),
      session_notes: formData.session_notes || '',
    },
    abort: {
      aborted_reason: formData.aborted_reason,
    },
  };

  const localUpdates = {
    start: {
      status: 'in_progress',
      actual_start_time: now,
      ...apiPayloads.start,
    },
    complete: {
      status: 'completed',
      actual_end_time: now,
      ...apiPayloads.complete,
    },
    abort: {
      status: 'aborted',
      actual_end_time: now,
      aborted_reason: formData.aborted_reason,
    },
  };

  if (isUuid(session.id) && navigator.onLine) {
    const response = await api.post(`/dialysis-sessions/${session.id}/${mode}`, apiPayloads[mode]);
    await putSessionLocal(response.data);
    return;
  }

  await putSessionLocal({
    ...(session.raw || {}),
    id: session.id,
    patient_id: session.patientId,
    machine_id: session.machineId,
    ...localUpdates[mode],
  });
}

async function putSessionLocal(session) {
  await db.dialysis_sessions.put({
    ...session,
    synced: false,
    updated_at: new Date().toISOString(),
  });
  window.dispatchEvent(new CustomEvent('dms-local-change', {
    detail: { entityType: 'dialysis_sessions', record: session },
  }));
}

function workflowTitle(mode) {
  const titles = {
    start: 'Start Dialysis Session',
    complete: 'Complete Dialysis Session',
    abort: 'Abort / Pause Dialysis Session',
    details: 'Session Report',
  };
  return titles[mode] || 'Session';
}

function workflowButtonLabel(mode) {
  const labels = {
    start: 'Start Session',
    complete: 'Complete Session',
    abort: 'Record Abort',
  };
  return labels[mode] || 'Save';
}

function normalizeSessionStatus(status) {
  return String(status || 'scheduled').replace('-', '_');
}

function humanizeStatus(status) {
  return String(status || '').replaceAll('_', ' ').replaceAll('-', ' ');
}

function dateOnly(value) {
  if (!value) return '';
  if (typeof value === 'string') return value.slice(0, 10);
  if (value.Time) return String(value.Time).slice(0, 10);
  return String(value).slice(0, 10);
}

function timeOnly(value) {
  if (!value) return '';
  if (typeof value === 'string') {
    if (value.includes('T')) return value.split('T')[1]?.slice(0, 5) || '';
    return value.slice(0, 5);
  }
  if (value.Microseconds !== undefined) {
    const totalSeconds = Math.floor(Number(value.Microseconds) / 1000000);
    const hours = String(Math.floor(totalSeconds / 3600)).padStart(2, '0');
    const minutes = String(Math.floor((totalSeconds % 3600) / 60)).padStart(2, '0');
    return `${hours}:${minutes}`;
  }
  return '';
}

function isUuid(value) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(String(value));
}

function calculateExpectedEnd(startTime, durationHours) {
  if (!startTime) return 'N/A';
  const [hour, minute] = startTime.split(':').map(Number);
  if (Number.isNaN(hour) || Number.isNaN(minute)) return 'N/A';
  const end = new Date();
  end.setHours(hour + durationHours, minute, 0, 0);
  return end.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}
