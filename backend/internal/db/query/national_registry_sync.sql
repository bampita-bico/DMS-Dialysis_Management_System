-- name: CreateRegistrySync :one
INSERT INTO national_registry_sync (
    hospital_id, patient_id, registry_name, registry_id, sync_type,
    payload, synced_by, sync_status, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetRegistrySync :one
SELECT * FROM national_registry_sync WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSyncsByPatient :many
SELECT * FROM national_registry_sync
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY synced_at DESC;

-- name: ListSyncsByRegistry :many
SELECT * FROM national_registry_sync
WHERE hospital_id = $1 AND registry_name = $2 AND deleted_at IS NULL
ORDER BY synced_at DESC;

-- name: ListPendingSyncs :many
SELECT * FROM national_registry_sync
WHERE hospital_id = $1 AND sync_status = 'pending' AND deleted_at IS NULL
ORDER BY synced_at ASC;

-- name: ListFailedSyncs :many
SELECT * FROM national_registry_sync
WHERE hospital_id = $1
  AND sync_status = 'failed'
  AND retry_count < 3
  AND deleted_at IS NULL
ORDER BY synced_at ASC;

-- name: UpdateSyncStatus :one
UPDATE national_registry_sync
SET sync_status = $2, error_message = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: IncrementRetryCount :one
UPDATE national_registry_sync
SET retry_count = retry_count + 1, last_retry_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteRegistrySync :exec
UPDATE national_registry_sync SET deleted_at = now() WHERE id = $1;
