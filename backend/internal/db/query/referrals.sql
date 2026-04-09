-- name: CreateReferral :one
INSERT INTO referrals (
    hospital_id, patient_id, direction, from_facility, to_facility, from_doctor,
    to_doctor_id, reason, urgency, referral_letter_url, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: ListReferralsByPatient :many
SELECT * FROM referrals
WHERE patient_id = $1 AND deleted_at IS NULL
ORDER BY referral_date DESC;

-- name: UpdateReferralStatus :one
UPDATE referrals
SET status = $2, response_date = $3, response_notes = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
