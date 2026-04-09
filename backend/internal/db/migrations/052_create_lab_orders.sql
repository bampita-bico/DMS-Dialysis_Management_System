-- +goose Up
CREATE TABLE lab_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    ordered_by UUID NOT NULL REFERENCES users(id),
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    order_time TIME NOT NULL DEFAULT CURRENT_TIME,
    priority lab_priority NOT NULL DEFAULT 'routine',
    status lab_status NOT NULL DEFAULT 'pending',
    clinical_notes TEXT,
    diagnosis_code VARCHAR(20),
    cancelled_by UUID REFERENCES users(id),
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE lab_orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_orders USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_orders_hospital ON lab_orders(hospital_id);
CREATE INDEX idx_lab_orders_patient ON lab_orders(patient_id);
CREATE INDEX idx_lab_orders_session ON lab_orders(session_id);
CREATE INDEX idx_lab_orders_status ON lab_orders(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_lab_orders_date ON lab_orders(hospital_id, order_date DESC);

CREATE TRIGGER trg_lab_orders_updated_at BEFORE UPDATE ON lab_orders FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_orders;
