-- +goose Up
CREATE TABLE session_staff_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    staff_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_by UUID NOT NULL REFERENCES users(id),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_staff_assignments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_staff_assignments USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_staff_assignments_session ON session_staff_assignments(session_id);
CREATE INDEX idx_staff_assignments_staff ON session_staff_assignments(staff_id);
CREATE INDEX idx_staff_assignments_active ON session_staff_assignments(hospital_id) WHERE started_at IS NOT NULL AND completed_at IS NULL;

CREATE TRIGGER trg_session_staff_assignments_updated_at BEFORE UPDATE ON session_staff_assignments FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_staff_assignments;
