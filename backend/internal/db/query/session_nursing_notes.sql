-- name: CreateSessionNursingNote :one
INSERT INTO session_nursing_notes (
    hospital_id, session_id, patient_id, nurse_id, note_type,
    recorded_at, content, is_flagged_for_doctor
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSessionNursingNote :one
SELECT * FROM session_nursing_notes WHERE id = $1 AND deleted_at IS NULL;

-- name: ListNursingNotesBySession :many
SELECT * FROM session_nursing_notes
WHERE session_id = $1 AND deleted_at IS NULL
ORDER BY recorded_at;

-- name: ListFlaggedNotes :many
SELECT * FROM session_nursing_notes
WHERE hospital_id = $1
  AND is_flagged_for_doctor = TRUE
  AND doctor_reviewed_at IS NULL
  AND deleted_at IS NULL
ORDER BY recorded_at;

-- name: MarkNoteReviewedByDoctor :one
UPDATE session_nursing_notes
SET doctor_reviewed_by = $2, doctor_reviewed_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSessionNursingNote :exec
UPDATE session_nursing_notes SET deleted_at = now() WHERE id = $1;
