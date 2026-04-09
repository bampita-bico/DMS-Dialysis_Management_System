-- +goose Up
CREATE TABLE session_complications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    reported_by UUID NOT NULL REFERENCES users(id),
    occurred_at TIMESTAMPTZ NOT NULL,
    complication_type VARCHAR(100) NOT NULL,
    severity complication_severity NOT NULL DEFAULT 'minor',
    symptoms TEXT NOT NULL,
    vital_signs_at_event JSONB DEFAULT '{}',
    immediate_action_taken TEXT,
    outcome TEXT,
    required_hospitalization BOOLEAN NOT NULL DEFAULT FALSE,
    was_session_terminated BOOLEAN NOT NULL DEFAULT FALSE,
    doctor_notified BOOLEAN NOT NULL DEFAULT FALSE,
    doctor_id UUID REFERENCES users(id),
    family_notified BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_complications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_complications USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_complications_session ON session_complications(session_id);
CREATE INDEX idx_complications_patient ON session_complications(patient_id);
CREATE INDEX idx_complications_severity ON session_complications(hospital_id, severity) WHERE severity IN ('severe', 'life_threatening');
CREATE INDEX idx_complications_hospitalization ON session_complications(hospital_id) WHERE required_hospitalization = TRUE;

CREATE TRIGGER trg_session_complications_updated_at BEFORE UPDATE ON session_complications FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_complications;
