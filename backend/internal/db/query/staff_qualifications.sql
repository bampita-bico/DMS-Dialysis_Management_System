-- name: CreateStaffQualification :one
INSERT INTO staff_qualifications (
    hospital_id, staff_id, qualification_type, qualification_name,
    institution, country, year_obtained, certificate_number,
    expiry_date, document_url, file_attachment_id, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetStaffQualification :one
SELECT * FROM staff_qualifications WHERE id = $1 AND deleted_at IS NULL;

-- name: ListQualificationsByStaff :many
SELECT * FROM staff_qualifications
WHERE staff_id = $1 AND deleted_at IS NULL
ORDER BY year_obtained DESC;

-- name: ListExpiringQualifications :many
SELECT * FROM staff_qualifications
WHERE hospital_id = $1
  AND expiry_date <= $2
  AND deleted_at IS NULL
ORDER BY expiry_date ASC;

-- name: VerifyQualification :one
UPDATE staff_qualifications
SET is_verified = TRUE, verified_by = $2, verified_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteStaffQualification :exec
UPDATE staff_qualifications SET deleted_at = now() WHERE id = $1;
