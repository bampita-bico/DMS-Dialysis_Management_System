-- +goose Up
-- Drop any partially created objects from previous failed attempt
DROP INDEX IF EXISTS idx_prescriptions_patient;
DROP INDEX IF EXISTS idx_prescriptions_hospital;
DROP INDEX IF EXISTS idx_prescriptions_session;
DROP INDEX IF EXISTS idx_prescriptions_status;
DROP INDEX IF EXISTS idx_prescriptions_date;
DROP TABLE IF EXISTS prescriptions CASCADE;

CREATE TABLE prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    prescribed_by UUID NOT NULL REFERENCES users(id),
    prescribed_date DATE NOT NULL DEFAULT CURRENT_DATE,
    prescribed_time TIME NOT NULL DEFAULT CURRENT_TIME,
    status prescription_status NOT NULL DEFAULT 'active',
    valid_from DATE NOT NULL DEFAULT CURRENT_DATE,
    valid_until DATE,
    diagnosis TEXT,
    clinical_notes TEXT,
    pharmacist_verified_by UUID REFERENCES users(id),
    pharmacist_verified_at TIMESTAMPTZ,
    dispensed_by UUID REFERENCES users(id),
    dispensed_at TIMESTAMPTZ,
    cancelled_by UUID REFERENCES users(id),
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE prescriptions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON prescriptions USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_prescriptions_hospital ON prescriptions(hospital_id);
CREATE INDEX idx_prescriptions_patient ON prescriptions(patient_id);
CREATE INDEX idx_prescriptions_session ON prescriptions(session_id);
CREATE INDEX idx_prescriptions_status ON prescriptions(hospital_id, status);
CREATE INDEX idx_prescriptions_date ON prescriptions(hospital_id, prescribed_date DESC);

CREATE TRIGGER trg_prescriptions_updated_at BEFORE UPDATE ON prescriptions FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS prescriptions;
