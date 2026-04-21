-- Seed: Insurance Schemes for East African Context
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 007_insurance_schemes.sql

INSERT INTO insurance_schemes (hospital_id, scheme_name, scheme_code, scheme_type, contact_person, contact_phone, contact_email, coverage_percentage, copay_amount, requires_preauthorization, claim_submission_method, reimbursement_period_days, is_active, deleted_at) VALUES
-- Kenya
(:hospital_id, 'National Hospital Insurance Fund (NHIF)', 'NHIF-KE', 'government', 'NHIF Contacts', '+254-20-2717000', 'customercare@nhif.or.ke', 100, 0.00, true, 'online', 60, true, NULL),
(:hospital_id, 'Britam Health Insurance', 'BRITAM', 'private', 'Claims Dept', '+254-711-066000', 'medical@britam.com', 80, 20.00, true, 'email', 30, true, NULL),
(:hospital_id, 'AAR Insurance Kenya', 'AAR', 'private', 'Medical Claims', '+254-709-910000', 'medical@aar-healthcare.com', 90, 10.00, true, 'online', 21, true, NULL),
(:hospital_id, 'CIC Insurance', 'CIC', 'private', 'Health Claims', '+254-703-007000', 'health@cic.co.ke', 75, 25.00, true, 'email', 45, true, NULL),
(:hospital_id, 'Jubilee Insurance', 'JUBILEE', 'private', 'Claims Unit', '+254-732-100000', 'claims@jubileekenya.com', 80, 20.00, true, 'online', 30, true, NULL),

-- Uganda
(:hospital_id, 'National Health Insurance Scheme (NHIS)', 'NHIS-UG', 'government', 'NHIS Office', '+256-414-233000', 'info@nhis.ug', 100, 0.00, true, 'online', 90, true, NULL),
(:hospital_id, 'IAA Uganda', 'IAA-UG', 'private', 'Medical Dept', '+256-312-260000', 'medical@iaa.co.ug', 85, 15.00, true, 'email', 45, true, NULL),
(:hospital_id, 'AAR Healthcare Uganda', 'AAR-UG', 'private', 'Claims Team', '+256-417-715000', 'uganda@aar-healthcare.com', 90, 10.00, true, 'online', 21, true, NULL),

-- Tanzania
(:hospital_id, 'National Health Insurance Fund (NHIF-TZ)', 'NHIF-TZ', 'government', 'NHIF Tanzania', '+255-22-2111992', 'nhif@nhif.or.tz', 100, 0.00, true, 'manual', 60, true, NULL),
(:hospital_id, 'Jubilee Insurance Tanzania', 'JUB-TZ', 'private', 'Medical Claims', '+255-22-2113480', 'claims@jubileetz.com', 80, 20.00, true, 'email', 30, true, NULL),

-- Employer Schemes
(:hospital_id, 'Kenya Power Staff Medical Scheme', 'KPLC-MED', 'employer', 'HR Benefits', '+254-711-000000', 'benefits@kplc.co.ke', 100, 0.00, false, 'email', 30, true, NULL),
(:hospital_id, 'Safaricom Staff Medical Cover', 'SCOM-MED', 'employer', 'Staff Benefits', '+254-722-000000', 'benefits@safaricom.co.ke', 100, 0.00, false, 'online', 21, true, NULL),

-- International/NGO
(:hospital_id, 'AMREF Flying Doctors', 'AMREF', 'private', 'Emergency Services', '+254-699-000000', 'medical@flydoc.org', 100, 0.00, true, 'online', 14, true, NULL),
(:hospital_id, 'International SOS', 'ISOS', 'private', 'Regional Office', '+254-20-2711900', 'operations.kenya@internationalsos.com', 100, 0.00, true, 'online', 30, true, NULL),

-- Self-Pay / Cash
(:hospital_id, 'Self Pay (Cash)', 'CASH', 'self_pay', NULL, NULL, NULL, 0, 0.00, false, 'manual', 0, true, NULL),
(:hospital_id, 'Corporate Credit Account', 'CORP-CREDIT', 'employer', NULL, NULL, NULL, 100, 0.00, false, 'manual', 30, true, NULL)

ON CONFLICT (hospital_id, scheme_code) DO NOTHING;

DO $$
BEGIN
  RAISE NOTICE 'Insurance schemes seeded successfully. Total schemes: 16';
END $$;
