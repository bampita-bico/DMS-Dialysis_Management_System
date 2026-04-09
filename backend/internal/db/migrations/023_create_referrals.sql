-- +goose Up
CREATE TYPE referral_direction AS ENUM ('inbound','outbound');
CREATE TYPE referral_status AS ENUM ('pending','accepted','rejected','completed','cancelled');
CREATE TYPE referral_urgency AS ENUM ('routine','urgent','emergency');

CREATE TABLE referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id UUID NOT NULL REFERENCES hospitals(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    direction referral_direction NOT NULL,
    from_facility VARCHAR(255),
    to_facility VARCHAR(255),
    from_doctor VARCHAR(255),
    to_doctor_id UUID REFERENCES users(id),
    reason TEXT NOT NULL,
    urgency referral_urgency NOT NULL DEFAULT 'routine',
    status referral_status NOT NULL DEFAULT 'pending',
    referral_date DATE NOT NULL DEFAULT CURRENT_DATE,
    response_date DATE,
    response_notes TEXT,
    referral_letter_url TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE referrals ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON referrals USING (hospital_id = current_setting('app.current_hospital_id')::UUID);

CREATE INDEX idx_referrals_patient ON referrals(patient_id);
CREATE INDEX idx_referrals_status ON referrals(status) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_referrals_updated_at BEFORE UPDATE ON referrals FOR EACH ROW EXECUTE FUNCTION dms_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS referrals;
DROP TYPE IF EXISTS referral_urgency;
DROP TYPE IF EXISTS referral_status;
DROP TYPE IF EXISTS referral_direction;
