-- +goose Up
CREATE TABLE consumables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    category consumable_category NOT NULL,
    unit VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(255),
    model VARCHAR(100),
    is_reusable BOOLEAN NOT NULL DEFAULT FALSE,
    max_reuse_count INTEGER,
    min_stock_level INTEGER,
    cost_per_unit DECIMAL(10,2),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_consumables_hospital ON consumables(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_category ON consumables(hospital_id, category) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_active ON consumables(hospital_id, is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_consumables_updated_at
BEFORE UPDATE ON consumables
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE consumables ENABLE ROW LEVEL SECURITY;
CREATE POLICY consumables_isolation ON consumables
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS consumables CASCADE;
