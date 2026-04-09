-- +goose Up
CREATE TABLE shift_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shift_date DATE NOT NULL,
    shift_type shift_type NOT NULL,
    shift_start_time TIME NOT NULL,
    shift_end_time TIME NOT NULL,
    machine_ids JSONB,
    assigned_by UUID REFERENCES users(id),
    clock_in_time TIMESTAMPTZ,
    clock_out_time TIMESTAMPTZ,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_shift_assignments_staff ON shift_assignments(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_shift_assignments_hospital ON shift_assignments(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_shift_assignments_date ON shift_assignments(hospital_id, shift_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_shift_assignments_staff_date ON shift_assignments(staff_id, shift_date) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_shift_assignments_updated_at
BEFORE UPDATE ON shift_assignments
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE shift_assignments ENABLE ROW LEVEL SECURITY;
CREATE POLICY shift_assignments_isolation ON shift_assignments
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS shift_assignments CASCADE;
