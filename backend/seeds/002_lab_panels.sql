-- Seed: Standard Laboratory Panels for Dialysis
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 002_lab_panels.sql

-- Pre-built lab panels for common dialysis workflows
INSERT INTO lab_panels (hospital_id, panel_code, panel_name, description, category, cost, turnaround_time_hours, deleted_at) VALUES
-- Essential Dialysis Panels
(:hospital_id, 'PRE_DIALYSIS', 'Pre-Dialysis Panel', 'Standard pre-dialysis assessment panel for adequacy monitoring', 'routine', 25.00, 4, NULL),
(:hospital_id, 'POST_DIALYSIS', 'Post-Dialysis Panel', 'Post-dialysis urea for URR/Kt/V calculation', 'routine', 10.00, 2, NULL),
(:hospital_id, 'MONTHLY_MON', 'Monthly Monitoring Panel', 'Comprehensive monthly assessment for dialysis patients', 'routine', 80.00, 6, NULL),
(:hospital_id, 'QUARTERLY_MON', 'Quarterly Monitoring Panel', 'Quarterly comprehensive metabolic and hematology panel', 'routine', 120.00, 24, NULL),

-- Specialized Panels
(:hospital_id, 'ANEMIA_PANEL', 'Anemia Workup Panel', 'Complete iron studies and CBC for anemia management', 'specialty', 45.00, 4, NULL),
(:hospital_id, 'BONE_PANEL', 'Bone & Mineral Panel', 'CKD-MBD (Mineral Bone Disorder) assessment', 'specialty', 60.00, 24, NULL),
(:hospital_id, 'HEPATITIS_SCREEN', 'Hepatitis Screening Panel', 'Hepatitis B and C screening for dialysis admission', 'serology', 35.00, 6, NULL),
(:hospital_id, 'INFECTION_SCREEN', 'Infection Screening Panel', 'Full infectious disease screening for new patients', 'serology', 50.00, 6, NULL),
(:hospital_id, 'ADEQUACY_PANEL', 'Dialysis Adequacy Panel', 'URR and Kt/V calculation with pre/post urea', 'routine', 15.00, 2, NULL),
(:hospital_id, 'ACCESS_COMPLIC', 'Vascular Access Complication Panel', 'Workup for access infection or thrombosis', 'urgent', 40.00, 4, NULL),

-- Emergency Panels
(:hospital_id, 'CARDIAC_PANEL', 'Cardiac Markers Panel', 'Troponin, CK-MB, ECG for chest pain/MI', 'urgent', 50.00, 2, NULL),
(:hospital_id, 'SEPSIS_PANEL', 'Sepsis Workup Panel', 'Blood culture, CBC, CRP, lactate for suspected sepsis', 'urgent', 60.00, 4, NULL)

ON CONFLICT (hospital_id, panel_code) DO NOTHING;

-- Link tests to panels (via junction table lab_panel_tests)
-- Note: This assumes lab_panel_tests table exists with schema: (panel_id, test_id, is_required)

-- PRE_DIALYSIS Panel
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'PRE_DIALYSIS'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('UREA', 'CREAT', 'SODIUM', 'POTASSIUM', 'BICARB', 'CALCIUM', 'PHOSPHATE')
ON CONFLICT DO NOTHING;

-- POST_DIALYSIS Panel
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'POST_DIALYSIS'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('UREA', 'POTASSIUM')
ON CONFLICT DO NOTHING;

-- MONTHLY_MON Panel
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  CASE WHEN t.code IN ('UREA', 'CREAT', 'HB', 'POTASSIUM', 'CALCIUM', 'PHOSPHATE', 'ALB') THEN true ELSE false END AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'MONTHLY_MON'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('UREA', 'CREAT', 'EGFR', 'SODIUM', 'POTASSIUM', 'CHLORIDE', 'BICARB', 'CALCIUM', 'PHOSPHATE', 'ALB', 'HB', 'HCT', 'WBC', 'PLT', 'GLUCOSE')
ON CONFLICT DO NOTHING;

-- QUARTERLY_MON Panel
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'QUARTERLY_MON'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('UREA', 'CREAT', 'EGFR', 'SODIUM', 'POTASSIUM', 'CHLORIDE', 'BICARB', 'CALCIUM', 'PHOSPHATE', 'ALP', 'PTH', 'HB', 'HCT', 'RBC', 'WBC', 'PLT', 'MCV', 'FERRITIN', 'TSAT', 'ALB', 'TP', 'ALT', 'AST', 'GLUCOSE', 'HBA1C', 'CHOL', 'TG', 'HDL', 'LDL')
ON CONFLICT DO NOTHING;

-- ANEMIA_PANEL
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'ANEMIA_PANEL'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('HB', 'HCT', 'RBC', 'MCV', 'MCH', 'MCHC', 'IRON', 'TIBC', 'FERRITIN', 'TSAT')
ON CONFLICT DO NOTHING;

-- BONE_PANEL (CKD-MBD)
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'BONE_PANEL'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('CALCIUM', 'CALCIUM_ION', 'PHOSPHATE', 'ALP', 'PTH', 'VIT_D', 'ALB')
ON CONFLICT DO NOTHING;

-- HEPATITIS_SCREEN
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'HEPATITIS_SCREEN'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('HEP_B', 'HEP_C', 'ALT', 'AST')
ON CONFLICT DO NOTHING;

-- INFECTION_SCREEN
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'INFECTION_SCREEN'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('HEP_B', 'HEP_C', 'HIV')
ON CONFLICT DO NOTHING;

-- ADEQUACY_PANEL
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'ADEQUACY_PANEL'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('UREA', 'URR', 'KTV')
ON CONFLICT DO NOTHING;

-- ACCESS_COMPLIC
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'ACCESS_COMPLIC'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('WBC', 'HB', 'PLT', 'BLOOD_CX', 'PT', 'APTT')
ON CONFLICT DO NOTHING;

-- CARDIAC_PANEL
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'CARDIAC_PANEL'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('TROP', 'POTASSIUM', 'SODIUM', 'CALCIUM')
ON CONFLICT DO NOTHING;

-- SEPSIS_PANEL
INSERT INTO lab_panel_tests (panel_id, test_id, is_required)
SELECT
  p.id AS panel_id,
  t.id AS test_id,
  true AS is_required
FROM lab_panels p
CROSS JOIN lab_test_catalog t
WHERE p.hospital_id = :hospital_id
  AND p.panel_code = 'SEPSIS_PANEL'
  AND t.hospital_id = :hospital_id
  AND t.code IN ('BLOOD_CX', 'WBC', 'POTASSIUM', 'CREAT', 'GLUCOSE')
ON CONFLICT DO NOTHING;

-- Display success message
DO $$
BEGIN
  RAISE NOTICE 'Lab panels seeded successfully. Total panels: 12';
END $$;
