-- name: CreateFileAttachment :one
INSERT INTO file_attachments (
  hospital_id, uploaded_by, entity_type, entity_id, file_name, file_path, mime_type, file_size_bytes, description, is_sensitive
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetFileAttachment :one
SELECT * FROM file_attachments
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListFileAttachmentsByEntity :many
SELECT * FROM file_attachments
WHERE entity_type = $1 AND entity_id = $2 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListFileAttachmentsByHospital :many
SELECT * FROM file_attachments
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteFileAttachment :exec
UPDATE file_attachments
SET deleted_at = now()
WHERE id = $1;
