-- name: CreateMortalityRecord :one
INSERT INTO mortality_records (
    hospital_id, patient_id, session_id, date_of_death, time_of_death,
    death_setting, session_related, primary_cause_of_death,
    contributing_factors, icd10_code, autopsy_performed, autopsy_findings,
    reported_by, certified_by, death_certificate_number, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: GetMortalityRecord :one
SELECT * FROM mortality_records WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMortalityByPatient :one
SELECT * FROM mortality_records
WHERE patient_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListMortalitiesByHospital :many
SELECT * FROM mortality_records
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY date_of_death DESC;

-- name: ListMortalitiesByPeriod :many
SELECT * FROM mortality_records
WHERE hospital_id = $1
  AND date_of_death BETWEEN $2 AND $3
  AND deleted_at IS NULL
ORDER BY date_of_death DESC;

-- name: ListSessionRelatedDeaths :many
SELECT * FROM mortality_records
WHERE hospital_id = $1
  AND session_related = TRUE
  AND deleted_at IS NULL
ORDER BY date_of_death DESC;

-- name: ListDeathsBySetting :many
SELECT * FROM mortality_records
WHERE hospital_id = $1 AND death_setting = $2 AND deleted_at IS NULL
ORDER BY date_of_death DESC;

-- name: CertifyDeath :one
UPDATE mortality_records
SET certified_by = $2, certified_at = now(), death_certificate_number = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMortalityRecord :exec
UPDATE mortality_records SET deleted_at = now() WHERE id = $1;
