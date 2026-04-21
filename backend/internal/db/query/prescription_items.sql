-- name: CreatePrescriptionItem :one
INSERT INTO prescription_items (
    hospital_id, prescription_id, medication_id, dose, frequency, route,
    duration_days, quantity_prescribed, instructions, start_date, end_date,
    is_prn, prn_indication, is_stat, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetPrescriptionItem :one
SELECT * FROM prescription_items WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPrescriptionItemsByPrescription :many
SELECT * FROM prescription_items
WHERE prescription_id = $1 AND deleted_at IS NULL
ORDER BY created_at;

-- name: UpdateDispensedQuantity :one
UPDATE prescription_items
SET quantity_dispensed = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ListActiveMedicationIDsForPatient :many
SELECT DISTINCT pi.medication_id
FROM prescription_items pi
JOIN prescriptions p ON p.id = pi.prescription_id
WHERE p.patient_id = $1
  AND p.status = 'active'
  AND p.deleted_at IS NULL
  AND pi.deleted_at IS NULL
  AND (pi.end_date IS NULL OR pi.end_date >= CURRENT_DATE);

-- name: DeletePrescriptionItem :exec
UPDATE prescription_items SET deleted_at = now() WHERE id = $1;
