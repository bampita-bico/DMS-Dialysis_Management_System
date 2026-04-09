-- +goose Up
CREATE TABLE session_fluid_balance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    session_id UUID NOT NULL REFERENCES dialysis_sessions(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    recorded_by UUID NOT NULL REFERENCES users(id),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uf_goal_ml NUMERIC(7,2) NOT NULL,
    uf_achieved_ml NUMERIC(7,2),
    uf_rate_ml_per_hr NUMERIC(6,2),
    fluid_intake_oral_ml NUMERIC(6,2) DEFAULT 0,
    fluid_intake_iv_ml NUMERIC(6,2) DEFAULT 0,
    fluid_output_urine_ml NUMERIC(6,2) DEFAULT 0,
    fluid_output_other_ml NUMERIC(6,2) DEFAULT 0,
    net_fluid_balance_ml NUMERIC(7,2),
    weight_change_kg NUMERIC(5,2),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE session_fluid_balance ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON session_fluid_balance USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_fluid_balance_session ON session_fluid_balance(session_id);
CREATE INDEX idx_fluid_balance_patient ON session_fluid_balance(patient_id);

CREATE TRIGGER trg_session_fluid_balance_updated_at BEFORE UPDATE ON session_fluid_balance FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS session_fluid_balance;
