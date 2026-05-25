-- +goose Up
CREATE TABLE IF NOT EXISTS user_login_aliases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(80) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_login_aliases_username_active
ON user_login_aliases (lower(username))
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_login_aliases_user
ON user_login_aliases (user_id)
WHERE deleted_at IS NULL;

DROP TRIGGER IF EXISTS trg_user_login_aliases_updated_at ON user_login_aliases;

CREATE TRIGGER trg_user_login_aliases_updated_at
BEFORE UPDATE ON user_login_aliases
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS user_login_aliases CASCADE;
