-- +goose Up
CREATE TABLE medication_administrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    prescription_item_id UUID NOT NULL REFERENCES prescription_items(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    administered_by UUID NOT NULL REFERENCES users(id),
    scheduled_time TIMESTAMPTZ NOT NULL,
    administered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    dose_given VARCHAR(100) NOT NULL,
    route medication_route NOT NULL,
    site VARCHAR(100),
    is_refused BOOLEAN NOT NULL DEFAULT FALSE,
    refusal_reason TEXT,
    is_omitted BOOLEAN NOT NULL DEFAULT FALSE,
    omission_reason TEXT,
    adverse_reaction BOOLEAN NOT NULL DEFAULT FALSE,
    adverse_reaction_details TEXT,
    witnessed_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE medication_administrations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON medication_administrations USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_medication_administrations_item ON medication_administrations(prescription_item_id);
CREATE INDEX idx_medication_administrations_patient ON medication_administrations(patient_id);
CREATE INDEX idx_medication_administrations_session ON medication_administrations(session_id);
CREATE INDEX idx_medication_administrations_date ON medication_administrations(hospital_id, administered_at DESC);
CREATE INDEX idx_medication_administrations_adverse ON medication_administrations(hospital_id) WHERE adverse_reaction = TRUE;

CREATE TRIGGER trg_medication_administrations_updated_at BEFORE UPDATE ON medication_administrations FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS medication_administrations;
