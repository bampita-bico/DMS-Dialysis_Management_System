-- name: CreateRole :one
INSERT INTO roles (
  hospital_id, name, description, permissions, is_system
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3, permissions = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteRole :exec
UPDATE roles
SET deleted_at = now()
WHERE id = $1;
