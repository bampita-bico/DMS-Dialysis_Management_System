-- name: CreateSessionFluidBalance :one
INSERT INTO session_fluid_balance (
    hospital_id, session_id, patient_id, recorded_by, recorded_at,
    uf_goal_ml, uf_achieved_ml, uf_rate_ml_per_hr, fluid_intake_oral_ml,
    fluid_intake_iv_ml, fluid_output_urine_ml, fluid_output_other_ml,
    net_fluid_balance_ml, weight_change_kg, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetSessionFluidBalance :one
SELECT * FROM session_fluid_balance WHERE id = $1 AND deleted_at IS NULL;

-- name: GetFluidBalanceBySession :one
SELECT * FROM session_fluid_balance WHERE session_id = $1 AND deleted_at IS NULL;

-- name: ListFluidBalancesByPatient :many
SELECT * FROM session_fluid_balance
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY recorded_at DESC;

-- name: UpdateSessionFluidBalance :one
UPDATE session_fluid_balance
SET uf_goal_ml = $2, uf_achieved_ml = $3, uf_rate_ml_per_hr = $4,
    fluid_intake_oral_ml = $5, fluid_intake_iv_ml = $6,
    fluid_output_urine_ml = $7, fluid_output_other_ml = $8,
    net_fluid_balance_ml = $9, weight_change_kg = $10, notes = $11
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionFluidBalance :exec
UPDATE session_fluid_balance SET deleted_at = now() WHERE id = $1;
