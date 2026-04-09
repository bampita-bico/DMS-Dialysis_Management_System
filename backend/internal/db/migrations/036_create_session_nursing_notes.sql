-- +goose Up
CREATE TABLE session_nursing_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    nurse_id UUID NOT NULL REFERENCES users(id),
    note_type VARCHAR(50) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    content TEXT NOT NULL,
    is_flagged_for_doctor BOOLEAN NOT NULL DEFAULT FALSE,
    doctor_reviewed_by UUID REFERENCES users(id),
    doctor_reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_nursing_notes ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_nursing_notes USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_nursing_notes_session ON session_nursing_notes(session_id, recorded_at);
CREATE INDEX idx_nursing_notes_patient ON session_nursing_notes(patient_id);
CREATE INDEX idx_nursing_notes_flagged ON session_nursing_notes(hospital_id) WHERE is_flagged_for_doctor = TRUE AND doctor_reviewed_at IS NULL;

CREATE TRIGGER trg_session_nursing_notes_updated_at BEFORE UPDATE ON session_nursing_notes FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_nursing_notes;
