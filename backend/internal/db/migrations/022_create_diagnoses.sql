-- +goose Up
CREATE TYPE diagnosis_type AS ENUM ('primary','secondary','differential','working','confirmed','ruled_out');

CREATE TABLE diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    icd10_code VARCHAR(20) NOT NULL,
    description VARCHAR(500) NOT NULL,
    diagnosis_type diagnosis_type NOT NULL DEFAULT 'working',
    diagnosed_by UUID NOT NULL REFERENCES users(id),
    diagnosed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    admission_id UUID,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE diagnoses ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON diagnoses USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_diagnoses_patient ON diagnoses(patient_id);
CREATE INDEX idx_diagnoses_icd10 ON diagnoses(icd10_code);

CREATE TRIGGER trg_diagnoses_updated_at BEFORE UPDATE ON diagnoses FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS diagnoses;
DROP TYPE IF EXISTS diagnosis_type;
