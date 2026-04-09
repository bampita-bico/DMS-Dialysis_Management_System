-- +goose Up
CREATE TABLE lab_critical_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    result_id UUID NOT NULL REFERENCES lab_results(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patients(id),
    test_name VARCHAR(255) NOT NULL,
    critical_value TEXT NOT NULL,
    reference_range TEXT,
    severity VARCHAR(50) NOT NULL DEFAULT 'high',
    alerted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMPTZ,
    action_taken TEXT,
    doctor_notified BOOLEAN NOT NULL DEFAULT FALSE,
    doctor_notified_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE lab_critical_alerts ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lab_critical_alerts USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_lab_critical_alerts_result ON lab_critical_alerts(result_id);
CREATE INDEX idx_lab_critical_alerts_patient ON lab_critical_alerts(patient_id);
CREATE INDEX idx_lab_critical_alerts_unacknowledged ON lab_critical_alerts(hospital_id) WHERE acknowledged_at IS NULL;
CREATE INDEX idx_lab_critical_alerts_alerted_at ON lab_critical_alerts(hospital_id, alerted_at DESC);

CREATE TRIGGER trg_lab_critical_alerts_updated_at BEFORE UPDATE ON lab_critical_alerts FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS lab_critical_alerts;
