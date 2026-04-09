-- name: CreateConsent :one
INSERT INTO consents (
    hospital_id, patient_id, consent_type, given_by, relationship, witness_id,
    expires_at, document_url, notes, recorded_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: CheckActiveConsent :one
SELECT * FROM consents
WHERE patient_id = $1 AND consent_type = $2 AND status = 'given' AND deleted_at IS NULL
    AND (expires_at IS NULL OR expires_at > now())
LIMIT 1;

-- name: GetConsentByType :one
SELECT * FROM consents
WHERE patient_id = $1 AND consent_type = $2 AND deleted_at IS NULL
ORDER BY signed_at DESC
LIMIT 1;

-- name: ListConsentsByPatient :many
SELECT * FROM consents
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY signed_at DESC;

-- name: WithdrawConsent :exec
UPDATE consents
SET status = 'withdrawn', withdrawn_at = now(), withdrawn_reason = $2
WHERE id = $1;
