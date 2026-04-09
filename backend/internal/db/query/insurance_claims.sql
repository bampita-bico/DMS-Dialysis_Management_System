-- name: CreateInsuranceClaim :one
INSERT INTO insurance_claims (
    hospital_id, invoice_id, scheme_id, patient_id, claim_number,
    claim_date, claimed_amount, status, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetInsuranceClaim :one
SELECT * FROM insurance_claims WHERE id = $1 AND deleted_at IS NULL;

-- name: GetClaimByNumber :one
SELECT * FROM insurance_claims
WHERE hospital_id = $1 AND claim_number = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListClaimsByInvoice :many
SELECT * FROM insurance_claims
WHERE invoice_id = $1 AND deleted_at IS NULL
ORDER BY claim_date DESC;

-- name: ListClaimsByScheme :many
SELECT * FROM insurance_claims
WHERE scheme_id = $1 AND deleted_at IS NULL
ORDER BY claim_date DESC;

-- name: ListClaimsByStatus :many
SELECT * FROM insurance_claims
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY claim_date DESC;

-- name: ListPendingClaims :many
SELECT * FROM insurance_claims
WHERE hospital_id = $1 AND status IN ('draft', 'submitted', 'under_review')
  AND deleted_at IS NULL
ORDER BY claim_date ASC;

-- name: SubmitClaim :one
UPDATE insurance_claims
SET status = 'submitted', submitted_by = $2, submitted_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ApproveClaim :one
UPDATE insurance_claims
SET status = 'approved', approved_amount = $2,
    approved_by_insurer = $3, approved_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RejectClaim :one
UPDATE insurance_claims
SET status = 'rejected', rejection_reason = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteInsuranceClaim :exec
UPDATE insurance_claims SET deleted_at = now() WHERE id = $1;
