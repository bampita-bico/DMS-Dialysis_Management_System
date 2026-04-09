-- +goose Up
CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    medication_id UUID NOT NULL REFERENCES medications(id),
    stock_id UUID REFERENCES pharmacy_stock(id),
    movement_type stock_movement_type NOT NULL,
    quantity INTEGER NOT NULL,
    quantity_before INTEGER NOT NULL,
    quantity_after INTEGER NOT NULL,
    unit_cost NUMERIC(10,2),
    total_cost NUMERIC(12,2),
    reference_type VARCHAR(100),
    reference_id UUID,
    performed_by UUID NOT NULL REFERENCES users(id),
    movement_date DATE NOT NULL DEFAULT CURRENT_DATE,
    movement_time TIME NOT NULL DEFAULT CURRENT_TIME,
    reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE stock_movements ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON stock_movements USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_stock_movements_medication ON stock_movements(medication_id);
CREATE INDEX idx_stock_movements_stock ON stock_movements(stock_id);
CREATE INDEX idx_stock_movements_type ON stock_movements(hospital_id, movement_type);
CREATE INDEX idx_stock_movements_date ON stock_movements(hospital_id, movement_date DESC);
CREATE INDEX idx_stock_movements_reference ON stock_movements(reference_type, reference_id);

CREATE TRIGGER trg_stock_movements_updated_at BEFORE UPDATE ON stock_movements FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS stock_movements;
