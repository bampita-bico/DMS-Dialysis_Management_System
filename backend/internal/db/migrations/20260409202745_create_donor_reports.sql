-- +goose Up
CREATE TABLE donor_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    report_name VARCHAR(255) NOT NULL,
    report_type report_type NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    recipient_organization VARCHAR(255) NOT NULL,
    data JSONB NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    generated_by UUID NOT NULL REFERENCES users(id),
    submitted_at TIMESTAMPTZ,
    submitted_by UUID REFERENCES users(id),
    submission_reference VARCHAR(100),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_donor_reports_hospital ON donor_reports(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_donor_reports_type ON donor_reports(hospital_id, report_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_donor_reports_period ON donor_reports(hospital_id, period_start, period_end) WHERE deleted_at IS NULL;
CREATE INDEX idx_donor_reports_submitted ON donor_reports(hospital_id, submitted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_donor_reports_updated_at
BEFORE UPDATE ON donor_reports
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE donor_reports ENABLE ROW LEVEL SECURITY;
CREATE POLICY donor_reports_isolation ON donor_reports
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS donor_reports CASCADE;
