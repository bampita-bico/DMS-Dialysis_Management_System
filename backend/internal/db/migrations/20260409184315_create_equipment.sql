-- +goose Up
CREATE TABLE equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    category equipment_category NOT NULL,
    serial_number VARCHAR(100),
    model VARCHAR(100),
    manufacturer VARCHAR(255),
    purchase_date DATE,
    purchase_cost DECIMAL(12,2),
    warranty_expiry_date DATE,
    status equipment_status NOT NULL DEFAULT 'operational',
    location VARCHAR(255),
    department_id UUID REFERENCES departments(id),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_equipment_serial_unique ON equipment(hospital_id, serial_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_hospital ON equipment(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_category ON equipment(hospital_id, category) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_status ON equipment(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_department ON equipment(department_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_equipment_updated_at
BEFORE UPDATE ON equipment
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE equipment ENABLE ROW LEVEL SECURITY;
CREATE POLICY equipment_isolation ON equipment
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS equipment CASCADE;
