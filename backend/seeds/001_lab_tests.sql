-- Seed: Standard Dialysis Laboratory Test Catalog
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 001_lab_tests.sql

-- Common Dialysis Lab Tests for African Dialysis Centers
INSERT INTO lab_test_catalog (hospital_id, code, name, category, specimen_type, turnaround_time_hours, cost, requires_fasting, instructions, deleted_at) VALUES
-- Renal Function Tests
(:hospital_id, 'UREA', 'Urea', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'CREAT', 'Creatinine', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'EGFR', 'eGFR (Estimated Glomerular Filtration Rate)', 'chemistry', 'serum', 2, 0.00, false, 'Calculated from creatinine, age, sex', NULL),
(:hospital_id, 'BUN', 'Blood Urea Nitrogen (BUN)', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'URR', 'Urea Reduction Ratio', 'chemistry', 'serum', 2, 0.00, false, 'Pre and post dialysis urea required', NULL),
(:hospital_id, 'KTV', 'Kt/V (Dialysis Adequacy)', 'chemistry', 'serum', 2, 0.00, false, 'Calculated from URR or direct measurement', NULL),

-- Electrolytes
(:hospital_id, 'SODIUM', 'Sodium (Na+)', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'POTASSIUM', 'Potassium (K+)', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw, avoid hemolysis', NULL),
(:hospital_id, 'CHLORIDE', 'Chloride (Cl-)', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'BICARB', 'Bicarbonate (HCO3-)', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),

-- Bone & Mineral Metabolism
(:hospital_id, 'CALCIUM', 'Calcium (Total)', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'CALCIUM_ION', 'Calcium (Ionized)', 'chemistry', 'serum', 2, 6.00, false, 'Standard blood draw, analyze immediately', NULL),
(:hospital_id, 'PHOSPHATE', 'Phosphate', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'ALP', 'Alkaline Phosphatase (ALP)', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'PTH', 'Parathyroid Hormone (PTH)', 'immunology', 'serum', 24, 25.00, false, 'Fasting preferred, morning sample', NULL),
(:hospital_id, 'VIT_D', 'Vitamin D (25-OH)', 'immunology', 'serum', 48, 30.00, false, 'Standard blood draw', NULL),

-- Complete Blood Count (CBC)
(:hospital_id, 'HB', 'Hemoglobin (Hb)', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'HCT', 'Hematocrit (Hct)', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'RBC', 'Red Blood Cell Count', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'WBC', 'White Blood Cell Count', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'PLT', 'Platelet Count', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'MCV', 'Mean Corpuscular Volume (MCV)', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'MCH', 'Mean Corpuscular Hemoglobin (MCH)', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),
(:hospital_id, 'MCHC', 'Mean Corpuscular Hb Concentration (MCHC)', 'hematology', 'whole_blood', 1, 3.00, false, 'EDTA tube', NULL),

-- Iron Studies
(:hospital_id, 'IRON', 'Serum Iron', 'chemistry', 'serum', 4, 8.00, true, 'Fasting required, morning sample', NULL),
(:hospital_id, 'TIBC', 'Total Iron Binding Capacity (TIBC)', 'chemistry', 'serum', 4, 8.00, true, 'Fasting required, morning sample', NULL),
(:hospital_id, 'FERRITIN', 'Ferritin', 'immunology', 'serum', 4, 12.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'TSAT', 'Transferrin Saturation', 'chemistry', 'serum', 4, 0.00, false, 'Calculated from iron and TIBC', NULL),

-- Liver Function
(:hospital_id, 'ALB', 'Albumin', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'TP', 'Total Protein', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'BILI_T', 'Bilirubin (Total)', 'chemistry', 'serum', 2, 4.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'ALT', 'Alanine Aminotransferase (ALT)', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'AST', 'Aspartate Aminotransferase (AST)', 'chemistry', 'serum', 2, 5.00, false, 'Standard blood draw', NULL),

-- Glucose & Lipids
(:hospital_id, 'GLUCOSE', 'Glucose (Random)', 'chemistry', 'serum', 1, 3.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'GLUCOSE_FAST', 'Glucose (Fasting)', 'chemistry', 'serum', 1, 3.00, true, 'Fasting 8-12 hours', NULL),
(:hospital_id, 'HBA1C', 'HbA1c (Glycated Hemoglobin)', 'chemistry', 'whole_blood', 2, 15.00, false, 'EDTA tube', NULL),
(:hospital_id, 'CHOL', 'Cholesterol (Total)', 'chemistry', 'serum', 2, 6.00, true, 'Fasting 12 hours', NULL),
(:hospital_id, 'TG', 'Triglycerides', 'chemistry', 'serum', 2, 6.00, true, 'Fasting 12 hours', NULL),
(:hospital_id, 'HDL', 'HDL Cholesterol', 'chemistry', 'serum', 2, 7.00, true, 'Fasting 12 hours', NULL),
(:hospital_id, 'LDL', 'LDL Cholesterol', 'chemistry', 'serum', 2, 7.00, true, 'Fasting 12 hours', NULL),

-- Infectious Disease Screening
(:hospital_id, 'HEP_B', 'Hepatitis B Surface Antigen (HBsAg)', 'serology', 'serum', 4, 15.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'HEP_C', 'Hepatitis C Antibody', 'serology', 'serum', 4, 18.00, false, 'Standard blood draw', NULL),
(:hospital_id, 'HIV', 'HIV 1 & 2 Antibody', 'serology', 'serum', 4, 12.00, false, 'Counseling required, consent needed', NULL),

-- Coagulation
(:hospital_id, 'PT', 'Prothrombin Time (PT)', 'hematology', 'plasma', 2, 6.00, false, 'Citrate tube', NULL),
(:hospital_id, 'INR', 'International Normalized Ratio (INR)', 'hematology', 'plasma', 2, 0.00, false, 'Calculated from PT', NULL),
(:hospital_id, 'APTT', 'Activated Partial Thromboplastin Time (aPTT)', 'hematology', 'plasma', 2, 6.00, false, 'Citrate tube', NULL),

-- Urinalysis
(:hospital_id, 'URINE_MICRO', 'Urine Microscopy', 'hematology', 'urine', 2, 4.00, false, 'Mid-stream clean catch', NULL),
(:hospital_id, 'URINE_PROTEIN', 'Urine Protein (24hr)', 'chemistry', 'urine', 4, 5.00, false, '24-hour collection required', NULL),
(:hospital_id, 'UPCR', 'Urine Protein/Creatinine Ratio', 'chemistry', 'urine', 2, 6.00, false, 'Spot urine sample', NULL),

-- Cardiac
(:hospital_id, 'TROP', 'Troponin I', 'immunology', 'serum', 1, 20.00, false, 'Urgent - cardiac marker', NULL),
(:hospital_id, 'BNP', 'B-type Natriuretic Peptide (BNP)', 'immunology', 'serum', 4, 25.00, false, 'Heart failure marker', NULL),

-- Microbiology (Culture)
(:hospital_id, 'BLOOD_CX', 'Blood Culture', 'microbiology', 'whole_blood', 72, 25.00, false, 'Sterile technique, before antibiotics', NULL),
(:hospital_id, 'URINE_CX', 'Urine Culture & Sensitivity', 'microbiology', 'urine', 72, 15.00, false, 'Mid-stream clean catch, sterile container', NULL),
(:hospital_id, 'CVC_CX', 'CVC Tip Culture', 'microbiology', 'tissue', 72, 20.00, false, 'Sterile removal and collection', NULL)

ON CONFLICT (hospital_id, code) DO NOTHING;

-- Display success message
DO $$
BEGIN
  RAISE NOTICE 'Lab test catalog seeded successfully. Total tests: 54';
END $$;
