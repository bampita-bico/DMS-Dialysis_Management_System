-- +goose Up
CREATE TABLE lab_panels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tests JSONB NOT NULL DEFAULT '[]',
    cost_amount NUMERIC(10,2),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(hospital_id, code)
);

ALTER TABLE lab_panels ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_panels USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_panels_hospital ON lab_panels(hospital_id);
CREATE INDEX idx_lab_panels_code ON lab_panels(hospital_id, code);

CREATE TRIGGER trg_lab_panels_updated_at BEFORE UPDATE ON lab_panels FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_panels;
