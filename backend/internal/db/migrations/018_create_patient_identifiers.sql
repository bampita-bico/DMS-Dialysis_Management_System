-- +goose Up
CREATE TABLE patient_identifiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    id_type id_type NOT NULL,
    id_value VARCHAR(255) NOT NULL,
    issuing_country VARCHAR(100) DEFAULT 'Uganda',
    issuing_authority VARCHAR(150),
    issued_date DATE,
    expiry_date DATE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(hospital_id, id_type, id_value)
);

ALTER TABLE patient_identifiers ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON patient_identifiers USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_patient_identifiers_patient ON patient_identifiers(patient_id);
CREATE INDEX idx_patient_identifiers_value ON patient_identifiers(id_type, id_value);

CREATE TRIGGER trg_patient_identifiers_updated_at BEFORE UPDATE ON patient_identifiers FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS patient_identifiers;
