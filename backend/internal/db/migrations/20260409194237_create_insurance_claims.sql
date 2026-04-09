-- +goose Up
CREATE TABLE insurance_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    scheme_id UUID NOT NULL REFERENCES insurance_schemes(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    claim_number VARCHAR(50) NOT NULL,
    claim_date DATE NOT NULL,
    claimed_amount DECIMAL(12,2) NOT NULL,
    approved_amount DECIMAL(12,2),
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    status claim_status NOT NULL DEFAULT 'draft',
    submitted_by UUID REFERENCES users(id),
    submitted_at TIMESTAMPTZ,
    approved_by_insurer VARCHAR(255),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    payment_reference VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_insurance_claims_number ON insurance_claims(hospital_id, claim_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_claims_invoice ON insurance_claims(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_claims_scheme ON insurance_claims(scheme_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_claims_patient ON insurance_claims(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_claims_hospital ON insurance_claims(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_claims_status ON insurance_claims(hospital_id, status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_insurance_claims_updated_at
BEFORE UPDATE ON insurance_claims
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE insurance_claims ENABLE ROW LEVEL SECURITY;
CREATE POLICY insurance_claims_isolation ON insurance_claims
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS insurance_claims CASCADE;
