-- name: CreateHospital :one
INSERT INTO hospitals (
  name, short_code, tier, region, country, address, phone, email, license_no, license_expiry, settings
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetHospital :one
SELECT * FROM hospitals
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListHospitals :many
SELECT * FROM hospitals
WHERE deleted_at IS NULL
ORDER BY name;

-- name: UpdateHospital :one
UPDATE hospitals
SET name = $2, tier = $3, region = $4, address = $5, phone = $6, email = $7
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteHospital :exec
UPDATE hospitals
SET deleted_at = now()
WHERE id = $1;

-- name: GetHospitalPlan :one
SELECT subscription_plan, enabled_modules FROM hospitals
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateHospitalPlan :exec
UPDATE hospitals
SET subscription_plan = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateEnabledModules :exec
UPDATE hospitals
SET enabled_modules = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListHospitalsByPlan :many
SELECT * FROM hospitals
WHERE subscription_plan = $1 AND deleted_at IS NULL
ORDER BY name;
