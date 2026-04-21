-- Seed: Lab Reference Ranges
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 003_lab_reference_ranges.sql

-- Standard reference ranges for dialysis lab tests
-- Ranges vary by age group and sex where applicable

INSERT INTO lab_reference_ranges (hospital_id, test_id, age_group, sex, min_value, max_value, unit, critical_low, critical_high, deleted_at)
SELECT
  :hospital_id AS hospital_id,
  t.id AS test_id,
  range_data.age_group,
  range_data.sex,
  range_data.min_value,
  range_data.max_value,
  range_data.unit,
  range_data.critical_low,
  range_data.critical_high,
  NULL AS deleted_at
FROM lab_test_catalog t
CROSS JOIN (VALUES
  -- Renal Function
  ('UREA', 'adult', 'all', 2.5, 7.8, 'mmol/L', 1.0, 35.0),
  ('CREAT', 'adult', 'male', 62, 115, 'µmol/L', 30, 800),
  ('CREAT', 'adult', 'female', 53, 97, 'µmol/L', 30, 800),

  -- Electrolytes
  ('SODIUM', 'adult', 'all', 135, 145, 'mmol/L', 120, 160),
  ('POTASSIUM', 'adult', 'all', 3.5, 5.0, 'mmol/L', 2.5, 6.5),
  ('CHLORIDE', 'adult', 'all', 98, 107, 'mmol/L', 80, 120),
  ('BICARB', 'adult', 'all', 22, 29, 'mmol/L', 10, 40),

  -- Bone & Mineral
  ('CALCIUM', 'adult', 'all', 2.10, 2.55, 'mmol/L', 1.5, 3.5),
  ('PHOSPHATE', 'adult', 'all', 0.80, 1.45, 'mmol/L', 0.3, 3.0),
  ('ALP', 'adult', 'all', 30, 130, 'U/L', NULL, 500),

  -- Hematology
  ('HB', 'adult', 'male', 130, 180, 'g/L', 70, 200),
  ('HB', 'adult', 'female', 115, 165, 'g/L', 70, 200),
  ('HCT', 'adult', 'male', 0.40, 0.54, 'L/L', 0.20, 0.60),
  ('HCT', 'adult', 'female', 0.36, 0.48, 'L/L', 0.20, 0.60),
  ('WBC', 'adult', 'all', 4.0, 11.0, '10^9/L', 1.0, 30.0),
  ('PLT', 'adult', 'all', 150, 400, '10^9/L', 50, 1000),

  -- Iron Studies
  ('IRON', 'adult', 'all', 10, 30, 'µmol/L', 5, 50),
  ('FERRITIN', 'adult', 'all', 200, 500, 'µg/L', 50, 1000),
  ('TSAT', 'adult', 'all', 20, 50, '%', 10, 80),

  -- Liver Function
  ('ALB', 'adult', 'all', 35, 50, 'g/L', 20, 60),
  ('ALT', 'adult', 'all', 0, 55, 'U/L', NULL, 300),
  ('AST', 'adult', 'all', 0, 48, 'U/L', NULL, 300),

  -- Glucose & Lipids
  ('GLUCOSE', 'adult', 'all', 3.5, 5.5, 'mmol/L', 2.2, 30.0),
  ('HBA1C', 'adult', 'all', 4.0, 6.0, '%', NULL, 15.0),
  ('CHOL', 'adult', 'all', 3.0, 5.2, 'mmol/L', NULL, 10.0)
) AS range_data(test_code, age_group, sex, min_value, max_value, unit, critical_low, critical_high)
WHERE t.hospital_id = :hospital_id AND t.code = range_data.test_code
ON CONFLICT DO NOTHING;

DO $$
BEGIN
  RAISE NOTICE 'Lab reference ranges seeded successfully';
END $$;
