-- name: CreateClinicalOutcome :one
INSERT INTO clinical_outcomes (
    hospital_id, patient_id, assessment_date, period_start, period_end,
    hemoglobin, hemoglobin_target_min, hemoglobin_target_max, kt_v, kt_v_target,
    urr, systolic_bp_avg, diastolic_bp_avg, bp_controlled, weight_gain_percent,
    albumin, phosphate, calcium, pth, quality_of_life_score, functional_status,
    adverse_events_count, hospitalizations_count, missed_sessions_count,
    outcome_trend, assessed_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
RETURNING *;

-- name: GetClinicalOutcome :one
SELECT * FROM clinical_outcomes WHERE id = $1 AND deleted_at IS NULL;

-- name: ListOutcomesByPatient :many
SELECT * FROM clinical_outcomes
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY assessment_date DESC;

-- name: GetLatestOutcome :one
SELECT * FROM clinical_outcomes
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY assessment_date DESC
LIMIT 1;

-- name: ListOutcomesByTrend :many
SELECT * FROM clinical_outcomes
WHERE hospital_id = $1 AND outcome_trend = $2 AND deleted_at IS NULL
ORDER BY assessment_date DESC;

-- name: ListDecliningPatients :many
SELECT * FROM clinical_outcomes
WHERE hospital_id = $1
  AND outcome_trend IN ('declining', 'critical')
  AND assessment_date >= $2
  AND deleted_at IS NULL
ORDER BY outcome_trend DESC, assessment_date DESC;

-- name: DeleteClinicalOutcome :exec
UPDATE clinical_outcomes SET deleted_at = now() WHERE id = $1;
