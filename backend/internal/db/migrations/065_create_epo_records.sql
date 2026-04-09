-- +goose Up
CREATE TABLE epo_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    administered_by UUID NOT NULL REFERENCES users(id),
    administered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    product_name VARCHAR(255) NOT NULL,
    dose_units INTEGER NOT NULL,
    route medication_route NOT NULL DEFAULT 'sc',
    injection_site VARCHAR(100),
    hb_at_time NUMERIC(4,1),
    hb_target_min NUMERIC(4,1),
    hb_target_max NUMERIC(4,1),
    ferritin_at_time NUMERIC(6,1),
    tsat_at_time NUMERIC(5,2),
    dose_adjustment_reason TEXT,
    next_dose_recommendation INTEGER,
    adverse_reaction BOOLEAN NOT NULL DEFAULT FALSE,
    adverse_reaction_details TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE epo_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON epo_records USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_epo_records_patient ON epo_records(patient_id);
CREATE INDEX idx_epo_records_session ON epo_records(session_id);
CREATE INDEX idx_epo_records_date ON epo_records(hospital_id, administered_at DESC);

CREATE TRIGGER trg_epo_records_updated_at BEFORE UPDATE ON epo_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS epo_records;
