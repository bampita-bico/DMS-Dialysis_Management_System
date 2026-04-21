import { useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import FormModal from '../../components/forms/FormModal';
import SessionForm from '../../components/forms/SessionForm';

export default function SessionSchedule() {
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [viewMode, setViewMode] = useState('day'); // day, week
  const [filterShift, setFilterShift] = useState('all');
  const [showScheduleModal, setShowScheduleModal] = useState(false);

  const { data: allSessions, loading } = useOfflineData('dialysis_sessions');
  const { data: patients } = useOfflineData('patients');

  const handleScheduleSession = async (sessionData) => {
    await offlineService.create('dialysis_sessions', sessionData, 8);
  };

  const handleUpdateSession = async (id, updates) => {
    await offlineService.update('dialysis_sessions', id, updates);
  };

  // Map database sessions to UI format
  const sessions = (allSessions || [])
    .filter(s => s.scheduled_date?.startsWith(selectedDate))
    .map(session => {
      const patient = patients?.find(p => p.id === session.patient_id);
      const fullName = patient ? `${patient.first_name || ''} ${patient.last_name || ''}`.trim() : 'Unknown';
      const schedTime = session.scheduled_date?.split('T')[1]?.substring(0, 5) || '00:00';
      const hour = parseInt(schedTime.split(':')[0]);
      const shift = hour < 12 ? 'morning' : hour < 17 ? 'afternoon' : 'evening';

      return {
        id: session.id,
        patientName: fullName,
        patientNumber: patient?.mrn || 'N/A',
        time: schedTime,
        duration: session.planned_duration_hours || 4,
        machine: session.machine_id || 'Unassigned',
        shift: shift,
        status: session.session_status || 'scheduled',
        nurse: session.nurse_assigned || 'Unassigned',
        startTime: session.actual_start_time?.split('T')[1]?.substring(0, 5) || null,
        endTimeExpected: null // Calculate from start + duration if needed
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
    inProgress: sessions.filter(s => s.status === 'in-progress').length,
    scheduled: sessions.filter(s => s.status === 'scheduled').length,
    completed: sessions.filter(s => s.status === 'completed').length
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
                        <SessionCard key={session.id} session={session} />
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
    </div>
  );
}

function SessionCard({ session }) {
  const statusConfig = {
    scheduled: {
      bg: 'bg-sky-50',
      border: 'border-sky-200',
      text: 'text-sky-700',
      badge: 'bg-sky-100 text-sky-700'
    },
    'in-progress': {
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

  const config = statusConfig[session.status];

  return (
    <div className={`${config.bg} rounded-xl border ${config.border} p-5 hover:shadow-md transition-shadow`}>
      {/* Time and Status */}
      <div className="flex items-start justify-between mb-3">
        <div>
          <div className="text-2xl font-bold text-gray-900">{session.time}</div>
          <div className="text-sm text-gray-600 mt-0.5">{session.duration}h session</div>
        </div>
        <span className={`px-3 py-1 rounded-full text-xs font-medium ${config.badge}`}>
          {session.status.replace('-', ' ').toUpperCase()}
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
      </div>

      {/* Actions */}
      <div className="mt-4 pt-3 border-t border-gray-200">
        <button className={`w-full py-2 rounded-lg font-medium transition-colors ${
          session.status === 'scheduled'
            ? 'bg-sky-600 text-white hover:bg-sky-700'
            : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-100'
        }`}>
          {session.status === 'scheduled' && 'Start Session'}
          {session.status === 'in-progress' && 'View Details'}
          {session.status === 'completed' && 'View Report'}
        </button>
      </div>
    </div>
  );
}
