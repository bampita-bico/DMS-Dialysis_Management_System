-- +goose Up
CREATE TABLE pharmacy_stock (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    medication_id UUID NOT NULL REFERENCES medications(id),
    batch_number VARCHAR(100) NOT NULL,
    quantity_current INTEGER NOT NULL DEFAULT 0,
    quantity_reserved INTEGER NOT NULL DEFAULT 0,
    quantity_available INTEGER NOT NULL DEFAULT 0,
    unit_cost NUMERIC(10,2),
    expiry_date DATE NOT NULL,
    received_date DATE NOT NULL DEFAULT CURRENT_DATE,
    supplier_name VARCHAR(255),
    storage_location VARCHAR(200),
    is_low_stock BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE pharmacy_stock ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON pharmacy_stock USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_pharmacy_stock_medication ON pharmacy_stock(medication_id);
CREATE INDEX idx_pharmacy_stock_batch ON pharmacy_stock(hospital_id, batch_number);
CREATE INDEX idx_pharmacy_stock_expiry ON pharmacy_stock(hospital_id, expiry_date);
CREATE INDEX idx_pharmacy_stock_low ON pharmacy_stock(hospital_id) WHERE is_low_stock = TRUE;

CREATE TRIGGER trg_pharmacy_stock_updated_at BEFORE UPDATE ON pharmacy_stock FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS pharmacy_stock;
