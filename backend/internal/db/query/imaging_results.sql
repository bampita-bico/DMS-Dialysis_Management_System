-- name: CreateImagingResult :one
INSERT INTO imaging_results (
    hospital_id, order_id, report_text, impression, recommendations,
    reported_by, report_date, report_time, image_count, image_urls,
    is_abnormal, is_critical, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetImagingResult :one
SELECT * FROM imaging_results WHERE id = $1 AND deleted_at IS NULL;

-- name: GetImagingResultByOrder :one
SELECT * FROM imaging_results WHERE order_id = $1 AND deleted_at IS NULL;

-- name: ListCriticalImagingResults :many
SELECT * FROM imaging_results
WHERE hospital_id = $1 AND is_critical = TRUE AND deleted_at IS NULL
ORDER BY report_date DESC, report_time DESC;

-- name: ListImagingResultsByDateRange :many
SELECT * FROM imaging_results
WHERE hospital_id = $1
  AND report_date >= $2
  AND report_date <= $3
  AND deleted_at IS NULL
ORDER BY report_date DESC, report_time DESC;

-- name: VerifyImagingResult :one
UPDATE imaging_results
SET verified_by = $2, verified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateImagingResult :one
UPDATE imaging_results
SET report_text = $2, impression = $3, recommendations = $4,
    is_abnormal = $5, is_critical = $6, notes = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteImagingResult :exec
UPDATE imaging_results SET deleted_at = now() WHERE id = $1;
