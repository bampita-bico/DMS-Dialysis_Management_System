-- +goose Up
CREATE TYPE comorbidity_status AS ENUM ('active','controlled','resolved','suspected');

CREATE TABLE comorbidities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    condition VARCHAR(255) NOT NULL,
    icd10_code VARCHAR(20),
    status comorbidity_status NOT NULL DEFAULT 'active',
    diagnosed_at DATE,
    diagnosed_by UUID REFERENCES users(id),
    resolved_at DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE comorbidities ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON comorbidities USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_comorbidities_patient ON comorbidities(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_comorbidities_icd10 ON comorbidities(icd10_code) WHERE icd10_code IS NOT NULL;

CREATE TRIGGER trg_comorbidities_updated_at BEFORE UPDATE ON comorbidities FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS comorbidities;
DROP TYPE IF EXISTS comorbidity_status;
