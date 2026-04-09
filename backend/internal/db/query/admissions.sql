-- name: CreateAdmission :one
INSERT INTO admissions (
    hospital_id, patient_id, admission_type, admitted_by, ward, bed_number,
    primary_diagnosis, admitting_doctor_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetCurrentAdmission :one
SELECT * FROM admissions
WHERE patient_id = $1 AND discharged_at IS NULL AND deleted_at IS NULL
LIMIT 1;

-- name: ListActiveAdmissions :many
SELECT * FROM admissions
WHERE hospital_id = $1 AND discharged_at IS NULL AND deleted_at IS NULL
ORDER BY admitted_at DESC;

-- name: DischargePatient :one
UPDATE admissions
SET discharged_at = now(), discharge_type = $2, discharge_summary = $3, discharged_by = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
