-- +goose Up
CREATE TABLE quality_indicators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    indicator_name VARCHAR(255) NOT NULL,
    indicator_code VARCHAR(50) NOT NULL,
    indicator_category VARCHAR(100),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    numerator INTEGER,
    denominator INTEGER,
    value DECIMAL(10,4) NOT NULL,
    unit VARCHAR(50),
    target_value DECIMAL(10,4),
    benchmark_value DECIMAL(10,4),
    meets_target BOOLEAN,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    calculated_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_quality_indicators_hospital ON quality_indicators(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_quality_indicators_code ON quality_indicators(hospital_id, indicator_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_quality_indicators_period ON quality_indicators(hospital_id, period_start, period_end) WHERE deleted_at IS NULL;
CREATE INDEX idx_quality_indicators_category ON quality_indicators(hospital_id, indicator_category) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_quality_indicators_updated_at
BEFORE UPDATE ON quality_indicators
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE quality_indicators ENABLE ROW LEVEL SECURITY;
CREATE POLICY quality_indicators_isolation ON quality_indicators
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS quality_indicators CASCADE;
