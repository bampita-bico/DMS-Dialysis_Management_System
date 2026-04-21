-- Seed: Standard Dialysis Medications
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 004_medications.sql

-- Essential medications for African dialysis centers (Kenya/Uganda focused)
INSERT INTO medications (hospital_id, generic_name, brand_name, medication_form, route, strength, unit, category, requires_prescription, is_controlled, reorder_level, storage_conditions, deleted_at) VALUES

-- Erythropoiesis-Stimulating Agents (ESAs)
(:hospital_id, 'Erythropoietin Alfa', 'Eprex', 'injection', 'iv', '4000', 'IU/ml', 'Erythropoiesis-stimulating agent', true, false, 20, 'Refrigerate 2-8°C, protect from light', NULL),
(:hospital_id, 'Erythropoietin Alfa', 'Eprex', 'injection', 'sc', '2000', 'IU/ml', 'Erythropoiesis-stimulating agent', true, false, 30, 'Refrigerate 2-8°C, protect from light', NULL),
(:hospital_id, 'Darbepoetin Alfa', 'Aranesp', 'injection', 'iv', '40', 'mcg/ml', 'Erythropoiesis-stimulating agent', true, false, 10, 'Refrigerate 2-8°C, protect from light', NULL),
(:hospital_id, 'Erythropoietin Beta', 'Recormon', 'injection', 'sc', '2000', 'IU', 'Erythropoiesis-stimulating agent', true, false, 25, 'Refrigerate 2-8°C', NULL),

-- Iron Preparations
(:hospital_id, 'Iron Sucrose', 'Venofer', 'infusion', 'iv', '100', 'mg/5ml', 'Iron supplement', true, false, 50, 'Room temperature, protect from light', NULL),
(:hospital_id, 'Iron Dextran', 'Imferon', 'injection', 'iv', '50', 'mg/ml', 'Iron supplement', true, false, 30, 'Room temperature', NULL),
(:hospital_id, 'Ferrous Sulphate', 'Fefol', 'tablet', 'oral', '200', 'mg', 'Iron supplement', false, false, 500, 'Room temperature, dry place', NULL),
(:hospital_id, 'Ferrous Fumarate', 'Galfer', 'capsule', 'oral', '305', 'mg', 'Iron supplement', false, false, 400, 'Room temperature, dry place', NULL),

-- Phosphate Binders
(:hospital_id, 'Calcium Acetate', 'PhosLo', 'tablet', 'oral', '667', 'mg', 'Phosphate binder', true, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Calcium Carbonate', 'Calcichew', 'tablet', 'oral', '1250', 'mg', 'Phosphate binder', false, false, 500, 'Room temperature, dry place', NULL),
(:hospital_id, 'Sevelamer Hydrochloride', 'Renagel', 'tablet', 'oral', '800', 'mg', 'Phosphate binder', true, false, 200, 'Room temperature, dry place', NULL),
(:hospital_id, 'Sevelamer Carbonate', 'Renvela', 'tablet', 'oral', '800', 'mg', 'Phosphate binder', true, false, 150, 'Room temperature, dry place', NULL),
(:hospital_id, 'Lanthanum Carbonate', 'Fosrenol', 'tablet', 'oral', '500', 'mg', 'Phosphate binder', true, false, 100, 'Room temperature, dry place', NULL),

-- Vitamin D Analogs
(:hospital_id, 'Calcitriol', 'Rocaltrol', 'capsule', 'oral', '0.25', 'mcg', 'Vitamin D analog', true, false, 200, 'Room temperature, protect from light', NULL),
(:hospital_id, 'Alfacalcidol', 'One-Alpha', 'capsule', 'oral', '0.25', 'mcg', 'Vitamin D analog', true, false, 150, 'Room temperature, protect from light', NULL),
(:hospital_id, 'Cholecalciferol', 'Vit D3', 'tablet', 'oral', '1000', 'IU', 'Vitamin D supplement', false, false, 300, 'Room temperature', NULL),

-- Antihypertensives - ACE Inhibitors
(:hospital_id, 'Enalapril', 'Renitec', 'tablet', 'oral', '10', 'mg', 'ACE inhibitor', true, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Lisinopril', 'Zestril', 'tablet', 'oral', '10', 'mg', 'ACE inhibitor', true, false, 250, 'Room temperature, dry place', NULL),
(:hospital_id, 'Ramipril', 'Tritace', 'capsule', 'oral', '5', 'mg', 'ACE inhibitor', true, false, 200, 'Room temperature, dry place', NULL),

-- Antihypertensives - ARBs
(:hospital_id, 'Losartan', 'Cozaar', 'tablet', 'oral', '50', 'mg', 'Angiotensin receptor blocker', true, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Valsartan', 'Diovan', 'tablet', 'oral', '80', 'mg', 'Angiotensin receptor blocker', true, false, 200, 'Room temperature, dry place', NULL),
(:hospital_id, 'Telmisartan', 'Micardis', 'tablet', 'oral', '40', 'mg', 'Angiotensin receptor blocker', true, false, 150, 'Room temperature, dry place', NULL),

-- Antihypertensives - Calcium Channel Blockers
(:hospital_id, 'Amlodipine', 'Norvasc', 'tablet', 'oral', '5', 'mg', 'Calcium channel blocker', true, false, 400, 'Room temperature, dry place', NULL),
(:hospital_id, 'Nifedipine SR', 'Adalat', 'tablet', 'oral', '20', 'mg', 'Calcium channel blocker', true, false, 300, 'Room temperature, dry place', NULL),

-- Antihypertensives - Beta Blockers
(:hospital_id, 'Atenolol', 'Tenormin', 'tablet', 'oral', '50', 'mg', 'Beta blocker', true, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Metoprolol', 'Lopressor', 'tablet', 'oral', '50', 'mg', 'Beta blocker', true, false, 250, 'Room temperature, dry place', NULL),
(:hospital_id, 'Carvedilol', 'Coreg', 'tablet', 'oral', '6.25', 'mg', 'Beta blocker', true, false, 200, 'Room temperature, dry place', NULL),

-- Diuretics
(:hospital_id, 'Furosemide', 'Lasix', 'tablet', 'oral', '40', 'mg', 'Loop diuretic', true, false, 500, 'Room temperature, dry place', NULL),
(:hospital_id, 'Furosemide', 'Lasix', 'injection', 'iv', '10', 'mg/ml', 'Loop diuretic', true, false, 100, 'Room temperature, protect from light', NULL),
(:hospital_id, 'Hydrochlorothiazide', 'HCTZ', 'tablet', 'oral', '25', 'mg', 'Thiazide diuretic', true, false, 300, 'Room temperature, dry place', NULL),

-- Anticoagulants
(:hospital_id, 'Heparin Sodium', 'Heparin', 'injection', 'iv', '5000', 'IU/ml', 'Anticoagulant', true, false, 100, 'Room temperature', NULL),
(:hospital_id, 'Enoxaparin', 'Clexane', 'injection', 'sc', '40', 'mg', 'Low molecular weight heparin', true, false, 50, 'Room temperature', NULL),
(:hospital_id, 'Warfarin', 'Coumadin', 'tablet', 'oral', '5', 'mg', 'Anticoagulant', true, false, 200, 'Room temperature, dry place', NULL),

-- Antibiotics (Common for Dialysis-Related Infections)
(:hospital_id, 'Ceftriaxone', 'Rocephin', 'injection', 'iv', '1', 'g', 'Antibiotic', true, false, 50, 'Room temperature before reconstitution', NULL),
(:hospital_id, 'Cefuroxime', 'Zinacef', 'injection', 'iv', '750', 'mg', 'Antibiotic', true, false, 50, 'Room temperature', NULL),
(:hospital_id, 'Vancomycin', 'Vancocin', 'infusion', 'iv', '500', 'mg', 'Antibiotic', true, false, 30, 'Room temperature, reconstituted refrigerate', NULL),
(:hospital_id, 'Gentamicin', 'Garamycin', 'injection', 'iv', '80', 'mg/2ml', 'Antibiotic', true, false, 40, 'Room temperature', NULL),
(:hospital_id, 'Ciprofloxacin', 'Cipro', 'tablet', 'oral', '500', 'mg', 'Antibiotic', true, false, 200, 'Room temperature, dry place', NULL),

-- Antifungals
(:hospital_id, 'Fluconazole', 'Diflucan', 'capsule', 'oral', '150', 'mg', 'Antifungal', true, false, 100, 'Room temperature, dry place', NULL),

-- Antiemetics
(:hospital_id, 'Metoclopramide', 'Reglan', 'tablet', 'oral', '10', 'mg', 'Antiemetic', true, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Ondansetron', 'Zofran', 'injection', 'iv', '4', 'mg/2ml', 'Antiemetic', true, false, 50, 'Room temperature, protect from light', NULL),

-- Antacids & GI Protection
(:hospital_id, 'Omeprazole', 'Losec', 'capsule', 'oral', '20', 'mg', 'Proton pump inhibitor', true, false, 400, 'Room temperature, dry place', NULL),
(:hospital_id, 'Ranitidine', 'Zantac', 'tablet', 'oral', '150', 'mg', 'H2 blocker', false, false, 300, 'Room temperature, dry place', NULL),
(:hospital_id, 'Aluminium Hydroxide', 'Alu-Cap', 'tablet', 'oral', '500', 'mg', 'Antacid', false, false, 400, 'Room temperature, dry place', NULL),

-- Antidiabetics
(:hospital_id, 'Insulin Regular', 'Actrapid', 'injection', 'sc', '100', 'IU/ml', 'Insulin', true, false, 50, 'Refrigerate 2-8°C', NULL),
(:hospital_id, 'Insulin Glargine', 'Lantus', 'injection', 'sc', '100', 'IU/ml', 'Long-acting insulin', true, false, 30, 'Refrigerate 2-8°C', NULL),
(:hospital_id, 'Metformin', 'Glucophage', 'tablet', 'oral', '500', 'mg', 'Antidiabetic', true, false, 500, 'Room temperature, dry place', NULL),

-- Analgesics
(:hospital_id, 'Paracetamol', 'Panadol', 'tablet', 'oral', '500', 'mg', 'Analgesic', false, false, 1000, 'Room temperature, dry place', NULL),
(:hospital_id, 'Tramadol', 'Tramal', 'capsule', 'oral', '50', 'mg', 'Opioid analgesic', true, true, 200, 'Room temperature, controlled substance', NULL),

-- Vitamins & Supplements
(:hospital_id, 'Folic Acid', 'Folic', 'tablet', 'oral', '5', 'mg', 'Vitamin supplement', false, false, 500, 'Room temperature, dry place', NULL),
(:hospital_id, 'Vitamin B Complex', 'B-Complex', 'tablet', 'oral', '1', 'tab', 'Vitamin supplement', false, false, 400, 'Room temperature, dry place', NULL),
(:hospital_id, 'Ascorbic Acid', 'Vitamin C', 'tablet', 'oral', '500', 'mg', 'Vitamin supplement', false, false, 300, 'Room temperature, dry place', NULL),

-- Emergency Medications
(:hospital_id, 'Adrenaline', 'Epinephrine', 'injection', 'iv', '1', 'mg/ml', 'Emergency drug', true, false, 20, 'Room temperature, protect from light', NULL),
(:hospital_id, 'Hydrocortisone', 'Solu-Cortef', 'injection', 'iv', '100', 'mg', 'Corticosteroid', true, false, 30, 'Room temperature', NULL),
(:hospital_id, 'Calcium Gluconate', 'Calcium Gluconate', 'injection', 'iv', '10', '%', 'Electrolyte', true, false, 50, 'Room temperature', NULL),

-- IV Fluids & Solutions
(:hospital_id, 'Normal Saline', 'NaCl 0.9%', 'infusion', 'iv', '1000', 'ml', 'IV fluid', false, false, 200, 'Room temperature', NULL),
(:hospital_id, 'Dextrose 5%', 'D5W', 'infusion', 'iv', '1000', 'ml', 'IV fluid', false, false, 150, 'Room temperature', NULL),
(:hospital_id, 'Sodium Bicarbonate', 'NaHCO3 8.4%', 'injection', 'iv', '50', 'ml', 'Electrolyte', true, false, 100, 'Room temperature', NULL)

ON CONFLICT (hospital_id, generic_name, strength, medication_form, route) DO NOTHING;

-- Display success message
DO $$
BEGIN
  RAISE NOTICE 'Medications catalog seeded successfully. Total medications: 68';
END $$;
