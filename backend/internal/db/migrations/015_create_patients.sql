-- +goose Up
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    mrn VARCHAR(50) NOT NULL,
    national_id VARCHAR(100),
    full_name VARCHAR(255) NOT NULL,
    preferred_name VARCHAR(100),
    date_of_birth DATE NOT NULL,
    sex sex_type NOT NULL,
    blood_type blood_type NOT NULL DEFAULT 'unknown',
    marital_status marital_status DEFAULT 'unknown',
    nationality VARCHAR(100) DEFAULT 'Ugandan',
    religion VARCHAR(100),
    occupation VARCHAR(150),
    education_level VARCHAR(100),
    photo_url TEXT,
    primary_language VARCHAR(50) DEFAULT 'English',
    interpreter_needed BOOLEAN NOT NULL DEFAULT FALSE,
    registration_date DATE NOT NULL DEFAULT CURRENT_DATE,
    registered_by UUID NOT NULL REFERENCES users(id),
    primary_doctor_id UUID REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    deceased_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(hospital_id, mrn)
);

ALTER TABLE patients ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON patients USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_patients_hospital ON patients(hospital_id);
CREATE INDEX idx_patients_mrn ON patients(hospital_id, mrn);
CREATE INDEX idx_patients_national_id ON patients(national_id) WHERE national_id IS NOT NULL;
CREATE INDEX idx_patients_name ON patients(full_name);
CREATE INDEX idx_patients_active ON patients(hospital_id, is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_patients_updated_at BEFORE UPDATE ON patients FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS patients;
