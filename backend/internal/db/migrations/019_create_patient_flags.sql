-- +goose Up
CREATE TABLE patient_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    flag_type flag_type NOT NULL,
    reason TEXT NOT NULL,
    flagged_by UUID NOT NULL REFERENCES users(id),
    flagged_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE patient_flags ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON patient_flags USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_patient_flags_patient ON patient_flags(patient_id, is_active) WHERE is_active = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_patient_flags_infectious ON patient_flags(flag_type) WHERE flag_type IN ('infectious','hiv_positive','hepatitis_b','hepatitis_c');

CREATE TRIGGER trg_patient_flags_updated_at BEFORE UPDATE ON patient_flags FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS patient_flags;
