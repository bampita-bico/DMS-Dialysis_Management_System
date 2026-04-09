-- +goose Up
CREATE TABLE imaging_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    ordered_by UUID NOT NULL REFERENCES users(id),
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    order_time TIME NOT NULL DEFAULT CURRENT_TIME,
    modality imaging_modality NOT NULL,
    body_part VARCHAR(200) NOT NULL,
    laterality VARCHAR(20),
    clinical_indication TEXT NOT NULL,
    priority lab_priority NOT NULL DEFAULT 'routine',
    status lab_status NOT NULL DEFAULT 'pending',
    scheduled_date DATE,
    scheduled_time TIME,
    performed_at TIMESTAMPTZ,
    performed_by VARCHAR(200),
    cancelled_by UUID REFERENCES users(id),
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE imaging_orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON imaging_orders USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_imaging_orders_hospital ON imaging_orders(hospital_id);
CREATE INDEX idx_imaging_orders_patient ON imaging_orders(patient_id);
CREATE INDEX idx_imaging_orders_session ON imaging_orders(session_id);
CREATE INDEX idx_imaging_orders_status ON imaging_orders(hospital_id, status);
CREATE INDEX idx_imaging_orders_modality ON imaging_orders(hospital_id, modality);
CREATE INDEX idx_imaging_orders_date ON imaging_orders(hospital_id, order_date DESC);

CREATE TRIGGER trg_imaging_orders_updated_at BEFORE UPDATE ON imaging_orders FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS imaging_orders;
