-- name: CreateDryWeightRecord :one
INSERT INTO dry_weight_records (
    hospital_id, patient_id, set_by, set_date, dry_weight_kg,
    assessment_method, clinical_indicators, bp_at_assessment_systolic,
    bp_at_assessment_diastolic, edema_present, edema_location,
    dyspnea_present, chest_xray_findings, effective_from, effective_until,
    is_current, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
RETURNING *;

-- name: GetDryWeightRecord :one
SELECT * FROM dry_weight_records WHERE id = $1 AND deleted_at IS NULL;

-- name: GetCurrentDryWeight :one
SELECT * FROM dry_weight_records
WHERE patient_id = $1
  AND is_current = TRUE
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListDryWeightRecordsByPatient :many
SELECT * FROM dry_weight_records
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY set_date DESC;

-- name: UpdateDryWeightRecord :one
UPDATE dry_weight_records
SET dry_weight_kg = $2, assessment_method = $3, clinical_indicators = $4,
    bp_at_assessment_systolic = $5, bp_at_assessment_diastolic = $6,
    edema_present = $7, edema_location = $8, dyspnea_present = $9,
    chest_xray_findings = $10, effective_until = $11, is_current = $12, notes = $13
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDryWeightRecord :exec
UPDATE dry_weight_records SET deleted_at = now() WHERE id = $1;
