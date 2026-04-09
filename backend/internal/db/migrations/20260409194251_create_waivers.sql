-- +goose Up
CREATE TABLE waivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    waiver_number VARCHAR(50) NOT NULL,
    waiver_amount DECIMAL(12,2) NOT NULL,
    waiver_percentage DECIMAL(5,2),
    waiver_type VARCHAR(100) NOT NULL,
    reason TEXT NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(id),
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    status waiver_status NOT NULL DEFAULT 'pending',
    supporting_documents TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_waivers_number ON waivers(hospital_id, waiver_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_waivers_invoice ON waivers(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_waivers_patient ON waivers(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_waivers_hospital ON waivers(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_waivers_status ON waivers(hospital_id, status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_waivers_updated_at
BEFORE UPDATE ON waivers
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE waivers ENABLE ROW LEVEL SECURITY;
CREATE POLICY waivers_isolation ON waivers
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS waivers CASCADE;
