import db from '../db/schema';

/**
 * Demo Data Seeder for Investor Presentations
 * Creates realistic African patient data in IndexedDB
 * DO NOT use in production - Kiruddu will have real patients
 */

const DEMO_PATIENTS = [
  {
    full_name: 'Nakato Sarah',
    mrn: 'KRD001',
    national_id: 'CM90012345678',
    date_of_birth: '1978-03-15',
    sex: 'female',
    blood_type: 'O+',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Teacher',
    education_level: 'tertiary',
  },
  {
    full_name: 'Okello James',
    mrn: 'KRD002',
    national_id: 'CM89012345679',
    date_of_birth: '1965-07-22',
    sex: 'male',
    blood_type: 'A+',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Protestant',
    occupation: 'Business Owner',
    education_level: 'secondary',
  },
  {
    full_name: 'Nambi Grace',
    mrn: 'KRD003',
    national_id: 'CM95012345680',
    date_of_birth: '1985-11-08',
    sex: 'female',
    blood_type: 'B+',
    marital_status: 'single',
    nationality: 'Ugandan',
    religion: 'Muslim',
    occupation: 'Nurse',
    education_level: 'tertiary',
  },
  {
    full_name: 'Musoke Daniel',
    mrn: 'KRD004',
    national_id: 'CM82012345681',
    date_of_birth: '1972-05-30',
    sex: 'male',
    blood_type: 'AB+',
    marital_status: 'widowed',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Accountant',
    education_level: 'postgraduate',
  },
  {
    full_name: 'Auma Betty',
    mrn: 'KRD005',
    national_id: 'CM92012345682',
    date_of_birth: '1980-09-12',
    sex: 'female',
    blood_type: 'O-',
    marital_status: 'divorced',
    nationality: 'Ugandan',
    religion: 'Protestant',
    occupation: 'Social Worker',
    education_level: 'tertiary',
  },
  {
    full_name: 'Ssemakula Patrick',
    mrn: 'KRD006',
    national_id: 'CM88012345683',
    date_of_birth: '1968-01-25',
    sex: 'male',
    blood_type: 'A-',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Engineer',
    education_level: 'tertiary',
  },
];

const DEMO_STAFF = [
  {
    full_name: 'Dr. Kiwanuka Robert',
    email: 'r.kiwanuka@kiruddu.go.ug',
    phone_number: '+256701234567',
    cadre: 'nephrologist',
    specialty: 'Nephrology',
    license_number: 'UMC-12345',
    is_active: true,
  },
  {
    full_name: 'Sr. Nalongo Mary',
    email: 'm.nalongo@kiruddu.go.ug',
    phone_number: '+256702345678',
    cadre: 'nurse',
    specialty: 'Dialysis Nursing',
    license_number: 'UNMC-23456',
    is_active: true,
  },
  {
    full_name: 'Sr. Akello Christine',
    email: 'c.akello@kiruddu.go.ug',
    phone_number: '+256703456789',
    cadre: 'nurse',
    specialty: 'Critical Care',
    license_number: 'UNMC-23457',
    is_active: true,
  },
];

const DEMO_MACHINES = [
  {
    machine_number: 'HD-001',
    manufacturer: 'Fresenius',
    model: '5008S',
    serial_number: 'FMC-5008-2023-001',
    operational_status: 'operational',
    location: 'Dialysis Unit A',
  },
  {
    machine_number: 'HD-002',
    manufacturer: 'Fresenius',
    model: '5008S',
    serial_number: 'FMC-5008-2023-002',
    operational_status: 'operational',
    location: 'Dialysis Unit A',
  },
  {
    machine_number: 'HD-003',
    manufacturer: 'Gambro',
    model: 'AK 200',
    serial_number: 'GMB-AK200-2022-003',
    operational_status: 'operational',
    location: 'Dialysis Unit B',
  },
];

const LAB_TESTS = [
  { code: 'K', name: 'Potassium', unit: 'mmol/L', low: 3.5, high: 5.5 },
  { code: 'PHOS', name: 'Phosphate', unit: 'mmol/L', low: 0.8, high: 1.5 },
  { code: 'ALB', name: 'Albumin / Protein', unit: 'g/L', low: 35, high: 50 },
  { code: 'HB', name: 'Hemoglobin', unit: 'g/dL', low: 10, high: 12 },
  { code: 'CR', name: 'Creatinine', unit: 'umol/L', low: 60, high: 110 },
  { code: 'UREA', name: 'Urea', unit: 'mmol/L', low: 2.5, high: 7.8 },
  { code: 'GLU', name: 'Glucose', unit: 'mmol/L', low: 4, high: 10 },
  { code: 'CA', name: 'Calcium', unit: 'mmol/L', low: 2.1, high: 2.6 },
];

const MEDICATIONS = [
  { id: 'med_amlodipine', generic_name: 'Amlodipine', drug_class: 'calcium channel blocker' },
  { id: 'med_erythropoietin', generic_name: 'Erythropoietin alfa', drug_class: 'erythropoiesis-stimulating agent' },
  { id: 'med_sevelamer', generic_name: 'Sevelamer carbonate', drug_class: 'phosphate binder' },
  { id: 'med_artesunate', generic_name: 'Artesunate', drug_class: 'antimalarial' },
  { id: 'med_labetalol', generic_name: 'Labetalol', drug_class: 'beta blocker' },
  { id: 'med_insulin', generic_name: 'Insulin soluble', drug_class: 'insulin' },
];

const CLINICAL_PROFILES = [
  {
    renalCourse: 'CKD on maintenance hemodialysis',
    cause: 'Hypertensive nephrosclerosis',
    indication: 'ESKD maintenance dialysis',
    malaria: 'not_applicable',
    pregnancy: false,
    comorbidities: ['Hypertension', 'Anemia of CKD'],
    meds: ['med_amlodipine', 'med_erythropoietin', 'med_sevelamer'],
  },
  {
    renalCourse: 'AKI requiring dialysis',
    cause: 'Severe malaria with blackwater fever',
    indication: 'Hyperkalemia and oliguric AKI',
    malaria: 'blackwater_fever',
    pregnancy: false,
    comorbidities: ['Severe malaria', 'Hemoglobinuria'],
    meds: ['med_artesunate', 'med_erythropoietin'],
  },
  {
    renalCourse: 'Pregnancy-associated AKI',
    cause: 'Antepartum preeclampsia with AKI',
    indication: 'Refractory hypertension and rising creatinine',
    malaria: 'not_applicable',
    pregnancy: true,
    gravida: 'G3',
    para: 'P2',
    ancVisits: 4,
    comorbidities: ['Preeclampsia', 'Pregnancy-associated AKI'],
    meds: ['med_labetalol', 'med_erythropoietin'],
  },
  {
    renalCourse: 'CKD on maintenance hemodialysis',
    cause: 'Diabetic kidney disease',
    indication: 'ESKD maintenance dialysis',
    malaria: 'not_applicable',
    pregnancy: false,
    comorbidities: ['Diabetes mellitus', 'Hypertension'],
    meds: ['med_insulin', 'med_amlodipine', 'med_sevelamer'],
  },
  {
    renalCourse: 'AKI on CKD watch',
    cause: 'Sepsis-associated AKI on suspected CKD',
    indication: 'Volume overload and uremic symptoms',
    malaria: 'non_blackwater_severe_malaria',
    pregnancy: false,
    comorbidities: ['Sepsis', 'Suspected CKD'],
    meds: ['med_artesunate', 'med_sevelamer'],
  },
  {
    renalCourse: 'CKD on maintenance hemodialysis',
    cause: 'Lupus nephritis',
    indication: 'ESKD maintenance dialysis',
    malaria: 'not_applicable',
    pregnancy: false,
    comorbidities: ['Systemic lupus erythematosus', 'Hypertension'],
    meds: ['med_amlodipine', 'med_erythropoietin', 'med_sevelamer'],
  },
];

function dateDaysAgo(days) {
  return new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
}

function timeDaysAgo(days) {
  return new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString();
}

function labValue(profileIndex, testCode, pointIndex) {
  const values = {
    K: [5.9, 5.4, 4.8, 4.6, 5.1, 4.9],
    PHOS: [1.9, 1.6, 1.3, 1.7, 2.0, 1.4],
    ALB: [32, 34, 38, 35, 30, 36],
    HB: [8.8, 9.5, 10.4, 9.8, 8.6, 10.2],
    CR: [880, 710, 620, 740, 690, 810],
    UREA: [28, 22, 18, 24, 20, 26],
    GLU: [5.8, 6.4, 5.1, 9.8, 7.2, 5.6],
    CA: [2.0, 2.2, 2.3, 2.1, 1.9, 2.2],
  };
  const base = values[testCode]?.[profileIndex % 6] || 0;
  const drift = pointIndex === 0 ? 0 : pointIndex === 1 ? -0.2 : pointIndex === 2 ? 0.15 : -0.1;
  return Number((base + drift).toFixed(testCode === 'CR' ? 0 : 1));
}

async function ensureClinicalDemoData(existingPatients) {
  const patients = existingPatients?.length ? existingPatients : await db.patients.toArray();
  if (patients.length === 0) return;

  const timestamp = new Date().toISOString();
  const hospitalId = patients[0].hospital_id || 'demo_hospital';
  const staff = await db.staff_profiles.toArray();
  const reviewerId = staff[0]?.id || 'staff_demo_reviewer';

  await db.lab_test_catalog.bulkPut(LAB_TESTS.map(test => ({
    id: `test_${test.code.toLowerCase()}`,
    hospital_id: hospitalId,
    code: test.code,
    name: test.name,
    category: 'dialysis_monitoring',
    default_unit: test.unit,
  })));

  await db.medications.bulkPut(MEDICATIONS.map(med => ({
    ...med,
    hospital_id: hospitalId,
    is_active: true,
    search_terms: [med.generic_name, med.drug_class],
    created_at: timestamp,
    updated_at: timestamp,
  })));

  for (let i = 0; i < patients.length; i++) {
    const patient = patients[i];
    const patientId = patient.id;
    const profile = CLINICAL_PROFILES[i % CLINICAL_PROFILES.length];
    const sessions = await db.dialysis_sessions.where('patient_id').equals(patientId).toArray();
    const accessRecords = await db.vascular_access.where('patient_id').equals(patientId).toArray();
    const primaryAccess = accessRecords.find(a => a.is_primary_access) || accessRecords[0];

    await db.patients.put({
      ...patient,
      renal_course: patient.renal_course || profile.renalCourse,
      kidney_disease_cause: patient.kidney_disease_cause || profile.cause,
      dialysis_indication: patient.dialysis_indication || profile.indication,
      symptom_onset_date: patient.symptom_onset_date || dateDaysAgo(120 + i * 12),
      diagnosis_date: patient.diagnosis_date || dateDaysAgo(100 + i * 9),
      dialysis_start_date: patient.dialysis_start_date || dateDaysAgo(92 + i * 18),
      malaria_aki_phenotype: patient.malaria_aki_phenotype || profile.malaria,
      pregnancy_related: patient.pregnancy_related ?? profile.pregnancy,
      gravida: patient.gravida || profile.gravida || '',
      para: patient.para || profile.para || '',
      anc_visits: patient.anc_visits ?? profile.ancVisits ?? '',
      tracking_notes: patient.tracking_notes || 'Monitor dialysis dose, access condition, 3P labs, hemoglobin, glucose, volume status, and regimen changes.',
      updated_at: timestamp,
    });

    await db.patient_clinical_profiles.put({
      id: `clinical_profile_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      renal_course: patient.renal_course || profile.renalCourse,
      kidney_disease_cause: patient.kidney_disease_cause || profile.cause,
      dialysis_indication: patient.dialysis_indication || profile.indication,
      symptom_onset_date: patient.symptom_onset_date || dateDaysAgo(120 + i * 12),
      diagnosis_date: patient.diagnosis_date || dateDaysAgo(100 + i * 9),
      dialysis_start_date: patient.dialysis_start_date || dateDaysAgo(92 + i * 18),
      aki_stage: profile.renalCourse.includes('AKI') ? 'KDIGO stage 3' : '',
      baseline_creatinine_mg_dl: profile.renalCourse.includes('AKI') ? 1.1 : null,
      highest_creatinine_mg_dl: profile.pregnancy ? 8.4 : profile.malaria !== 'not_applicable' ? 9.6 : 11.2,
      urine_output_ml_day: profile.renalCourse.includes('AKI') ? 420 : 120,
      residual_kidney_function: profile.renalCourse.includes('AKI') ? 'Daily urine output and creatinine recovery watch' : 'Maintenance dialysis residual function review',
      reversible_cause: profile.renalCourse.includes('AKI') ? profile.cause : '',
      recovery_plan: profile.renalCourse.includes('AKI') ? 'Review dialysis frequency weekly and reassess AKI-to-CKD status at 3 months.' : 'Monthly adequacy, access, anemia, MBD, and transplant/modality review.',
      aki_to_ckd_reassessment_due: profile.renalCourse.includes('AKI') ? dateDaysAgo(i === 1 ? -12 : 4) : '',
      malaria_aki_phenotype: patient.malaria_aki_phenotype || profile.malaria,
      malaria_symptom_onset_date: profile.malaria !== 'not_applicable' ? dateDaysAgo(18 + i) : '',
      fever_onset_date: profile.malaria !== 'not_applicable' ? dateDaysAgo(17 + i) : '',
      dark_urine_onset_date: profile.malaria === 'blackwater_fever' ? dateDaysAgo(16 + i) : '',
      malaria_test_date: profile.malaria !== 'not_applicable' ? dateDaysAgo(15 + i) : '',
      first_antimalarial_dose_at: profile.malaria !== 'not_applicable' ? timeDaysAgo(15 + i) : '',
      definitive_antimalarial_regimen: profile.malaria !== 'not_applicable' ? 'IV artesunate then oral ACT completion course' : '',
      parasite_clearance_date: profile.malaria !== 'not_applicable' ? dateDaysAgo(11 + i) : '',
      pregnancy_related: patient.pregnancy_related ?? profile.pregnancy,
      gestational_age_weeks: profile.pregnancy ? 31 : null,
      expected_delivery_date: profile.pregnancy ? dateDaysAgo(-46) : '',
      gravida: patient.gravida || profile.gravida || '',
      para: patient.para || profile.para || '',
      anc_visits: patient.anc_visits ?? profile.ancVisits ?? '',
      bp_proteinuria_summary: profile.pregnancy ? 'Severe-range BP with proteinuria before dialysis start' : '',
      preeclampsia_severity: profile.pregnancy ? 'severe_features' : '',
      magnesium_sulfate_given: profile.pregnancy ? true : null,
      pregnancy_antihypertensive_regimen: profile.pregnancy ? 'Labetalol titrated with obstetric review' : '',
      postpartum_renal_recovery: profile.pregnancy ? 'Postpartum renal recovery not yet confirmed' : '',
      transplant_referral_status: profile.renalCourse.includes('CKD') ? 'screening_needed' : 'not_yet_applicable',
      modality_plan: profile.renalCourse.includes('CKD') ? 'Maintenance HD with AV access planning and transplant discussion' : 'AKI recovery watch before chronic modality decision',
      conservative_care_flag: false,
      tracking_notes: patient.tracking_notes || 'Monitor dialysis dose, access condition, 3P labs, hemoglobin, glucose, volume status, and regimen changes.',
      created_by: reviewerId,
      updated_by: reviewerId,
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    await db.diagnoses.put({
      id: `diagnosis_primary_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      icd10_code: profile.renalCourse.includes('AKI') ? 'N17.9' : 'N18.6',
      description: profile.cause,
      diagnosis_type: 'primary',
      diagnosed_by: reviewerId,
      diagnosed_at: timeDaysAgo(100 + i * 9),
      notes: profile.indication,
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    for (let c = 0; c < profile.comorbidities.length; c++) {
      const condition = profile.comorbidities[c];
      await db.comorbidities.put({
        id: `comorbidity_${patientId}_${c}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        condition,
        status: condition.includes('Suspected') ? 'suspected' : 'active',
        diagnosed_at: dateDaysAgo(90 + c * 7),
        notes: 'Tracked for dialysis planning, medication review, and session risk flags.',
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
    }

    await db.consents.put({
      id: `consent_dialysis_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      consent_type: 'dialysis_treatment',
      status: i === 4 ? 'expired' : 'given',
      given_by: patient.full_name,
      relationship: 'self',
      signed_at: timeDaysAgo(i === 4 ? 420 : 30),
      expires_at: i === 4 ? timeDaysAgo(10) : timeDaysAgo(-335),
      recorded_by: reviewerId,
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    for (let t = 0; t < LAB_TESTS.length; t++) {
      const test = LAB_TESTS[t];
      for (let p = 0; p < 4; p++) {
        const phase = ['pre_dialysis', 'post_dialysis', 'routine', 'intra_complication'][p];
        const value = labValue(i, test.code, p);
        await db.lab_results.put({
          id: `lab_${patientId}_${test.code}_${p}`,
          patient_id: patientId,
          hospital_id: hospitalId,
          order_item_id: `order_item_${patientId}_${test.code}_${p}`,
          test_id: `test_${test.code.toLowerCase()}`,
          test_code: test.code,
          test_name: test.name,
          value_numeric: value,
          unit: test.unit,
          reference_range: `${test.low}-${test.high}`,
          result_phase: phase,
          status: 'verified',
          result_status: 'completed',
          is_abnormal: value < test.low || value > test.high,
          is_critical: test.code === 'K' && value >= 6,
          result_date: dateDaysAgo(p * 9 + t),
          result_time: '08:00:00',
          entered_by: reviewerId,
          notes: `${phase.replaceAll('_', ' ')} dialysis monitoring value`,
          synced: false,
          created_at: timestamp,
          updated_at: timestamp,
        });
      }
    }

    const latestPotassium = labValue(i, 'K', 0);
    if (latestPotassium >= 5.8) {
      await db.lab_critical_alerts.put({
        id: `alert_k_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        test_id: 'test_k',
        test_name: 'Potassium',
        critical_value: `${latestPotassium} mmol/L`,
        reference_range: '3.5-5.5 mmol/L',
        severity: latestPotassium >= 6 ? 'critical_high' : 'high_risk',
        acknowledged: false,
        created_at: timestamp,
        synced: false,
      });
    }

    await db.clinical_alerts.put({
      id: `clinical_alert_${patientId}_primary`,
      patient_id: patientId,
      hospital_id: hospitalId,
      session_id: sessions[0]?.id || null,
      alert_type: latestPotassium >= 5.8 ? 'hyperkalemia' : i === 4 ? 'consent_expired' : 'monthly_review',
      severity: latestPotassium >= 6 ? 'critical' : latestPotassium >= 5.8 || i === 4 ? 'high' : 'moderate',
      title: latestPotassium >= 5.8 ? 'Potassium safety review before dialysis' : i === 4 ? 'Dialysis consent needs renewal' : 'Monthly consultant review due',
      triggering_value: latestPotassium >= 5.8 ? `${latestPotassium} mmol/L` : '',
      threshold: latestPotassium >= 5.8 ? 'K >= 5.8 mmol/L local safety threshold' : '',
      source_table: latestPotassium >= 5.8 ? 'lab_results' : i === 4 ? 'consents' : 'clinical_review',
      patient_context: {
        renal_course: profile.renalCourse,
        cause: profile.cause,
        access_type: primaryAccess?.access_type || 'not_recorded',
      },
      suggested_action: latestPotassium >= 5.8 ? 'Confirm pre-dialysis potassium, ECG symptoms, dialysate potassium bath, and physician review.' : 'Complete review and document action.',
      governance_note: 'Clinician acknowledgement required; this alert does not prescribe autonomously.',
      status: 'open',
      created_by: reviewerId,
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    for (const session of sessions) {
      const blockedReasons = [];
      if (latestPotassium >= 5.8) blockedReasons.push('pre_dialysis_potassium_high');
      if (i === 4) blockedReasons.push('dialysis_consent_expired');
      if (primaryAccess && primaryAccess.access_type === 'tunneled_catheter' && i === 1) blockedReasons.push('cvc_infection_review_needed');

      await db.session_safety_checks.put({
        id: `safety_${session.id}`,
        session_id: session.id,
        patient_id: patientId,
        hospital_id: hospitalId,
        checked_by: reviewerId,
        checked_at: timeDaysAgo(1),
        check_status: blockedReasons.length ? 'blocked' : 'cleared',
        risk_score: blockedReasons.length * 35 + (profile.renalCourse.includes('AKI') ? 10 : 0),
        checklist: {
          potassium: latestPotassium < 5.8 ? 'cleared' : 'blocked',
          blood_pressure: i === 2 ? 'doctor_review' : 'cleared',
          access_assessment: primaryAccess ? 'present' : 'missing',
          machine_status: 'available',
          water_quality: i === 5 ? 'retest_pending' : 'passed',
          disinfection_log: 'present',
          hepatitis_isolation: 'checked',
          consent: i === 4 ? 'expired' : 'valid',
        },
        hard_stop_reasons: blockedReasons,
        source_summary: {
          potassium: latestPotassium,
          access_type: primaryAccess?.access_type || null,
          consent_status: i === 4 ? 'expired' : 'given',
        },
        override_required: blockedReasons.length > 0,
        notes: blockedReasons.length ? 'Safety gate requires documented clinician review before start.' : 'Pre-session gate cleared from local demo data.',
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.session_vitals.put({
        id: `vitals_${session.id}_pre`,
        session_id: session.id,
        patient_id: patientId,
        hospital_id: hospitalId,
        recorded_by: reviewerId,
        recorded_at: timeDaysAgo(1),
        time_on_dialysis_mins: 0,
        bp_systolic: 128 + i * 6,
        bp_diastolic: 78 + i * 3,
        heart_rate: 78 + i,
        temperature: Number((36.5 + i * 0.1).toFixed(1)),
        spo2: 97,
        blood_flow_actual: 300,
        dialysate_flow_actual: 500,
        has_hypotension_alert: false,
        has_hypertension_alert: i === 2,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.treatment_telemetry.put({
        id: `telemetry_${session.id}`,
        session_id: session.id,
        patient_id: patientId,
        hospital_id: hospitalId,
        recorded_by: reviewerId,
        recorded_at: timeDaysAgo(1),
        minutes_on_dialysis: session.status === 'completed' ? 240 : 0,
        blood_flow_actual: 285 + i * 5,
        dialysate_flow_actual: 500,
        tmp_mmhg: 92 + i * 4,
        venous_pressure_mmhg: 130 + i * 6,
        arterial_pressure_mmhg: -120 - i * 4,
        conductivity_ms_cm: 14.1,
        temperature_celsius: 36.5,
        blood_volume_processed_l: session.status === 'completed' ? 72 + i * 2 : 0,
        access_recirculation_percent: i === 4 ? 12 : 4,
        delivered_minutes: session.status === 'completed' ? (i === 4 ? 210 : 240) : 0,
        final_delivered_dose: session.status === 'completed' ? (i === 4 ? 'shortened' : 'full') : 'pending',
        early_termination_reason: i === 4 && session.status === 'completed' ? 'Access pain and high venous pressure' : '',
        alarms: i === 4 ? [{ type: 'venous_pressure', minute: 118, action: 'access review' }] : [],
        interruptions: i === 4 ? [{ reason: 'access pain', minutes: 12 }] : [],
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.dialysate_records.put({
        id: `dialysate_${session.id}`,
        session_id: session.id,
        patient_id: patientId,
        hospital_id: hospitalId,
        recorded_by: reviewerId,
        recorded_at: timeDaysAgo(1),
        sodium_meq_l: 138,
        potassium_meq_l: i === 1 ? 2 : 3,
        bicarbonate_meq_l: 35,
        calcium_meq_l: 1.25,
        magnesium_meq_l: i === 5 ? 0.75 : 0.5,
        glucose_mg_dl: i === 3 ? 200 : 100,
        conductivity_ms_cm: 14.1,
        temperature_celsius: 36.5,
        composition_verified: i !== 4,
        deviations_noted: i === 4 ? 'Composition awaiting second-person verification.' : '',
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
    }

    if (primaryAccess) {
      await db.vascular_access_assessments.put({
        id: `access_assessment_${patientId}_latest`,
        access_id: primaryAccess.id,
        patient_id: patientId,
        session_id: sessions[0]?.id,
        hospital_id: hospitalId,
        assessed_by: reviewerId,
        assessed_at: timeDaysAgo(1),
        has_thrill: primaryAccess.access_type !== 'tunneled_catheter',
        has_bruit: primaryAccess.access_type !== 'tunneled_catheter',
        has_redness: i === 1,
        has_swelling: i === 1,
        has_discharge: false,
        has_bleeding: false,
        has_pain: i === 1 || i === 4,
        appearance_normal: i !== 1 && i !== 4,
        flow_rate_ml_min: primaryAccess.access_type === 'tunneled_catheter' ? 260 : 360,
        venous_pressure_mmhg: 130 + i * 6,
        arterial_pressure_mmhg: -120 - i * 4,
        recirculation_percent: i === 4 ? 12 : 4,
        requires_intervention: i === 1 || i === 4,
        intervention_type: i === 1 ? 'CVC infection review' : i === 4 ? 'Access flow assessment' : '',
        intervention_urgency: i === 1 ? 'same_day' : i === 4 ? 'routine' : '',
        notes: i === 1 ? 'CVC site tenderness and swelling. Review for infection and catheter change decision.' : 'Routine access surveillance entry.',
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.access_lifecycle_events.put({
        id: `access_lifecycle_${patientId}`,
        patient_id: patientId,
        access_id: primaryAccess.id,
        hospital_id: hospitalId,
        event_type: primaryAccess.access_type === 'tunneled_catheter' ? 'cvc_review' : 'access_maturation_review',
        event_date: dateDaysAgo(1),
        operator_id: reviewerId,
        operator_name: 'Dr. Kiwanuka Robert',
        insertion_attempts: primaryAccess.access_type === 'tunneled_catheter' ? 2 : 1,
        ultrasound_used: primaryAccess.access_type === 'tunneled_catheter',
        side_site: `${primaryAccess.site_side || 'unknown'} ${primaryAccess.access_site || 'site'}`,
        catheter_length_cm: primaryAccess.access_type === 'tunneled_catheter' ? 15 : null,
        tip_position_confirmation: primaryAccess.access_type === 'tunneled_catheter' ? 'documented at insertion' : '',
        immediate_complications: i === 1 ? ['difficult_cannulation'] : [],
        long_term_complications: i === 1 ? ['suspected_exit_site_infection'] : i === 4 ? ['high_recirculation'] : [],
        exit_site_condition: i === 1 ? 'Tender with swelling' : 'Clean',
        tunnel_cuff_status: primaryAccess.access_type === 'tunneled_catheter' ? 'stable cuff' : '',
        lock_solution: primaryAccess.access_type === 'tunneled_catheter' ? 'heparin lock' : '',
        dressing_date: dateDaysAgo(1),
        hub_scrub_compliant: i !== 1,
        catheter_free_plan: primaryAccess.access_type === 'tunneled_catheter' ? 'Permanent access planning required' : 'Continue AV access surveillance',
        culture_result: i === 1 ? 'culture pending' : '',
        antibiotic_course: i === 1 ? 'review after cultures' : '',
        notes: primaryAccess.access_type === 'tunneled_catheter' ? 'CVC lifecycle visible for infection, blockage, replacement, and catheter-free planning.' : 'AV access life-plan review.',
        created_by: reviewerId,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      if (i === 1) {
        await db.infection_surveillance_events.put({
          id: `infection_event_${patientId}`,
          patient_id: patientId,
          session_id: sessions[0]?.id,
          access_id: primaryAccess.id,
          hospital_id: hospitalId,
          event_type: 'suspected_cvc_bsi',
          event_date: dateDaysAgo(1),
          iv_antimicrobial_started: true,
          positive_blood_culture: false,
          access_pus_redness_swelling: true,
          suspected_source: 'tunneled_cvc',
          organism: 'culture pending',
          culture_collected_at: timeDaysAgo(1),
          antimicrobial_started_at: timeDaysAgo(1),
          hospitalized: false,
          death_related: false,
          recurrence_window_notes: 'Track 21-day recurrence window after culture result.',
          reported_to_registry: false,
          reported_by: reviewerId,
          notes: 'Access tenderness and swelling captured as a reportable surveillance object.',
          synced: false,
          created_at: timestamp,
          updated_at: timestamp,
        });
      }
    }

    await db.adequacy_reviews.put({
      id: `adequacy_review_${patientId}`,
      patient_id: patientId,
      session_id: sessions[0]?.id,
      hospital_id: hospitalId,
      review_month: new Date().toISOString().slice(0, 7) + '-01',
      reviewed_by: reviewerId,
      pre_urea_mg_dl: 88 + i * 4,
      post_urea_mg_dl: 28 + i * 2,
      urr_percent: i === 4 ? 61 : 68 + i,
      sp_kt_v: i === 4 ? 1.12 : 1.28 + i * 0.02,
      target_kt_v: 1.4,
      residual_kidney_function: profile.renalCourse.includes('AKI') ? 'recovering urine output under review' : 'low residual function',
      urine_volume_ml_day: profile.renalCourse.includes('AKI') ? 420 : 120,
      missed_treatments: i === 4 ? 1 : 0,
      shortened_treatments: i === 4 ? 2 : 0,
      interdialytic_weight_gain_kg: 2.1 + i * 0.4,
      normalized_protein_catabolic_rate: 0.95,
      uf_rate_ml_kg_hr: i === 4 ? 13.5 : 8.5 + i * 0.4,
      adequacy_status: i === 4 ? 'doctor_review_required' : 'adequate',
      doctor_review_required: i === 4,
      recommendations: i === 4 ? 'Review shortened sessions, access flow, UF tolerance, and delivered dose.' : 'Continue monthly adequacy trend review.',
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    await db.prescriptions.put({
      id: `prescription_active_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      status: 'active',
      prescribed_at: timeDaysAgo(14),
      prescribed_by: reviewerId,
      regimen_notes: 'Medication list is tracked for later regimen change review and interaction checks.',
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    for (let m = 0; m < profile.meds.length; m++) {
      await db.prescription_items.put({
        id: `rx_item_${patientId}_${m}`,
        prescription_id: `prescription_active_${patientId}`,
        patient_id: patientId,
        medication_id: profile.meds[m],
        dose: m === 0 ? '5 mg' : 'standard renal dose',
        route: 'oral',
        frequency: m === 1 ? 'weekly' : 'daily',
        start_date: dateDaysAgo(14),
        changed_by: reviewerId,
        change_reason: m === 2 ? '3P control: phosphate binder review' : 'active regimen tracking',
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
    }

    await db.medication_reconciliation_reviews.put({
      id: `med_rec_${patientId}`,
      patient_id: patientId,
      session_id: sessions[0]?.id,
      hospital_id: hospitalId,
      review_type: profile.pregnancy ? 'pregnancy_dialysis_review' : 'monthly_dialysis_review',
      review_date: dateDaysAgo(2),
      reviewed_by: reviewerId,
      medications: profile.meds.map((medId, index) => ({
        medication_id: medId,
        dose: index === 0 ? '5 mg' : 'standard renal dose',
        timing: index === 1 ? 'after dialysis / weekly' : 'daily',
        changed_by: reviewerId,
        change_reason: index === 2 ? 'phosphate control' : 'active regimen',
      })),
      renal_dosing_flags: profile.meds.includes('med_artesunate') ? ['severe_malaria_aki_follow_up'] : [],
      dialysis_timing_flags: profile.meds.includes('med_erythropoietin') ? ['administer_on_dialysis_day'] : [],
      pregnancy_cautions: profile.pregnancy ? ['avoid_unreviewed_nephrotoxins', 'obstetric_joint_review'] : [],
      nephrotoxin_flags: profile.cause.includes('Drug') ? ['toxin_review_required'] : [],
      adherence_notes: 'Medication list reconciled against dialysis status and clinical pathway.',
      regimen_change_reason: 'Monthly dialysis medication review',
      recommendations: 'Document dose changes, changed-by, indication, and dialysis-day timing.',
      status: 'open',
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    await db.patient_reported_events.put({
      id: `patient_reported_${patientId}`,
      patient_id: patientId,
      session_id: sessions[0]?.id,
      hospital_id: hospitalId,
      event_type: i === 4 ? 'missed_session_follow_up' : 'routine_symptom_review',
      reported_at: timeDaysAgo(1),
      reported_by: reviewerId,
      symptoms: {
        cramps: i === 4,
        dizziness: i === 4,
        dyspnea: false,
        pruritus: i === 0,
        post_dialysis_fatigue: true,
        access_symptoms: i === 1 ? 'CVC tenderness' : '',
      },
      education_topics: ['fluid restriction', 'potassium warning signs', 'access care', 'medication adherence'],
      teach_back_completed: i !== 4,
      transport_reliability: i === 4 ? 'unreliable' : 'reliable',
      no_show_reason: i === 4 ? 'transport and payment barrier' : '',
      payment_barrier: i === 4,
      food_insecurity: i === 5,
      caregiver_support: i === 4 ? 'limited' : 'available',
      sms_whatsapp_consent: true,
      follow_up_channel: 'whatsapp',
      follow_up_due_at: timeDaysAgo(-2),
      notes: 'Patient-centered dialysis tracking with symptoms, education, transport, and follow-up.',
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    await db.interoperability_exports.put({
      id: `interop_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      export_type: 'fhir_readiness',
      fhir_resource_type: 'Bundle',
      local_table: 'patients',
      local_id: patientId,
      coding_system: 'ICD-10 / LOINC / UCUM mapping queue',
      coding_version: 'local-demo',
      code_value: profile.renalCourse.includes('AKI') ? 'N17.9' : 'N18.6',
      ucum_unit: 'mmol/L',
      fhir_payload: {
        patient: `Patient/${patientId}`,
        observations: ['potassium', 'phosphate', 'albumin', 'hemoglobin', 'glucose'],
        conditions: [profile.cause],
      },
      export_status: i === 5 ? 'missing_codes' : 'draft',
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });

    await db.ontology_relationships.put({
      id: `ontology_${patientId}`,
      patient_id: patientId,
      hospital_id: hospitalId,
      source_type: 'diagnosis',
      source_id: `diagnosis_primary_${patientId}`,
      relation_type: 'drives_dialysis_indication',
      target_type: 'dialysis_course',
      target_id: sessions[0]?.id,
      relation_context: {
        cause: profile.cause,
        renal_course: profile.renalCourse,
        note: 'Relational tables remain source of truth; graph edges support relationship queries.',
      },
      provenance: 'relational_source_of_truth',
      confidence: 0.92,
      neo4j_node_ref: `Patient:${patientId}`,
      is_active: true,
      created_by: reviewerId,
      synced: false,
      created_at: timestamp,
      updated_at: timestamp,
    });
  }

  await db.staff_attendance_verifications.bulkPut(staff.map((member, index) => ({
    id: `attendance_${member.id}`,
    staff_id: member.id,
    user_id: member.user_id || null,
    hospital_id: hospitalId,
    verification_date: new Date().toISOString().split('T')[0],
    verification_time: timeDaysAgo(0),
    biometric_method: 'fingerprint',
    biometric_device_id: 'demo-biometric-01',
    biometric_template_ref: `template_ref_${member.id}`,
    verification_result: index === 2 ? 'late_exception' : 'verified',
    station_assignment: index === 0 ? 'doctor_review' : 'dialysis_station_a',
    patient_assignment: patients.slice(index, index + 2).map(p => p.id),
    handover_accepted: index !== 2,
    competency_status: {
      catheter_care: index !== 2 ? 'current' : 'renewal_due',
      machine_competency: 'current',
      emergency_drill: index === 0 ? 'current' : 'due_this_quarter',
    },
    exception_reason: index === 2 ? 'Late arrival requires supervisor review' : '',
    synced: false,
    created_at: timestamp,
    updated_at: timestamp,
  })));

  await db.unit_safety_events.put({
    id: 'unit_safety_water_retest_demo',
    hospital_id: hospitalId,
    event_type: 'water_retest_pending',
    severity: 'high',
    event_date: new Date().toISOString().split('T')[0],
    affected_sessions: [],
    affected_patients: [],
    consumable_lots: [{ item: 'dialyzer', lot: 'DX-UG-2026-05', status: 'in_use' }],
    immediate_action: 'Keep retest queue visible before first shift sign-off.',
    root_cause: 'Demo retest event for water safety workflow',
    corrective_action: 'Document retest clearance before session start.',
    closure_due_date: dateDaysAgo(-1),
    reported_by: reviewerId,
    assigned_to: reviewerId,
    notes: 'Unit-level safety event links water, machine, consumables, outage, and QI workflows.',
    synced: false,
    created_at: timestamp,
    updated_at: timestamp,
  });
}

export async function seedDemoData() {
  console.log('🌱 Starting demo data seeding...');

  try {
    // Check if already seeded
    const existingPatients = await db.patients.toArray();
    if (existingPatients.length > 0) {
      await ensureClinicalDemoData(existingPatients);
      console.log('✅ Demo data already exists. Clinical tracking data checked.');
      return;
    }

    const hospitalId = 'demo_hospital';
    const timestamp = new Date().toISOString();

    // 1. Seed Staff (needed for doctor assignments)
    const staffIds = [];
    for (const staff of DEMO_STAFF) {
      const staffId = `staff_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      await db.staff_profiles.put({
        id: staffId,
        ...staff,
        hospital_id: hospitalId,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
      staffIds.push(staffId);
    }

    // 2. Seed Machines
    const machineIds = [];
    for (const machine of DEMO_MACHINES) {
      const machineId = `machine_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      await db.dialysis_machines.put({
        id: machineId,
        ...machine,
        hospital_id: hospitalId,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
      machineIds.push(machineId);
    }

    // 3. Seed Patients with full data
    for (let i = 0; i < DEMO_PATIENTS.length; i++) {
      const patient = DEMO_PATIENTS[i];
      const patientId = `patient_${Date.now()}_${i}`;

      // Create patient
      await db.patients.put({
        id: patientId,
        ...patient,
        hospital_id: hospitalId,
        registered_by: staffIds[0], // First doctor
        primary_doctor_id: staffIds[0],
        registration_date: new Date().toISOString().split('T')[0],
        primary_language: 'English',
        interpreter_needed: false,
        is_active: true,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add contacts
      await db.patient_contacts.put({
        id: `contact_phone_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'phone',
        value: `+25670${7000000 + i}`,
        label: 'Primary Phone',
        is_primary: true,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.patient_contacts.put({
        id: `contact_address_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'address',
        value: `Kampala, Nakawa Division, Plot ${100 + i}`,
        label: 'Home Address',
        is_primary: true,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add emergency contact
      await db.patient_contacts.put({
        id: `contact_emergency_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'emergency',
        value: `Emergency Contact ${i}|+25677${1000000 + i}|Spouse`,
        label: 'Emergency Contact',
        is_primary: false,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add vascular access
      const accessTypes = ['avf', 'avg', 'tunneled_catheter'];
      const accessType = accessTypes[i % 3];
      await db.vascular_access.put({
        id: `access_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        access_type: accessType,
        access_site: accessType === 'tunneled_catheter' ? 'internal_jugular' : 'forearm',
        site_side: i % 2 === 0 ? 'left' : 'right',
        insertion_date: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        first_use_date: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        status: 'active',
        is_primary_access: true,
        inserted_by: staffIds[0],
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add 2-3 recent sessions for each patient
      for (let j = 0; j < 3; j++) {
        const sessionDate = new Date(Date.now() - j * 2 * 24 * 60 * 60 * 1000);
        const sessionId = `session_${patientId}_${j}`;
        
        await db.dialysis_sessions.put({
          id: sessionId,
          patient_id: patientId,
          hospital_id: hospitalId,
          machine_id: machineIds[i % machineIds.length],
          scheduled_date: sessionDate.toISOString().split('T')[0],
          scheduled_start_time: '08:00',
          shift: 'morning',
          prescribed_duration_mins: 240,
          modality: 'hd',
          status: j === 0 ? 'scheduled' : 'completed',
          primary_nurse_id: staffIds[1], // First nurse
          supervising_doctor_id: staffIds[0], // Doctor
          was_patient_reviewed: j === 0 ? false : true,
          synced: false,
          created_at: timestamp,
          updated_at: timestamp,
        });
      }
    }

    console.log('✅ Demo data seeded successfully!');
    console.log(`   - ${DEMO_PATIENTS.length} patients`);
    console.log(`   - ${DEMO_STAFF.length} staff members`);
    console.log(`   - ${DEMO_MACHINES.length} machines`);
    console.log(`   - ${DEMO_PATIENTS.length * 3} contacts per patient`);
    console.log(`   - ${DEMO_PATIENTS.length} vascular access records`);
    console.log(`   - ${DEMO_PATIENTS.length * 3} dialysis sessions`);

    await ensureClinicalDemoData(await db.patients.toArray());
    console.log('   - Clinical tracking, lab trends, vitals, access surveillance, and medication histories');

  } catch (error) {
    console.error('❌ Error seeding demo data:', error);
    throw error;
  }
}

// Function to clear demo data (for re-seeding)
export async function clearDemoData() {
  console.log('🗑️  Clearing demo data...');
  await db.patients.clear();
  await db.patient_contacts.clear();
  await db.vascular_access.clear();
  await db.vascular_access_assessments.clear();
  await db.patient_clinical_profiles.clear();
  await db.session_safety_checks.clear();
  await db.clinical_alerts.clear();
  await db.treatment_telemetry.clear();
  await db.access_lifecycle_events.clear();
  await db.infection_surveillance_events.clear();
  await db.adequacy_reviews.clear();
  await db.medication_reconciliation_reviews.clear();
  await db.patient_reported_events.clear();
  await db.staff_attendance_verifications.clear();
  await db.unit_safety_events.clear();
  await db.interoperability_exports.clear();
  await db.ontology_relationships.clear();
  await db.dialysis_sessions.clear();
  await db.session_vitals.clear();
  await db.dialysate_records.clear();
  await db.lab_results.clear();
  await db.lab_test_catalog.clear();
  await db.lab_critical_alerts.clear();
  await db.diagnoses.clear();
  await db.comorbidities.clear();
  await db.consents.clear();
  await db.prescriptions.clear();
  await db.prescription_items.clear();
  await db.medications.clear();
  await db.staff_profiles.clear();
  await db.dialysis_machines.clear();
  console.log('✅ Demo data cleared');
}
