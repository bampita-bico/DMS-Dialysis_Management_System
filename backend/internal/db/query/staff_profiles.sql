-- name: CreateStaffProfile :one
INSERT INTO staff_profiles (
    hospital_id, user_id, department_id, cadre, license_number,
    license_expiry_date, registration_body, specialization,
    years_of_experience, hire_date, employee_number,
    emergency_contact_name, emergency_contact_phone, blood_type, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetStaffProfile :one
SELECT * FROM staff_profiles WHERE id = $1 AND deleted_at IS NULL;

-- name: GetStaffProfileByUser :one
SELECT * FROM staff_profiles
WHERE hospital_id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListStaffProfilesByHospital :many
SELECT * FROM staff_profiles
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListStaffByCadre :many
SELECT * FROM staff_profiles
WHERE hospital_id = $1 AND cadre = $2 AND deleted_at IS NULL
ORDER BY user_id;

-- name: ListActiveStaff :many
SELECT * FROM staff_profiles
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY user_id;

-- name: ListStaffByDepartment :many
SELECT * FROM staff_profiles
WHERE department_id = $1 AND deleted_at IS NULL
ORDER BY user_id;

-- name: ListExpiringLicenses :many
SELECT * FROM staff_profiles
WHERE hospital_id = $1
  AND license_expiry_date <= $2
  AND is_active = TRUE
  AND deleted_at IS NULL
ORDER BY license_expiry_date ASC;

-- name: UpdateStaffProfile :one
UPDATE staff_profiles
SET cadre = $2, license_number = $3, license_expiry_date = $4,
    specialization = $5, years_of_experience = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteStaffProfile :exec
UPDATE staff_profiles SET deleted_at = now() WHERE id = $1;
