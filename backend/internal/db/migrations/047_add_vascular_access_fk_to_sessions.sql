-- +goose Up
ALTER TABLE dialysis_sessions
ADD CONSTRAINT fk_sessions_access FOREIGN KEY (access_id) REFERENCES vascular_access(id);

CREATE INDEX idx_sessions_access ON dialysis_sessions(access_id);

-- +goose Down
DROP INDEX IF EXISTS idx_sessions_access;
ALTER TABLE dialysis_sessions DROP CONSTRAINT IF EXISTS fk_sessions_access;
