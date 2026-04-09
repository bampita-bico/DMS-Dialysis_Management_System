-- name: CreatePatient :one
INSERT INTO patients (
    hospital_id, mrn, national_id, full_name, preferred_name, date_of_birth, sex, blood_type,
    marital_status, nationality, religion, occupation, education_level, photo_url,
    primary_language, interpreter_needed, registration_date, registered_by, primary_doctor_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetPatient :one
SELECT * FROM patients
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPatientByMRN :one
SELECT * FROM patients
WHERE hospital_id = $1 AND mrn = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPatientByNationalID :one
SELECT * FROM patients
WHERE national_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: SearchPatientsByName :many
SELECT * FROM patients
WHERE hospital_id = $1 AND full_name ILIKE $2 AND deleted_at IS NULL
ORDER BY full_name
LIMIT $3 OFFSET $4;

-- name: ListActivePatients :many
SELECT * FROM patients
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY full_name
LIMIT $2 OFFSET $3;

-- name: UpdatePatient :one
UPDATE patients
SET full_name = $2, preferred_name = $3, blood_type = $4, marital_status = $5,
    nationality = $6, religion = $7, occupation = $8, education_level = $9,
    primary_language = $10, interpreter_needed = $11, primary_doctor_id = $12
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeletePatient :exec
UPDATE patients
SET deleted_at = now()
WHERE id = $1;

-- name: MarkPatientDeceased :exec
UPDATE patients
SET deceased_at = now(), is_active = FALSE
WHERE id = $1;
