-- +goose Up
CREATE TABLE insurance_schemes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    short_code VARCHAR(50),
    country VARCHAR(100),
    covers_dialysis BOOLEAN NOT NULL DEFAULT TRUE,
    reimbursement_rate DECIMAL(5,2),
    requires_pre_authorization BOOLEAN NOT NULL DEFAULT FALSE,
    claim_submission_url TEXT,
    contact_person VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_insurance_schemes_hospital ON insurance_schemes(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_schemes_active ON insurance_schemes(hospital_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_schemes_country ON insurance_schemes(country) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_insurance_schemes_updated_at
BEFORE UPDATE ON insurance_schemes
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE insurance_schemes ENABLE ROW LEVEL SECURITY;
CREATE POLICY insurance_schemes_isolation ON insurance_schemes
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS insurance_schemes CASCADE;
