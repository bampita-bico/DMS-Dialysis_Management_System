-- name: CreateUser :one
INSERT INTO users (
  hospital_id, email, phone, password_hash, full_name, is_active
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE hospital_id = $1 AND email = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListActiveUsers :many
SELECT * FROM users
WHERE hospital_id = $1 AND is_active = true AND deleted_at IS NULL
ORDER BY full_name;

-- name: UpdateUser :one
UPDATE users
SET email = $2, phone = $3, full_name = $4, is_active = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, password_reset_at = now()
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = now()
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;
