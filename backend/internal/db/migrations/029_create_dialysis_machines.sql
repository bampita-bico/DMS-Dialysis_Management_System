-- +goose Up
CREATE TABLE dialysis_machines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    machine_code VARCHAR(50) NOT NULL,
    serial_number VARCHAR(100) NOT NULL,
    model VARCHAR(150) NOT NULL,
    manufacturer VARCHAR(150) NOT NULL,
    manufacture_year INTEGER,
    installation_date DATE,
    location VARCHAR(100),
    status machine_status NOT NULL DEFAULT 'available',
    is_hbv_dedicated BOOLEAN NOT NULL DEFAULT FALSE,
    last_service_date DATE,
    next_service_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(hospital_id, machine_code)
);

ALTER TABLE dialysis_machines ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON dialysis_machines USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_machines_hospital ON dialysis_machines(hospital_id);
CREATE INDEX idx_machines_status ON dialysis_machines(hospital_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_machines_hbv ON dialysis_machines(hospital_id, is_hbv_dedicated) WHERE is_hbv_dedicated = TRUE;

CREATE TRIGGER trg_dialysis_machines_updated_at BEFORE UPDATE ON dialysis_machines FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS dialysis_machines;
