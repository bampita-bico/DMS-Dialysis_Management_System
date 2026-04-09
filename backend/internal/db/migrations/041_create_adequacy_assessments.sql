-- +goose Up
CREATE TABLE adequacy_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    assessed_by UUID NOT NULL REFERENCES users(id),
    assessment_date DATE NOT NULL,
    kt_v NUMERIC(4,2),
    kt_v_method VARCHAR(50),
    urr_percent NUMERIC(5,2),
    pre_bun_mg_dl NUMERIC(6,2),
    post_bun_mg_dl NUMERIC(6,2),
    pre_creatinine_mg_dl NUMERIC(6,2),
    post_creatinine_mg_dl NUMERIC(6,2),
    dialysis_duration_mins INTEGER,
    blood_flow_rate INTEGER,
    dialyzer_clearance INTEGER,
    body_water_volume_liters NUMERIC(5,2),
    is_adequate BOOLEAN NOT NULL DEFAULT TRUE,
    recommendations TEXT,
    next_assessment_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE adequacy_assessments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON adequacy_assessments USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_adequacy_patient ON adequacy_assessments(patient_id, assessment_date DESC);
CREATE INDEX idx_adequacy_inadequate ON adequacy_assessments(hospital_id) WHERE is_adequate = FALSE;

CREATE TRIGGER trg_adequacy_assessments_updated_at BEFORE UPDATE ON adequacy_assessments FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS adequacy_assessments;
