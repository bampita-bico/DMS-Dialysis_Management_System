-- +goose Up
CREATE TABLE prescription_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    prescription_id UUID NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    medication_id UUID NOT NULL REFERENCES medications(id),
    dose VARCHAR(100) NOT NULL,
    frequency VARCHAR(100) NOT NULL,
    route medication_route NOT NULL,
    duration_days INTEGER,
    quantity_prescribed INTEGER NOT NULL,
    quantity_dispensed INTEGER DEFAULT 0,
    instructions TEXT,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE,
    is_prn BOOLEAN NOT NULL DEFAULT FALSE,
    prn_indication TEXT,
    is_stat BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE prescription_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON prescription_items USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_prescription_items_prescription ON prescription_items(prescription_id);
CREATE INDEX idx_prescription_items_medication ON prescription_items(medication_id);
CREATE INDEX idx_prescription_items_dates ON prescription_items(start_date, end_date);

CREATE TRIGGER trg_prescription_items_updated_at BEFORE UPDATE ON prescription_items FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS prescription_items;
