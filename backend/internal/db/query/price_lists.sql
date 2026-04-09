-- name: CreatePriceList :one
INSERT INTO price_lists (
    hospital_id, service_name, service_code, service_category,
    unit_price, scheme_id, effective_from, effective_until, description
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPriceList :one
SELECT * FROM price_lists WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPriceByServiceCode :one
SELECT * FROM price_lists
WHERE hospital_id = $1
  AND service_code = $2
  AND (scheme_id = $3 OR scheme_id IS NULL)
  AND effective_from <= CURRENT_DATE
  AND (effective_until IS NULL OR effective_until >= CURRENT_DATE)
  AND deleted_at IS NULL
ORDER BY scheme_id DESC NULLS LAST
LIMIT 1;

-- name: ListPricesByHospital :many
SELECT * FROM price_lists
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY service_name;

-- name: ListPricesByScheme :many
SELECT * FROM price_lists
WHERE hospital_id = $1 AND scheme_id = $2 AND deleted_at IS NULL
ORDER BY service_name;

-- name: ListActivePrices :many
SELECT * FROM price_lists
WHERE hospital_id = $1
  AND effective_from <= CURRENT_DATE
  AND (effective_until IS NULL OR effective_until >= CURRENT_DATE)
  AND deleted_at IS NULL
ORDER BY service_name;

-- name: UpdatePriceList :one
UPDATE price_lists
SET service_name = $2, unit_price = $3, effective_from = $4, effective_until = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePriceList :exec
UPDATE price_lists SET deleted_at = now() WHERE id = $1;
