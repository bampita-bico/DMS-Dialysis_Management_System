-- +goose Up
CREATE TABLE mortality_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    session_id UUID REFERENCES dialysis_sessions(id),
    date_of_death DATE NOT NULL,
    time_of_death TIME,
    death_setting death_setting NOT NULL,
    session_related BOOLEAN NOT NULL DEFAULT FALSE,
    primary_cause_of_death TEXT NOT NULL,
    contributing_factors TEXT,
    icd10_code VARCHAR(10),
    autopsy_performed BOOLEAN NOT NULL DEFAULT FALSE,
    autopsy_findings TEXT,
    reported_by UUID NOT NULL REFERENCES users(id),
    reported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    certified_by UUID REFERENCES users(id),
    certified_at TIMESTAMPTZ,
    death_certificate_number VARCHAR(100),
    family_notified_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_mortality_records_patient ON mortality_records(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_mortality_records_hospital ON mortality_records(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_mortality_records_date ON mortality_records(hospital_id, date_of_death) WHERE deleted_at IS NULL;
CREATE INDEX idx_mortality_records_setting ON mortality_records(hospital_id, death_setting) WHERE deleted_at IS NULL;
CREATE INDEX idx_mortality_records_session_related ON mortality_records(hospital_id, session_related) WHERE deleted_at IS NULL AND session_related = TRUE;

CREATE TRIGGER trg_mortality_records_updated_at
BEFORE UPDATE ON mortality_records
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE mortality_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY mortality_records_isolation ON mortality_records
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS mortality_records CASCADE;
