-- +goose Up
CREATE TYPE transfer_reason AS ENUM ('higher_level_care','specialist_referral','patient_request','bed_availability','geographic_convenience','transplant_workup','other');
CREATE TYPE transfer_status AS ENUM ('requested','approved','in_transit','completed','cancelled');

CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    from_hospital_id UUID REFERENCES hospitals(id),
    to_hospital_id UUID REFERENCES hospitals(id),
    from_department VARCHAR(150),
    to_department VARCHAR(150),
    reason transfer_reason NOT NULL,
    reason_notes TEXT,
    status transfer_status NOT NULL DEFAULT 'requested',
    requested_by UUID NOT NULL REFERENCES users(id),
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    transferred_at TIMESTAMPTZ,
    received_by VARCHAR(255),
    transfer_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE transfers ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON transfers USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_transfers_patient ON transfers(patient_id);
CREATE INDEX idx_transfers_status ON transfers(status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_transfers_updated_at BEFORE UPDATE ON transfers FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS transfers;
DROP TYPE IF EXISTS transfer_status;
DROP TYPE IF EXISTS transfer_reason;
