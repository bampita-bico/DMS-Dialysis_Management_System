-- +goose Up
CREATE TABLE equipment_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
    certification_type certification_type NOT NULL,
    certificate_number VARCHAR(100),
    issued_by VARCHAR(255) NOT NULL,
    issued_date DATE,
    valid_from DATE NOT NULL,
    valid_until DATE NOT NULL,
    document_url TEXT,
    file_attachment_id UUID REFERENCES file_attachments(id),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_equipment_certifications_equipment ON equipment_certifications(equipment_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_certifications_hospital ON equipment_certifications(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_certifications_valid_until ON equipment_certifications(hospital_id, valid_until) WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_equipment_certifications_type ON equipment_certifications(hospital_id, certification_type) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_equipment_certifications_updated_at
BEFORE UPDATE ON equipment_certifications
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE equipment_certifications ENABLE ROW LEVEL SECURITY;
CREATE POLICY equipment_certifications_isolation ON equipment_certifications
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS equipment_certifications CASCADE;
