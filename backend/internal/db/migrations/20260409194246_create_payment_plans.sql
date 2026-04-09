-- +goose Up
CREATE TABLE payment_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES billing_accounts(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    plan_number VARCHAR(50) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    down_payment DECIMAL(12,2) NOT NULL DEFAULT 0,
    installment_amount DECIMAL(12,2) NOT NULL,
    installment_frequency VARCHAR(50) NOT NULL,
    number_of_installments INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    amount_paid DECIMAL(12,2) NOT NULL DEFAULT 0,
    balance_remaining DECIMAL(12,2) NOT NULL,
    status payment_plan_status NOT NULL DEFAULT 'active',
    approved_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_payment_plans_number ON payment_plans(hospital_id, plan_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_plans_account ON payment_plans(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_plans_patient ON payment_plans(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_plans_hospital ON payment_plans(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_plans_status ON payment_plans(hospital_id, status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_payment_plans_updated_at
BEFORE UPDATE ON payment_plans
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE payment_plans ENABLE ROW LEVEL SECURITY;
CREATE POLICY payment_plans_isolation ON payment_plans
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS payment_plans CASCADE;
