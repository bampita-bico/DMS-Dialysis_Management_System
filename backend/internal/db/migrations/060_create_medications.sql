-- +goose Up
CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    generic_name VARCHAR(255) NOT NULL,
    brand_names JSONB DEFAULT '[]',
    drug_class VARCHAR(100),
    form medication_form NOT NULL,
    strength VARCHAR(100),
    unit VARCHAR(50),
    is_controlled BOOLEAN NOT NULL DEFAULT FALSE,
    requires_prescription BOOLEAN NOT NULL DEFAULT TRUE,
    is_essential_who BOOLEAN NOT NULL DEFAULT FALSE,
    storage_conditions TEXT,
    cost_per_unit NUMERIC(10,2),
    reorder_level INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE medications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON medications USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_medications_hospital ON medications(hospital_id);
CREATE INDEX idx_medications_generic_name ON medications(hospital_id, generic_name);
CREATE INDEX idx_medications_drug_class ON medications(hospital_id, drug_class);
CREATE INDEX idx_medications_active ON medications(hospital_id, is_active) WHERE is_active = TRUE;

CREATE TRIGGER trg_medications_updated_at BEFORE UPDATE ON medications FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS medications;
