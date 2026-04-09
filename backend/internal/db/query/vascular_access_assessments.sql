-- name: CreateVascularAccessAssessment :one
INSERT INTO vascular_access_assessments (
    hospital_id, access_id, patient_id, session_id, assessed_by, assessed_at,
    has_thrill, has_bruit, has_redness, has_swelling, has_discharge,
    has_bleeding, has_pain, appearance_normal, flow_rate_ml_min,
    venous_pressure_mmhg, arterial_pressure_mmhg, recirculation_percent,
    requires_intervention, intervention_type, intervention_urgency, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
RETURNING *;

-- name: GetVascularAccessAssessment :one
SELECT * FROM vascular_access_assessments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAssessmentsByAccess :many
SELECT * FROM vascular_access_assessments
WHERE access_id = $1 AND deleted_at IS NULL
ORDER BY assessed_at DESC;

-- name: ListAssessmentsBySession :many
SELECT * FROM vascular_access_assessments
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY assessed_at;

-- name: ListAssessmentsRequiringIntervention :many
SELECT * FROM vascular_access_assessments
WHERE hospital_id = $1
  AND requires_intervention = TRUE
  AND deleted_at IS NULL
ORDER BY assessed_at DESC;

-- name: DeleteVascularAccessAssessment :exec
UPDATE vascular_access_assessments SET deleted_at = now() WHERE id = $1;
