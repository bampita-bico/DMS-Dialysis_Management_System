-- name: CreateMedicationAdministration :one
INSERT INTO medication_administrations (
    hospital_id, prescription_item_id, patient_id, session_id, administered_by,
    scheduled_time, dose_given, route, site, is_refused, refusal_reason,
    is_omitted, omission_reason, adverse_reaction, adverse_reaction_details,
    witnessed_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
RETURNING *;

-- name: GetMedicationAdministration :one
SELECT * FROM medication_administrations WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAdministrationsByPatient :many
SELECT * FROM medication_administrations
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY administered_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAdministrationsBySession :many
SELECT * FROM medication_administrations
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY administered_at;

-- name: ListAdverseReactions :many
SELECT * FROM medication_administrations
WHERE hospital_id = $1 AND adverse_reaction = TRUE AND deleted_at IS NULL
ORDER BY administered_at DESC;

-- name: DeleteMedicationAdministration :exec
UPDATE medication_administrations SET deleted_at = now() WHERE id = $1;
