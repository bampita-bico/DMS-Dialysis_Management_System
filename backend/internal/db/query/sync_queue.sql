-- name: EnqueueSync :one
INSERT INTO sync_queue (
  hospital_id, user_id, entity_type, entity_id, operation, payload, priority
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPendingSyncItems :many
SELECT * FROM sync_queue
WHERE hospital_id = $1 AND status = 'pending' AND deleted_at IS NULL
ORDER BY priority DESC, created_at
LIMIT $2;

-- name: MarkSyncSynced :exec
UPDATE sync_queue
SET status = 'synced', synced_at = now()
WHERE id = $1;

-- name: MarkSyncFailed :exec
UPDATE sync_queue
SET status = 'failed', attempts = attempts + 1, last_attempt_at = now(), error_message = $2
WHERE id = $1;

-- name: RequeueFailedSyncs :exec
UPDATE sync_queue
SET status = 'pending', attempts = 0, error_message = NULL
WHERE status = 'failed' AND attempts < 3 AND deleted_at IS NULL;
