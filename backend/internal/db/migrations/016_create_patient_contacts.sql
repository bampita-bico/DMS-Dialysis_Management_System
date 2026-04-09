-- +goose Up
CREATE TABLE patient_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    contact_type contact_type NOT NULL,
    value VARCHAR(255) NOT NULL,
    label VARCHAR(100),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE patient_contacts ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON patient_contacts USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_patient_contacts_patient ON patient_contacts(patient_id);
CREATE INDEX idx_patient_contacts_primary ON patient_contacts(patient_id, is_primary) WHERE is_primary = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_patient_contacts_updated_at BEFORE UPDATE ON patient_contacts FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS patient_contacts;
