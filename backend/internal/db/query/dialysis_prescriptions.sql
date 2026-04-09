-- name: CreateDialysisPrescription :one
INSERT INTO dialysis_prescriptions (
    hospital_id, session_id, patient_id, prescribed_by, modality, duration_mins,
    target_uf_ml, blood_flow_rate, dialysate_flow_rate, membrane_type,
    membrane_surface_area, dialysate_temp, conductivity_target, pd_params,
    anticoagulant, anticoag_route, loading_dose_units, maintenance_dose_units,
    maintenance_rate, sodium_target, potassium_target, bicarbonate_target,
    calcium_target, glucose_target, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
RETURNING *;

-- name: GetDialysisPrescription :one
SELECT * FROM dialysis_prescriptions WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPrescriptionBySession :one
SELECT * FROM dialysis_prescriptions WHERE session_id = $1 AND deleted_at IS NULL;

-- name: ListPrescriptionsByPatient :many
SELECT * FROM dialysis_prescriptions
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateDialysisPrescription :one
UPDATE dialysis_prescriptions
SET modality = $2, duration_mins = $3, target_uf_ml = $4, blood_flow_rate = $5,
    dialysate_flow_rate = $6, membrane_type = $7, membrane_surface_area = $8,
    dialysate_temp = $9, conductivity_target = $10, pd_params = $11,
    anticoagulant = $12, anticoag_route = $13, loading_dose_units = $14,
    maintenance_dose_units = $15, maintenance_rate = $16, sodium_target = $17,
    potassium_target = $18, bicarbonate_target = $19, calcium_target = $20,
    glucose_target = $21, notes = $22
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDialysisPrescription :exec
UPDATE dialysis_prescriptions SET deleted_at = now() WHERE id = $1;
