-- +goose Up
CREATE TABLE national_registry_sync (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    registry_name VARCHAR(255) NOT NULL,
    registry_id VARCHAR(100),
    sync_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_by UUID REFERENCES users(id),
    sync_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_retry_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_national_registry_sync_patient ON national_registry_sync(patient_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_national_registry_sync_hospital ON national_registry_sync(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_national_registry_sync_status ON national_registry_sync(hospital_id, sync_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_national_registry_sync_registry ON national_registry_sync(hospital_id, registry_name) WHERE deleted_at IS NULL;
CREATE INDEX idx_national_registry_sync_date ON national_registry_sync(hospital_id, synced_at) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_national_registry_sync_updated_at
BEFORE UPDATE ON national_registry_sync
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE national_registry_sync ENABLE ROW LEVEL SECURITY;
CREATE POLICY national_registry_sync_isolation ON national_registry_sync
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS national_registry_sync CASCADE;
