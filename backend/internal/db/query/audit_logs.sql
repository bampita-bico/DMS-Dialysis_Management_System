-- name: CreateAuditLog :one
INSERT INTO audit_logs (
  hospital_id, user_id, action, table_name, record_id, old_data, new_data, ip_address, user_agent
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetAuditLog :one
SELECT * FROM audit_logs
WHERE id = $1
LIMIT 1;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
WHERE hospital_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByUser :many
SELECT * FROM audit_logs
WHERE hospital_id = $1 AND user_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAuditLogsByTable :many
SELECT * FROM audit_logs
WHERE hospital_id = $1 AND table_name = $2 AND record_id = $3
ORDER BY created_at DESC;
