import { useMemo, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import useOfflineData from '../../hooks/useOfflineData';
import FormModal from '../../components/forms/FormModal';
import PatientForm from '../../components/forms/PatientForm';
import { htmlTable, keyValueTable, printHtml } from '../../utils/print';

const LAB_TRACKING_GROUPS = [
  {
    id: 'potassium',
    label: 'Potassium',
    shortLabel: 'K',
    keywords: ['potassium', 'k'],
    unit: 'mmol/L',
    low: 3.5,
    high: 5.5,
    criticalHigh: 6,
  },
  {
    id: 'phosphate',
    label: 'Phosphate',
    shortLabel: 'Phosphate',
    keywords: ['phosphate', 'phosphorus', 'phos'],
    unit: 'mmol/L',
    high: 1.8,
  },
  {
    id: 'protein',
    label: 'Protein / Albumin',
    shortLabel: 'Protein',
    keywords: ['protein', 'albumin', 'alb'],
    unit: 'g/L',
    low: 35,
    criticalLow: 30,
  },
  {
    id: 'hemoglobin',
    label: 'Hemoglobin',
    shortLabel: 'Hb',
    keywords: ['hemoglobin', 'haemoglobin', 'hb'],
    unit: 'g/dL',
    low: 10,
    criticalLow: 8.5,
  },
  {
    id: 'creatinine',
    label: 'Creatinine',
    shortLabel: 'Creatinine',
    keywords: ['creatinine', 'cr'],
    unit: 'umol/L',
  },
  {
    id: 'glucose',
    label: 'Glucose',
    shortLabel: 'Glucose',
    keywords: ['glucose', 'blood sugar', 'glu'],
    unit: 'mmol/L',
    low: 4,
    high: 10,
  },
];

const PHASE_LABELS = {
  pre_dialysis: 'Pre-dialysis',
  routine: 'Routine',
  intra_dialysis: 'Intra-dialysis',
  intra_complication: 'Intra-dialysis complication',
  post_dialysis: 'Post-dialysis',
  off_session: 'Off-session',
};

export default function PatientDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { data: patients, loading } = useOfflineData('patients');
  const { data: sessions } = useOfflineData('dialysis_sessions');
  const { data: contacts } = useOfflineData('patient_contacts');
  const { data: vascularAccess } = useOfflineData('vascular_access');
  const { data: accessAssessments } = useOfflineData('vascular_access_assessments');
  const { data: labResults } = useOfflineData('lab_results');
  const { data: labTests } = useOfflineData('lab_test_catalog');
  const { data: vitals } = useOfflineData('session_vitals');
  const { data: dialysateRecords } = useOfflineData('dialysate_records');
  const { data: diagnoses } = useOfflineData('diagnoses');
  const { data: comorbidities } = useOfflineData('comorbidities');
  const { data: consents } = useOfflineData('consents');
  const { data: criticalAlerts } = useOfflineData('lab_critical_alerts');
  const { data: prescriptions } = useOfflineData('prescriptions');
  const { data: prescriptionItems } = useOfflineData('prescription_items');
  const { data: medications } = useOfflineData('medications');
  const { data: clinicalProfiles } = useOfflineData('patient_clinical_profiles');
  const { data: accessLifecycleEvents } = useOfflineData('access_lifecycle_events');
  const { data: infectionEvents } = useOfflineData('infection_surveillance_events');

  const [activeTab, setActiveTab] = useState('clinical');
  const [showEditModal, setShowEditModal] = useState(false);

  const patient = useMemo(
    () => patients?.find(p => idMatches(p.id, id)),
    [patients, id]
  );
  const clinicalProfile = useMemo(
    () => (clinicalProfiles || []).find(profile => idMatches(profile.patient_id, id)),
    [clinicalProfiles, id]
  );
  const patientRecord = useMemo(
    () => patient ? { ...patient, ...(clinicalProfile || {}) } : null,
    [patient, clinicalProfile]
  );

  const patientSessions = useMemo(
    () => sortByDate((sessions || []).filter(s => idMatches(s.patient_id, id)), 'scheduled_date'),
    [sessions, id]
  );
  const patientContacts = useMemo(
    () => (contacts || []).filter(c => idMatches(c.patient_id, id)),
    [contacts, id]
  );
  const patientAccess = useMemo(
    () => sortByDate((vascularAccess || []).filter(v => idMatches(v.patient_id, id)), 'insertion_date'),
    [vascularAccess, id]
  );
  const patientAccessAssessments = useMemo(
    () => sortByDate((accessAssessments || []).filter(a => idMatches(a.patient_id, id)), 'assessed_at'),
    [accessAssessments, id]
  );
  const patientAccessLifecycle = useMemo(
    () => sortByDate((accessLifecycleEvents || []).filter(a => idMatches(a.patient_id, id)), 'event_date'),
    [accessLifecycleEvents, id]
  );
  const patientInfections = useMemo(
    () => sortByDate((infectionEvents || []).filter(a => idMatches(a.patient_id, id)), 'event_date'),
    [infectionEvents, id]
  );
  const patientVitals = useMemo(
    () => sortByDate((vitals || []).filter(v => idMatches(v.patient_id, id)), 'recorded_at'),
    [vitals, id]
  );
  const patientDialysate = useMemo(
    () => sortByDate((dialysateRecords || []).filter(d => idMatches(d.patient_id, id)), 'recorded_at'),
    [dialysateRecords, id]
  );
  const patientLabs = useMemo(
    () => buildLabTimeline((labResults || []).filter(r => idMatches(r.patient_id, id)), labTests || []),
    [labResults, labTests, id]
  );
  const patientDiagnoses = useMemo(
    () => sortByDate((diagnoses || []).filter(d => idMatches(d.patient_id, id)), 'diagnosed_at'),
    [diagnoses, id]
  );
  const patientComorbidities = useMemo(
    () => sortByDate((comorbidities || []).filter(c => idMatches(c.patient_id, id)), 'diagnosed_at'),
    [comorbidities, id]
  );
  const patientConsents = useMemo(
    () => sortByDate((consents || []).filter(c => idMatches(c.patient_id, id)), 'signed_at'),
    [consents, id]
  );
  const patientAlerts = useMemo(
    () => sortByDate((criticalAlerts || []).filter(a => idMatches(a.patient_id, id)), 'created_at'),
    [criticalAlerts, id]
  );
  const patientPrescriptions = useMemo(
    () => sortByDate((prescriptions || []).filter(p => idMatches(p.patient_id, id)), 'prescribed_at'),
    [prescriptions, id]
  );
  const patientPrescriptionItems = useMemo(
    () => (prescriptionItems || []).filter(item => {
      if (idMatches(item.patient_id, id)) return true;
      return patientPrescriptions.some(rx => idMatches(rx.id, item.prescription_id));
    }),
    [prescriptionItems, patientPrescriptions, id]
  );

  const snapshot = useMemo(() => {
    if (!patientRecord) return null;
    return buildClinicalSnapshot({
      patient: patientRecord,
      sessions: patientSessions,
      contacts: patientContacts,
      access: patientAccess,
      assessments: patientAccessAssessments,
      labs: patientLabs,
      vitals: patientVitals,
      dialysate: patientDialysate,
      diagnoses: patientDiagnoses,
      comorbidities: patientComorbidities,
      consents: patientConsents,
      alerts: patientAlerts,
    });
  }, [
    patientRecord,
    patientSessions,
    patientContacts,
    patientAccess,
    patientAccessAssessments,
    patientLabs,
    patientVitals,
    patientDialysate,
    patientDiagnoses,
    patientComorbidities,
    patientConsents,
    patientAlerts,
  ]);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-600"></div>
      </div>
    );
  }

  if (!patient) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Patient Not Found</h2>
          <p className="text-gray-600 mb-4">The patient record is not available in local storage.</p>
          <Link to="/patients" className="text-sky-600 hover:underline">
            Back to Patients
          </Link>
        </div>
      </div>
    );
  }

  const age = calculateAge(patientRecord.date_of_birth);

  const tabs = [
    { id: 'clinical', label: 'Clinical Tracking' },
    { id: 'overview', label: 'Overview' },
    { id: 'sessions', label: 'Sessions' },
    { id: 'labs', label: 'Labs' },
    { id: 'medications', label: 'Medications' },
    { id: 'access', label: 'Access / CVC' },
  ];

  const handlePrintPatientHistory = () => {
    printHtml({
      title: 'Patient History',
      subtitle: `${patientRecord.full_name || 'Unknown Patient'} | MRN ${patientRecord.mrn || 'N/A'}`,
      body: buildPatientHistoryPrintBody({
        patient: patientRecord,
        snapshot,
        contacts: patientContacts,
        sessions: patientSessions,
        labs: patientLabs,
        diagnoses: patientDiagnoses,
        comorbidities: patientComorbidities,
        accessData: patientAccess,
        prescriptionItems: patientPrescriptionItems,
        medications: medications || [],
      }),
      footer: 'Printed from DMS patient clinical record.',
    });
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div className="flex items-start gap-4">
              <button
                onClick={() => navigate('/patients')}
                className="px-3 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Back
              </button>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">
                  {patientRecord.full_name || 'Unknown Patient'}
                </h1>
                <p className="text-sm text-gray-600 mt-1">
                  MRN: {patientRecord.mrn || 'N/A'} | Age: {age} | Sex: {patientRecord.sex || 'N/A'}
                </p>
                <div className="mt-3 flex flex-wrap gap-2">
                  <StatusPill tone={patientRecord.is_active ? 'green' : 'gray'}>
                    {patientRecord.is_active ? 'Active' : 'Inactive'}
                  </StatusPill>
                  <StatusPill tone={snapshot?.courseTone || 'gray'}>
                    {snapshot?.renalCourse || 'Renal course not staged'}
                  </StatusPill>
                  {snapshot?.smsReady ? (
                    <StatusPill tone="blue">SMS contact ready</StatusPill>
                  ) : (
                    <StatusPill tone="yellow">SMS contact missing</StatusPill>
                  )}
                </div>
              </div>
            </div>
            <div className="flex flex-wrap gap-2 self-start">
              <button
                onClick={handlePrintPatientHistory}
                className="px-4 py-2 border border-gray-300 bg-white text-gray-700 rounded-lg hover:bg-gray-50"
              >
                Print History
              </button>
              <button
                onClick={() => setShowEditModal(true)}
                className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700"
              >
                Edit Patient
              </button>
            </div>
          </div>

          <div className="mt-6 flex flex-wrap gap-1 border-b border-gray-200">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`px-3 sm:px-4 py-2 text-sm font-medium transition-colors ${
                  activeTab === tab.id
                    ? 'text-sky-700 border-b-2 border-sky-600 bg-sky-50'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8">
        {activeTab === 'clinical' && (
          <ClinicalTrackingTab
            snapshot={snapshot}
            patient={patientRecord}
            labs={patientLabs}
            diagnoses={patientDiagnoses}
            comorbidities={patientComorbidities}
          />
        )}
        {activeTab === 'overview' && (
          <OverviewTab
            patient={patientRecord}
            contacts={patientContacts}
            diagnoses={patientDiagnoses}
            comorbidities={patientComorbidities}
          />
        )}
        {activeTab === 'sessions' && (
          <SessionsTab
            sessions={patientSessions}
            vitals={patientVitals}
            dialysate={patientDialysate}
          />
        )}
        {activeTab === 'labs' && <LabsTab labs={patientLabs} />}
        {activeTab === 'medications' && (
          <MedicationsTab
            prescriptions={patientPrescriptions}
            prescriptionItems={patientPrescriptionItems}
            medications={medications || []}
          />
        )}
        {activeTab === 'access' && (
          <AccessTab
            accessData={patientAccess}
            assessments={patientAccessAssessments}
            lifecycle={patientAccessLifecycle}
            infections={patientInfections}
          />
        )}
      </div>

      <FormModal
        isOpen={showEditModal}
        onClose={() => setShowEditModal(false)}
        title="Edit Patient"
        size="xl"
      >
        <PatientForm
          patient={patientRecord}
          onSuccess={() => {
            setShowEditModal(false);
          }}
          onCancel={() => setShowEditModal(false)}
        />
      </FormModal>
    </div>
  );
}

function ClinicalTrackingTab({ snapshot, patient, labs, diagnoses, comorbidities }) {
  if (!snapshot) return null;

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <MetricCard label="Renal course" value={snapshot.renalCourse} detail={snapshot.courseDetail} tone={snapshot.courseTone} />
        <MetricCard label="Dialysis exposure" value={snapshot.dialysisVintage} detail={snapshot.dialysisDecision} tone={snapshot.dialysisTone} />
        <MetricCard label="Primary access" value={snapshot.primaryAccessLabel} detail={snapshot.accessDetail} tone={snapshot.accessTone} />
        <MetricCard label="Latest BP" value={snapshot.latestBp} detail={snapshot.vitalsDetail} tone={snapshot.vitalsTone} />
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <Panel title="Session Safety Gates" className="xl:col-span-2">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {snapshot.safetyGates.map(gate => (
              <SafetyGate key={gate.label} gate={gate} />
            ))}
          </div>
        </Panel>

        <Panel title="Risk Flags">
          {snapshot.riskFlags.length === 0 ? (
            <p className="text-sm text-gray-600">No active risk flags from local records.</p>
          ) : (
            <div className="space-y-3">
              {snapshot.riskFlags.map(flag => (
                <div key={flag.label} className={`rounded-lg border p-3 ${toneClasses(flag.tone).soft}`}>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className={`text-sm font-semibold ${toneClasses(flag.tone).text}`}>{flag.label}</p>
                      <p className="text-xs text-gray-700 mt-1">{flag.detail}</p>
                    </div>
                    <StatusPill tone={flag.tone}>{flag.status}</StatusPill>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <Panel title="3P, Hb, Creatinine, Glucose Trends" className="xl:col-span-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {LAB_TRACKING_GROUPS.map(group => (
              <LabTrendCard key={group.id} group={group} labs={labs} compact />
            ))}
          </div>
        </Panel>

        <Panel title="Disease Chain">
          <TimelineItem label="Symptom onset" value={formatDate(patient.symptom_onset_date)} />
          <TimelineItem label="Diagnosis date" value={formatDate(patient.diagnosis_date)} />
          <TimelineItem label="Dialysis start" value={formatDate(patient.dialysis_start_date)} />
          <TimelineItem label="AKI-to-CKD review due" value={formatDate(patient.aki_to_ckd_reassessment_due)} />
          <TimelineItem label="Kidney cause" value={patient.kidney_disease_cause || firstValue(diagnoses, 'description') || 'Not recorded'} />
          <TimelineItem label="Dialysis indication" value={patient.dialysis_indication || firstValue(diagnoses, 'notes') || 'Not recorded'} />
          <TimelineItem label="Creatinine / urine output" value={`${patient.baseline_creatinine_mg_dl || '--'} baseline | ${patient.highest_creatinine_mg_dl || '--'} highest | urine ${patient.urine_output_ml_day || '--'} ml/day`} />
          <TimelineItem label="Recovery plan" value={patient.recovery_plan || patient.residual_kidney_function || 'Not recorded'} />
          {patient.pregnancy_related && (
            <>
              <TimelineItem label="Pregnancy context" value={`${patient.gravida || 'G?'} ${patient.para || 'P?'} | GA ${patient.gestational_age_weeks || '?'} weeks | ANC visits: ${patient.anc_visits || 'not recorded'}`} />
              <TimelineItem label="Preeclampsia / delivery" value={`${humanize(patient.preeclampsia_severity || 'not recorded')} | EDD ${formatDate(patient.expected_delivery_date)} | delivery ${formatDate(patient.delivery_date)}`} />
              <TimelineItem label="Maternal-fetal follow-up" value={patient.maternal_fetal_follow_up || patient.postpartum_renal_recovery || 'Not recorded'} />
            </>
          )}
          {patient.malaria_aki_phenotype && patient.malaria_aki_phenotype !== 'not_applicable' && (
            <>
              <TimelineItem label="Malaria phenotype" value={humanize(patient.malaria_aki_phenotype)} />
              <TimelineItem label="Malaria dates" value={`symptoms ${formatDate(patient.malaria_symptom_onset_date)} | test ${formatDate(patient.malaria_test_date)} | first dose ${formatDate(patient.first_antimalarial_dose_at)}`} />
              <TimelineItem label="Definitive antimalarial regimen" value={patient.definitive_antimalarial_regimen || 'Not recorded'} />
            </>
          )}
          <div className="pt-3 border-t border-gray-200">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Comorbidities</p>
            <div className="flex flex-wrap gap-2">
              {comorbidities.length > 0 ? comorbidities.map(item => (
                <StatusPill key={item.id || item.condition} tone={item.status === 'controlled' ? 'green' : 'gray'}>
                  {item.condition}
                </StatusPill>
              )) : <span className="text-sm text-gray-600">No comorbidities recorded.</span>}
            </div>
          </div>
        </Panel>
      </div>
    </div>
  );
}

function OverviewTab({ patient, contacts, diagnoses, comorbidities }) {
  const age = calculateAge(patient.date_of_birth);
  const phoneContact = findContact(contacts, 'phone');
  const emailContact = findContact(contacts, 'email');
  const addressContact = findContact(contacts, 'address');
  const emergencyContact = findContact(contacts, 'emergency');
  const emergencyParts = emergencyContact?.value?.split('|') || [];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Panel title="Demographics">
        <InfoList
          rows={[
            ['Full name', patient.full_name],
            ['Preferred name', patient.preferred_name],
            ['MRN', patient.mrn],
            ['National ID', patient.national_id],
            ['Date of birth', formatDate(patient.date_of_birth)],
            ['Age', age === 'N/A' ? 'N/A' : `${age} years`],
            ['Sex', humanize(patient.sex)],
            ['Blood type', patient.blood_type || 'Unknown'],
            ['Marital status', humanize(patient.marital_status)],
            ['Nationality', patient.nationality],
            ['Occupation', patient.occupation],
          ]}
        />
      </Panel>

      <Panel title="Contacts And Alerts">
        <InfoList
          rows={[
            ['Phone', phoneContact?.value || 'Not provided'],
            ['Email', emailContact?.value || 'Not provided'],
            ['Address', addressContact?.value || 'Not provided'],
            ['Emergency contact', emergencyParts[0] || 'Not provided'],
            ['Emergency phone', emergencyParts[1] || 'Not provided'],
            ['SMS readiness', phoneContact?.value ? 'Ready' : 'Missing phone contact'],
          ]}
        />
      </Panel>

      <Panel title="Medical Tracking">
        <InfoList
          rows={[
            ['Renal course', patient.renal_course],
            ['Kidney disease cause', patient.kidney_disease_cause || firstValue(diagnoses, 'description')],
            ['Dialysis indication', patient.dialysis_indication || firstValue(diagnoses, 'notes')],
            ['Symptom onset', formatDate(patient.symptom_onset_date)],
            ['Diagnosis date', formatDate(patient.diagnosis_date)],
            ['Dialysis start', formatDate(patient.dialysis_start_date)],
            ['Registration date', formatDate(patient.registration_date)],
            ['Status', patient.is_active ? 'Active' : 'Inactive'],
          ]}
        />
        <div className="mt-4 pt-4 border-t border-gray-200">
          <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Known conditions</p>
          <div className="flex flex-wrap gap-2">
            {comorbidities.length > 0 ? comorbidities.map(condition => (
              <StatusPill key={condition.id || condition.condition} tone="gray">
                {condition.condition}
              </StatusPill>
            )) : <span className="text-sm text-gray-600">No conditions recorded.</span>}
          </div>
        </div>
      </Panel>
    </div>
  );
}

function SessionsTab({ sessions, vitals, dialysate }) {
  return (
    <Panel title={`Dialysis Sessions (${sessions.length})`}>
      {sessions.length === 0 ? (
        <p className="text-gray-600">No sessions recorded yet.</p>
      ) : (
        <div className="space-y-4">
          {sessions.map(session => {
            const sessionVitals = vitals.filter(v => idMatches(v.session_id, session.id));
            const latestVitals = sessionVitals[0];
            const sessionDialysate = dialysate.find(d => idMatches(d.session_id, session.id));
            const status = session.status || session.session_status || 'scheduled';

            return (
              <div key={session.id} className="border border-gray-200 rounded-lg p-4 bg-white">
                <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                  <div>
                    <p className="font-semibold text-gray-900">
                      {formatDate(session.scheduled_date)}
                    </p>
                    <p className="text-sm text-gray-600 mt-1">
                      {session.scheduled_start_time || timeFromIso(session.actual_start_time) || 'Time not set'} | {session.prescribed_duration_mins || session.actual_duration_mins || 240} mins | {humanize(session.modality || 'hd')}
                    </p>
                  </div>
                  <StatusPill tone={statusTone(status)}>{humanize(status)}</StatusPill>
                </div>

                <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-3 text-sm">
                  <MiniReading label="Latest BP" value={latestVitals ? `${latestVitals.bp_systolic || '--'}/${latestVitals.bp_diastolic || '--'}` : 'Not recorded'} />
                  <MiniReading label="Blood flow" value={latestVitals?.blood_flow_actual ? `${latestVitals.blood_flow_actual} ml/min` : 'Not recorded'} />
                  <MiniReading label="Dialysate Mg" value={sessionDialysate?.magnesium_meq_l ? `${sessionDialysate.magnesium_meq_l} mEq/L` : 'Not recorded'} />
                  <MiniReading label="Dialysate Na" value={sessionDialysate?.sodium_meq_l ? `${sessionDialysate.sodium_meq_l} mEq/L` : 'Not recorded'} />
                  <MiniReading label="Dialysate K" value={sessionDialysate?.potassium_meq_l ? `${sessionDialysate.potassium_meq_l} mEq/L` : 'Not recorded'} />
                  <MiniReading label="Composition check" value={sessionDialysate?.composition_verified ? 'Verified' : 'Needs verification'} tone={sessionDialysate?.composition_verified ? 'green' : 'yellow'} />
                </div>
              </div>
            );
          })}
        </div>
      )}
    </Panel>
  );
}

function LabsTab({ labs }) {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {LAB_TRACKING_GROUPS.map(group => (
          <LabTrendCard key={group.id} group={group} labs={labs} />
        ))}
      </div>

      <Panel title="All Lab Values By Clinical State">
        {labs.length === 0 ? (
          <p className="text-sm text-gray-600">No lab results recorded yet.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
              <thead>
                <tr className="text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
                  <th className="py-3 pr-4">Date</th>
                  <th className="py-3 pr-4">State</th>
                  <th className="py-3 pr-4">Test</th>
                  <th className="py-3 pr-4">Value</th>
                  <th className="py-3 pr-4">Status</th>
                  <th className="py-3 pr-4">Notes</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {labs.map(result => {
                  const tone = result.is_critical ? 'red' : result.is_abnormal ? 'yellow' : 'green';
                  return (
                    <tr key={result.id}>
                      <td className="py-3 pr-4 text-gray-700">{formatDate(result.result_date || result.created_at)}</td>
                      <td className="py-3 pr-4">
                        <StatusPill tone="blue">{PHASE_LABELS[result.result_phase] || humanize(result.result_phase || 'routine')}</StatusPill>
                      </td>
                      <td className="py-3 pr-4 font-medium text-gray-900">{result.test_name}</td>
                      <td className="py-3 pr-4 text-gray-900">{formatLabValue(result)}</td>
                      <td className="py-3 pr-4"><StatusPill tone={tone}>{result.is_critical ? 'Critical' : result.is_abnormal ? 'Abnormal' : 'In range'}</StatusPill></td>
                      <td className="py-3 pr-4 text-gray-600">{result.notes || 'No notes'}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </Panel>
    </div>
  );
}

function MedicationsTab({ prescriptions, prescriptionItems, medications }) {
  const medicationById = new Map(medications.map(med => [String(med.id), med]));

  return (
    <Panel title="Medication And Regimen Tracking">
      {prescriptions.length === 0 && prescriptionItems.length === 0 ? (
        <p className="text-sm text-gray-600">No active medication history recorded yet.</p>
      ) : (
        <div className="space-y-4">
          {prescriptions.map(rx => {
            const items = prescriptionItems.filter(item => idMatches(item.prescription_id, rx.id));
            return (
              <div key={rx.id} className="border border-gray-200 rounded-lg p-4">
                <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
                  <div>
                    <p className="font-semibold text-gray-900">{humanize(rx.status || 'active')} prescription</p>
                    <p className="text-sm text-gray-600">Started {formatDate(rx.prescribed_at || rx.created_at)}</p>
                  </div>
                  <StatusPill tone={rx.status === 'active' ? 'green' : 'gray'}>{humanize(rx.status || 'active')}</StatusPill>
                </div>
                {rx.regimen_notes && <p className="mt-3 text-sm text-gray-700">{rx.regimen_notes}</p>}
                <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-3">
                  {items.map(item => {
                    const med = medicationById.get(String(item.medication_id));
                    return (
                      <div key={item.id} className="rounded-lg border border-gray-200 bg-gray-50 p-3">
                        <p className="font-medium text-gray-900">{med?.generic_name || item.medication_name || 'Medication'}</p>
                        <p className="text-xs text-gray-600 mt-1">
                          {item.dose || 'Dose not recorded'} | {item.frequency || 'Frequency not recorded'} | {item.route || 'Route not recorded'}
                        </p>
                        <p className="text-xs text-gray-700 mt-2">
                          Change reason: {item.change_reason || 'Not recorded'}
                        </p>
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </Panel>
  );
}

function AccessTab({ accessData, assessments, lifecycle = [], infections = [] }) {
  return (
    <div className="space-y-6">
      <Panel title="Vascular Access And CVC Registry">
        {accessData.length === 0 ? (
          <p className="text-gray-600">No vascular access records found.</p>
        ) : (
          <div className="space-y-4">
            {accessData.map(access => {
              const accessAssessments = assessments.filter(a => idMatches(a.access_id, access.id));
              const latest = accessAssessments[0];
              return (
                <div key={access.id} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                    <div>
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="text-lg font-semibold text-gray-900">
                          {humanize(access.access_type)}
                        </span>
                        <StatusPill tone={access.status === 'active' ? 'green' : 'gray'}>{humanize(access.status)}</StatusPill>
                        {access.is_primary_access && <StatusPill tone="blue">Primary</StatusPill>}
                      </div>
                      <p className="text-sm text-gray-600 mt-1">
                        {humanize(access.access_site)} | {humanize(access.site_side)} | inserted {formatDate(access.insertion_date)}
                      </p>
                    </div>
                    <StatusPill tone={latest?.requires_intervention ? 'red' : 'green'}>
                      {latest?.requires_intervention ? 'Intervention needed' : 'No intervention flag'}
                    </StatusPill>
                  </div>

                  <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-3 text-sm">
                    <MiniReading label="Inserted by" value={access.inserted_by || 'Not recorded'} />
                    <MiniReading label="Catheter type" value={access.catheter_type || 'N/A'} />
                    <MiniReading label="First use" value={formatDate(access.first_use_date)} />
                    <MiniReading label="Thrill / bruit" value={latest ? `${yesNo(latest.has_thrill)} / ${yesNo(latest.has_bruit)}` : 'Not assessed'} />
                    <MiniReading label="Flow" value={latest?.flow_rate_ml_min ? `${latest.flow_rate_ml_min} ml/min` : 'Not recorded'} />
                    <MiniReading label="Recirculation" value={latest?.recirculation_percent ? `${latest.recirculation_percent}%` : 'Not recorded'} />
                  </div>

                  {latest && (
                    <div className="mt-4 rounded-lg bg-gray-50 border border-gray-200 p-3">
                      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Latest access assessment</p>
                      <div className="flex flex-wrap gap-2">
                        <StatusPill tone={latest.has_redness ? 'red' : 'green'}>Redness: {yesNo(latest.has_redness)}</StatusPill>
                        <StatusPill tone={latest.has_swelling ? 'red' : 'green'}>Swelling: {yesNo(latest.has_swelling)}</StatusPill>
                        <StatusPill tone={latest.has_discharge ? 'red' : 'green'}>Discharge: {yesNo(latest.has_discharge)}</StatusPill>
                        <StatusPill tone={latest.has_pain ? 'yellow' : 'green'}>Pain: {yesNo(latest.has_pain)}</StatusPill>
                      </div>
                      <p className="text-sm text-gray-700 mt-3">{latest.notes || 'No notes recorded.'}</p>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </Panel>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <Panel title="CVC / Access Lifecycle Events">
          {lifecycle.length === 0 ? (
            <p className="text-sm text-gray-600">No insertion, replacement, blockage, complication, or catheter-free plan events recorded.</p>
          ) : (
            <div className="space-y-3">
              {lifecycle.map(event => (
                <div key={event.id} className="rounded-lg border border-gray-200 bg-white p-4">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-semibold text-gray-900">{humanize(event.event_type)}</p>
                      <p className="text-sm text-gray-600">{formatDate(event.event_date)} | {event.operator_name || 'Operator not recorded'}</p>
                    </div>
                    <StatusPill tone={event.culture_result || includesText(event.long_term_complications, 'infection') ? 'red' : 'blue'}>
                      {event.culture_result ? 'Culture watch' : 'Lifecycle'}
                    </StatusPill>
                  </div>
                  <div className="mt-3 grid grid-cols-1 md:grid-cols-3 gap-2">
                    <MiniReading label="Attempts / pricks" value={event.insertion_attempts || 'Not recorded'} />
                    <MiniReading label="Immediate complications" value={arrayText(event.immediate_complications) || 'None recorded'} />
                    <MiniReading label="Long-term complications" value={arrayText(event.long_term_complications) || 'None recorded'} />
                    <MiniReading label="Exit site" value={event.exit_site_condition || 'Not recorded'} />
                    <MiniReading label="Lock solution" value={event.lock_solution || 'Not recorded'} />
                    <MiniReading label="Catheter-free plan" value={event.catheter_free_plan || 'Pending'} />
                  </div>
                </div>
              ))}
            </div>
          )}
        </Panel>

        <Panel title="Infection Surveillance">
          {infections.length === 0 ? (
            <p className="text-sm text-gray-600">No dialysis-event infection surveillance records for this patient.</p>
          ) : (
            <div className="space-y-3">
              {infections.map(event => (
                <div key={event.id} className="rounded-lg border border-gray-200 bg-white p-4">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-semibold text-gray-900">{humanize(event.event_type)}</p>
                      <p className="text-sm text-gray-600">{formatDate(event.event_date)} | source: {event.suspected_source || 'not recorded'}</p>
                    </div>
                    <StatusPill tone={event.reported_to_registry ? 'green' : 'red'}>
                      {event.reported_to_registry ? 'Registry sent' : 'Registry pending'}
                    </StatusPill>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    <StatusPill tone={event.iv_antimicrobial_started ? 'yellow' : 'gray'}>IV antimicrobial: {yesNo(event.iv_antimicrobial_started)}</StatusPill>
                    <StatusPill tone={event.positive_blood_culture ? 'red' : 'gray'}>Blood culture: {yesNo(event.positive_blood_culture)}</StatusPill>
                    <StatusPill tone={event.access_pus_redness_swelling ? 'red' : 'gray'}>Pus/redness/swelling: {yesNo(event.access_pus_redness_swelling)}</StatusPill>
                  </div>
                  <p className="text-sm text-gray-700 mt-3">Organism: {event.organism || 'not recorded'} | {event.notes || 'No notes'}</p>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>
    </div>
  );
}

function LabTrendCard({ group, labs, compact = false }) {
  const points = trendForGroup(labs, group);
  const latest = points[points.length - 1];
  const tone = labTone(latest, group);
  const value = latest ? `${latest.value} ${latest.unit || group.unit || ''}`.trim() : 'No data';

  return (
    <div className={`rounded-lg border p-4 ${toneClasses(tone).soft}`}>
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-gray-900">{group.label}</p>
          <p className={`mt-1 ${compact ? 'text-xl' : 'text-2xl'} font-bold ${toneClasses(tone).text}`}>{value}</p>
          <p className="text-xs text-gray-600 mt-1">
            {latest ? `${PHASE_LABELS[latest.result_phase] || humanize(latest.result_phase || 'routine')} | ${formatDate(latest.result_date)}` : 'Awaiting first result'}
          </p>
        </div>
        <StatusPill tone={tone}>{labStatus(latest, group)}</StatusPill>
      </div>
      <LabSparkline points={points} group={group} />
    </div>
  );
}

function LabSparkline({ points, group }) {
  if (points.length === 0) {
    return <div className="mt-4 h-12 rounded bg-white/70 border border-dashed border-gray-300"></div>;
  }

  const values = points.map(point => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = Math.max(max - min, 1);

  return (
    <div className="mt-4 h-16 flex items-end gap-1">
      {points.map(point => {
        const height = 22 + ((point.value - min) / span) * 38;
        const tone = labTone(point, group);
        return (
          <div
            key={`${point.id}-${point.result_phase}`}
            title={`${formatDate(point.result_date)}: ${point.value} ${point.unit || group.unit || ''}`}
            className={`flex-1 min-w-[8px] rounded-t ${toneClasses(tone).bar}`}
            style={{ height: `${height}px` }}
          ></div>
        );
      })}
    </div>
  );
}

function Panel({ title, children, className = '' }) {
  return (
    <section className={`bg-white rounded-lg border border-gray-200 p-5 ${className}`}>
      <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
      {children}
    </section>
  );
}

function MetricCard({ label, value, detail, tone = 'gray' }) {
  return (
    <div className={`rounded-lg border p-4 ${toneClasses(tone).soft}`}>
      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">{label}</p>
      <p className={`mt-2 text-xl font-bold ${toneClasses(tone).text}`}>{value || 'Not recorded'}</p>
      <p className="text-sm text-gray-700 mt-1">{detail || 'No detail recorded'}</p>
    </div>
  );
}

function SafetyGate({ gate }) {
  return (
    <div className={`rounded-lg border p-3 ${toneClasses(gate.tone).soft}`}>
      <div className="flex items-start justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-gray-900">{gate.label}</p>
          <p className="text-xs text-gray-700 mt-1">{gate.detail}</p>
        </div>
        <StatusPill tone={gate.tone}>{gate.status}</StatusPill>
      </div>
    </div>
  );
}

function StatusPill({ tone = 'gray', children }) {
  const classes = toneClasses(tone);
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold ${classes.pill}`}>
      {children}
    </span>
  );
}

function MiniReading({ label, value, tone = 'gray' }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-gray-50 p-3">
      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">{label}</p>
      <p className={`mt-1 font-medium ${toneClasses(tone).text}`}>{value || 'Not recorded'}</p>
    </div>
  );
}

function InfoList({ rows }) {
  return (
    <dl className="space-y-3">
      {rows.filter(([, value]) => value !== undefined && value !== null && value !== '').map(([label, value]) => (
        <div key={label}>
          <dt className="text-sm text-gray-600">{label}</dt>
          <dd className="text-base font-medium text-gray-900">{value}</dd>
        </div>
      ))}
    </dl>
  );
}

function TimelineItem({ label, value }) {
  return (
    <div className="pb-3">
      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">{label}</p>
      <p className="text-sm font-medium text-gray-900 mt-1">{value || 'Not recorded'}</p>
    </div>
  );
}

function buildPatientHistoryPrintBody({
  patient,
  snapshot,
  contacts,
  sessions,
  labs,
  diagnoses,
  comorbidities,
  accessData,
  prescriptionItems,
  medications,
}) {
  const medicationById = new Map(medications.map(med => [String(med.id), med]));
  const phoneContact = findContact(contacts, 'phone');
  const emergencyContact = findContact(contacts, 'emergency');
  const emergencyParts = emergencyContact?.value?.split('|') || [];

  const demographics = keyValueTable([
    ['Full name', patient.full_name],
    ['MRN', patient.mrn],
    ['National ID', patient.national_id],
    ['Date of birth', formatDate(patient.date_of_birth)],
    ['Age', calculateAge(patient.date_of_birth)],
    ['Sex', humanize(patient.sex)],
    ['Phone', phoneContact?.value || 'Not provided'],
    ['Emergency contact', emergencyParts[0] || 'Not provided'],
    ['Emergency phone', emergencyParts[1] || 'Not provided'],
  ]);

  const clinicalSummary = keyValueTable([
    ['Renal course', snapshot?.renalCourse || patient.renal_course],
    ['Dialysis exposure', snapshot?.dialysisVintage],
    ['Kidney disease cause', patient.kidney_disease_cause || firstValue(diagnoses, 'description')],
    ['Dialysis indication', patient.dialysis_indication || firstValue(diagnoses, 'notes')],
    ['Dialysis start', formatDate(patient.dialysis_start_date)],
    ['Primary access', snapshot?.primaryAccessLabel],
    ['Latest BP', snapshot?.latestBp],
    ['Hepatitis isolation', patient.hepatitis_isolation_status],
    ['Recovery plan', patient.recovery_plan || patient.residual_kidney_function],
  ]);

  const diagnosisRows = diagnoses.map(item => [
    formatDate(item.diagnosed_at || item.created_at),
    item.icd10_code || 'N/A',
    item.description || item.diagnosis || humanize(item.diagnosis_type),
    item.notes || '',
  ]);

  const comorbidityRows = comorbidities.map(item => [
    formatDate(item.diagnosed_at || item.created_at),
    item.condition || 'Condition',
    humanize(item.status || 'active'),
    item.notes || '',
  ]);

  const sessionRows = sessions.slice(0, 12).map(session => [
    formatDate(session.scheduled_date),
    session.scheduled_start_time || timeFromIso(session.actual_start_time) || 'N/A',
    humanize(session.status || session.session_status || 'scheduled'),
    humanize(session.modality || 'hd'),
    session.prescribed_duration_mins || session.actual_duration_mins || 'N/A',
  ]);

  const labRows = labs.slice(0, 30).map(result => [
    formatDate(result.result_date || result.created_at),
    PHASE_LABELS[result.result_phase] || humanize(result.result_phase || 'routine'),
    result.test_name,
    formatLabValue(result),
    result.is_critical ? 'Critical' : result.is_abnormal ? 'Abnormal' : 'In range',
    result.notes || '',
  ]);

  const accessRows = accessData.map(access => [
    humanize(access.access_type),
    humanize(access.access_site),
    humanize(access.site_side),
    formatDate(access.insertion_date),
    access.is_primary_access ? 'Yes' : 'No',
    humanize(access.status || access.access_status),
  ]);

  const medicationRows = prescriptionItems.map(item => {
    const med = medicationById.get(String(item.medication_id));
    return [
      med?.generic_name || item.medication_name || 'Medication',
      item.dose || 'N/A',
      item.frequency || 'N/A',
      item.route || 'N/A',
      item.change_reason || '',
    ];
  });

  return `
    <h2>Demographics</h2>
    ${demographics}
    <h2>Clinical Summary</h2>
    ${clinicalSummary}
    <h2>Diagnoses</h2>
    ${diagnosisRows.length ? htmlTable(['Date', 'ICD-10', 'Diagnosis', 'Notes'], diagnosisRows) : '<p>No diagnoses recorded.</p>'}
    <h2>Comorbidities</h2>
    ${comorbidityRows.length ? htmlTable(['Date', 'Condition', 'Status', 'Notes'], comorbidityRows) : '<p>No comorbidities recorded.</p>'}
    <h2>Recent Dialysis Sessions</h2>
    ${sessionRows.length ? htmlTable(['Date', 'Time', 'Status', 'Modality', 'Duration mins'], sessionRows) : '<p>No sessions recorded.</p>'}
    <h2>Recent Lab Results</h2>
    ${labRows.length ? htmlTable(['Date', 'Clinical state', 'Test', 'Value', 'Status', 'Notes'], labRows) : '<p>No lab results recorded.</p>'}
    <h2>Vascular Access</h2>
    ${accessRows.length ? htmlTable(['Type', 'Site', 'Side', 'Inserted', 'Primary', 'Status'], accessRows) : '<p>No access records recorded.</p>'}
    <h2>Medication History</h2>
    ${medicationRows.length ? htmlTable(['Medication', 'Dose', 'Frequency', 'Route', 'Change reason'], medicationRows) : '<p>No medication records recorded.</p>'}
    <div class="signatures">
      <div class="line">Prepared by</div>
      <div class="line">Reviewed by</div>
    </div>
  `;
}

function buildClinicalSnapshot({
  patient,
  sessions,
  contacts,
  access,
  assessments,
  labs,
  vitals,
  dialysate,
  diagnoses,
  comorbidities,
  consents,
  alerts,
}) {
  const latestSession = sessions[0];
  const latestVitals = vitals[0];
  const primaryAccess = access.find(item => item.is_primary_access) || access[0];
  const latestAssessment = assessments[0];
  const latestDialysate = dialysate[0];
  const dialysisStart = patient.dialysis_start_date || latestSession?.scheduled_date;
  const daysOnDialysis = daysSince(dialysisStart);
  const phone = findContact(contacts, 'phone');
  const latestConsent = consents.find(c => c.consent_type === 'dialysis_treatment') || consents[0];
  const potassium = latestForGroup(labs, LAB_TRACKING_GROUPS[0]);
  const hemoglobin = latestForGroup(labs, LAB_TRACKING_GROUPS[3]);
  const phosphate = latestForGroup(labs, LAB_TRACKING_GROUPS[1]);
  const glucose = latestForGroup(labs, LAB_TRACKING_GROUPS[5]);
  const activeAlerts = alerts.filter(alert => alert.acknowledged !== true);
  const renalCourse = patient.renal_course || inferRenalCourse(patient, diagnoses, daysOnDialysis);
  const isCvc = primaryAccess?.access_type?.includes('catheter');
  const consentActive = latestConsent?.status === 'given' && !isExpired(latestConsent?.expires_at);
  const hasPregnancy = Boolean(patient.pregnancy_related || comorbidities.some(item => item.condition?.toLowerCase().includes('pregnancy') || item.condition?.toLowerCase().includes('preeclampsia')));
  const hasMalaria = patient.malaria_aki_phenotype && patient.malaria_aki_phenotype !== 'not_applicable';

  const safetyGates = [
    {
      label: 'Latest potassium',
      status: labStatus(potassium, LAB_TRACKING_GROUPS[0]),
      detail: potassium ? `${potassium.value} ${potassium.unit || 'mmol/L'} on ${formatDate(potassium.result_date)}` : 'No potassium result found.',
      tone: labTone(potassium, LAB_TRACKING_GROUPS[0]),
    },
    {
      label: 'Blood pressure',
      status: latestVitals ? bpStatus(latestVitals) : 'Missing',
      detail: latestVitals ? `${latestVitals.bp_systolic || '--'}/${latestVitals.bp_diastolic || '--'} at ${formatDateTime(latestVitals.recorded_at)}` : 'No vitals attached to recent sessions.',
      tone: latestVitals ? bpTone(latestVitals) : 'yellow',
    },
    {
      label: 'Access assessment',
      status: latestAssessment?.requires_intervention ? 'Review' : latestAssessment ? 'Checked' : 'Missing',
      detail: latestAssessment ? accessAssessmentSummary(latestAssessment) : 'No thrill/bruit, CVC site, or infection marker check found.',
      tone: latestAssessment?.requires_intervention ? 'red' : latestAssessment ? 'green' : 'yellow',
    },
    {
      label: 'Dialysate composition',
      status: latestDialysate?.composition_verified ? 'Verified' : latestDialysate ? 'Pending' : 'Missing',
      detail: latestDialysate ? `Na ${latestDialysate.sodium_meq_l || '--'}, K ${latestDialysate.potassium_meq_l || '--'}, Mg ${latestDialysate.magnesium_meq_l || '--'} mEq/L` : 'No dialysate sodium, potassium, or magnesium record.',
      tone: latestDialysate?.composition_verified ? 'green' : 'yellow',
    },
    {
      label: 'Consent',
      status: consentActive ? 'Active' : 'Review',
      detail: latestConsent ? `${humanize(latestConsent.status)} consent signed ${formatDate(latestConsent.signed_at)}` : 'No dialysis consent record found.',
      tone: consentActive ? 'green' : 'red',
    },
    {
      label: 'Hepatitis isolation',
      status: patient.hepatitis_isolation_status ? humanize(patient.hepatitis_isolation_status) : 'Missing',
      detail: patient.hepatitis_isolation_status ? 'Isolation state recorded on patient profile.' : 'Record HBV/HCV status and machine/isolation requirement before session start.',
      tone: patient.hepatitis_isolation_status ? 'green' : 'yellow',
    },
  ];

  const riskFlags = [];
  if (activeAlerts.length > 0) {
    riskFlags.push({
      label: 'Unacknowledged clinical alert',
      status: `${activeAlerts.length} open`,
      detail: activeAlerts.map(alert => alert.test_name || alert.severity || 'Alert').join(', '),
      tone: 'red',
    });
  }
  if (isCvc) {
    riskFlags.push({
      label: 'CVC infection and replacement watch',
      status: latestAssessment?.requires_intervention ? 'Review' : 'Track',
      detail: 'Central venous catheters need visible insertion, site, complication, blockage, infection, and change history.',
      tone: latestAssessment?.requires_intervention ? 'red' : 'yellow',
    });
  }
  if (daysOnDialysis >= 90 && renalCourse.toLowerCase().includes('aki')) {
    riskFlags.push({
      label: 'AKI to CKD reassessment',
      status: 'Due',
      detail: `${daysOnDialysis} days on dialysis. Review recovery, CKD transition, or ESKD classification.`,
      tone: 'yellow',
    });
  }
  if (!phone?.value) {
    riskFlags.push({
      label: 'SMS follow-up blocked',
      status: 'Missing',
      detail: 'No patient phone number is available for session reminders or clinical callbacks.',
      tone: 'yellow',
    });
  }
  if (hasPregnancy) {
    riskFlags.push({
      label: 'Pregnancy-associated kidney risk',
      status: 'Track dates',
      detail: `ANC visits: ${patient.anc_visits || 'not recorded'}; ${patient.gravida || 'G?'} ${patient.para || 'P?'}.`,
      tone: 'blue',
    });
  }
  if (hasMalaria) {
    riskFlags.push({
      label: 'Malaria-associated AKI',
      status: humanize(patient.malaria_aki_phenotype),
      detail: 'Track symptom onset, first dose, diagnosis date, definitive therapy, hemoglobinuria, and renal recovery.',
      tone: 'blue',
    });
  }

  return {
    renalCourse,
    courseTone: renalCourse.toLowerCase().includes('aki') ? 'yellow' : 'blue',
    courseDetail: patient.kidney_disease_cause || firstValue(diagnoses, 'description') || 'Primary cause not recorded.',
    dialysisVintage: daysOnDialysis >= 0 ? `${daysOnDialysis} days` : 'Unknown',
    dialysisDecision: daysOnDialysis >= 90 ? 'Reassess recovery, CKD transition, and dialysis dose.' : 'Continue close AKI/CKD recovery tracking.',
    dialysisTone: daysOnDialysis >= 90 ? 'yellow' : 'green',
    primaryAccessLabel: primaryAccess ? humanize(primaryAccess.access_type) : 'No access',
    accessDetail: primaryAccess ? `${humanize(primaryAccess.access_site)} ${humanize(primaryAccess.site_side)} | ${formatDate(primaryAccess.insertion_date)}` : 'Create access record before treatment.',
    accessTone: latestAssessment?.requires_intervention ? 'red' : primaryAccess ? 'green' : 'yellow',
    latestBp: latestVitals ? `${latestVitals.bp_systolic || '--'}/${latestVitals.bp_diastolic || '--'}` : 'Missing',
    vitalsDetail: latestVitals ? formatDateTime(latestVitals.recorded_at) : 'No session vitals recorded.',
    vitalsTone: latestVitals ? bpTone(latestVitals) : 'yellow',
    smsReady: Boolean(phone?.value),
    safetyGates,
    riskFlags,
    potassium,
    hemoglobin,
    phosphate,
    glucose,
  };
}

function buildLabTimeline(results, labTests) {
  const testById = new Map(labTests.map(test => [String(test.id), test]));
  return results
    .map(result => {
      const test = testById.get(String(result.test_id));
      const value = parseNumeric(result.value_numeric ?? result.result_value ?? result.value);
      const testName = result.test_name || test?.name || result.name || result.test_code || 'Lab result';
      return {
        ...result,
        test_name: testName,
        test_code: result.test_code || test?.code || '',
        value,
        unit: result.unit || test?.default_unit || '',
        result_date: result.result_date || result.created_at || result.updated_at,
      };
    })
    .filter(result => result.value !== null || result.value_text)
    .sort((a, b) => toTime(b.result_date) - toTime(a.result_date));
}

function trendForGroup(labs, group) {
  return labs
    .filter(lab => labMatchesGroup(lab, group) && lab.value !== null)
    .slice(0, 6)
    .reverse();
}

function latestForGroup(labs, group) {
  return labs.find(lab => labMatchesGroup(lab, group) && lab.value !== null);
}

function labMatchesGroup(lab, group) {
  const text = `${lab.test_name || ''} ${lab.test_code || ''}`.toLowerCase();
  return group.keywords.some(keyword => text === keyword || text.includes(keyword));
}

function labTone(lab, group) {
  if (!lab) return 'yellow';
  if (lab.is_critical) return 'red';
  if (group.criticalHigh !== undefined && lab.value >= group.criticalHigh) return 'red';
  if (group.criticalLow !== undefined && lab.value <= group.criticalLow) return 'red';
  if (group.high !== undefined && lab.value > group.high) return 'yellow';
  if (group.low !== undefined && lab.value < group.low) return 'yellow';
  if (lab.is_abnormal) return 'yellow';
  return 'green';
}

function labStatus(lab, group) {
  const tone = labTone(lab, group);
  if (!lab) return 'Missing';
  if (tone === 'red') return 'Critical';
  if (tone === 'yellow') return 'Review';
  return 'In range';
}

function bpTone(vitals) {
  const systolic = Number(vitals.bp_systolic);
  const diastolic = Number(vitals.bp_diastolic);
  if (vitals.has_hypotension_alert || systolic < 90) return 'red';
  if (vitals.has_hypertension_alert || systolic >= 180 || diastolic >= 110) return 'red';
  if (systolic >= 160 || systolic < 100) return 'yellow';
  return 'green';
}

function bpStatus(vitals) {
  const tone = bpTone(vitals);
  if (tone === 'red') return 'Alert';
  if (tone === 'yellow') return 'Watch';
  return 'OK';
}

function accessAssessmentSummary(assessment) {
  const markers = [];
  if (assessment.has_redness) markers.push('redness');
  if (assessment.has_swelling) markers.push('swelling');
  if (assessment.has_discharge) markers.push('discharge');
  if (assessment.has_pain) markers.push('pain');
  if (assessment.requires_intervention) markers.push(assessment.intervention_type || 'intervention');
  return markers.length > 0 ? markers.join(', ') : 'No infection or intervention marker recorded.';
}

function inferRenalCourse(patient, diagnoses, daysOnDialysis) {
  const text = `${patient.kidney_disease_cause || ''} ${diagnoses.map(d => d.description).join(' ')}`.toLowerCase();
  if (text.includes('aki')) return daysOnDialysis >= 90 ? 'AKI on CKD review' : 'AKI requiring dialysis';
  if (text.includes('esrd') || text.includes('eskd')) return 'ESKD on maintenance hemodialysis';
  if (text.includes('ckd')) return 'CKD on hemodialysis';
  return 'Renal course not staged';
}

function idMatches(a, b) {
  return String(a) === String(b);
}

function findContact(contacts, type) {
  return contacts.find(contact => contact.contact_type === type);
}

function firstValue(rows, field) {
  return rows.find(row => row?.[field])?.[field];
}

function sortByDate(rows, field) {
  return [...rows].sort((a, b) => toTime(b[field] || b.created_at) - toTime(a[field] || a.created_at));
}

function toTime(value) {
  if (!value) return 0;
  const time = new Date(value).getTime();
  return Number.isNaN(time) ? 0 : time;
}

function calculateAge(dateOfBirth) {
  if (!dateOfBirth) return 'N/A';
  const birth = new Date(dateOfBirth);
  if (Number.isNaN(birth.getTime())) return 'N/A';
  return Math.floor((new Date() - birth) / 31557600000);
}

function daysSince(value) {
  if (!value) return -1;
  const time = new Date(value).getTime();
  if (Number.isNaN(time)) return -1;
  return Math.max(0, Math.floor((Date.now() - time) / 86400000));
}

function isExpired(value) {
  if (!value) return false;
  const time = new Date(value).getTime();
  return !Number.isNaN(time) && time < Date.now();
}

function parseNumeric(value) {
  if (value === null || value === undefined || value === '') return null;
  if (typeof value === 'number') return value;
  if (typeof value === 'string') {
    const parsed = Number(value);
    return Number.isNaN(parsed) ? null : parsed;
  }
  if (typeof value === 'object') {
    if (value.value !== undefined) return parseNumeric(value.value);
    if (value.String !== undefined) return parseNumeric(value.String);
    if (value.Int !== undefined && value.Exp !== undefined) {
      return Number(value.Int) * Math.pow(10, Number(value.Exp));
    }
  }
  return null;
}

function formatDate(value) {
  if (!value) return 'Not recorded';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleDateString();
}

function formatDateTime(value) {
  if (!value) return 'Not recorded';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString();
}

function timeFromIso(value) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function formatLabValue(result) {
  if (result.value_text) return result.value_text;
  if (result.value === null || result.value === undefined) return 'Not recorded';
  return `${result.value} ${result.unit || ''}`.trim();
}

function humanize(value) {
  if (!value) return 'Not recorded';
  return String(value)
    .replaceAll('_', ' ')
    .replaceAll('-', ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .replace(/\b\w/g, char => char.toUpperCase());
}

function yesNo(value) {
  if (value === null || value === undefined) return 'N/A';
  return value ? 'Yes' : 'No';
}

function arrayText(value) {
  if (!value) return '';
  if (Array.isArray(value)) return value.map(humanize).join(', ');
  if (typeof value === 'string') return humanize(value);
  if (typeof value === 'object') return Object.values(value).map(humanize).join(', ');
  return String(value);
}

function includesText(value, needle) {
  return arrayText(value).toLowerCase().includes(String(needle).toLowerCase());
}

function statusTone(status) {
  const normalized = String(status || '').toLowerCase();
  if (['completed', 'reviewed'].includes(normalized)) return 'green';
  if (['in_progress', 'in-progress', 'priming', 'pre_assessment'].includes(normalized)) return 'blue';
  if (['aborted', 'failed'].includes(normalized)) return 'red';
  if (['paused', 'scheduled', 'checked_in'].includes(normalized)) return 'yellow';
  return 'gray';
}

function toneClasses(tone) {
  const classes = {
    green: {
      soft: 'bg-emerald-50 border-emerald-200',
      text: 'text-emerald-800',
      pill: 'bg-emerald-100 text-emerald-800',
      bar: 'bg-emerald-500',
    },
    red: {
      soft: 'bg-red-50 border-red-200',
      text: 'text-red-800',
      pill: 'bg-red-100 text-red-800',
      bar: 'bg-red-500',
    },
    yellow: {
      soft: 'bg-amber-50 border-amber-200',
      text: 'text-amber-800',
      pill: 'bg-amber-100 text-amber-800',
      bar: 'bg-amber-500',
    },
    blue: {
      soft: 'bg-sky-50 border-sky-200',
      text: 'text-sky-800',
      pill: 'bg-sky-100 text-sky-800',
      bar: 'bg-sky-500',
    },
    gray: {
      soft: 'bg-gray-50 border-gray-200',
      text: 'text-gray-800',
      pill: 'bg-gray-100 text-gray-800',
      bar: 'bg-gray-400',
    },
  };
  return classes[tone] || classes.gray;
}
