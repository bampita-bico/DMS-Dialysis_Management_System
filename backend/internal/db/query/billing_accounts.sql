-- name: CreateBillingAccount :one
INSERT INTO billing_accounts (
    hospital_id, patient_id, guarantor_id, account_number,
    account_status, credit_limit, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBillingAccount :one
SELECT * FROM billing_accounts WHERE id = $1 AND deleted_at IS NULL;

-- name: GetBillingAccountByPatient :one
SELECT * FROM billing_accounts
WHERE hospital_id = $1 AND patient_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListBillingAccountsByHospital :many
SELECT * FROM billing_accounts
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListAccountsByStatus :many
SELECT * FROM billing_accounts
WHERE hospital_id = $1 AND account_status = $2 AND deleted_at IS NULL
ORDER BY current_balance DESC;

-- name: UpdateAccountBalance :one
UPDATE billing_accounts
SET current_balance = $2, total_billed = $3, total_paid = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateAccountStatus :one
UPDATE billing_accounts
SET account_status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteBillingAccount :exec
UPDATE billing_accounts SET deleted_at = now() WHERE id = $1;
