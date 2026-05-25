import { useMemo, useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import api from '../../services/api';
import db from '../../db/schema';
import { authService } from '../../services/auth';

const QUEUES = [
  { id: 'safety', label: 'Safety' },
  { id: 'delivery', label: 'Delivery' },
  { id: 'access', label: 'Access' },
  { id: 'meds', label: 'Meds' },
  { id: 'ops', label: 'Ops' },
  { id: 'interop', label: 'Interop' },
];

export default function ClinicalCommandCenter() {
  const [activeQueue, setActiveQueue] = useState('safety');
  const [actionModal, setActionModal] = useState({ type: null, record: null });
  const [actionError, setActionError] = useState('');

  const { data: patients, loading: patientsLoading } = useOfflineData('patients');
  const { data: sessions } = useOfflineData('dialysis_sessions');
  const { data: profiles } = useOfflineData('patient_clinical_profiles');
  const { data: safetyChecks } = useOfflineData('session_safety_checks');
  const { data: clinicalAlerts } = useOfflineData('clinical_alerts');
  const { data: telemetry } = useOfflineData('treatment_telemetry');
  const { data: accessLifecycle } = useOfflineData('access_lifecycle_events');
  const { data: infections } = useOfflineData('infection_surveillance_events');
  const { data: adequacyReviews } = useOfflineData('adequacy_reviews');
  const { data: medReviews } = useOfflineData('medication_reconciliation_reviews');
  const { data: patientReports } = useOfflineData('patient_reported_events');
  const { data: attendance } = useOfflineData('staff_attendance_verifications');
  const { data: unitEvents } = useOfflineData('unit_safety_events');
  const { data: exportsData } = useOfflineData('interoperability_exports');
  const { data: graphEdges } = useOfflineData('ontology_relationships');
  const { data: labResults } = useOfflineData('lab_results');
  const { data: vascularAccess } = useOfflineData('vascular_access');
  const { data: dialysateRecords } = useOfflineData('dialysate_records');

  const context = useMemo(() => {
    const patientMap = new Map((patients || []).map(patient => [String(patient.id), patient]));
    const sessionMap = new Map((sessions || []).map(session => [String(session.id), session]));
    const profileMap = new Map((profiles || []).map(profile => [String(profile.patient_id), profile]));

    const blockedChecks = (safetyChecks || []).filter(check =>
      check.override_required || ['blocked', 'failed'].includes(check.check_status)
    );
    const openAlerts = (clinicalAlerts || []).filter(alert =>
      !['acknowledged', 'closed', 'resolved'].includes(alert.status)
    );
    const accessRisks = (accessLifecycle || []).filter(event =>
      includesAny(event.long_term_complications, ['infection', 'recirculation', 'blockage']) ||
      includesAny(event.immediate_complications, ['difficult', 'bleeding']) ||
      event.culture_result
    );
    const adequacyRisks = (adequacyReviews || []).filter(review =>
      review.doctor_review_required || ['doctor_review_required', 'inadequate', 'needs_review'].includes(review.adequacy_status)
    );
    const staffExceptions = (attendance || []).filter(item =>
      !['verified', 'accepted'].includes(item.verification_result)
    );
    const openUnitEvents = (unitEvents || []).filter(event => !event.closed_at);
    const dueAKI = (profiles || []).filter(profile =>
      profile.aki_to_ckd_reassessment_due && new Date(profile.aki_to_ckd_reassessment_due) <= new Date()
    );
    const exportGaps = (exportsData || []).filter(item =>
      !['exported', 'ready'].includes(item.export_status)
    );

    return {
      patientMap,
      sessionMap,
      profileMap,
      blockedChecks,
      openAlerts,
      accessRisks,
      adequacyRisks,
      staffExceptions,
      openUnitEvents,
      dueAKI,
      exportGaps,
    };
  }, [
    patients,
    sessions,
    profiles,
    safetyChecks,
    clinicalAlerts,
    accessLifecycle,
    adequacyReviews,
    attendance,
    unitEvents,
    exportsData,
  ]);

  const loading = patientsLoading;

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 tracking-tight">Clinical Command Center</h1>
              <p className="mt-1 text-sm text-gray-600">
                {formatDate(new Date().toISOString())} | {patients?.length || 0} patients | {sessions?.length || 0} sessions
              </p>
            </div>
            <div className="flex flex-wrap gap-2">
              {QUEUES.map(queue => (
                <button
                  key={queue.id}
                  onClick={() => setActiveQueue(queue.id)}
                  className={`h-10 min-w-20 rounded-lg border px-3 text-sm font-semibold transition-colors ${
                    activeQueue === queue.id
                      ? 'border-sky-700 bg-sky-700 text-white'
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  {queue.label}
                </button>
              ))}
            </div>
          </div>

          <div className="mt-6 grid grid-cols-2 md:grid-cols-4 xl:grid-cols-8 gap-3">
            <KpiCard label="Blocked Starts" value={context.blockedChecks.length} tone={context.blockedChecks.length ? 'red' : 'green'} />
            <KpiCard label="Open Alerts" value={context.openAlerts.length} tone={context.openAlerts.length ? 'yellow' : 'green'} />
            <KpiCard label="Access Risks" value={context.accessRisks.length + (infections?.length || 0)} tone={context.accessRisks.length || infections?.length ? 'red' : 'green'} />
            <KpiCard label="Adequacy" value={context.adequacyRisks.length} tone={context.adequacyRisks.length ? 'yellow' : 'green'} />
            <KpiCard label="AKI Review" value={context.dueAKI.length} tone={context.dueAKI.length ? 'purple' : 'green'} />
            <KpiCard label="Staff Exceptions" value={context.staffExceptions.length} tone={context.staffExceptions.length ? 'yellow' : 'green'} />
            <KpiCard label="Unit Events" value={context.openUnitEvents.length} tone={context.openUnitEvents.length ? 'red' : 'green'} />
            <KpiCard label="Export Gaps" value={context.exportGaps.length} tone={context.exportGaps.length ? 'blue' : 'green'} />
          </div>
        </div>
      </div>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-6">
        {activeQueue === 'safety' && (
          <SafetyQueue
            checks={safetyChecks || []}
            alerts={clinicalAlerts || []}
            patients={context.patientMap}
            sessions={context.sessionMap}
            labs={labResults || []}
            dialysate={dialysateRecords || []}
            onAction={(type, record) => setActionModal({ type, record })}
          />
        )}
        {activeQueue === 'delivery' && (
          <DeliveryQueue
            telemetry={telemetry || []}
            adequacy={adequacyReviews || []}
            profiles={context.profileMap}
            patients={context.patientMap}
          />
        )}
        {activeQueue === 'access' && (
          <AccessQueue
            access={vascularAccess || []}
            lifecycle={accessLifecycle || []}
            infections={infections || []}
            patients={context.patientMap}
            onAction={(type, record) => setActionModal({ type, record })}
          />
        )}
        {activeQueue === 'meds' && (
          <MedicationQueue
            reviews={medReviews || []}
            reports={patientReports || []}
            patients={context.patientMap}
            onAction={(type, record) => setActionModal({ type, record })}
          />
        )}
        {activeQueue === 'ops' && (
          <OperationsQueue
            attendance={attendance || []}
            unitEvents={unitEvents || []}
            onAction={(type, record) => setActionModal({ type, record })}
          />
        )}
        {activeQueue === 'interop' && (
          <InteropQueue
            exportsData={exportsData || []}
            graphEdges={graphEdges || []}
            patients={context.patientMap}
            onAction={(type, record) => setActionModal({ type, record })}
          />
        )}

        <DataQualityPanel
          patients={patients || []}
          profiles={profiles || []}
          access={vascularAccess || []}
          safetyChecks={safetyChecks || []}
          adequacy={adequacyReviews || []}
          exportsData={exportsData || []}
        />
      </main>

      <ClinicalActionModal
        key={`${actionModal.type || 'none'}-${actionModal.record?.id || 'new'}`}
        action={actionModal}
        patients={patients || []}
        access={vascularAccess || []}
        error={actionError}
        onClose={() => {
          setActionModal({ type: null, record: null });
          setActionError('');
        }}
        onSubmit={async (payload) => {
          setActionError('');
          try {
            await runClinicalAction(actionModal, payload);
            setActionModal({ type: null, record: null });
          } catch (error) {
            setActionError(error.response?.data?.error || error.message || 'Clinical action failed');
          }
        }}
      />
    </div>
  );
}

function SafetyQueue({ checks, alerts, patients, sessions, labs, dialysate, onAction }) {
  const rows = checks
    .slice()
    .sort((a, b) => Number(b.override_required) - Number(a.override_required) || compareDates(b.checked_at, a.checked_at))
    .slice(0, 12);

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="Pre-Session Safety Gates" className="xl:col-span-2">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="text-left text-xs uppercase tracking-wide text-gray-500">
                <th className="py-3 pr-4">Patient</th>
                <th className="py-3 pr-4">Session</th>
                <th className="py-3 pr-4">Risk</th>
                <th className="py-3 pr-4">Hard Stops</th>
                <th className="py-3 pr-4">Status</th>
                <th className="py-3 pr-4">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {rows.map(check => {
                const session = sessions.get(String(check.session_id));
                return (
                  <tr key={check.id}>
                    <td className="py-3 pr-4 font-medium text-gray-900">{patientName(patients, check.patient_id)}</td>
                    <td className="py-3 pr-4 text-gray-700">{formatDate(session?.scheduled_date || check.checked_at)}</td>
                    <td className="py-3 pr-4"><StatusPill tone={check.risk_score >= 60 ? 'red' : check.risk_score > 0 ? 'yellow' : 'green'}>{check.risk_score || 0}</StatusPill></td>
                    <td className="py-3 pr-4 text-gray-700">{arrayText(check.hard_stop_reasons) || 'None'}</td>
                    <td className="py-3 pr-4"><StatusPill tone={check.override_required ? 'red' : statusTone(check.check_status)}>{humanize(check.check_status)}</StatusPill></td>
                    <td className="py-3 pr-4">
                      <button
                        type="button"
                        onClick={() => onAction(check.override_required ? 'override-safety' : 'clear-safety', check)}
                        className="rounded-lg border border-sky-200 bg-sky-50 px-3 py-1.5 text-xs font-semibold text-sky-700 hover:bg-sky-100"
                      >
                        {check.override_required ? 'Override' : 'Clear'}
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </Panel>

      <Panel title="Alert And Source Signals">
        <div className="space-y-3">
          {alerts.slice(0, 5).map(alert => (
            <SignalRow
              key={alert.id}
              title={alert.title}
              detail={`${patientName(patients, alert.patient_id)} | ${alert.triggering_value || humanize(alert.alert_type)}`}
              tone={alert.severity === 'critical' || alert.severity === 'high' ? 'red' : 'yellow'}
              action={(
                <button
                  type="button"
                  onClick={() => onAction('ack-alert', alert)}
                  className="rounded-lg border border-amber-200 bg-white px-3 py-1.5 text-xs font-semibold text-amber-800 hover:bg-amber-50"
                >
                  Acknowledge
                </button>
              )}
            />
          ))}
          <SignalRow title="Latest lab values" detail={`${labs.length} values across pre, intra, post, routine, and off-session states`} tone="blue" />
          <SignalRow title="Dialysate composition" detail={`${dialysate.filter(item => item.composition_verified).length}/${dialysate.length} verified`} tone={dialysate.some(item => !item.composition_verified) ? 'yellow' : 'green'} />
        </div>
      </Panel>
    </div>
  );
}

function DeliveryQueue({ telemetry, adequacy, profiles, patients }) {
  const reviewByPatient = new Map(adequacy.map(review => [String(review.patient_id), review]));

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="Delivered Treatment Telemetry" className="xl:col-span-2">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {telemetry.slice(0, 8).map(record => {
            const review = reviewByPatient.get(String(record.patient_id));
            return (
              <div key={record.id} className="rounded-lg border border-gray-200 bg-white p-4">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="font-semibold text-gray-900">{patientName(patients, record.patient_id)}</p>
                    <p className="text-xs text-gray-600 mt-1">{formatDate(record.recorded_at)}</p>
                  </div>
                  <StatusPill tone={record.delivered_minutes >= 240 ? 'green' : 'yellow'}>
                    {record.delivered_minutes || 0} min
                  </StatusPill>
                </div>
                <div className="mt-4 grid grid-cols-3 gap-2 text-sm">
                  <Mini label="BFR" value={record.blood_flow_actual ? `${record.blood_flow_actual}` : '--'} />
                  <Mini label="DFR" value={record.dialysate_flow_actual ? `${record.dialysate_flow_actual}` : '--'} />
                  <Mini label="TMP" value={record.tmp_mmhg || '--'} />
                  <Mini label="VP" value={record.venous_pressure_mmhg || '--'} />
                  <Mini label="AP" value={record.arterial_pressure_mmhg || '--'} />
                  <Mini label="BV" value={record.blood_volume_processed_l || '--'} />
                </div>
                <p className="mt-3 text-xs text-gray-700">
                  Kt/V {review?.sp_kt_v || '--'} | URR {review?.urr_percent || '--'}% | UF rate {review?.uf_rate_ml_kg_hr || '--'}
                </p>
              </div>
            );
          })}
        </div>
      </Panel>

      <Panel title="AKI To CKD And Adequacy Review">
        <div className="space-y-3">
          {adequacy.slice(0, 6).map(review => {
            const profile = profiles.get(String(review.patient_id));
            return (
              <SignalRow
                key={review.id}
                title={patientName(patients, review.patient_id)}
                detail={`${profile?.renal_course || 'Renal course pending'} | ${humanize(review.adequacy_status)} | ${review.recommendations || 'No recommendation'}`}
                tone={review.doctor_review_required ? 'yellow' : 'green'}
              />
            );
          })}
        </div>
      </Panel>
    </div>
  );
}

function AccessQueue({ access, lifecycle, infections, patients, onAction }) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="CVC / AVF / AVG Lifecycle" className="xl:col-span-2">
        <div className="mb-4 flex flex-wrap gap-2">
          <button type="button" onClick={() => onAction('access-event', null)}
            className="rounded-lg bg-sky-600 px-3 py-2 text-sm font-semibold text-white hover:bg-sky-700">
            Log CVC / Access Event
          </button>
          <button type="button" onClick={() => onAction('infection-event', null)}
            className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm font-semibold text-red-700 hover:bg-red-100">
            Report Infection Event
          </button>
        </div>
        <div className="space-y-4">
          {lifecycle.slice(0, 10).map(event => (
            <div key={event.id} className="rounded-lg border border-gray-200 bg-white p-4">
              <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
                <div>
                  <p className="font-semibold text-gray-900">{patientName(patients, event.patient_id)}</p>
                  <p className="text-sm text-gray-600">{humanize(event.event_type)} | {formatDate(event.event_date)}</p>
                </div>
                <StatusPill tone={event.culture_result || includesAny(event.long_term_complications, ['infection']) ? 'red' : 'green'}>
                  {event.culture_result ? 'Culture watch' : 'Tracked'}
                </StatusPill>
              </div>
              <div className="mt-3 grid grid-cols-1 md:grid-cols-4 gap-2 text-sm">
                <Mini label="Attempts" value={event.insertion_attempts || '--'} />
                <Mini label="Ultrasound" value={event.ultrasound_used ? 'Yes' : 'No'} />
                <Mini label="Hub scrub" value={event.hub_scrub_compliant ? 'Yes' : 'No'} />
                <Mini label="Plan" value={event.catheter_free_plan || 'Pending'} />
              </div>
            </div>
          ))}
        </div>
      </Panel>

      <Panel title="Infection Surveillance">
        <div className="space-y-3">
          {infections.length === 0 ? (
            <p className="text-sm text-gray-600">No active infection surveillance events.</p>
          ) : infections.map(event => (
            <SignalRow
              key={event.id}
              title={patientName(patients, event.patient_id)}
              detail={`${humanize(event.event_type)} | culture: ${event.organism || 'not recorded'} | registry: ${event.reported_to_registry ? 'sent' : 'pending'}`}
              tone={event.reported_to_registry ? 'green' : 'red'}
            />
          ))}
          <SignalRow title="Access inventory" detail={`${access.length} access records in local tracking`} tone="blue" />
        </div>
      </Panel>
    </div>
  );
}

function MedicationQueue({ reviews, reports, patients, onAction }) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="Medication Reconciliation" className="xl:col-span-2">
        <div className="mb-4 flex flex-wrap gap-2">
          <button type="button" onClick={() => onAction('med-review', null)}
            className="rounded-lg bg-sky-600 px-3 py-2 text-sm font-semibold text-white hover:bg-sky-700">
            Add Medication Review
          </button>
          <button type="button" onClick={() => onAction('patient-followup', null)}
            className="rounded-lg border border-sky-200 bg-sky-50 px-3 py-2 text-sm font-semibold text-sky-700 hover:bg-sky-100">
            Add Follow-Up
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {reviews.slice(0, 8).map(review => (
            <div key={review.id} className="rounded-lg border border-gray-200 bg-white p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="font-semibold text-gray-900">{patientName(patients, review.patient_id)}</p>
                  <p className="text-xs text-gray-600">{humanize(review.review_type)} | {formatDate(review.review_date)}</p>
                </div>
                <StatusPill tone={review.status === 'closed' ? 'green' : 'yellow'}>{humanize(review.status)}</StatusPill>
              </div>
              <div className="mt-3 space-y-2 text-sm text-gray-700">
                <p>Renal dose flags: {arrayText(review.renal_dosing_flags) || 'None'}</p>
                <p>Dialysis timing: {arrayText(review.dialysis_timing_flags) || 'None'}</p>
                <p>Pregnancy cautions: {arrayText(review.pregnancy_cautions) || 'None'}</p>
              </div>
            </div>
          ))}
        </div>
      </Panel>

      <Panel title="Patient Follow-Up">
        <div className="space-y-3">
          {reports.slice(0, 6).map(report => (
            <SignalRow
              key={report.id}
              title={patientName(patients, report.patient_id)}
              detail={`${humanize(report.event_type)} | ${report.follow_up_channel || 'no channel'} | ${report.no_show_reason || 'routine'}`}
              tone={report.payment_barrier || report.transport_reliability === 'unreliable' ? 'yellow' : 'green'}
            />
          ))}
        </div>
      </Panel>
    </div>
  );
}

function OperationsQueue({ attendance, unitEvents, onAction }) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="Biometric Shift Verification" className="xl:col-span-2">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {attendance.map(item => (
            <div key={item.id} className="rounded-lg border border-gray-200 bg-white p-4">
              <p className="font-semibold text-gray-900">{item.station_assignment || 'Station pending'}</p>
              <p className="text-xs text-gray-600 mt-1">{formatDate(item.verification_date)} | {item.biometric_method || 'biometric'}</p>
              <div className="mt-3">
                <StatusPill tone={['verified', 'accepted'].includes(item.verification_result) ? 'green' : 'yellow'}>
                  {humanize(item.verification_result)}
                </StatusPill>
              </div>
              <p className="mt-3 text-xs text-gray-700">{item.exception_reason || 'No exception recorded'}</p>
              {!['verified', 'accepted'].includes(item.verification_result) && (
                <button type="button" onClick={() => onAction('accept-attendance', item)}
                  className="mt-3 w-full rounded-lg border border-amber-200 bg-white px-3 py-2 text-xs font-semibold text-amber-800 hover:bg-amber-50">
                  Accept Exception
                </button>
              )}
            </div>
          ))}
        </div>
      </Panel>

      <Panel title="Water, Machine, Consumables, QI">
        <div className="space-y-3">
          {unitEvents.map(event => (
            <SignalRow
              key={event.id}
              title={humanize(event.event_type)}
              detail={`${humanize(event.severity)} | due ${formatDate(event.closure_due_date)} | ${event.corrective_action || event.immediate_action || 'Action pending'}`}
              tone={event.closed_at ? 'green' : event.severity === 'high' ? 'red' : 'yellow'}
              action={!event.closed_at ? (
                <button type="button" onClick={() => onAction('close-unit-event', event)}
                  className="rounded-lg border border-red-200 bg-white px-3 py-1.5 text-xs font-semibold text-red-700 hover:bg-red-50">
                  Close
                </button>
              ) : null}
            />
          ))}
        </div>
      </Panel>
    </div>
  );
}

function InteropQueue({ exportsData, graphEdges, patients, onAction }) {
  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Panel title="FHIR / Coding Export Readiness" className="xl:col-span-2">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {exportsData.map(item => (
            <div key={item.id} className="rounded-lg border border-gray-200 bg-white p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="font-semibold text-gray-900">{patientName(patients, item.patient_id)}</p>
                  <p className="text-xs text-gray-600">{item.fhir_resource_type || 'FHIR'} | {item.coding_system || 'coding pending'}</p>
                </div>
                <StatusPill tone={item.export_status === 'ready' || item.export_status === 'exported' ? 'green' : 'blue'}>
                  {humanize(item.export_status)}
                </StatusPill>
              </div>
              <p className="mt-3 text-sm text-gray-700">Code: {item.code_value || 'pending'} | UCUM: {item.ucum_unit || 'pending'}</p>
              <div className="mt-3 flex flex-wrap gap-2">
                <button type="button" onClick={() => onAction('mark-export-ready', item)}
                  className="rounded-lg border border-sky-200 bg-sky-50 px-3 py-1.5 text-xs font-semibold text-sky-700 hover:bg-sky-100">
                  Mark Ready
                </button>
                {item.patient_id && (
                  <button type="button" onClick={() => onAction('view-fhir', item)}
                    className="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-xs font-semibold text-gray-700 hover:bg-gray-50">
                    FHIR Bundle
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </Panel>

      <Panel title="Ontology / Neo4j Boundary">
        <div className="space-y-3">
          {graphEdges.slice(0, 8).map(edge => (
            <SignalRow
              key={edge.id}
              title={`${humanize(edge.source_type)} ${humanize(edge.relation_type)}`}
              detail={`${patientName(patients, edge.patient_id)} | target: ${humanize(edge.target_type)} | source: ${edge.provenance}`}
              tone={edge.is_active ? 'green' : 'gray'}
            />
          ))}
        </div>
      </Panel>
    </div>
  );
}

function DataQualityPanel({ patients, profiles, access, safetyChecks, adequacy, exportsData }) {
  const profileIds = new Set(profiles.map(item => String(item.patient_id)));
  const accessIds = new Set(access.map(item => String(item.patient_id)));
  const safetyIds = new Set(safetyChecks.map(item => String(item.patient_id)));
  const adequacyIds = new Set(adequacy.map(item => String(item.patient_id)));
  const exportIds = new Set(exportsData.map(item => String(item.patient_id)));

  const rows = patients.map(patient => {
    const missing = [];
    if (!profileIds.has(String(patient.id))) missing.push('profile');
    if (!accessIds.has(String(patient.id))) missing.push('access');
    if (!safetyIds.has(String(patient.id))) missing.push('safety');
    if (!adequacyIds.has(String(patient.id))) missing.push('adequacy');
    if (!exportIds.has(String(patient.id))) missing.push('FHIR');
    if (!patient.mrn) missing.push('MRN');
    if (!patient.date_of_birth) missing.push('DOB');
    return { patient, missing };
  }).filter(row => row.missing.length > 0);

  return (
    <Panel title="Registry Data Quality">
      {rows.length === 0 ? (
        <p className="text-sm text-gray-600">No local registry gaps detected.</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
          {rows.slice(0, 9).map(row => (
            <div key={row.patient.id} className="rounded-lg border border-amber-200 bg-amber-50 p-3">
              <p className="font-semibold text-gray-900">{row.patient.full_name}</p>
              <p className="mt-1 text-xs text-amber-800">{row.missing.join(', ')}</p>
            </div>
          ))}
        </div>
      )}
    </Panel>
  );
}

function ClinicalActionModal({ action, patients, access, error, onClose, onSubmit }) {
  const [formData, setFormData] = useState(() => initialActionData(action));
  const [saving, setSaving] = useState(false);

  if (!action.type) return null;

  const handleChange = (event) => {
    const { name, value, type, checked } = event.target;
    setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setSaving(true);
    try {
      await onSubmit(formData);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="fixed inset-0 bg-gray-900/70" onClick={onClose}></div>
      <div className="flex min-h-full items-center justify-center p-4">
        <form onSubmit={handleSubmit} className="relative w-full max-w-2xl rounded-lg bg-white shadow-2xl">
          <div className="flex items-center justify-between border-b border-gray-200 px-6 py-4">
            <h3 className="text-lg font-semibold text-gray-900">{actionTitle(action.type)}</h3>
            <button type="button" onClick={onClose} className="text-gray-500 hover:text-gray-800">Close</button>
          </div>

          <div className="space-y-4 px-6 py-5">
            {error && <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>}
            {action.record && (
              <div className="rounded-lg border border-gray-200 bg-gray-50 p-3 text-sm text-gray-700">
                {actionSummary(action)}
              </div>
            )}

            {requiresPatient(action.type) && (
              <ActionSelect label="Patient" name="patient_id" value={formData.patient_id} onChange={handleChange} required
                options={patients.map(patient => ({ value: patient.id, label: `${patient.full_name} (${patient.mrn})` }))} />
            )}

            {action.type === 'override-safety' && (
              <ActionTextArea label="Override reason" name="override_reason" value={formData.override_reason} onChange={handleChange} required />
            )}

            {action.type === 'ack-alert' && (
              <>
                <ActionSelect label="Alert status" name="status" value={formData.status} onChange={handleChange}
                  options={[
                    { value: 'acknowledged', label: 'Acknowledged' },
                    { value: 'resolved', label: 'Resolved' },
                    { value: 'closed', label: 'Closed' },
                  ]} />
                <ActionTextArea label="Action / override note" name="override_reason" value={formData.override_reason} onChange={handleChange} />
              </>
            )}

            {action.type === 'access-event' && (
              <>
                <ActionSelect label="Access record" name="access_id" value={formData.access_id} onChange={handleChange}
                  options={access.filter(item => !formData.patient_id || String(item.patient_id) === String(formData.patient_id)).map(item => ({ value: item.id, label: `${humanize(item.access_type)} ${humanize(item.access_site)}` }))} />
                <ActionSelect label="Event type" name="event_type" value={formData.event_type} onChange={handleChange}
                  options={[
                    { value: 'cvc_insertion', label: 'CVC insertion' },
                    { value: 'cvc_replacement', label: 'CVC replacement' },
                    { value: 'cvc_removal', label: 'CVC removal' },
                    { value: 'avf_maturation_check', label: 'AVF maturation check' },
                    { value: 'access_complication', label: 'Access complication' },
                  ]} />
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <ActionInput label="Event date" name="event_date" type="date" value={formData.event_date} onChange={handleChange} />
                  <ActionInput label="Attempts / pricks" name="insertion_attempts" type="number" value={formData.insertion_attempts} onChange={handleChange} />
                  <ActionInput label="Operator" name="operator_name" value={formData.operator_name} onChange={handleChange} />
                </div>
                <ActionInput label="Immediate complications" name="immediate_complications" value={formData.immediate_complications} onChange={handleChange} placeholder="bleeding, difficult insertion" />
                <ActionInput label="Long-term complications" name="long_term_complications" value={formData.long_term_complications} onChange={handleChange} placeholder="infection, blockage, thrombosis" />
                <ActionInput label="Culture result" name="culture_result" value={formData.culture_result} onChange={handleChange} />
                <ActionTextArea label="Catheter-free plan / notes" name="catheter_free_plan" value={formData.catheter_free_plan} onChange={handleChange} />
              </>
            )}

            {action.type === 'infection-event' && (
              <>
                <ActionSelect label="Event type" name="event_type" value={formData.event_type} onChange={handleChange}
                  options={[
                    { value: 'iv_antimicrobial_start', label: 'IV antimicrobial start' },
                    { value: 'positive_blood_culture', label: 'Positive blood culture' },
                    { value: 'access_site_pus_redness_swelling', label: 'Access pus/redness/swelling' },
                    { value: 'suspected_bsi', label: 'Suspected bloodstream infection' },
                  ]} />
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <ActionInput label="Event date" name="event_date" type="date" value={formData.event_date} onChange={handleChange} />
                  <ActionInput label="Suspected source" name="suspected_source" value={formData.suspected_source} onChange={handleChange} />
                  <ActionInput label="Organism" name="organism" value={formData.organism} onChange={handleChange} />
                </div>
                <ActionChecks fields={[
                  ['iv_antimicrobial_started', 'IV antimicrobial started'],
                  ['positive_blood_culture', 'Positive blood culture'],
                  ['access_pus_redness_swelling', 'Pus/redness/swelling'],
                  ['hospitalized', 'Hospitalized'],
                  ['death_related', 'Death related'],
                  ['reported_to_registry', 'Reported to registry'],
                ]} values={formData} onChange={handleChange} />
                <ActionTextArea label="Notes" name="notes" value={formData.notes} onChange={handleChange} />
              </>
            )}

            {action.type === 'med-review' && (
              <>
                <ActionSelect label="Review type" name="review_type" value={formData.review_type} onChange={handleChange}
                  options={[
                    { value: 'admission', label: 'Admission' },
                    { value: 'monthly', label: 'Monthly' },
                    { value: 'post_hospitalization', label: 'Post-hospitalization' },
                    { value: 'prescription_change', label: 'Prescription change' },
                  ]} />
                <ActionInput label="Renal dosing flags" name="renal_dosing_flags" value={formData.renal_dosing_flags} onChange={handleChange} />
                <ActionInput label="Dialysis timing flags" name="dialysis_timing_flags" value={formData.dialysis_timing_flags} onChange={handleChange} />
                <ActionInput label="Regimen change reason" name="regimen_change_reason" value={formData.regimen_change_reason} onChange={handleChange} />
                <ActionTextArea label="Recommendations" name="recommendations" value={formData.recommendations} onChange={handleChange} />
              </>
            )}

            {action.type === 'patient-followup' && (
              <>
                <ActionInput label="Symptoms" name="symptoms" value={formData.symptoms} onChange={handleChange} placeholder="cramps, dizziness, access pain" />
                <ActionSelect label="Transport reliability" name="transport_reliability" value={formData.transport_reliability} onChange={handleChange}
                  options={[
                    { value: 'reliable', label: 'Reliable' },
                    { value: 'unreliable', label: 'Unreliable' },
                    { value: 'unknown', label: 'Unknown' },
                  ]} />
                <ActionSelect label="Follow-up channel" name="follow_up_channel" value={formData.follow_up_channel} onChange={handleChange}
                  options={[
                    { value: 'sms', label: 'SMS' },
                    { value: 'whatsapp', label: 'WhatsApp' },
                    { value: 'phone', label: 'Phone' },
                    { value: 'clinic', label: 'Clinic review' },
                  ]} />
                <ActionChecks fields={[
                  ['sms_whatsapp_consent', 'SMS/WhatsApp consent'],
                  ['payment_barrier', 'Payment barrier'],
                  ['food_insecurity', 'Food insecurity'],
                  ['teach_back_completed', 'Teach-back completed'],
                ]} values={formData} onChange={handleChange} />
                <ActionTextArea label="Notes" name="notes" value={formData.notes} onChange={handleChange} />
              </>
            )}

            {['clear-safety', 'accept-attendance', 'close-unit-event', 'mark-export-ready', 'view-fhir'].includes(action.type) && (
              <ActionTextArea label="Note" name="notes" value={formData.notes} onChange={handleChange} />
            )}
          </div>

          <div className="flex justify-end gap-3 border-t border-gray-200 px-6 py-4">
            <button type="button" onClick={onClose} disabled={saving} className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 disabled:opacity-50">Cancel</button>
            <button type="submit" disabled={saving} className="rounded-lg bg-sky-600 px-4 py-2 text-sm font-semibold text-white hover:bg-sky-700 disabled:opacity-50">
              {saving ? 'Saving...' : actionButton(action.type)}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

function ActionInput({ label, name, value, onChange, type = 'text', required = false, placeholder = '' }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}{required && <span className="text-red-500 ml-1">*</span>}</span>
      <input name={name} type={type} value={value} onChange={onChange} required={required} placeholder={placeholder}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500" />
    </label>
  );
}

function ActionTextArea({ label, name, value, onChange, required = false }) {
  return (
    <label className="block">
      <span className="block text-sm font-medium text-gray-700 mb-1">{label}{required && <span className="text-red-500 ml-1">*</span>}</span>
      <textarea name={name} value={value} onChange={onChange} required={required} rows={3}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-sky-500" />
    </label>
  );
}

function ActionSelect({ label, name, value, onChange, options, required = false }) {
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

function ActionChecks({ fields, values, onChange }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
      {fields.map(([name, label]) => (
        <label key={name} className="flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-700">
          <input type="checkbox" name={name} checked={Boolean(values[name])} onChange={onChange} className="h-4 w-4 rounded text-sky-600" />
          {label}
        </label>
      ))}
    </div>
  );
}

async function runClinicalAction(action, payload) {
  const { type, record } = action;
  const currentUser = authService.getCurrentUser();
  const now = new Date().toISOString();
  const hospitalId = currentUser?.hospital_id || record?.hospital_id || 'demo_hospital';

  if (type === 'ack-alert') {
    const updates = {
      status: payload.status || 'acknowledged',
      acknowledged_by: currentUser?.id,
      acknowledged_at: now,
      override_reason: payload.override_reason || '',
    };
    if (isUuid(record?.id) && navigator.onLine) {
      const response = await api.post(`/clinical-alerts/${record.id}/acknowledge`, updates);
      await putLocal('clinical_alerts', response.data);
      return;
    }
    await putLocal('clinical_alerts', { ...record, ...updates });
    return;
  }

  if (type === 'override-safety') {
    const updates = {
      check_status: 'overridden',
      override_required: false,
      override_reason: payload.override_reason,
      override_approved_by: currentUser?.id,
      override_approved_at: now,
    };
    if (isUuid(record?.id) && navigator.onLine) {
      const response = await api.post(`/session-safety-checks/${record.id}/override`, updates);
      await putLocal('session_safety_checks', response.data);
      return;
    }
    await putLocal('session_safety_checks', { ...record, ...updates });
    return;
  }

  if (type === 'clear-safety') {
    await putLocal('session_safety_checks', {
      ...record,
      check_status: 'cleared',
      override_required: false,
      notes: payload.notes || record?.notes || 'Cleared from command center.',
      updated_at: now,
    });
    return;
  }

  if (type === 'access-event') {
    await offlineService.create('access_lifecycle_events', {
      hospital_id: hospitalId,
      patient_id: payload.patient_id,
      access_id: payload.access_id || null,
      event_type: payload.event_type || 'access_complication',
      event_date: payload.event_date || now.slice(0, 10),
      operator_name: payload.operator_name || '',
      insertion_attempts: payload.insertion_attempts ? Number(payload.insertion_attempts) : null,
      immediate_complications: csvList(payload.immediate_complications),
      long_term_complications: csvList(payload.long_term_complications),
      culture_result: payload.culture_result || '',
      catheter_free_plan: payload.catheter_free_plan || '',
      created_by: currentUser?.id,
      notes: payload.notes || '',
    }, 3);
    return;
  }

  if (type === 'infection-event') {
    await offlineService.create('infection_surveillance_events', {
      hospital_id: hospitalId,
      patient_id: payload.patient_id,
      event_type: payload.event_type || 'suspected_bsi',
      event_date: payload.event_date || now.slice(0, 10),
      iv_antimicrobial_started: Boolean(payload.iv_antimicrobial_started),
      positive_blood_culture: Boolean(payload.positive_blood_culture),
      access_pus_redness_swelling: Boolean(payload.access_pus_redness_swelling),
      suspected_source: payload.suspected_source || '',
      organism: payload.organism || '',
      hospitalized: Boolean(payload.hospitalized),
      death_related: Boolean(payload.death_related),
      reported_to_registry: Boolean(payload.reported_to_registry),
      reported_by: currentUser?.id,
      notes: payload.notes || '',
    }, 2);
    return;
  }

  if (type === 'med-review') {
    await offlineService.create('medication_reconciliation_reviews', {
      hospital_id: hospitalId,
      patient_id: payload.patient_id,
      review_type: payload.review_type || 'monthly',
      review_date: now.slice(0, 10),
      reviewed_by: currentUser?.id,
      renal_dosing_flags: csvList(payload.renal_dosing_flags),
      dialysis_timing_flags: csvList(payload.dialysis_timing_flags),
      regimen_change_reason: payload.regimen_change_reason || '',
      recommendations: payload.recommendations || '',
      status: 'open',
    }, 4);
    return;
  }

  if (type === 'patient-followup') {
    await offlineService.create('patient_reported_events', {
      hospital_id: hospitalId,
      patient_id: payload.patient_id,
      event_type: 'follow_up',
      reported_at: now,
      reported_by: currentUser?.id,
      symptoms: csvList(payload.symptoms),
      transport_reliability: payload.transport_reliability || 'unknown',
      sms_whatsapp_consent: Boolean(payload.sms_whatsapp_consent),
      follow_up_channel: payload.follow_up_channel || 'phone',
      payment_barrier: Boolean(payload.payment_barrier),
      food_insecurity: Boolean(payload.food_insecurity),
      teach_back_completed: Boolean(payload.teach_back_completed),
      notes: payload.notes || '',
    }, 5);
    return;
  }

  if (type === 'accept-attendance') {
    await putLocal('staff_attendance_verifications', {
      ...record,
      verification_result: 'accepted',
      exception_reason: payload.notes || record?.exception_reason || 'Accepted by command center reviewer.',
      updated_at: now,
    });
    return;
  }

  if (type === 'close-unit-event') {
    await putLocal('unit_safety_events', {
      ...record,
      closed_at: now,
      corrective_action: payload.notes || record?.corrective_action || 'Closed by command center reviewer.',
      updated_at: now,
    });
    return;
  }

  if (type === 'mark-export-ready' || type === 'view-fhir') {
    let fhirPayload = record?.fhir_payload || null;
    if (type === 'view-fhir' && record?.patient_id && isUuid(record.patient_id) && navigator.onLine) {
      const response = await api.get(`/patients/${record.patient_id}/fhir-summary`);
      fhirPayload = response.data;
    }
    await putLocal('interoperability_exports', {
      ...record,
      export_status: 'ready',
      fhir_payload: fhirPayload,
      error_message: '',
      updated_at: now,
    });
  }
}

async function putLocal(entityType, record) {
  if (!record?.id) return;
  await db.table(entityType).put({
    ...record,
    synced: false,
    updated_at: new Date().toISOString(),
  });
  window.dispatchEvent(new CustomEvent('dms-local-change', {
    detail: { entityType, record },
  }));
}

function initialActionData(action) {
  const today = new Date().toISOString().slice(0, 10);
  const patientId = action.record?.patient_id || '';
  return {
    patient_id: patientId,
    status: 'acknowledged',
    override_reason: '',
    notes: '',
    access_id: action.record?.access_id || '',
    event_type: '',
    event_date: today,
    insertion_attempts: '',
    operator_name: '',
    immediate_complications: '',
    long_term_complications: '',
    culture_result: '',
    catheter_free_plan: '',
    suspected_source: '',
    organism: '',
    review_type: 'monthly',
    renal_dosing_flags: '',
    dialysis_timing_flags: '',
    regimen_change_reason: '',
    recommendations: '',
    symptoms: '',
    transport_reliability: 'unknown',
    follow_up_channel: 'phone',
    sms_whatsapp_consent: true,
    reported_to_registry: false,
  };
}

function actionTitle(type) {
  const titles = {
    'override-safety': 'Override Safety Gate',
    'clear-safety': 'Clear Safety Gate',
    'ack-alert': 'Acknowledge Clinical Alert',
    'access-event': 'Log CVC / Access Event',
    'infection-event': 'Report Infection Event',
    'med-review': 'Medication Reconciliation Review',
    'patient-followup': 'Patient Follow-Up',
    'accept-attendance': 'Accept Attendance Exception',
    'close-unit-event': 'Close Unit Safety Event',
    'mark-export-ready': 'Mark Export Ready',
    'view-fhir': 'Prepare FHIR Bundle',
  };
  return titles[type] || 'Clinical Action';
}

function actionButton(type) {
  if (type === 'view-fhir') return 'Prepare Bundle';
  if (type === 'mark-export-ready') return 'Mark Ready';
  if (type === 'close-unit-event') return 'Close Event';
  if (type === 'accept-attendance') return 'Accept';
  if (type === 'ack-alert') return 'Acknowledge';
  return 'Save';
}

function actionSummary(action) {
  const record = action.record || {};
  if (action.type === 'ack-alert') return `${record.title || humanize(record.alert_type)} | ${record.triggering_value || 'no triggering value'}`;
  if (action.type === 'override-safety' || action.type === 'clear-safety') return `Safety check ${humanize(record.check_status)} | hard stops: ${arrayText(record.hard_stop_reasons) || 'none'}`;
  if (action.type === 'close-unit-event') return `${humanize(record.event_type)} | ${record.corrective_action || record.immediate_action || 'action pending'}`;
  return record.id ? `Record ${record.id}` : 'Create a new command-center record.';
}

function requiresPatient(type) {
  return ['access-event', 'infection-event', 'med-review', 'patient-followup'].includes(type);
}

function csvList(value) {
  if (!value) return [];
  return String(value).split(',').map(item => item.trim()).filter(Boolean);
}

function isUuid(value) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(String(value));
}

function KpiCard({ label, value, tone }) {
  return (
    <div className={`rounded-lg border p-3 ${toneClasses(tone).soft}`}>
      <p className="text-[11px] font-semibold uppercase tracking-wide text-gray-600">{label}</p>
      <p className={`mt-1 text-2xl font-bold ${toneClasses(tone).text}`}>{value}</p>
    </div>
  );
}

function Panel({ title, children, className = '' }) {
  return (
    <section className={`rounded-lg border border-gray-200 bg-white p-5 ${className}`}>
      <h2 className="text-lg font-semibold text-gray-900 mb-4">{title}</h2>
      {children}
    </section>
  );
}

function SignalRow({ title, detail, tone = 'gray', action = null }) {
  return (
    <div className={`rounded-lg border p-3 ${toneClasses(tone).soft}`}>
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-gray-900">{title || 'Untitled'}</p>
          <p className="text-xs text-gray-700 mt-1">{detail || 'No detail'}</p>
        </div>
        {action || <span className={`mt-0.5 h-2.5 w-2.5 rounded-full ${toneClasses(tone).dot}`}></span>}
      </div>
    </div>
  );
}

function Mini({ label, value }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 min-h-16">
      <p className="text-[11px] font-semibold uppercase tracking-wide text-gray-500">{label}</p>
      <p className="mt-1 text-sm font-semibold text-gray-900 break-words">{value}</p>
    </div>
  );
}

function StatusPill({ tone = 'gray', children }) {
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold ${toneClasses(tone).pill}`}>
      {children}
    </span>
  );
}

function toneClasses(tone) {
  const tones = {
    green: { soft: 'bg-emerald-50 border-emerald-200', text: 'text-emerald-700', pill: 'bg-emerald-100 text-emerald-800', dot: 'bg-emerald-500' },
    red: { soft: 'bg-red-50 border-red-200', text: 'text-red-700', pill: 'bg-red-100 text-red-800', dot: 'bg-red-500' },
    yellow: { soft: 'bg-amber-50 border-amber-200', text: 'text-amber-700', pill: 'bg-amber-100 text-amber-800', dot: 'bg-amber-500' },
    blue: { soft: 'bg-sky-50 border-sky-200', text: 'text-sky-700', pill: 'bg-sky-100 text-sky-800', dot: 'bg-sky-500' },
    purple: { soft: 'bg-violet-50 border-violet-200', text: 'text-violet-700', pill: 'bg-violet-100 text-violet-800', dot: 'bg-violet-500' },
    gray: { soft: 'bg-gray-50 border-gray-200', text: 'text-gray-700', pill: 'bg-gray-100 text-gray-800', dot: 'bg-gray-400' },
  };
  return tones[tone] || tones.gray;
}

function patientName(patientMap, patientId) {
  const patient = patientMap.get(String(patientId));
  return patient?.full_name || patient?.mrn || 'Unknown patient';
}

function statusTone(status) {
  if (['cleared', 'completed', 'verified', 'adequate', 'ready'].includes(status)) return 'green';
  if (['blocked', 'failed', 'critical', 'high'].includes(status)) return 'red';
  if (['pending', 'open', 'needs_review', 'doctor_review_required'].includes(status)) return 'yellow';
  return 'gray';
}

function humanize(value) {
  if (value === null || value === undefined || value === '') return 'Not recorded';
  return String(value).replaceAll('_', ' ').replaceAll('-', ' ');
}

function formatDate(value) {
  if (!value) return 'N/A';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value).slice(0, 10);
  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: '2-digit' });
}

function compareDates(a, b) {
  return new Date(a || 0).getTime() - new Date(b || 0).getTime();
}

function arrayText(value) {
  if (!value) return '';
  if (Array.isArray(value)) return value.map(humanize).join(', ');
  if (typeof value === 'string') return humanize(value);
  return Object.values(value).map(humanize).join(', ');
}

function includesAny(value, needles) {
  const text = arrayText(value).toLowerCase();
  return needles.some(needle => text.includes(needle));
}
