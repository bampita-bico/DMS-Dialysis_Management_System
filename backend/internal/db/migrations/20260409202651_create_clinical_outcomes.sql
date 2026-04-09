-- +goose Up
CREATE TABLE clinical_outcomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    assessment_date DATE NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    hemoglobin DECIMAL(5,2),
    hemoglobin_target_min DECIMAL(5,2),
    hemoglobin_target_max DECIMAL(5,2),
    kt_v DECIMAL(5,3),
    kt_v_target DECIMAL(5,3),
    urr DECIMAL(5,2),
    systolic_bp_avg DECIMAL(5,2),
    diastolic_bp_avg DECIMAL(5,2),
    bp_controlled BOOLEAN,
    weight_gain_percent DECIMAL(5,2),
    albumin DECIMAL(5,2),
    phosphate DECIMAL(5,2),
    calcium DECIMAL(5,2),
    pth DECIMAL(7,2),
    quality_of_life_score DECIMAL(5,2),
    functional_status VARCHAR(100),
    adverse_events_count INTEGER NOT NULL DEFAULT 0,
    hospitalizations_count INTEGER NOT NULL DEFAULT 0,
    missed_sessions_count INTEGER NOT NULL DEFAULT 0,
    outcome_trend outcome_trend,
    assessed_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_clinical_outcomes_patient ON clinical_outcomes(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_clinical_outcomes_hospital ON clinical_outcomes(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_clinical_outcomes_date ON clinical_outcomes(hospital_id, assessment_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_clinical_outcomes_trend ON clinical_outcomes(hospital_id, outcome_trend) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_clinical_outcomes_updated_at
BEFORE UPDATE ON clinical_outcomes
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE clinical_outcomes ENABLE ROW LEVEL SECURITY;
CREATE POLICY clinical_outcomes_isolation ON clinical_outcomes
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS clinical_outcomes CASCADE;
