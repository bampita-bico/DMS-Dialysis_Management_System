-- +goose Up
CREATE TABLE hospitalizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    admission_date DATE NOT NULL,
    admission_time TIME,
    discharge_date DATE,
    discharge_time TIME,
    length_of_stay_days INTEGER,
    admission_reason TEXT NOT NULL,
    admission_diagnosis TEXT,
    icd10_codes TEXT,
    admitting_facility VARCHAR(255),
    ward_name VARCHAR(100),
    dialysis_related BOOLEAN NOT NULL DEFAULT FALSE,
    access_related BOOLEAN NOT NULL DEFAULT FALSE,
    infection_related BOOLEAN NOT NULL DEFAULT FALSE,
    treatment_given TEXT,
    procedures_performed TEXT,
    outcome hospitalization_outcome,
    discharge_destination VARCHAR(255),
    follow_up_required BOOLEAN NOT NULL DEFAULT FALSE,
    follow_up_date DATE,
    recorded_by UUID NOT NULL REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_hospitalizations_patient ON hospitalizations(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_hospitalizations_hospital ON hospitalizations(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_hospitalizations_admission ON hospitalizations(hospital_id, admission_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_hospitalizations_dialysis_related ON hospitalizations(hospital_id, dialysis_related) WHERE deleted_at IS NULL AND dialysis_related = TRUE;
CREATE INDEX idx_hospitalizations_outcome ON hospitalizations(hospital_id, outcome) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_hospitalizations_updated_at
BEFORE UPDATE ON hospitalizations
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE hospitalizations ENABLE ROW LEVEL SECURITY;
CREATE POLICY hospitalizations_isolation ON hospitalizations
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS hospitalizations CASCADE;
