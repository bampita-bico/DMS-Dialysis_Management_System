-- +goose Up
CREATE TABLE drug_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    medication_a_id UUID NOT NULL REFERENCES medications(id),
    medication_b_id UUID NOT NULL REFERENCES medications(id),
    severity VARCHAR(50) NOT NULL,
    interaction_type VARCHAR(100),
    description TEXT NOT NULL,
    clinical_effect TEXT,
    management_recommendation TEXT,
    evidence_level VARCHAR(50),
    reference_source TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    CONSTRAINT check_different_medications CHECK (medication_a_id != medication_b_id)
);

ALTER TABLE drug_interactions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON drug_interactions USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_drug_interactions_med_a ON drug_interactions(medication_a_id);
CREATE INDEX idx_drug_interactions_med_b ON drug_interactions(medication_b_id);
CREATE INDEX idx_drug_interactions_severity ON drug_interactions(hospital_id, severity);
CREATE UNIQUE INDEX idx_drug_interactions_pair ON drug_interactions(hospital_id, medication_a_id, medication_b_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_drug_interactions_updated_at BEFORE UPDATE ON drug_interactions FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS drug_interactions;
