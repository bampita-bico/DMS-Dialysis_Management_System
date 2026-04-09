-- +goose Up
CREATE TABLE community_health_workers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    alt_phone VARCHAR(50),
    region VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    village VARCHAR(100),
    catchment_area TEXT,
    chw_id_number VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    registered_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE community_health_workers ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON community_health_workers USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE TABLE chw_patient_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    chw_id UUID NOT NULL REFERENCES community_health_workers(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(chw_id, patient_id)
);

ALTER TABLE chw_patient_assignments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chw_patient_assignments USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_chw_patient_chw ON chw_patient_assignments(chw_id);
CREATE INDEX idx_chw_patient_patient ON chw_patient_assignments(patient_id);

CREATE TRIGGER trg_community_health_workers_updated_at BEFORE UPDATE ON community_health_workers FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();
CREATE TRIGGER trg_chw_patient_assignments_updated_at BEFORE UPDATE ON chw_patient_assignments FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS chw_patient_assignments;
DROP TABLE IF EXISTS community_health_workers;
