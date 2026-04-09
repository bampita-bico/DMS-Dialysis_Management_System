-- +goose Up
CREATE TABLE next_of_kin (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    full_name VARCHAR(255) NOT NULL,
    relationship VARCHAR(100) NOT NULL,
    phone_primary VARCHAR(50) NOT NULL,
    phone_secondary VARCHAR(50),
    address TEXT,
    national_id VARCHAR(100),
    is_legal_guardian BOOLEAN NOT NULL DEFAULT FALSE,
    is_emergency_contact BOOLEAN NOT NULL DEFAULT TRUE,
    priority_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE next_of_kin ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON next_of_kin USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_next_of_kin_patient ON next_of_kin(patient_id);

CREATE TRIGGER trg_next_of_kin_updated_at BEFORE UPDATE ON next_of_kin FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS next_of_kin;
