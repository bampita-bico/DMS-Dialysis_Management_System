-- +goose Up
CREATE TABLE consumables_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    consumable_id UUID NOT NULL REFERENCES consumables(id) ON DELETE RESTRICT,
    batch_number VARCHAR(100),
    quantity_current INTEGER NOT NULL DEFAULT 0,
    quantity_reserved INTEGER NOT NULL DEFAULT 0,
    quantity_available INTEGER NOT NULL DEFAULT 0,
    unit_cost DECIMAL(10,2),
    expiry_date DATE,
    received_date DATE,
    supplier_name VARCHAR(255),
    storage_location VARCHAR(255),
    is_low_stock BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_consumables_inventory_consumable ON consumables_inventory(consumable_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_inventory_hospital ON consumables_inventory(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_consumables_inventory_expiry ON consumables_inventory(hospital_id, expiry_date) WHERE deleted_at IS NULL AND quantity_current > 0;
CREATE INDEX idx_consumables_inventory_low_stock ON consumables_inventory(hospital_id, is_low_stock) WHERE deleted_at IS NULL AND is_low_stock = TRUE;

CREATE TRIGGER trg_consumables_inventory_updated_at
BEFORE UPDATE ON consumables_inventory
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE consumables_inventory ENABLE ROW LEVEL SECURITY;
CREATE POLICY consumables_inventory_isolation ON consumables_inventory
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS consumables_inventory CASCADE;
