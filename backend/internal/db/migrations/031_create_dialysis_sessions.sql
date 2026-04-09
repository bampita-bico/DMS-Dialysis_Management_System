-- +goose Up
CREATE TABLE dialysis_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    schedule_id UUID REFERENCES session_schedules(id),
    machine_id UUID NOT NULL REFERENCES dialysis_machines(id),
    access_id UUID,
    modality dialysis_modality NOT NULL DEFAULT 'hd',
    shift shift_type NOT NULL,
    status session_status NOT NULL DEFAULT 'scheduled',
    scheduled_date DATE NOT NULL,
    scheduled_start_time TIME NOT NULL,
    actual_start_time TIMESTAMPTZ,
    actual_end_time TIMESTAMPTZ,
    prescribed_duration_mins INTEGER NOT NULL DEFAULT 240,
    actual_duration_mins INTEGER,
    pre_weight_kg NUMERIC(5,2),
    pre_bp_systolic INTEGER,
    pre_bp_diastolic INTEGER,
    pre_hr INTEGER,
    pre_temp NUMERIC(4,1),
    post_weight_kg NUMERIC(5,2),
    post_bp_systolic INTEGER,
    post_bp_diastolic INTEGER,
    post_hr INTEGER,
    session_notes TEXT,
    aborted_reason TEXT,
    was_patient_reviewed BOOLEAN NOT NULL DEFAULT FALSE,
    reviewed_by UUID REFERENCES users(id),
    primary_nurse_id UUID REFERENCES users(id),
    supervising_doctor_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE dialysis_sessions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON dialysis_sessions USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_sessions_patient ON dialysis_sessions(patient_id);
CREATE INDEX idx_sessions_machine ON dialysis_sessions(machine_id);
CREATE INDEX idx_sessions_date ON dialysis_sessions(hospital_id, scheduled_date);
CREATE INDEX idx_sessions_status ON dialysis_sessions(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_sessions_active ON dialysis_sessions(hospital_id) WHERE status = 'in_progress';

CREATE TRIGGER trg_dialysis_sessions_updated_at BEFORE UPDATE ON dialysis_sessions FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS dialysis_sessions;
