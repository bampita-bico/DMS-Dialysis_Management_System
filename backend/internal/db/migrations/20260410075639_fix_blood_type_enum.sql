-- +goose Up
-- Fix blood_type enum to have SQL-friendly values without special characters
-- This resolves sqlc duplicate constant generation issue

-- Create new enum with proper names
CREATE TYPE blood_type_new AS ENUM (
    'a_positive',
    'a_negative',
    'b_positive',
    'b_negative',
    'ab_positive',
    'ab_negative',
    'o_positive',
    'o_negative',
    'unknown'
);

-- Update patients table
-- First drop the default
ALTER TABLE patients ALTER COLUMN blood_type DROP DEFAULT;

-- Then change the type
ALTER TABLE patients
    ALTER COLUMN blood_type TYPE blood_type_new
    USING CASE blood_type::text
        WHEN 'A+' THEN 'a_positive'::blood_type_new
        WHEN 'A-' THEN 'a_negative'::blood_type_new
        WHEN 'B+' THEN 'b_positive'::blood_type_new
        WHEN 'B-' THEN 'b_negative'::blood_type_new
        WHEN 'AB+' THEN 'ab_positive'::blood_type_new
        WHEN 'AB-' THEN 'ab_negative'::blood_type_new
        WHEN 'O+' THEN 'o_positive'::blood_type_new
        WHEN 'O-' THEN 'o_negative'::blood_type_new
        ELSE 'unknown'::blood_type_new
    END;

-- Set the new default
ALTER TABLE patients ALTER COLUMN blood_type SET DEFAULT 'unknown'::blood_type_new;

-- Update staff_profiles table (uses VARCHAR, but for consistency)
-- First convert VARCHAR to the new enum
ALTER TABLE staff_profiles
    ALTER COLUMN blood_type TYPE blood_type_new
    USING CASE blood_type
        WHEN 'A+' THEN 'a_positive'::blood_type_new
        WHEN 'A-' THEN 'a_negative'::blood_type_new
        WHEN 'B+' THEN 'b_positive'::blood_type_new
        WHEN 'B-' THEN 'b_negative'::blood_type_new
        WHEN 'AB+' THEN 'ab_positive'::blood_type_new
        WHEN 'AB-' THEN 'ab_negative'::blood_type_new
        WHEN 'O+' THEN 'o_positive'::blood_type_new
        WHEN 'O-' THEN 'o_negative'::blood_type_new
        ELSE 'unknown'::blood_type_new
    END;

-- Drop old enum
DROP TYPE blood_type;

-- Rename new enum to blood_type
ALTER TYPE blood_type_new RENAME TO blood_type;

COMMENT ON TYPE blood_type IS 'Blood type enum with SQL-friendly identifiers (a_positive, b_negative, etc.)';

-- +goose Down
-- Restore original enum with special characters
CREATE TYPE blood_type_old AS ENUM ('A+','A-','B+','B-','AB+','AB-','O+','O-','unknown');

-- Revert patients table
ALTER TABLE patients
    ALTER COLUMN blood_type TYPE blood_type_old
    USING CASE blood_type::text
        WHEN 'a_positive' THEN 'A+'::blood_type_old
        WHEN 'a_negative' THEN 'A-'::blood_type_old
        WHEN 'b_positive' THEN 'B+'::blood_type_old
        WHEN 'b_negative' THEN 'B-'::blood_type_old
        WHEN 'ab_positive' THEN 'AB+'::blood_type_old
        WHEN 'ab_negative' THEN 'AB-'::blood_type_old
        WHEN 'o_positive' THEN 'O+'::blood_type_old
        WHEN 'o_negative' THEN 'O-'::blood_type_old
        ELSE 'unknown'::blood_type_old
    END;

-- Revert staff_profiles table
ALTER TABLE staff_profiles
    ALTER COLUMN blood_type TYPE VARCHAR(10)
    USING CASE blood_type::text
        WHEN 'a_positive' THEN 'A+'
        WHEN 'a_negative' THEN 'A-'
        WHEN 'b_positive' THEN 'B+'
        WHEN 'b_negative' THEN 'B-'
        WHEN 'ab_positive' THEN 'AB+'
        WHEN 'ab_negative' THEN 'AB-'
        WHEN 'o_positive' THEN 'O+'
        WHEN 'o_negative' THEN 'O-'
        ELSE 'unknown'
    END;

-- Drop new enum and restore old one
DROP TYPE blood_type;
ALTER TYPE blood_type_old RENAME TO blood_type;
