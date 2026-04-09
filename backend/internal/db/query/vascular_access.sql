-- name: CreateVascularAccess :one
INSERT INTO vascular_access (
    hospital_id, patient_id, access_type, access_site, site_side,
    insertion_date, inserted_by, insertion_location, status, maturation_date,
    first_use_date, catheter_type, catheter_length_cm, catheter_position,
    fistula_vein, fistula_artery, graft_material, is_primary_access, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING *;

-- name: GetVascularAccess :one
SELECT * FROM vascular_access WHERE id = $1 AND deleted_at IS NULL;

-- name: ListVascularAccessByPatient :many
SELECT * FROM vascular_access
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY insertion_date DESC;

-- name: GetPrimaryAccessForPatient :one
SELECT * FROM vascular_access
WHERE patient_id = $1
  AND is_primary_access = TRUE
  AND status = 'active'
  AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateVascularAccess :one
UPDATE vascular_access
SET access_type = $2, access_site = $3, site_side = $4, status = $5,
    maturation_date = $6, first_use_date = $7, abandonment_date = $8,
    abandonment_reason = $9, catheter_type = $10, catheter_length_cm = $11,
    catheter_position = $12, fistula_vein = $13, fistula_artery = $14,
    graft_material = $15, is_primary_access = $16, notes = $17
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: AbandonAccess :one
UPDATE vascular_access
SET status = 'abandoned', abandonment_date = CURRENT_DATE, abandonment_reason = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteVascularAccess :exec
UPDATE vascular_access SET deleted_at = now() WHERE id = $1;
