-- name: CreateMicrobiologyResult :one
INSERT INTO microbiology_results (
    hospital_id, order_item_id, culture_date, culture_time, growth_detected,
    organism, organism_count, sensitivity, antibiotic_recommendations,
    result_date, result_time, reported_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetMicrobiologyResult :one
SELECT * FROM microbiology_results WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMicrobiologyResultByOrderItem :one
SELECT * FROM microbiology_results WHERE order_item_id = $1 AND deleted_at IS NULL;

-- name: ListMicrobiologyResultsWithGrowth :many
SELECT * FROM microbiology_results
WHERE hospital_id = $1 AND growth_detected = TRUE AND deleted_at IS NULL
ORDER BY result_date DESC, result_time DESC;

-- name: ListMicrobiologyResultsByDateRange :many
SELECT * FROM microbiology_results
WHERE hospital_id = $1
  AND result_date >= $2
  AND result_date <= $3
  AND deleted_at IS NULL
ORDER BY result_date DESC, result_time DESC;

-- name: VerifyMicrobiologyResult :one
UPDATE microbiology_results
SET verified_by = $2, verified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateMicrobiologyResult :one
UPDATE microbiology_results
SET growth_detected = $2, organism = $3, organism_count = $4,
    sensitivity = $5, antibiotic_recommendations = $6, notes = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMicrobiologyResult :exec
UPDATE microbiology_results SET deleted_at = now() WHERE id = $1;
