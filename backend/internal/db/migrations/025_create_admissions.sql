-- +goose Up
CREATE TYPE admission_type AS ENUM ('elective','emergency','transfer');
CREATE TYPE discharge_type AS ENUM ('recovered','referred','absconded','deceased','against_advice');

CREATE TABLE admissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    admission_type admission_type NOT NULL DEFAULT 'elective',
    admitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    admitted_by UUID NOT NULL REFERENCES users(id),
    ward VARCHAR(100),
    bed_number VARCHAR(20),
    primary_diagnosis TEXT,
    admitting_doctor_id UUID REFERENCES users(id),
    discharged_at TIMESTAMPTZ,
    discharge_type discharge_type,
    discharge_summary TEXT,
    discharged_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE admissions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON admissions USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

ALTER TABLE diagnoses ADD CONSTRAINT fk_diagnoses_admission FOREIGN KEY (admission_id) REFERENCES admissions(id);

CREATE INDEX idx_admissions_patient ON admissions(patient_id);
CREATE INDEX idx_admissions_active ON admissions(hospital_id) WHERE discharged_at IS NULL AND deleted_at IS NULL;

CREATE TRIGGER trg_admissions_updated_at BEFORE UPDATE ON admissions FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
ALTER TABLE diagnoses DROP CONSTRAINT IF EXISTS fk_diagnoses_admission;
DROP TABLE IF EXISTS admissions;
DROP TYPE IF EXISTS discharge_type;
DROP TYPE IF EXISTS admission_type;
