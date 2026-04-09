-- +goose Up
CREATE TYPE sync_operation AS ENUM ('CREATE', 'UPDATE', 'DELETE');
CREATE TYPE sync_status AS ENUM ('pending', 'syncing', 'synced', 'failed', 'conflict');

CREATE TABLE sync_queue (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id UUID NOT NULL REFERENCES hospitals(id),
  user_id UUID NOT NULL REFERENCES users(id),
  entity_type VARCHAR(100) NOT NULL,
  entity_id UUID NOT NULL,
  operation sync_operation NOT NULL,
  payload JSONB NOT NULL,
  priority INTEGER NOT NULL DEFAULT 5 CHECK (priority BETWEEN 1 AND 10),
  status sync_status NOT NULL DEFAULT 'pending',
  attempts INTEGER NOT NULL DEFAULT 0,
  last_attempt_at TIMESTAMPTZ,
  synced_at TIMESTAMPTZ,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE sync_queue ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sync_queue
  USING (hospital_id = current_setting('app.current_hospital_id', true)::UUID);

CREATE INDEX idx_sync_queue_status ON sync_queue(status, priority DESC);
CREATE INDEX idx_sync_queue_entity ON sync_queue(entity_type, entity_id);

CREATE TRIGGER trg_sync_queue_updated_at
BEFORE UPDATE ON sync_queue
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS sync_queue;
DROP TYPE IF EXISTS sync_status;
DROP TYPE IF EXISTS sync_operation;
