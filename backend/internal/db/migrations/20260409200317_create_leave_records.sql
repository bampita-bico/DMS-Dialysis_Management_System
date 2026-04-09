-- +goose Up
CREATE TABLE leave_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    leave_type leave_type NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days_requested INTEGER NOT NULL,
    days_approved INTEGER,
    reason TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    status leave_status NOT NULL DEFAULT 'pending',
    relief_staff_id UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_leave_records_staff ON leave_records(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_leave_records_hospital ON leave_records(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_leave_records_type ON leave_records(hospital_id, leave_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_leave_records_status ON leave_records(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_leave_records_dates ON leave_records(hospital_id, start_date, end_date) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_leave_records_updated_at
BEFORE UPDATE ON leave_records
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE leave_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY leave_records_isolation ON leave_records
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS leave_records CASCADE;
