-- +goose Up
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    price_list_id UUID REFERENCES price_lists(id),
    service_name VARCHAR(255) NOT NULL,
    service_code VARCHAR(50),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    is_claimed BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_invoice_items_invoice ON invoice_items(invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_items_hospital ON invoice_items(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_items_price_list ON invoice_items(price_list_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_items_claimed ON invoice_items(hospital_id, is_claimed) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_invoice_items_updated_at
BEFORE UPDATE ON invoice_items
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE invoice_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY invoice_items_isolation ON invoice_items
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS invoice_items CASCADE;
