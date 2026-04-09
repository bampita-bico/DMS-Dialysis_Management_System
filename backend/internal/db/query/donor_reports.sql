-- name: CreateDonorReport :one
INSERT INTO donor_reports (
    hospital_id, report_name, report_type, period_start, period_end,
    recipient_organization, data, generated_by, notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetDonorReport :one
SELECT * FROM donor_reports WHERE id = $1 AND deleted_at IS NULL;

-- name: ListReportsByHospital :many
SELECT * FROM donor_reports
WHERE hospital_id = $1 AND deleted_at IS NULL
ORDER BY period_start DESC, report_type;

-- name: ListReportsByType :many
SELECT * FROM donor_reports
WHERE hospital_id = $1 AND report_type = $2 AND deleted_at IS NULL
ORDER BY period_start DESC;

-- name: ListReportsByPeriod :many
SELECT * FROM donor_reports
WHERE hospital_id = $1
  AND period_start >= $2
  AND period_end <= $3
  AND deleted_at IS NULL
ORDER BY period_start DESC;

-- name: ListUnsubmittedReports :many
SELECT * FROM donor_reports
WHERE hospital_id = $1 AND submitted_at IS NULL AND deleted_at IS NULL
ORDER BY generated_at DESC;

-- name: SubmitReport :one
UPDATE donor_reports
SET submitted_at = now(), submitted_by = $2, submission_reference = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ApproveReport :one
UPDATE donor_reports
SET approved_by = $2, approved_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDonorReport :exec
UPDATE donor_reports SET deleted_at = now() WHERE id = $1;
