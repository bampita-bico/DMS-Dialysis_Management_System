-- +goose Up
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES billing_accounts(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    session_id UUID REFERENCES dialysis_sessions(id),
    invoice_number VARCHAR(50) NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE,
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    net_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    balance_due DECIMAL(12,2) NOT NULL DEFAULT 0,
    status invoice_status NOT NULL DEFAULT 'draft',
    issued_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_invoices_number ON invoices(hospital_id, invoice_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_account ON invoices(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_patient ON invoices(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_session ON invoices(session_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_hospital ON invoices(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_status ON invoices(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_date ON invoices(hospital_id, invoice_date) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_invoices_updated_at
BEFORE UPDATE ON invoices
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
CREATE POLICY invoices_isolation ON invoices
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS invoices CASCADE;
