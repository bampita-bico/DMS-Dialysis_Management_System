-- +goose Up
CREATE TABLE anticoagulation_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    administered_by UUID NOT NULL REFERENCES users(id),
    anticoagulant VARCHAR(100) NOT NULL,
    route anticoag_route NOT NULL DEFAULT 'systemic',
    loading_dose_units NUMERIC(8,2),
    loading_dose_time TIMESTAMPTZ,
    maintenance_dose_units NUMERIC(8,2),
    maintenance_rate VARCHAR(50),
    total_dose_units NUMERIC(8,2),
    reversal_agent_given BOOLEAN NOT NULL DEFAULT FALSE,
    reversal_agent VARCHAR(100),
    reversal_dose_units NUMERIC(8,2),
    reversal_time TIMESTAMPTZ,
    bleeding_complications BOOLEAN NOT NULL DEFAULT FALSE,
    clotting_observed BOOLEAN NOT NULL DEFAULT FALSE,
    clotting_location VARCHAR(200),
    aptt_pre NUMERIC(5,1),
    aptt_post NUMERIC(5,1),
    act_pre INTEGER,
    act_post INTEGER,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE anticoagulation_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON anticoagulation_records USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_anticoag_session ON anticoagulation_records(session_id);
CREATE INDEX idx_anticoag_patient ON anticoagulation_records(patient_id);
CREATE INDEX idx_anticoag_complications ON anticoagulation_records(hospital_id) WHERE bleeding_complications = TRUE OR clotting_observed = TRUE;

CREATE TRIGGER trg_anticoagulation_records_updated_at BEFORE UPDATE ON anticoagulation_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS anticoagulation_records;
