-- name: CreateHospitalization :one
INSERT INTO hospitalizations (
    hospital_id, patient_id, admission_date, admission_time, discharge_date,
    discharge_time, length_of_stay_days, admission_reason, admission_diagnosis,
    icd10_codes, admitting_facility, ward_name, dialysis_related, access_related,
    infection_related, treatment_given, procedures_performed, outcome,
    discharge_destination, follow_up_required, follow_up_date, recorded_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
RETURNING *;

-- name: GetHospitalization :one
SELECT * FROM hospitalizations WHERE id = $1 AND deleted_at IS NULL;

-- name: ListHospitalizationsByPatient :many
SELECT * FROM hospitalizations
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY admission_date DESC;

-- name: ListHospitalizationsByPeriod :many
SELECT * FROM hospitalizations
WHERE hospital_id = $1
  AND admission_date BETWEEN $2 AND $3
  AND deleted_at IS NULL
ORDER BY admission_date DESC;

-- name: ListDialysisRelatedHospitalizations :many
SELECT * FROM hospitalizations
WHERE hospital_id = $1
  AND dialysis_related = TRUE
  AND deleted_at IS NULL
ORDER BY admission_date DESC;

-- name: ListAccessRelatedHospitalizations :many
SELECT * FROM hospitalizations
WHERE hospital_id = $1
  AND access_related = TRUE
  AND deleted_at IS NULL
ORDER BY admission_date DESC;

-- name: UpdateHospitalizationDischarge :one
UPDATE hospitalizations
SET discharge_date = $2, discharge_time = $3, length_of_stay_days = $4,
    outcome = $5, discharge_destination = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteHospitalization :exec
UPDATE hospitalizations SET deleted_at = now() WHERE id = $1;
