-- +goose Up
CREATE TABLE consumables_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id) ON DELETE RESTRICT,
    consumable_id UUID NOT NULL REFERENCES consumables(id) ON DELETE RESTRICT,
    inventory_id UUID REFERENCES consumables_inventory(id),
    quantity_used INTEGER NOT NULL,
    reuse_number INTEGER,
    recorded_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_consumables_usage_session ON consumables_usage(session_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_usage_consumable ON consumables_usage(consumable_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_usage_hospital ON consumables_usage(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_usage_inventory ON consumables_usage(inventory_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_consumables_usage_updated_at
BEFORE UPDATE ON consumables_usage
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE consumables_usage ENABLE ROW LEVEL SECURITY;
CREATE POLICY consumables_usage_isolation ON consumables_usage
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS consumables_usage CASCADE;
