-- name: CreateHospitalSetting :one
INSERT INTO hospital_settings (
  hospital_id, key, value, data_type, description
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetHospitalSetting :one
SELECT * FROM hospital_settings
WHERE hospital_id = $1 AND key = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListHospitalSettings :many
SELECT * FROM hospital_settings
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY key;

-- name: UpdateHospitalSetting :one
UPDATE hospital_settings
SET value = $3, data_type = $4, description = $5
WHERE hospital_id = $1 AND key = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteHospitalSetting :exec
UPDATE hospital_settings
SET deleted_at = now()
WHERE hospital_id = $1 AND key = $2;
