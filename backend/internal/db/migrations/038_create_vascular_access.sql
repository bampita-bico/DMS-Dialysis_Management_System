-- +goose Up
CREATE TABLE vascular_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    access_type access_type NOT NULL,
    access_site access_site NOT NULL,
    site_side VARCHAR(10) NOT NULL,
    insertion_date DATE NOT NULL,
    inserted_by VARCHAR(200),
    insertion_location VARCHAR(200),
    status access_status NOT NULL DEFAULT 'active',
    maturation_date DATE,
    first_use_date DATE,
    abandonment_date DATE,
    abandonment_reason TEXT,
    catheter_type VARCHAR(100),
    catheter_length_cm NUMERIC(4,1),
    catheter_position VARCHAR(100),
    fistula_vein VARCHAR(100),
    fistula_artery VARCHAR(100),
    graft_material VARCHAR(100),
    is_primary_access BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE vascular_access ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON vascular_access USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_vascular_access_patient ON vascular_access(patient_id);
CREATE INDEX idx_vascular_access_status ON vascular_access(patient_id, status) WHERE status = 'active';
CREATE INDEX idx_vascular_access_primary ON vascular_access(patient_id) WHERE is_primary_access = TRUE AND deleted_at IS NULL;

CREATE TRIGGER trg_vascular_access_updated_at BEFORE UPDATE ON vascular_access FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS vascular_access;
