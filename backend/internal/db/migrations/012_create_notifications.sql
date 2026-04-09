-- +goose Up
CREATE TYPE notification_priority AS ENUM ('low','medium','high','critical');

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hospital_id UUID NOT NULL REFERENCES hospitals(id),
  user_id UUID NOT NULL REFERENCES users(id),
  type VARCHAR(100) NOT NULL,
  title VARCHAR(255) NOT NULL,
  message TEXT NOT NULL,
  priority notification_priority NOT NULL DEFAULT 'medium',
  entity_type VARCHAR(100),
  entity_id UUID,
  read_at TIMESTAMPTZ,
  actioned_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON notifications
  USING (hospital_id = current_setting('app.current_hospital_id', true)::UUID);

CREATE INDEX idx_notifications_user_unread ON notifications(user_id, read_at)
  WHERE read_at IS NULL AND deleted_at IS NULL;
CREATE INDEX idx_notifications_priority ON notifications(priority, created_at DESC);

CREATE TRIGGER trg_notifications_updated_at
BEFORE UPDATE ON notifications
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_priority;
