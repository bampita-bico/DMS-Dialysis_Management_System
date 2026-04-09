-- +goose Up
CREATE TABLE lab_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    order_id UUID NOT NULL REFERENCES lab_orders(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES lab_test_catalog(id),
    specimen_type specimen_type NOT NULL DEFAULT 'serum',
    specimen_collected_by UUID REFERENCES users(id),
    specimen_collected_at TIMESTAMPTZ,
    specimen_barcode VARCHAR(100),
    specimen_quality VARCHAR(50),
    specimen_rejected BOOLEAN NOT NULL DEFAULT FALSE,
    rejection_reason TEXT,
    status lab_status NOT NULL DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE lab_order_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_order_items USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_order_items_order ON lab_order_items(order_id);
CREATE INDEX idx_lab_order_items_test ON lab_order_items(test_id);
CREATE INDEX idx_lab_order_items_status ON lab_order_items(hospital_id, status);
CREATE INDEX idx_lab_order_items_barcode ON lab_order_items(hospital_id, specimen_barcode) WHERE specimen_barcode IS NOT NULL;

CREATE TRIGGER trg_lab_order_items_updated_at BEFORE UPDATE ON lab_order_items FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_order_items;
