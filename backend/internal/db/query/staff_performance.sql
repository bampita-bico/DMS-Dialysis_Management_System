-- name: CreateStaffPerformance :one
INSERT INTO staff_performance (
    hospital_id, staff_id, review_period_start, review_period_end,
    review_date, appraised_by, overall_score, technical_competence_score,
    communication_score, teamwork_score, punctuality_score, patient_care_score,
    strengths, areas_for_improvement, goals_next_period,
    training_recommendations, promotion_recommended, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
RETURNING *;

-- name: GetStaffPerformance :one
SELECT * FROM staff_performance WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPerformanceByStaff :many
SELECT * FROM staff_performance
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY review_period_start DESC;

-- name: GetLatestPerformance :one
SELECT * FROM staff_performance
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY review_period_end DESC
LIMIT 1;

-- name: ListPerformanceByPeriod :many
SELECT * FROM staff_performance
WHERE hospital_id = $1
  AND review_period_start >= $2
  AND review_period_end <= $3
  AND deleted_at IS NULL
ORDER BY review_period_start;

-- name: ListTopPerformers :many
SELECT * FROM staff_performance
WHERE hospital_id = $1
  AND review_period_start >= $2
  AND overall_score IS NOT NULL
  AND deleted_at IS NULL
ORDER BY overall_score DESC
LIMIT $3;

-- name: DeleteStaffPerformance :exec
UPDATE staff_performance SET deleted_at = now() WHERE id = $1;
