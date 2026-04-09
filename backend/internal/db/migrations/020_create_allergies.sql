-- +goose Up
CREATE TYPE allergy_category AS ENUM ('drug','food','contrast','latex','environmental','other');
CREATE TYPE allergy_reaction AS ENUM ('rash','urticaria','angioedema','anaphylaxis','bronchospasm','nausea','other');

CREATE TABLE allergies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    allergen VARCHAR(255) NOT NULL,
    category allergy_category NOT NULL,
    reaction allergy_reaction NOT NULL,
    severity severity_level NOT NULL,
    onset_date DATE,
    notes TEXT,
    recorded_by UUID NOT NULL REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE allergies ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON allergies USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_allergies_patient ON allergies(patient_id) WHERE is_active = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_allergies_updated_at BEFORE UPDATE ON allergies FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS allergies;
DROP TYPE IF EXISTS allergy_reaction;
DROP TYPE IF EXISTS allergy_category;
