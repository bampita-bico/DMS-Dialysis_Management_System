-- name: CreateEquipmentCertification :one
INSERT INTO equipment_certifications (
    hospital_id, equipment_id, certification_type, certificate_number,
    issued_by, issued_date, valid_from, valid_until, document_url,
    file_attachment_id, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetEquipmentCertification :one
SELECT * FROM equipment_certifications WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCertificationsByEquipment :many
SELECT * FROM equipment_certifications
WHERE equipment_id = $1 AND deleted_at IS NULL
ORDER BY valid_until DESC;

-- name: ListActiveCertifications :many
SELECT * FROM equipment_certifications
WHERE hospital_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY valid_until ASC;

-- name: ListExpiringCertifications :many
SELECT c.*, e.name AS equipment_name
FROM equipment_certifications c
JOIN equipment e ON c.equipment_id = e.id
WHERE c.hospital_id = $1
  AND c.valid_until <= $2
  AND c.is_active = TRUE
  AND c.deleted_at IS NULL
ORDER BY c.valid_until ASC;

-- name: ListCertificationsByType :many
SELECT * FROM equipment_certifications
WHERE hospital_id = $1 AND certification_type = $2 AND deleted_at IS NULL
ORDER BY valid_until DESC;

-- name: DeleteEquipmentCertification :exec
UPDATE equipment_certifications SET deleted_at = now() WHERE id = $1;
