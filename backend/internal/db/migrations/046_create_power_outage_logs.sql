-- +goose Up
CREATE TABLE power_outage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    outage_start TIMESTAMPTZ NOT NULL,
    outage_end TIMESTAMPTZ,
    duration_mins INTEGER,
    affected_sessions UUID[],
    sessions_terminated_count INTEGER DEFAULT 0,
    sessions_paused_count INTEGER DEFAULT 0,
    generator_available BOOLEAN NOT NULL DEFAULT FALSE,
    generator_start_delay_mins INTEGER,
    backup_power_duration_mins INTEGER,
    logged_by UUID NOT NULL REFERENCES users(id),
    incident_severity VARCHAR(50) NOT NULL,
    patient_safety_impact TEXT,
    equipment_damage TEXT,
    data_loss_occurred BOOLEAN NOT NULL DEFAULT FALSE,
    recovery_actions TEXT,
    utility_company_notified BOOLEAN NOT NULL DEFAULT FALSE,
    incident_report_filed BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE power_outage_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON power_outage_logs USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_power_outage_hospital_date ON power_outage_logs(hospital_id, outage_start DESC);
CREATE INDEX idx_power_outage_ongoing ON power_outage_logs(hospital_id) WHERE outage_end IS NULL;
CREATE INDEX idx_power_outage_severe ON power_outage_logs(hospital_id) WHERE incident_severity IN ('severe', 'life_threatening');

CREATE TRIGGER trg_power_outage_logs_updated_at BEFORE UPDATE ON power_outage_logs FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS power_outage_logs;
