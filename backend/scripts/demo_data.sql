-- Demo Hospital
INSERT INTO hospitals (id, name, short_code, tier, region, country, address, phone, email, subscription_plan)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'Demo Dialysis Center',
    'DEMO',
    'private',
    'Central',
    'Uganda',
    '123 Medical Plaza, Kampala',
    '+256700000000',
    'admin@demodialysis.com',
    'enterprise'
);

-- Demo User (Doctor)
INSERT INTO users (id, hospital_id, email, phone, password_hash, full_name, is_active, is_verified)
VALUES (
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'doctor@demo.com',
    '+256701234567',
    '$2a$10$F0.jeR5ZS/ToEgKN4GVx7ON9gKdUH/i/Jq5RBEqwwLWg3Shr4hSUW',
    'Dr. Demo User',
    true,
    true
);

-- Role for Doctor
INSERT INTO roles (id, hospital_id, name, description, permissions)
VALUES (
    'e0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'Doctor',
    'Medical Doctor with full clinical access',
    '{"patients": ["read", "write"], "sessions": ["read", "write"], "prescriptions": ["read", "write"]}'::jsonb
);

-- User Role Assignment
INSERT INTO user_roles (id, user_id, hospital_id, role_id)
VALUES (
    'c0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'e0000000-0000-0000-0000-000000000001'
);

-- Staff Profile for the doctor
INSERT INTO staff_profiles (id, user_id, hospital_id, employee_number, cadre, specialization, is_active)
VALUES (
    'd0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'DOC001',
    'doctor',
    'Nephrology',
    true
);
