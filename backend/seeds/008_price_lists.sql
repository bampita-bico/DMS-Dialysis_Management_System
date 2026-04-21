-- Seed: Standard Price List for Dialysis Services
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 008_price_lists.sql

INSERT INTO price_lists (hospital_id, service_code, service_name, service_category, unit_price, currency, effective_date, expiry_date, is_active, deleted_at) VALUES
-- Dialysis Sessions
(:hospital_id, 'HD-SESS', 'Hemodialysis Session (4 hours)', 'dialysis', 8000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'HDF-SESS', 'Hemodiafiltration Session (4 hours)', 'dialysis', 10000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'PD-CAPD', 'CAPD Exchange (per day)', 'dialysis', 3000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'PD-APD', 'APD Session (overnight)', 'dialysis', 4500.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Dialysis-Related Procedures
(:hospital_id, 'CVC-INSERT', 'Central Venous Catheter Insertion', 'procedure', 15000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'AVF-CREATION', 'AV Fistula Creation', 'procedure', 35000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'AVG-INSERT', 'AV Graft Insertion', 'procedure', 45000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'PD-CATH', 'Peritoneal Dialysis Catheter Insertion', 'procedure', 40000.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Consultations
(:hospital_id, 'CONS-NEPH', 'Nephrologist Consultation', 'consultation', 3000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'CONS-DOC', 'Doctor Consultation', 'consultation', 2000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'CONS-NURSE', 'Nurse Consultation', 'consultation', 1000.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Lab Tests (Common)
(:hospital_id, 'LAB-CBC', 'Complete Blood Count', 'laboratory', 800.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-UREA', 'Urea & Electrolytes', 'laboratory', 1500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-CREAT', 'Creatinine', 'laboratory', 500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-HB', 'Hemoglobin', 'laboratory', 300.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-IRON', 'Iron Studies (Full)', 'laboratory', 2500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-PTH', 'Parathyroid Hormone (PTH)', 'laboratory', 5000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-HBA1C', 'HbA1c', 'laboratory', 1500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-HEPB', 'Hepatitis B Surface Antigen', 'laboratory', 1500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-HEPC', 'Hepatitis C Antibody', 'laboratory', 2000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'LAB-HIV', 'HIV Test', 'laboratory', 1000.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Imaging
(:hospital_id, 'IMG-XRAY', 'X-Ray (Single View)', 'imaging', 1500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'IMG-US', 'Ultrasound Scan', 'imaging', 3000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'IMG-ECHO', 'Echocardiogram', 'imaging', 5000.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'IMG-FISTULOGRAM', 'Fistulogram', 'imaging', 8000.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Medications (Samples - per unit)
(:hospital_id, 'MED-EPO-4K', 'Erythropoietin 4000 IU Injection', 'medication', 1200.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'MED-IRON-IV', 'Iron Sucrose 100mg IV', 'medication', 800.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'MED-HEPARIN', 'Heparin 5000 IU Vial', 'medication', 300.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'MED-PHOSLO', 'Calcium Acetate 667mg (30 tabs)', 'medication', 1500.00, 'KES', CURRENT_DATE, NULL, true, NULL),

-- Administrative Fees
(:hospital_id, 'ADM-REG', 'Patient Registration Fee', 'administrative', 500.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'ADM-FILE', 'File Opening Fee', 'administrative', 200.00, 'KES', CURRENT_DATE, NULL, true, NULL),
(:hospital_id, 'ADM-REPORT', 'Medical Report', 'administrative', 1000.00, 'KES', CURRENT_DATE, NULL, true, NULL)

ON CONFLICT (hospital_id, service_code) DO NOTHING;

DO $$
BEGIN
  RAISE NOTICE 'Price list seeded successfully. Total services: 33';
END $$;
