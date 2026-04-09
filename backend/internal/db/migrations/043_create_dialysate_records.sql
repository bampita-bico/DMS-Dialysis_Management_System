-- +goose Up
CREATE TABLE dialysate_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    recorded_by UUID NOT NULL REFERENCES users(id),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    batch_number VARCHAR(100),
    sodium_meq_l NUMERIC(5,1),
    potassium_meq_l NUMERIC(4,1),
    bicarbonate_meq_l NUMERIC(5,1),
    calcium_meq_l NUMERIC(4,2),
    magnesium_meq_l NUMERIC(4,2),
    chloride_meq_l NUMERIC(5,1),
    glucose_mg_dl NUMERIC(5,1),
    acetate_meq_l NUMERIC(5,1),
    conductivity_ms_cm NUMERIC(4,1),
    ph_level NUMERIC(4,2),
    temperature_celsius NUMERIC(4,1),
    flow_rate_ml_min INTEGER,
    total_volume_liters NUMERIC(6,2),
    composition_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    deviations_noted TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE dialysate_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON dialysate_records USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_dialysate_session ON dialysate_records(session_id);
CREATE INDEX idx_dialysate_patient ON dialysate_records(patient_id);
CREATE INDEX idx_dialysate_batch ON dialysate_records(hospital_id, batch_number);

CREATE TRIGGER trg_dialysate_records_updated_at BEFORE UPDATE ON dialysate_records FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS dialysate_records;
