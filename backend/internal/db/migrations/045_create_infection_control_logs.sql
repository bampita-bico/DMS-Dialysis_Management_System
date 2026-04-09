-- +goose Up
CREATE TABLE infection_control_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    machine_id UUID NOT NULL REFERENCES dialysis_machines(id),
    session_id UUID REFERENCES dialysis_sessions(id),
    activity_type VARCHAR(100) NOT NULL,
    performed_by UUID NOT NULL REFERENCES users(id),
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    disinfectant_used VARCHAR(200),
    disinfectant_concentration VARCHAR(100),
    contact_time_mins INTEGER,
    rinse_cycles_count INTEGER,
    bacterial_test_done BOOLEAN NOT NULL DEFAULT FALSE,
    bacterial_result water_test_result DEFAULT 'pending',
    cfu_count NUMERIC(8,2),
    was_hbv_patient BOOLEAN NOT NULL DEFAULT FALSE,
    bleach_disinfection_done BOOLEAN NOT NULL DEFAULT FALSE,
    external_surfaces_cleaned BOOLEAN NOT NULL DEFAULT FALSE,
    chair_cleaned BOOLEAN NOT NULL DEFAULT FALSE,
    machine_ready_for_use BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE infection_control_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON infection_control_logs USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_infection_logs_machine ON infection_control_logs(machine_id, performed_at DESC);
CREATE INDEX idx_infection_logs_session ON infection_control_logs(session_id);
CREATE INDEX idx_infection_logs_hbv ON infection_control_logs(hospital_id) WHERE was_hbv_patient = TRUE;
CREATE INDEX idx_infection_logs_unverified ON infection_control_logs(hospital_id) WHERE machine_ready_for_use = FALSE;

CREATE TRIGGER trg_infection_control_logs_updated_at BEFORE UPDATE ON infection_control_logs FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS infection_control_logs;
