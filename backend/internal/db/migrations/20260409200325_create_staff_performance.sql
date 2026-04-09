-- +goose Up
CREATE TABLE staff_performance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    review_period_start DATE NOT NULL,
    review_period_end DATE NOT NULL,
    review_date DATE NOT NULL,
    appraised_by UUID NOT NULL REFERENCES users(id),
    overall_score DECIMAL(5,2),
    technical_competence_score DECIMAL(5,2),
    communication_score DECIMAL(5,2),
    teamwork_score DECIMAL(5,2),
    punctuality_score DECIMAL(5,2),
    patient_care_score DECIMAL(5,2),
    strengths TEXT,
    areas_for_improvement TEXT,
    goals_next_period TEXT,
    training_recommendations TEXT,
    staff_comments TEXT,
    promotion_recommended BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_staff_performance_staff ON staff_performance(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_performance_hospital ON staff_performance(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_performance_period ON staff_performance(hospital_id, review_period_start, review_period_end) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_staff_performance_updated_at
BEFORE UPDATE ON staff_performance
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE staff_performance ENABLE ROW LEVEL SECURITY;
CREATE POLICY staff_performance_isolation ON staff_performance
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS staff_performance CASCADE;
