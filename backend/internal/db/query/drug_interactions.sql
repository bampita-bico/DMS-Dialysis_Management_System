-- name: CreateDrugInteraction :one
INSERT INTO drug_interactions (
    hospital_id, medication_a_id, medication_b_id, severity, description,
    clinical_effect, management_recommendation
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetDrugInteraction :one
SELECT * FROM drug_interactions WHERE id = $1 AND deleted_at IS NULL;

-- name: CheckInteraction :one
SELECT * FROM drug_interactions
WHERE ((medication_a_id = $1 AND medication_b_id = $2)
   OR (medication_a_id = $2 AND medication_b_id = $1))
AND deleted_at IS NULL
LIMIT 1;

-- name: ListInteractionsForMedication :many
SELECT * FROM drug_interactions
WHERE (medication_a_id = $1 OR medication_b_id = $1)
AND deleted_at IS NULL
ORDER BY severity DESC;

-- name: ListSevereInteractions :many
SELECT * FROM drug_interactions
WHERE hospital_id = $1 AND severity IN ('severe', 'contraindicated')
AND deleted_at IS NULL
ORDER BY severity DESC, medication_a_id;

-- name: DeleteDrugInteraction :exec
UPDATE drug_interactions SET deleted_at = now() WHERE id = $1;
