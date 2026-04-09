-- +goose Up
CREATE TABLE lab_test_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    unit VARCHAR(50),
    turnaround_hrs INTEGER NOT NULL DEFAULT 24,
    specimen_type specimen_type NOT NULL DEFAULT 'serum',
    specimen_volume_ml NUMERIC(5,2),
    requires_fasting BOOLEAN NOT NULL DEFAULT FALSE,
    cost_amount NUMERIC(10,2),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(hospital_id, code)
);

ALTER TABLE lab_test_catalog ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_test_catalog USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_test_catalog_hospital ON lab_test_catalog(hospital_id);
CREATE INDEX idx_lab_test_catalog_code ON lab_test_catalog(hospital_id, code);
CREATE INDEX idx_lab_test_catalog_category ON lab_test_catalog(hospital_id, category);

CREATE TRIGGER trg_lab_test_catalog_updated_at BEFORE UPDATE ON lab_test_catalog FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_test_catalog;
