-- +goose Up
CREATE TABLE training_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    training_name VARCHAR(255) NOT NULL,
    training_category VARCHAR(100),
    training_provider VARCHAR(255),
    training_start_date DATE NOT NULL,
    training_end_date DATE,
    duration_hours DECIMAL(5,2),
    completed_at TIMESTAMPTZ,
    score DECIMAL(5,2),
    pass_status VARCHAR(50),
    certificate_url TEXT,
    certificate_number VARCHAR(100),
    certificate_expiry_date DATE,
    cpd_points DECIMAL(5,2),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_training_records_staff ON training_records(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_training_records_hospital ON training_records(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_training_records_category ON training_records(hospital_id, training_category) WHERE deleted_at IS NULL;
CREATE INDEX idx_training_records_expiry ON training_records(hospital_id, certificate_expiry_date) WHERE deleted_at IS NULL AND certificate_expiry_date IS NOT NULL;

CREATE TRIGGER trg_training_records_updated_at
BEFORE UPDATE ON training_records
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE training_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY training_records_isolation ON training_records
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS training_records CASCADE;
