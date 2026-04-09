-- name: CreateAdequacyAssessment :one
INSERT INTO adequacy_assessments (
    hospital_id, patient_id, session_id, assessed_by, assessment_date,
    kt_v, kt_v_method, urr_percent, pre_bun_mg_dl, post_bun_mg_dl,
    pre_creatinine_mg_dl, post_creatinine_mg_dl, dialysis_duration_mins,
    blood_flow_rate, dialyzer_clearance, body_water_volume_liters,
    is_adequate, recommendations, next_assessment_date, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING *;

-- name: GetAdequacyAssessment :one
SELECT * FROM adequacy_assessments WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAdequacyAssessmentsByPatient :many
SELECT * FROM adequacy_assessments
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY assessment_date DESC;

-- name: ListInadequateAssessments :many
SELECT * FROM adequacy_assessments
WHERE hospital_id = $1
  AND is_adequate = FALSE
  AND deleted_at IS NULL
ORDER BY assessment_date DESC;

-- name: UpdateAdequacyAssessment :one
UPDATE adequacy_assessments
SET kt_v = $2, kt_v_method = $3, urr_percent = $4, pre_bun_mg_dl = $5,
    post_bun_mg_dl = $6, pre_creatinine_mg_dl = $7, post_creatinine_mg_dl = $8,
    dialysis_duration_mins = $9, blood_flow_rate = $10, dialyzer_clearance = $11,
    body_water_volume_liters = $12, is_adequate = $13, recommendations = $14,
    next_assessment_date = $15, notes = $16
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAdequacyAssessment :exec
UPDATE adequacy_assessments SET deleted_at = now() WHERE id = $1;
