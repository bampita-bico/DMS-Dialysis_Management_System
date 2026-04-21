-- Seed: Critical Drug Interactions for Dialysis
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 005_drug_interactions.sql

-- Common critical drug interactions in dialysis patients
INSERT INTO drug_interactions (hospital_id, drug_a_id, drug_b_id, severity, interaction_type, clinical_effect, management_recommendation, deleted_at)
SELECT
  :hospital_id AS hospital_id,
  a.id AS drug_a_id,
  b.id AS drug_b_id,
  int_data.severity,
  int_data.interaction_type,
  int_data.clinical_effect,
  int_data.management_recommendation,
  NULL AS deleted_at
FROM medications a
CROSS JOIN medications b
CROSS JOIN (VALUES
  -- ACE Inhibitors + Potassium
  ('Enalapril', 'Calcium Acetate', 'moderate', 'pharmacodynamic', 'Increased risk of hyperkalemia in dialysis patients', 'Monitor potassium levels closely, especially post-dialysis'),
  ('Lisinopril', 'Calcium Acetate', 'moderate', 'pharmacodynamic', 'Risk of hyperkalemia', 'Regular potassium monitoring required'),

  -- Warfarin Interactions
  ('Warfarin', 'Ciprofloxacin', 'major', 'pharmacokinetic', 'Ciprofloxacin increases warfarin effect, risk of bleeding', 'Monitor INR closely, may need warfarin dose reduction'),
  ('Warfarin', 'Omeprazole', 'moderate', 'pharmacokinetic', 'Omeprazole may increase warfarin levels', 'Monitor INR when starting/stopping omeprazole'),
  ('Warfarin', 'Metronidazole', 'major', 'pharmacokinetic', 'Metronidazole significantly increases warfarin effect', 'Avoid combination if possible, otherwise reduce warfarin dose'),

  -- NSAIDs (if added) + Anticoagulants
  ('Warfarin', 'Aspirin', 'major', 'pharmacodynamic', 'Increased bleeding risk', 'Use together only if essential, monitor closely'),

  -- Antibiotics + Medications
  ('Gentamicin', 'Furosemide', 'major', 'pharmacodynamic', 'Increased risk of ototoxicity and nephrotoxicity', 'Avoid combination, monitor renal function and hearing'),
  ('Vancomycin', 'Gentamicin', 'major', 'pharmacodynamic', 'Increased nephrotoxicity risk', 'Monitor drug levels and renal function closely'),

  -- Calcium Supplements + Phosphate Binders
  ('Calcium Carbonate', 'Calcium Acetate', 'moderate', 'pharmacodynamic', 'Risk of hypercalcemia', 'Monitor calcium levels, avoid excessive calcium intake'),

  -- Antihypertensives
  ('Amlodipine', 'Atenolol', 'moderate', 'pharmacodynamic', 'Additive hypotensive effect, risk of bradycardia', 'Monitor BP and heart rate regularly'),
  ('Enalapril', 'Furosemide', 'moderate', 'pharmacodynamic', 'Risk of hypotension, especially first dose', 'Start with low doses, monitor BP'),

  -- Insulin + Oral Hypoglycemics
  ('Insulin Regular', 'Metformin', 'minor', 'pharmacodynamic', 'Additive glucose-lowering effect', 'Monitor blood glucose, adjust doses as needed'),

  -- Iron + Other Medications
  ('Iron Sucrose', 'Calcium Carbonate', 'minor', 'pharmacokinetic', 'Calcium may reduce iron absorption', 'Separate administration by 2 hours'),

  -- ESA + Iron
  ('Erythropoietin Alfa', 'Iron Sucrose', 'beneficial', 'pharmacodynamic', 'Iron supplementation enhances ESA response', 'Ensure adequate iron stores for optimal ESA efficacy'),

  -- Sevelamer Interactions
  ('Sevelamer Hydrochloride', 'Ciprofloxacin', 'moderate', 'pharmacokinetic', 'Sevelamer may reduce ciprofloxacin absorption', 'Give ciprofloxacin 2 hours before sevelamer')

) AS int_data(drug_a_name, drug_b_name, severity, interaction_type, clinical_effect, management_recommendation)
WHERE a.hospital_id = :hospital_id
  AND b.hospital_id = :hospital_id
  AND a.generic_name = int_data.drug_a_name
  AND b.generic_name = int_data.drug_b_name
  AND a.id < b.id  -- Avoid duplicate pairs
ON CONFLICT DO NOTHING;

DO $$
BEGIN
  RAISE NOTICE 'Drug interactions seeded successfully';
END $$;
