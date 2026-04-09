-- +goose Up
CREATE TYPE staff_cadre AS ENUM (
    'doctor',
    'nephrologist',
    'nurse',
    'dialysis_technician',
    'pharmacist',
    'lab_technician',
    'radiologist',
    'social_worker',
    'nutritionist',
    'admin',
    'receptionist',
    'accountant',
    'other'
);

-- Note: shift_type already exists from migration 028 (clinical enums)
-- Adding 'on_call' value to existing shift_type enum
ALTER TYPE shift_type ADD VALUE IF NOT EXISTS 'on_call';

CREATE TYPE schedule_type AS ENUM (
    'fixed',
    'rotating',
    'on_call',
    'part_time'
);

CREATE TYPE leave_type AS ENUM (
    'annual',
    'sick',
    'maternity',
    'paternity',
    'compassionate',
    'study',
    'unpaid'
);

CREATE TYPE leave_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'cancelled'
);

CREATE TYPE transport_status AS ENUM (
    'scheduled',
    'en_route_pickup',
    'picked_up',
    'en_route_hospital',
    'arrived',
    'en_route_dropoff',
    'completed',
    'cancelled'
);

-- +goose Down
DROP TYPE IF EXISTS transport_status;
DROP TYPE IF EXISTS leave_status;
DROP TYPE IF EXISTS leave_type;
DROP TYPE IF EXISTS schedule_type;
-- shift_type is owned by migration 028, not dropped here
DROP TYPE IF EXISTS staff_cadre;
