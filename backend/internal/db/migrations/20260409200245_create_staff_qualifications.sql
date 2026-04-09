-- +goose Up
CREATE TABLE staff_qualifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    staff_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    qualification_type VARCHAR(100) NOT NULL,
    qualification_name VARCHAR(255) NOT NULL,
    institution VARCHAR(255),
    country VARCHAR(100),
    year_obtained INTEGER,
    certificate_number VARCHAR(100),
    expiry_date DATE,
    document_url TEXT,
    file_attachment_id UUID REFERENCES file_attachments(id),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_staff_qualifications_staff ON staff_qualifications(staff_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_qualifications_hospital ON staff_qualifications(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_qualifications_type ON staff_qualifications(hospital_id, qualification_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_qualifications_expiry ON staff_qualifications(hospital_id, expiry_date) WHERE deleted_at IS NULL AND expiry_date IS NOT NULL;

CREATE TRIGGER trg_staff_qualifications_updated_at
BEFORE UPDATE ON staff_qualifications
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE staff_qualifications ENABLE ROW LEVEL SECURITY;
CREATE POLICY staff_qualifications_isolation ON staff_qualifications
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS staff_qualifications CASCADE;
