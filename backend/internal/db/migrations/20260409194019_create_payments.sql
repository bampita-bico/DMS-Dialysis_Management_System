-- +goose Up
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    invoice_id UUID REFERENCES invoices(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES billing_accounts(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    payment_date DATE NOT NULL,
    payment_time TIME NOT NULL DEFAULT CURRENT_TIME,
    amount DECIMAL(12,2) NOT NULL,
    payment_method payment_method NOT NULL,
    reference_number VARCHAR(100),
    mobile_money_number VARCHAR(50),
    bank_name VARCHAR(100),
    cheque_number VARCHAR(50),
    card_last_four VARCHAR(4),
    received_by UUID NOT NULL REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_payments_invoice ON payments(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payments_account ON payments(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payments_patient ON payments(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payments_hospital ON payments(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payments_date ON payments(hospital_id, payment_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_payments_method ON payments(hospital_id, payment_method) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
CREATE POLICY payments_isolation ON payments
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS payments CASCADE;
