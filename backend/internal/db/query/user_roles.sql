-- name: AssignUserRole :one
INSERT INTO user_roles (
  hospital_id, user_id, role_id, department_id, assigned_by, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserRoles :many
SELECT ur.*, r.name as role_name, r.permissions
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1 AND ur.deleted_at IS NULL AND r.deleted_at IS NULL;

-- name: GetUsersInRole :many
SELECT ur.*, u.full_name, u.email
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE ur.role_id = $1 AND ur.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY u.full_name;

-- name: RevokeUserRole :exec
UPDATE user_roles
SET deleted_at = now()
WHERE user_id = $1 AND role_id = $2 AND department_id = $3;
