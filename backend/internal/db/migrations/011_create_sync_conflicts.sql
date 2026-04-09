-- +goose Up
CREATE TYPE conflict_resolution AS ENUM ('pending','server_wins','client_wins','manual');

CREATE TABLE sync_conflicts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id UUID NOT NULL REFERENCES hospitals(id),
  queue_id UUID NOT NULL REFERENCES sync_queue(id),
  entity_type VARCHAR(100) NOT NULL,
  entity_id UUID NOT NULL,
  local_data JSONB NOT NULL,
  server_data JSONB NOT NULL,
  resolution conflict_resolution NOT NULL DEFAULT 'pending',
  resolved_by UUID REFERENCES users(id),
  resolved_at TIMESTAMPTZ,
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE sync_conflicts ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sync_conflicts
  USING (hospital_id = current_setting('app.current_hospital_id', true)::UUID);

CREATE INDEX idx_sync_conflicts_pending ON sync_conflicts(resolution)
  WHERE resolution = 'pending';

CREATE TRIGGER trg_sync_conflicts_updated_at
BEFORE UPDATE ON sync_conflicts
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS sync_conflicts;
DROP TYPE IF EXISTS conflict_resolution;
