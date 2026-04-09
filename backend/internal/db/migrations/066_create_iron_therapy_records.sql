-- +goose Up
CREATE TABLE iron_therapy_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    administered_by UUID NOT NULL REFERENCES users(id),
    administered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    product VARCHAR(255) NOT NULL,
    dose_mg INTEGER NOT NULL,
    route medication_route NOT NULL DEFAULT 'iv',
    infusion_duration_mins INTEGER,
    dilution VARCHAR(200),
    ferritin_at_time NUMERIC(6,1),
    ferritin_target_min NUMERIC(6,1),
    ferritin_target_max NUMERIC(6,1),
    tsat_at_time NUMERIC(5,2),
    hb_at_time NUMERIC(4,1),
    adverse_reaction BOOLEAN NOT NULL DEFAULT FALSE,
    adverse_reaction_type VARCHAR(200),
    adverse_reaction_severity VARCHAR(50),
    treatment_given TEXT,
    test_dose_given BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE iron_therapy_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON iron_therapy_records USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_iron_therapy_records_patient ON iron_therapy_records(patient_id);
CREATE INDEX idx_iron_therapy_records_session ON iron_therapy_records(session_id);
CREATE INDEX idx_iron_therapy_records_date ON iron_therapy_records(hospital_id, administered_at DESC);
CREATE INDEX idx_iron_therapy_records_adverse ON iron_therapy_records(hospital_id) WHERE adverse_reaction = TRUE;

CREATE TRIGGER trg_iron_therapy_records_updated_at BEFORE UPDATE ON iron_therapy_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS iron_therapy_records;
