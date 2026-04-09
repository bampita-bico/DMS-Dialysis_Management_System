-- +goose Up
CREATE TABLE water_treatment_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    test_date DATE NOT NULL,
    test_time TIME NOT NULL,
    tested_by UUID NOT NULL REFERENCES users(id),
    sample_location VARCHAR(200) NOT NULL,
    bacterial_count_cfu_ml NUMERIC(8,2),
    endotoxin_level_eu_ml NUMERIC(6,3),
    chlorine_level_ppm NUMERIC(5,2),
    chloramine_level_ppm NUMERIC(5,2),
    ph_level NUMERIC(4,2),
    conductivity_us_cm NUMERIC(8,2),
    hardness_mg_l NUMERIC(6,2),
    bacteria_result water_test_result NOT NULL DEFAULT 'pending',
    endotoxin_result water_test_result NOT NULL DEFAULT 'pending',
    chlorine_result water_test_result NOT NULL DEFAULT 'pending',
    overall_result water_test_result NOT NULL DEFAULT 'pending',
    out_of_spec_parameters TEXT,
    corrective_action_taken TEXT,
    corrective_action_by UUID REFERENCES users(id),
    retest_required BOOLEAN NOT NULL DEFAULT FALSE,
    retest_date DATE,
    systems_shut_down BOOLEAN NOT NULL DEFAULT FALSE,
    shutdown_time TIMESTAMPTZ,
    resumed_time TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE water_treatment_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON water_treatment_logs USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_water_logs_hospital_date ON water_treatment_logs(hospital_id, test_date DESC);
CREATE INDEX idx_water_logs_failed ON water_treatment_logs(hospital_id) WHERE overall_result = 'fail';
CREATE INDEX idx_water_logs_pending ON water_treatment_logs(hospital_id) WHERE overall_result = 'pending';

CREATE TRIGGER trg_water_treatment_logs_updated_at BEFORE UPDATE ON water_treatment_logs FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS water_treatment_logs;
