-- name: CreateDialysateRecord :one
INSERT INTO dialysate_records (
    hospital_id, session_id, patient_id, recorded_by, recorded_at,
    batch_number, sodium_meq_l, potassium_meq_l, bicarbonate_meq_l,
    calcium_meq_l, magnesium_meq_l, chloride_meq_l, glucose_mg_dl,
    acetate_meq_l, conductivity_ms_cm, ph_level, temperature_celsius,
    flow_rate_ml_min, total_volume_liters, composition_verified,
    verified_by, deviations_noted, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
RETURNING *;

-- name: GetDialysateRecord :one
SELECT * FROM dialysate_records WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDialysateBySession :one
SELECT * FROM dialysate_records WHERE session_id = $1 AND deleted_at IS NULL;

-- name: ListDialysateRecordsByBatch :many
SELECT * FROM dialysate_records
WHERE hospital_id = $1 AND batch_number = $2 AND deleted_at IS NULL
ORDER BY recorded_at;

-- name: VerifyDialysateComposition :one
UPDATE dialysate_records
SET composition_verified = TRUE, verified_by = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDialysateRecord :exec
UPDATE dialysate_records SET deleted_at = now() WHERE id = $1;
