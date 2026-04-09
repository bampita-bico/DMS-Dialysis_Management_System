-- +goose Up
CREATE TABLE dry_weight_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    set_by UUID NOT NULL REFERENCES users(id),
    set_date DATE NOT NULL DEFAULT CURRENT_DATE,
    dry_weight_kg NUMERIC(5,2) NOT NULL,
    assessment_method VARCHAR(100),
    clinical_indicators TEXT,
    bp_at_assessment_systolic INTEGER,
    bp_at_assessment_diastolic INTEGER,
    edema_present BOOLEAN NOT NULL DEFAULT FALSE,
    edema_location VARCHAR(200),
    dyspnea_present BOOLEAN NOT NULL DEFAULT FALSE,
    chest_xray_findings TEXT,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_until DATE,
    is_current BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE dry_weight_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON dry_weight_records USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_dry_weight_patient ON dry_weight_records(patient_id, set_date DESC);
CREATE INDEX idx_dry_weight_current ON dry_weight_records(patient_id) WHERE is_current = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_dry_weight_records_updated_at BEFORE UPDATE ON dry_weight_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS dry_weight_records;
