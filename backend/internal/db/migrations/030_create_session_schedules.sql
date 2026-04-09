-- +goose Up
CREATE TABLE session_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    machine_id UUID REFERENCES dialysis_machines(id),
    shift shift_type NOT NULL,
    days_of_week INTEGER[] NOT NULL,
    frequency_weeks INTEGER NOT NULL DEFAULT 1,
    modality dialysis_modality NOT NULL DEFAULT 'hd',
    prescribed_duration_mins INTEGER NOT NULL DEFAULT 240,
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_until DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID NOT NULL REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_schedules ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_schedules USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_schedules_patient ON session_schedules(patient_id) WHERE is_active = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_schedules_machine ON session_schedules(machine_id) WHERE is_active = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_session_schedules_updated_at BEFORE UPDATE ON session_schedules FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_schedules;
