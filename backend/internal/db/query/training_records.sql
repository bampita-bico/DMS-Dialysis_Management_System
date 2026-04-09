-- name: CreateTrainingRecord :one
INSERT INTO training_records (
    hospital_id, staff_id, training_name, training_category,
    training_provider, training_start_date, training_end_date,
    duration_hours, cpd_points, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetTrainingRecord :one
SELECT * FROM training_records WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTrainingByStaff :many
SELECT * FROM training_records
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY training_start_date DESC;

-- name: ListTrainingByCategory :many
SELECT * FROM training_records
WHERE hospital_id = $1 AND training_category = $2 AND deleted_at IS NULL
ORDER BY training_start_date DESC;

-- name: CompleteTraining :one
UPDATE training_records
SET completed_at = now(), score = $2, pass_status = $3,
    certificate_url = $4, certificate_number = $5, certificate_expiry_date = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ListExpiringCertificates :many
SELECT * FROM training_records
WHERE hospital_id = $1
  AND certificate_expiry_date <= $2
  AND deleted_at IS NULL
ORDER BY certificate_expiry_date ASC;

-- name: GetTotalCPDPoints :one
SELECT COALESCE(SUM(cpd_points), 0) AS total_cpd_points
FROM training_records
WHERE staff_id = $1
  AND training_start_date BETWEEN $2 AND $3
  AND deleted_at IS NULL;

-- name: DeleteTrainingRecord :exec
UPDATE training_records SET deleted_at = now() WHERE id = $1;
