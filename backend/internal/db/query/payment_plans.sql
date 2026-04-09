-- name: CreatePaymentPlan :one
INSERT INTO payment_plans (
    hospital_id, account_id, patient_id, plan_number, total_amount,
    down_payment, installment_amount, installment_frequency,
    number_of_installments, start_date, end_date, balance_remaining,
    approved_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetPaymentPlan :one
SELECT * FROM payment_plans WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPlanByNumber :one
SELECT * FROM payment_plans
WHERE hospital_id = $1 AND plan_number = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListPlansByPatient :many
SELECT * FROM payment_plans
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListPlansByAccount :many
SELECT * FROM payment_plans
WHERE account_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListActivePlans :many
SELECT * FROM payment_plans
WHERE hospital_id = $1 AND status = 'active' AND deleted_at IS NULL
ORDER BY start_date ASC;

-- name: UpdatePlanPayment :one
UPDATE payment_plans
SET amount_paid = $2, balance_remaining = total_amount - $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CompletePlan :one
UPDATE payment_plans
SET status = 'completed'
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePaymentPlan :exec
UPDATE payment_plans SET deleted_at = now() WHERE id = $1;
