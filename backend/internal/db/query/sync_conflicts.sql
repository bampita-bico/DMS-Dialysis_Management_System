-- name: CreateSyncConflict :one
INSERT INTO sync_conflicts (
  hospital_id, queue_id, entity_type, entity_id, local_data, server_data
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSyncConflict :one
SELECT * FROM sync_conflicts
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListPendingConflicts :many
SELECT * FROM sync_conflicts
WHERE hospital_id = $1 AND resolution = 'pending' AND deleted_at IS NULL
ORDER BY created_at;

-- name: ResolveConflict :exec
UPDATE sync_conflicts
SET resolution = $2, resolved_by = $3, resolved_at = now(), notes = $4
WHERE id = $1;
