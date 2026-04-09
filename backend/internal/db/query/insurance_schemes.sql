-- name: CreateInsuranceScheme :one
INSERT INTO insurance_schemes (
    hospital_id, name, short_code, country, covers_dialysis,
    reimbursement_rate, requires_pre_authorization,
    claim_submission_url, contact_person, contact_phone,
    contact_email, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetInsuranceScheme :one
SELECT * FROM insurance_schemes WHERE id = $1 AND deleted_at IS NULL;

-- name: ListInsuranceSchemesByHospital :many
SELECT * FROM insurance_schemes
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListActiveInsuranceSchemes :many
SELECT * FROM insurance_schemes
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY name;

-- name: ListSchemesByCountry :many
SELECT * FROM insurance_schemes
WHERE hospital_id = $1 AND country = $2 AND deleted_at IS NULL
ORDER BY name;

-- name: UpdateInsuranceScheme :one
UPDATE insurance_schemes
SET name = $2, covers_dialysis = $3, reimbursement_rate = $4,
    contact_person = $5, contact_phone = $6, contact_email = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteInsuranceScheme :exec
UPDATE insurance_schemes SET deleted_at = now() WHERE id = $1;
