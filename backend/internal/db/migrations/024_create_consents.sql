-- +goose Up
CREATE TYPE consent_type AS ENUM ('dialysis_treatment','vascular_access_procedure','blood_transfusion','surgery','hiv_testing','data_sharing','research','photography','general_treatment');
CREATE TYPE consent_status AS ENUM ('given','refused','withdrawn','expired');

CREATE TABLE consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    consent_type consent_type NOT NULL,
    status consent_status NOT NULL DEFAULT 'given',
    given_by VARCHAR(255) NOT NULL,
    relationship VARCHAR(100),
    witness_id UUID REFERENCES users(id),
    signed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ,
    withdrawn_at TIMESTAMPTZ,
    withdrawn_reason TEXT,
    document_url TEXT,
    notes TEXT,
    recorded_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE consents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON consents USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_consents_patient ON consents(patient_id);
CREATE INDEX idx_consents_active ON consents(patient_id, consent_type, status) WHERE status = 'given' AND deleted_at IS NULL;

CREATE TRIGGER trg_consents_updated_at BEFORE UPDATE ON consents FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS consents;
DROP TYPE IF EXISTS consent_status;
DROP TYPE IF EXISTS consent_type;
