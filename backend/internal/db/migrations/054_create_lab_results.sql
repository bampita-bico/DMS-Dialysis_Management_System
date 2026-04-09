-- +goose Up
CREATE TABLE lab_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    order_item_id UUID NOT NULL REFERENCES lab_order_items(id) ON DELETE CASCADE,
    value_text TEXT,
    value_numeric NUMERIC(15,4),
    unit VARCHAR(50),
    reference_range TEXT,
    is_abnormal BOOLEAN NOT NULL DEFAULT FALSE,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    status result_status NOT NULL DEFAULT 'pending',
    result_date DATE NOT NULL DEFAULT CURRENT_DATE,
    result_time TIME NOT NULL DEFAULT CURRENT_TIME,
    entered_by UUID NOT NULL REFERENCES users(id),
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    source VARCHAR(50) DEFAULT 'internal',
    external_lab_name VARCHAR(200),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE lab_results ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_results USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_results_order_item ON lab_results(order_item_id);
CREATE INDEX idx_lab_results_status ON lab_results(hospital_id, status);
CREATE INDEX idx_lab_results_critical ON lab_results(hospital_id) WHERE is_critical = TRUE;
CREATE INDEX idx_lab_results_date ON lab_results(hospital_id, result_date DESC);

CREATE TRIGGER trg_lab_results_updated_at BEFORE UPDATE ON lab_results FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_results;
