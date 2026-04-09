-- name: CreateQualityIndicator :one
INSERT INTO quality_indicators (
    hospital_id, indicator_name, indicator_code, indicator_category,
    period_start, period_end, numerator, denominator, value, unit,
    target_value, benchmark_value, meets_target, calculated_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetQualityIndicator :one
SELECT * FROM quality_indicators WHERE id = $1 AND deleted_at IS NULL;

-- name: ListIndicatorsByHospital :many
SELECT * FROM quality_indicators
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY period_start DESC, indicator_name;

-- name: ListIndicatorsByCode :many
SELECT * FROM quality_indicators
WHERE hospital_id = $1 AND indicator_code = $2 AND deleted_at IS NULL
ORDER BY period_start DESC;

-- name: ListIndicatorsByPeriod :many
SELECT * FROM quality_indicators
WHERE hospital_id = $1
  AND period_start >= $2
  AND period_end <= $3
  AND deleted_at IS NULL
ORDER BY indicator_category, indicator_name;

-- name: ListIndicatorsBelowTarget :many
SELECT * FROM quality_indicators
WHERE hospital_id = $1
  AND meets_target = FALSE
  AND period_start >= $2
  AND deleted_at IS NULL
ORDER BY indicator_category, indicator_name;

-- name: GetLatestIndicatorByCode :one
SELECT * FROM quality_indicators
WHERE hospital_id = $1 AND indicator_code = $2 AND deleted_at IS NULL
ORDER BY period_end DESC
LIMIT 1;

-- name: DeleteQualityIndicator :exec
UPDATE quality_indicators SET deleted_at = now() WHERE id = $1;
