-- +goose Up
CREATE TABLE hospital_settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
  key VARCHAR(100) NOT NULL,
  value TEXT NOT NULL,
  data_type VARCHAR(20) NOT NULL CHECK (data_type IN ('string','integer','boolean','json')),
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  UNIQUE(hospital_id, key)
);

ALTER TABLE hospital_settings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON hospital_settings
  USING (hospital_id = current_setting('app.current_hospital_id', true)::UUID);

CREATE INDEX idx_hospital_settings_hospital ON hospital_settings(hospital_id);

CREATE TRIGGER trg_hospital_settings_updated_at
BEFORE UPDATE ON hospital_settings
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS hospital_settings;
