-- +goose Up
CREATE TABLE lab_reference_ranges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    test_id UUID NOT NULL REFERENCES lab_test_catalog(id),
    age_group VARCHAR(50),
    sex VARCHAR(20),
    min_value NUMERIC(15,4),
    max_value NUMERIC(15,4),
    critical_low NUMERIC(15,4),
    critical_high NUMERIC(15,4),
    reference_text TEXT,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE lab_reference_ranges ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_reference_ranges USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_reference_ranges_test ON lab_reference_ranges(test_id);
CREATE INDEX idx_lab_reference_ranges_demographics ON lab_reference_ranges(test_id, age_group, sex);

CREATE TRIGGER trg_lab_reference_ranges_updated_at BEFORE UPDATE ON lab_reference_ranges FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_reference_ranges;
