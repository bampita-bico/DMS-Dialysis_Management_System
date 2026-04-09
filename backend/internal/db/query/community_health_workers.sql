-- name: CreateCHW :one
INSERT INTO community_health_workers (
    hospital_id, full_name, phone, alt_phone, region, district, village,
    catchment_area, chw_id_number, registered_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetPatientsByCHW :many
SELECT p.* FROM patients p
INNER JOIN chw_patient_assignments cpa ON cpa.patient_id = p.id
WHERE cpa.chw_id = $1 AND cpa.deleted_at IS NULL AND p.deleted_at IS NULL
ORDER BY p.full_name;

-- name: GetCHWByPatient :one
SELECT chw.* FROM community_health_workers chw
INNER JOIN chw_patient_assignments cpa ON cpa.chw_id = chw.id
WHERE cpa.patient_id = $1 AND cpa.deleted_at IS NULL
LIMIT 1;

-- name: ListCHWsByRegion :many
SELECT * FROM community_health_workers
WHERE hospital_id = $1 AND region = $2 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY full_name;

-- name: AssignCHWToPatient :one
INSERT INTO chw_patient_assignments (
    hospital_id, chw_id, patient_id, notes
) VALUES (
    $1, $2, $3, $4
) RETURNING *;
