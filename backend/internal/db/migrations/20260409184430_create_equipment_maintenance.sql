-- +goose Up
CREATE TABLE equipment_maintenance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
    maintenance_type maintenance_type NOT NULL,
    scheduled_date DATE,
    performed_date DATE,
    performed_by UUID REFERENCES users(id),
    technician_name VARCHAR(255),
    technician_company VARCHAR(255),
    next_due_date DATE,
    cost DECIMAL(10,2),
    findings TEXT,
    actions_taken TEXT,
    parts_replaced TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_equipment_maintenance_equipment ON equipment_maintenance(equipment_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_maintenance_hospital ON equipment_maintenance(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_maintenance_next_due ON equipment_maintenance(hospital_id, next_due_date) WHERE deleted_at IS NULL AND is_completed = FALSE;
CREATE INDEX idx_equipment_maintenance_performed ON equipment_maintenance(performed_by) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_equipment_maintenance_updated_at
BEFORE UPDATE ON equipment_maintenance
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE equipment_maintenance ENABLE ROW LEVEL SECURITY;
CREATE POLICY equipment_maintenance_isolation ON equipment_maintenance
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS equipment_maintenance CASCADE;
