-- +goose Up
CREATE TABLE staff_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id),
    cadre staff_cadre NOT NULL,
    license_number VARCHAR(100),
    license_expiry_date DATE,
    registration_body VARCHAR(255),
    specialization VARCHAR(255),
    years_of_experience INTEGER,
    hire_date DATE,
    contract_end_date DATE,
    employee_number VARCHAR(50),
    emergency_contact_name VARCHAR(255),
    emergency_contact_phone VARCHAR(50),
    blood_type VARCHAR(10),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_staff_profiles_user ON staff_profiles(hospital_id, user_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_staff_profiles_employee ON staff_profiles(hospital_id, employee_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_profiles_hospital ON staff_profiles(hospital_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_profiles_department ON staff_profiles(department_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_profiles_cadre ON staff_profiles(hospital_id, cadre) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_profiles_active ON staff_profiles(hospital_id, is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_staff_profiles_updated_at
BEFORE UPDATE ON staff_profiles
FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

ALTER TABLE staff_profiles ENABLE ROW LEVEL SECURITY;
CREATE POLICY staff_profiles_isolation ON staff_profiles
  USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

-- +goose Down
DROP TABLE IF EXISTS staff_profiles CASCADE;
