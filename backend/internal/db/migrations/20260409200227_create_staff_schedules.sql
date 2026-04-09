-- +goose Up
CREATE TABLE staff_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    schedule_name VARCHAR(255) NOT NULL,
    schedule_type schedule_type NOT NULL,
    effective_from DATE NOT NULL,
    effective_until DATE,
    monday_shift VARCHAR(50),
    tuesday_shift VARCHAR(50),
    wednesday_shift VARCHAR(50),
    thursday_shift VARCHAR(50),
    friday_shift VARCHAR(50),
    saturday_shift VARCHAR(50),
    sunday_shift VARCHAR(50),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_staff_schedules_staff ON staff_schedules(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_schedules_hospital ON staff_schedules(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_schedules_effective ON staff_schedules(hospital_id, effective_from, effective_until) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_schedules_active ON staff_schedules(hospital_id, is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_staff_schedules_updated_at
BEFORE UPDATE ON staff_schedules
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE staff_schedules ENABLE ROW LEVEL SECURITY;
CREATE POLICY staff_schedules_isolation ON staff_schedules
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS staff_schedules CASCADE;
