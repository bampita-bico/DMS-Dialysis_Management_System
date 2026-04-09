-- +goose Up
CREATE TABLE equipment_faults (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
    reported_by UUID NOT NULL REFERENCES users(id),
    reported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    fault_description TEXT NOT NULL,
    severity fault_severity NOT NULL,
    is_equipment_unusable BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    resolution_description TEXT,
    downtime_minutes INTEGER,
    cost DECIMAL(10,2),
    is_resolved BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_equipment_faults_equipment ON equipment_faults(equipment_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_faults_hospital ON equipment_faults(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_faults_unresolved ON equipment_faults(hospital_id, is_resolved) WHERE deleted_at IS NULL AND is_resolved = FALSE;
CREATE INDEX idx_equipment_faults_severity ON equipment_faults(hospital_id, severity) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_equipment_faults_updated_at
BEFORE UPDATE ON equipment_faults
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE equipment_faults ENABLE ROW LEVEL SECURITY;
CREATE POLICY equipment_faults_isolation ON equipment_faults
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS equipment_faults CASCADE;
