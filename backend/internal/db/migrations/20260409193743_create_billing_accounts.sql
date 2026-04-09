-- +goose Up
CREATE TABLE billing_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    guarantor_id UUID REFERENCES patients(id),
    account_number VARCHAR(50) NOT NULL,
    account_status account_status NOT NULL DEFAULT 'active',
    credit_limit DECIMAL(12,2),
    current_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_billed DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_paid DECIMAL(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_billing_accounts_number ON billing_accounts(hospital_id, account_number) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_billing_accounts_patient ON billing_accounts(hospital_id, patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_billing_accounts_hospital ON billing_accounts(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_billing_accounts_status ON billing_accounts(hospital_id, account_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_billing_accounts_guarantor ON billing_accounts(guarantor_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_billing_accounts_updated_at
BEFORE UPDATE ON billing_accounts
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE billing_accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY billing_accounts_isolation ON billing_accounts
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS billing_accounts CASCADE;
