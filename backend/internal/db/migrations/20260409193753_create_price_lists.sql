-- +goose Up
CREATE TABLE price_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    service_name VARCHAR(255) NOT NULL,
    service_code VARCHAR(50),
    service_category VARCHAR(100),
    unit_price DECIMAL(12,2) NOT NULL,
    scheme_id UUID REFERENCES insurance_schemes(id),
    effective_from DATE NOT NULL,
    effective_until DATE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_price_lists_hospital ON price_lists(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_lists_service ON price_lists(hospital_id, service_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_lists_scheme ON price_lists(scheme_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_lists_effective ON price_lists(hospital_id, effective_from, effective_until) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_price_lists_updated_at
BEFORE UPDATE ON price_lists
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE price_lists ENABLE ROW LEVEL SECURITY;
CREATE POLICY price_lists_isolation ON price_lists
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS price_lists CASCADE;
