-- +goose Up
CREATE TABLE patient_transport (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    session_id UUID REFERENCES dialysis_sessions(id),
    transport_date DATE NOT NULL,
    pickup_location TEXT NOT NULL,
    pickup_time TIME,
    dropoff_location TEXT,
    dropoff_time TIME,
    vehicle_registration VARCHAR(50),
    driver_name VARCHAR(255),
    driver_phone VARCHAR(50),
    distance_km DECIMAL(8,2),
    cost DECIMAL(10,2),
    status transport_status NOT NULL DEFAULT 'scheduled',
    arranged_by UUID NOT NULL REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_patient_transport_patient ON patient_transport(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_patient_transport_session ON patient_transport(session_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_patient_transport_hospital ON patient_transport(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_patient_transport_date ON patient_transport(hospital_id, transport_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_patient_transport_status ON patient_transport(hospital_id, status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_patient_transport_updated_at
BEFORE UPDATE ON patient_transport
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE patient_transport ENABLE ROW LEVEL SECURITY;
CREATE POLICY patient_transport_isolation ON patient_transport
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS patient_transport CASCADE;
