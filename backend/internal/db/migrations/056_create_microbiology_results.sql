-- +goose Up
CREATE TABLE microbiology_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    order_item_id UUID NOT NULL REFERENCES lab_order_items(id) ON DELETE CASCADE,
    culture_date DATE NOT NULL,
    culture_time TIME NOT NULL,
    growth_detected BOOLEAN NOT NULL DEFAULT FALSE,
    organism VARCHAR(255),
    organism_count VARCHAR(50),
    sensitivity JSONB DEFAULT '{}',
    antibiotic_recommendations TEXT,
    result_date DATE NOT NULL DEFAULT CURRENT_DATE,
    result_time TIME NOT NULL DEFAULT CURRENT_TIME,
    reported_by UUID NOT NULL REFERENCES users(id),
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE microbiology_results ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON microbiology_results USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_microbiology_results_order_item ON microbiology_results(order_item_id);
CREATE INDEX idx_microbiology_results_growth ON microbiology_results(hospital_id) WHERE growth_detected = TRUE;
CREATE INDEX idx_microbiology_results_date ON microbiology_results(hospital_id, result_date DESC);

CREATE TRIGGER trg_microbiology_results_updated_at BEFORE UPDATE ON microbiology_results FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS microbiology_results;
