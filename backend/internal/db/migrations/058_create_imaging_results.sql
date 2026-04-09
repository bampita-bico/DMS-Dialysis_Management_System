-- +goose Up
CREATE TABLE imaging_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    order_id UUID NOT NULL REFERENCES imaging_orders(id) ON DELETE CASCADE,
    report_text TEXT NOT NULL,
    impression TEXT,
    recommendations TEXT,
    reported_by UUID NOT NULL REFERENCES users(id),
    report_date DATE NOT NULL DEFAULT CURRENT_DATE,
    report_time TIME NOT NULL DEFAULT CURRENT_TIME,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    image_count INTEGER DEFAULT 0,
    image_urls JSONB DEFAULT '[]',
    is_abnormal BOOLEAN NOT NULL DEFAULT FALSE,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE imaging_results ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON imaging_results USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_imaging_results_order ON imaging_results(order_id);
CREATE INDEX idx_imaging_results_critical ON imaging_results(hospital_id) WHERE is_critical = TRUE;
CREATE INDEX idx_imaging_results_date ON imaging_results(hospital_id, report_date DESC);

CREATE TRIGGER trg_imaging_results_updated_at BEFORE UPDATE ON imaging_results FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS imaging_results;
