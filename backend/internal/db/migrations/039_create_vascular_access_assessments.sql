-- +goose Up
CREATE TABLE vascular_access_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    access_id UUID NOT NULL REFERENCES vascular_access(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    assessed_by UUID NOT NULL REFERENCES users(id),
    assessed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    has_thrill BOOLEAN,
    has_bruit BOOLEAN,
    has_redness BOOLEAN NOT NULL DEFAULT FALSE,
    has_swelling BOOLEAN NOT NULL DEFAULT FALSE,
    has_discharge BOOLEAN NOT NULL DEFAULT FALSE,
    has_bleeding BOOLEAN NOT NULL DEFAULT FALSE,
    has_pain BOOLEAN NOT NULL DEFAULT FALSE,
    appearance_normal BOOLEAN NOT NULL DEFAULT TRUE,
    flow_rate_ml_min INTEGER,
    venous_pressure_mmhg INTEGER,
    arterial_pressure_mmhg INTEGER,
    recirculation_percent NUMERIC(4,1),
    requires_intervention BOOLEAN NOT NULL DEFAULT FALSE,
    intervention_type VARCHAR(100),
    intervention_urgency VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE vascular_access_assessments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON vascular_access_assessments USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_access_assessments_access ON vascular_access_assessments(access_id, assessed_at);
CREATE INDEX idx_access_assessments_patient ON vascular_access_assessments(patient_id);
CREATE INDEX idx_access_assessments_session ON vascular_access_assessments(session_id);
CREATE INDEX idx_access_assessments_intervention ON vascular_access_assessments(hospital_id) WHERE requires_intervention = TRUE;

CREATE TRIGGER trg_vascular_access_assessments_updated_at BEFORE UPDATE ON vascular_access_assessments FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS vascular_access_assessments;
