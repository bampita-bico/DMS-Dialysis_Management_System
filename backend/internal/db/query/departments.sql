-- name: CreateDepartment :one
INSERT INTO departments (
  hospital_id, name, code, description, head_user_id, is_active
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetDepartment :one
SELECT * FROM departments
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListDepartments :many
SELECT * FROM departments
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: UpdateDepartment :one
UPDATE departments
SET name = $2, code = $3, description = $4, head_user_id = $5, is_active = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteDepartment :exec
UPDATE departments
SET deleted_at = now()
WHERE id = $1;
