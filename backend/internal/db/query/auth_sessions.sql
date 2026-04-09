-- name: CreateAuthSession :one
INSERT INTO auth_sessions (
  hospital_id, user_id, token_hash, device_info, ip_address, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAuthSession :one
SELECT * FROM auth_sessions
WHERE token_hash = $1 AND expires_at > now() AND revoked_at IS NULL AND deleted_at IS NULL
LIMIT 1;

-- name: ListUserSessions :many
SELECT * FROM auth_sessions
WHERE user_id = $1 AND expires_at > now() AND revoked_at IS NULL AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeSession :exec
UPDATE auth_sessions
SET revoked_at = now()
WHERE token_hash = $1;

-- name: RevokeAllUserSessions :exec
UPDATE auth_sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: CleanupExpiredSessions :exec
UPDATE auth_sessions
SET deleted_at = now()
WHERE expires_at < now() - INTERVAL '30 days' AND deleted_at IS NULL;
