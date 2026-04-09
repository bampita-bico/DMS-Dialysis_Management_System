-- +goose Up
CREATE TABLE hospitals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  short_code VARCHAR(20) NOT NULL UNIQUE,
  tier VARCHAR(50) NOT NULL CHECK (tier IN ('national','regional','district','private')),
  region VARCHAR(100) NOT NULL,
  country VARCHAR(100) NOT NULL DEFAULT 'Uganda',
  address TEXT,
  phone VARCHAR(50),
  email VARCHAR(255),
  license_no VARCHAR(100),
  license_expiry DATE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  settings JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_hospitals_country ON hospitals(country);
CREATE INDEX idx_hospitals_deleted ON hospitals(deleted_at);

CREATE TRIGGER trg_hospitals_updated_at
BEFORE UPDATE ON hospitals
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS hospitals;
