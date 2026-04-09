-- name: CreateNotification :one
INSERT INTO notifications (
  hospital_id, user_id, type, title, message, priority, entity_type, entity_id, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetNotification :one
SELECT * FROM notifications
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUnreadNotifications :many
SELECT * FROM notifications
WHERE user_id = $1 AND read_at IS NULL AND deleted_at IS NULL
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY priority DESC, created_at DESC;

-- name: ListCriticalUnread :many
SELECT * FROM notifications
WHERE user_id = $1 AND priority = 'critical' AND read_at IS NULL AND deleted_at IS NULL
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC;

-- name: MarkNotificationRead :exec
UPDATE notifications
SET read_at = now()
WHERE id = $1;

-- name: MarkNotificationActioned :exec
UPDATE notifications
SET actioned_at = now()
WHERE id = $1;
