-- +goose Up
-- Consultant-grade dialysis tracking extensions.
-- These tables add the workflow objects that existing dialysis CRUD surfaces were missing:
-- safety gates, delivered treatment telemetry, CVC/access lifecycle, infection surveillance,
-- adequacy review, patient pathways, medication reconciliation, staff attendance, exports,
-- and ontology edges. Existing core tables remain the source of truth.

CREATE TABLE patient_clinical_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    renal_course VARCHAR(120),
    kidney_disease_cause VARCHAR(200),
    dialysis_indication TEXT,
    symptom_onset_date DATE,
    diagnosis_date DATE,
    dialysis_start_date DATE,
    aki_stage VARCHAR(30),
    baseline_creatinine_mg_dl NUMERIC(7,2),
    highest_creatinine_mg_dl NUMERIC(7,2),
    urine_output_ml_day NUMERIC(9,2),
    residual_kidney_function TEXT,
    reversible_cause TEXT,
    recovery_plan TEXT,
    stop_dialysis_trial_date DATE,
    aki_to_ckd_reassessment_due DATE,
    malaria_aki_phenotype VARCHAR(100),
    malaria_symptom_onset_date DATE,
    fever_onset_date DATE,
    dark_urine_onset_date DATE,
    malaria_test_date DATE,
    parasitemia VARCHAR(100),
    first_antimalarial_dose_at TIMESTAMPTZ,
    definitive_antimalarial_regimen TEXT,
    parasite_clearance_date DATE,
    pregnancy_related BOOLEAN NOT NULL DEFAULT FALSE,
    gestational_age_weeks NUMERIC(4,1),
    expected_delivery_date DATE,
    gravida VARCHAR(20),
    para VARCHAR(20),
    anc_visits INTEGER,
    bp_proteinuria_summary TEXT,
    preeclampsia_severity VARCHAR(100),
    delivery_date DATE,
    fetal_outcome TEXT,
    magnesium_sulfate_given BOOLEAN,
    pregnancy_antihypertensive_regimen TEXT,
    postpartum_renal_recovery TEXT,
    maternal_fetal_follow_up TEXT,
    transplant_referral_status VARCHAR(120),
    modality_plan VARCHAR(120),
    conservative_care_flag BOOLEAN NOT NULL DEFAULT FALSE,
    tracking_notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(patient_id)
);

CREATE TABLE session_safety_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    checked_by UUID REFERENCES users(id),
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    check_status VARCHAR(40) NOT NULL DEFAULT 'pending',
    risk_score INTEGER NOT NULL DEFAULT 0,
    checklist JSONB NOT NULL DEFAULT '{}'::jsonb,
    hard_stop_reasons JSONB NOT NULL DEFAULT '[]'::jsonb,
    source_summary JSONB NOT NULL DEFAULT '{}'::jsonb,
    override_required BOOLEAN NOT NULL DEFAULT FALSE,
    override_reason TEXT,
    override_approved_by UUID REFERENCES users(id),
    override_approved_at TIMESTAMPTZ,
    audit_log_id UUID REFERENCES audit_logs(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE clinical_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    alert_type VARCHAR(100) NOT NULL,
    severity VARCHAR(40) NOT NULL DEFAULT 'info',
    title VARCHAR(200) NOT NULL,
    triggering_value TEXT,
    threshold TEXT,
    source_table VARCHAR(100),
    source_id UUID,
    patient_context JSONB NOT NULL DEFAULT '{}'::jsonb,
    suggested_action TEXT,
    governance_note TEXT,
    status VARCHAR(40) NOT NULL DEFAULT 'open',
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMPTZ,
    override_reason TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE treatment_telemetry_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    recorded_by UUID REFERENCES users(id),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    minutes_on_dialysis INTEGER,
    blood_flow_actual INTEGER,
    dialysate_flow_actual INTEGER,
    tmp_mmhg NUMERIC(7,2),
    venous_pressure_mmhg INTEGER,
    arterial_pressure_mmhg INTEGER,
    conductivity_ms_cm NUMERIC(5,2),
    temperature_celsius NUMERIC(4,1),
    blood_volume_processed_l NUMERIC(7,2),
    access_recirculation_percent NUMERIC(5,2),
    delivered_minutes INTEGER,
    final_delivered_dose VARCHAR(100),
    early_termination_reason TEXT,
    alarms JSONB NOT NULL DEFAULT '[]'::jsonb,
    interruptions JSONB NOT NULL DEFAULT '[]'::jsonb,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE vascular_access_lifecycle_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    access_id UUID REFERENCES vascular_access(id),
    event_type VARCHAR(100) NOT NULL,
    event_date DATE NOT NULL DEFAULT CURRENT_DATE,
    event_time TIME,
    operator_id UUID REFERENCES users(id),
    operator_name VARCHAR(200),
    insertion_attempts INTEGER,
    ultrasound_used BOOLEAN,
    side_site VARCHAR(120),
    catheter_length_cm NUMERIC(5,1),
    tip_position_confirmation VARCHAR(160),
    immediate_complications JSONB NOT NULL DEFAULT '[]'::jsonb,
    long_term_complications JSONB NOT NULL DEFAULT '[]'::jsonb,
    exit_site_condition TEXT,
    tunnel_cuff_status TEXT,
    lock_solution VARCHAR(120),
    dressing_date DATE,
    hub_scrub_compliant BOOLEAN,
    removal_date DATE,
    removal_reason TEXT,
    replacement_reason TEXT,
    catheter_free_plan TEXT,
    culture_result TEXT,
    antibiotic_course TEXT,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE infection_surveillance_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    access_id UUID REFERENCES vascular_access(id),
    event_type VARCHAR(100) NOT NULL,
    event_date DATE NOT NULL DEFAULT CURRENT_DATE,
    event_time TIME,
    iv_antimicrobial_started BOOLEAN NOT NULL DEFAULT FALSE,
    positive_blood_culture BOOLEAN NOT NULL DEFAULT FALSE,
    access_pus_redness_swelling BOOLEAN NOT NULL DEFAULT FALSE,
    suspected_source VARCHAR(120),
    organism TEXT,
    culture_collected_at TIMESTAMPTZ,
    antimicrobial_started_at TIMESTAMPTZ,
    hospitalized BOOLEAN NOT NULL DEFAULT FALSE,
    death_related BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_window_notes TEXT,
    reported_to_registry BOOLEAN NOT NULL DEFAULT FALSE,
    reported_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE dialysis_adequacy_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    review_month DATE NOT NULL,
    reviewed_by UUID REFERENCES users(id),
    pre_urea_mg_dl NUMERIC(7,2),
    post_urea_mg_dl NUMERIC(7,2),
    urr_percent NUMERIC(5,2),
    sp_kt_v NUMERIC(5,2),
    target_kt_v NUMERIC(5,2),
    residual_kidney_function TEXT,
    urine_volume_ml_day NUMERIC(9,2),
    missed_treatments INTEGER NOT NULL DEFAULT 0,
    shortened_treatments INTEGER NOT NULL DEFAULT 0,
    interdialytic_weight_gain_kg NUMERIC(6,2),
    normalized_protein_catabolic_rate NUMERIC(5,2),
    uf_rate_ml_kg_hr NUMERIC(6,2),
    adequacy_status VARCHAR(60) NOT NULL DEFAULT 'needs_review',
    doctor_review_required BOOLEAN NOT NULL DEFAULT FALSE,
    recommendations TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE medication_reconciliation_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    review_type VARCHAR(80) NOT NULL,
    review_date DATE NOT NULL DEFAULT CURRENT_DATE,
    reviewed_by UUID REFERENCES users(id),
    medications JSONB NOT NULL DEFAULT '[]'::jsonb,
    renal_dosing_flags JSONB NOT NULL DEFAULT '[]'::jsonb,
    dialysis_timing_flags JSONB NOT NULL DEFAULT '[]'::jsonb,
    pregnancy_cautions JSONB NOT NULL DEFAULT '[]'::jsonb,
    nephrotoxin_flags JSONB NOT NULL DEFAULT '[]'::jsonb,
    adherence_notes TEXT,
    regimen_change_reason TEXT,
    recommendations TEXT,
    status VARCHAR(60) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE patient_reported_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    event_type VARCHAR(80) NOT NULL,
    reported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reported_by UUID REFERENCES users(id),
    symptoms JSONB NOT NULL DEFAULT '{}'::jsonb,
    education_topics JSONB NOT NULL DEFAULT '[]'::jsonb,
    teach_back_completed BOOLEAN,
    transport_reliability VARCHAR(80),
    no_show_reason TEXT,
    payment_barrier BOOLEAN,
    food_insecurity BOOLEAN,
    caregiver_support VARCHAR(120),
    sms_whatsapp_consent BOOLEAN,
    follow_up_channel VARCHAR(80),
    follow_up_due_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE staff_attendance_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    staff_id UUID REFERENCES staff_profiles(id),
    user_id UUID REFERENCES users(id),
    shift_assignment_id UUID REFERENCES shift_assignments(id),
    verification_date DATE NOT NULL DEFAULT CURRENT_DATE,
    verification_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    biometric_method VARCHAR(80),
    biometric_device_id VARCHAR(120),
    biometric_template_ref VARCHAR(200),
    verification_result VARCHAR(60) NOT NULL DEFAULT 'pending',
    station_assignment VARCHAR(120),
    patient_assignment JSONB NOT NULL DEFAULT '[]'::jsonb,
    handover_accepted BOOLEAN,
    competency_status JSONB NOT NULL DEFAULT '{}'::jsonb,
    exception_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE unit_safety_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL DEFAULT 'info',
    event_date DATE NOT NULL DEFAULT CURRENT_DATE,
    event_time TIME,
    affected_sessions JSONB NOT NULL DEFAULT '[]'::jsonb,
    affected_patients JSONB NOT NULL DEFAULT '[]'::jsonb,
    machine_id UUID REFERENCES dialysis_machines(id),
    equipment_id UUID REFERENCES equipment(id),
    water_test_id UUID REFERENCES water_treatment_logs(id),
    consumable_lots JSONB NOT NULL DEFAULT '[]'::jsonb,
    immediate_action TEXT,
    root_cause TEXT,
    corrective_action TEXT,
    closure_due_date DATE,
    closed_at TIMESTAMPTZ,
    reported_by UUID REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE interoperability_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID REFERENCES patients(id),
    export_type VARCHAR(80) NOT NULL,
    fhir_resource_type VARCHAR(80),
    local_table VARCHAR(100),
    local_id UUID,
    coding_system VARCHAR(160),
    coding_version VARCHAR(80),
    code_value VARCHAR(120),
    ucum_unit VARCHAR(80),
    fhir_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    export_status VARCHAR(60) NOT NULL DEFAULT 'draft',
    exported_by UUID REFERENCES users(id),
    exported_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE ontology_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID REFERENCES patients(id),
    source_type VARCHAR(100) NOT NULL,
    source_id UUID,
    relation_type VARCHAR(100) NOT NULL,
    target_type VARCHAR(100) NOT NULL,
    target_id UUID,
    relation_context JSONB NOT NULL DEFAULT '{}'::jsonb,
    provenance VARCHAR(160) NOT NULL DEFAULT 'relational_source_of_truth',
    confidence NUMERIC(5,2),
    neo4j_node_ref VARCHAR(200),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE patient_clinical_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE session_safety_checks ENABLE ROW LEVEL SECURITY;
ALTER TABLE clinical_alerts ENABLE ROW LEVEL SECURITY;
ALTER TABLE treatment_telemetry_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE vascular_access_lifecycle_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE infection_surveillance_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE dialysis_adequacy_reviews ENABLE ROW LEVEL SECURITY;
ALTER TABLE medication_reconciliation_reviews ENABLE ROW LEVEL SECURITY;
ALTER TABLE patient_reported_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE staff_attendance_verifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE unit_safety_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE interoperability_exports ENABLE ROW LEVEL SECURITY;
ALTER TABLE ontology_relationships ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_patient_clinical_profiles ON patient_clinical_profiles
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_session_safety_checks ON session_safety_checks
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_clinical_alerts ON clinical_alerts
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_treatment_telemetry_records ON treatment_telemetry_records
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_vascular_access_lifecycle_events ON vascular_access_lifecycle_events
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_infection_surveillance_events ON infection_surveillance_events
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_dialysis_adequacy_reviews ON dialysis_adequacy_reviews
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_medication_reconciliation_reviews ON medication_reconciliation_reviews
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_patient_reported_events ON patient_reported_events
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_staff_attendance_verifications ON staff_attendance_verifications
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_unit_safety_events ON unit_safety_events
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_interoperability_exports ON interoperability_exports
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());
CREATE POLICY tenant_isolation_ontology_relationships ON ontology_relationships
    FOR ALL USING (hospital_id = dms_current_hospital_id()) WITH CHECK (hospital_id = dms_current_hospital_id());

CREATE INDEX idx_patient_clinical_profiles_patient ON patient_clinical_profiles(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_session_safety_checks_session ON session_safety_checks(session_id, checked_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_session_safety_checks_queue ON session_safety_checks(hospital_id, check_status, override_required) WHERE deleted_at IS NULL;
CREATE INDEX idx_clinical_alerts_patient ON clinical_alerts(patient_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_clinical_alerts_queue ON clinical_alerts(hospital_id, status, severity) WHERE deleted_at IS NULL;
CREATE INDEX idx_treatment_telemetry_session ON treatment_telemetry_records(session_id, recorded_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_access_lifecycle_patient ON vascular_access_lifecycle_events(patient_id, event_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_infection_surveillance_patient ON infection_surveillance_events(patient_id, event_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_infection_surveillance_queue ON infection_surveillance_events(hospital_id, reported_to_registry, event_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_adequacy_reviews_patient ON dialysis_adequacy_reviews(patient_id, review_month DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_med_rec_patient ON medication_reconciliation_reviews(patient_id, review_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_patient_reported_patient ON patient_reported_events(patient_id, reported_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_attendance_date ON staff_attendance_verifications(hospital_id, verification_date, verification_result) WHERE deleted_at IS NULL;
CREATE INDEX idx_unit_safety_events_queue ON unit_safety_events(hospital_id, severity, closed_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_interop_patient ON interoperability_exports(patient_id, export_type, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_ontology_patient ON ontology_relationships(patient_id, relation_type) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_patient_clinical_profiles_updated_at BEFORE UPDATE ON patient_clinical_profiles FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_session_safety_checks_updated_at BEFORE UPDATE ON session_safety_checks FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_clinical_alerts_updated_at BEFORE UPDATE ON clinical_alerts FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_treatment_telemetry_records_updated_at BEFORE UPDATE ON treatment_telemetry_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_vascular_access_lifecycle_events_updated_at BEFORE UPDATE ON vascular_access_lifecycle_events FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_infection_surveillance_events_updated_at BEFORE UPDATE ON infection_surveillance_events FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_dialysis_adequacy_reviews_updated_at BEFORE UPDATE ON dialysis_adequacy_reviews FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_medication_reconciliation_reviews_updated_at BEFORE UPDATE ON medication_reconciliation_reviews FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_patient_reported_events_updated_at BEFORE UPDATE ON patient_reported_events FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_staff_attendance_verifications_updated_at BEFORE UPDATE ON staff_attendance_verifications FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_unit_safety_events_updated_at BEFORE UPDATE ON unit_safety_events FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_interoperability_exports_updated_at BEFORE UPDATE ON interoperability_exports FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_ontology_relationships_updated_at BEFORE UPDATE ON ontology_relationships FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS ontology_relationships;
DROP TABLE IF EXISTS interoperability_exports;
DROP TABLE IF EXISTS unit_safety_events;
DROP TABLE IF EXISTS staff_attendance_verifications;
DROP TABLE IF EXISTS patient_reported_events;
DROP TABLE IF EXISTS medication_reconciliation_reviews;
DROP TABLE IF EXISTS dialysis_adequacy_reviews;
DROP TABLE IF EXISTS infection_surveillance_events;
DROP TABLE IF EXISTS vascular_access_lifecycle_events;
DROP TABLE IF EXISTS treatment_telemetry_records;
DROP TABLE IF EXISTS clinical_alerts;
DROP TABLE IF EXISTS session_safety_checks;
DROP TABLE IF EXISTS patient_clinical_profiles;
