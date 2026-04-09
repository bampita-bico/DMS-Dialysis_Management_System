-- name: CreatePatientTransport :one
INSERT INTO patient_transport (
    hospital_id, patient_id, session_id, transport_date, pickup_location,
    pickup_time, dropoff_location, vehicle_registration, driver_name,
    driver_phone, distance_km, cost, status, arranged_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetPatientTransport :one
SELECT * FROM patient_transport WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTransportByPatient :many
SELECT * FROM patient_transport
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY transport_date DESC;

-- name: ListTransportByDate :many
SELECT * FROM patient_transport
WHERE hospital_id = $1 AND transport_date = $2 AND deleted_at IS NULL
ORDER BY pickup_time NULLS LAST;

-- name: ListTransportByStatus :many
SELECT * FROM patient_transport
WHERE hospital_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY transport_date, pickup_time NULLS LAST;

-- name: UpdateTransportStatus :one
UPDATE patient_transport
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePatientTransport :exec
UPDATE patient_transport SET deleted_at = now() WHERE id = $1;
