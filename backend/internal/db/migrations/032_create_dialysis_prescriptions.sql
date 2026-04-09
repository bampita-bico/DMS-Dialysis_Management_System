-- +goose Up
CREATE TABLE dialysis_prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    prescribed_by UUID NOT NULL REFERENCES users(id),
    modality dialysis_modality NOT NULL DEFAULT 'hd',
    duration_mins INTEGER NOT NULL DEFAULT 240,
    target_uf_ml NUMERIC(7,2),
    blood_flow_rate INTEGER,
    dialysate_flow_rate INTEGER,
    membrane_type VARCHAR(100),
    membrane_surface_area NUMERIC(4,2),
    dialysate_temp NUMERIC(4,1),
    conductivity_target NUMERIC(4,1),
    pd_params JSONB DEFAULT '{}',
    anticoagulant VARCHAR(100),
    anticoag_route anticoag_route DEFAULT 'systemic',
    loading_dose_units NUMERIC(8,2),
    maintenance_dose_units NUMERIC(8,2),
    maintenance_rate VARCHAR(50),
    sodium_target NUMERIC(5,1),
    potassium_target NUMERIC(4,1),
    bicarbonate_target NUMERIC(5,1),
    calcium_target NUMERIC(4,2),
    glucose_target NUMERIC(5,1),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE dialysis_prescriptions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON dialysis_prescriptions USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_prescriptions_session ON dialysis_prescriptions(session_id);
CREATE INDEX idx_prescriptions_patient ON dialysis_prescriptions(patient_id);

CREATE TRIGGER trg_dialysis_prescriptions_updated_at BEFORE UPDATE ON dialysis_prescriptions FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS dialysis_prescriptions;
