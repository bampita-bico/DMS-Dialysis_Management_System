-- +goose Up
CREATE TABLE session_vitals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    recorded_by UUID NOT NULL REFERENCES users(id),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    time_on_dialysis_mins INTEGER,
    bp_systolic INTEGER,
    bp_diastolic INTEGER,
    heart_rate INTEGER,
    temperature NUMERIC(4,1),
    spo2 INTEGER,
    respiratory_rate INTEGER,
    blood_flow_actual INTEGER,
    dialysate_flow_actual INTEGER,
    venous_pressure INTEGER,
    arterial_pressure INTEGER,
    tmp INTEGER,
    uf_removed_so_far NUMERIC(7,2),
    conductivity_actual NUMERIC(4,1),
    has_hypotension_alert BOOLEAN NOT NULL DEFAULT FALSE,
    has_hypertension_alert BOOLEAN NOT NULL DEFAULT FALSE,
    has_tachycardia_alert BOOLEAN NOT NULL DEFAULT FALSE,
    alert_acknowledged_by UUID REFERENCES users(id),
    alert_acknowledged_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_vitals ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_vitals USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_vitals_session ON session_vitals(session_id, recorded_at);
CREATE INDEX idx_vitals_alerts ON session_vitals(session_id) WHERE has_hypotension_alert = TRUE OR has_hypertension_alert = TRUE;

CREATE TRIGGER trg_session_vitals_updated_at BEFORE UPDATE ON session_vitals FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_vitals;
