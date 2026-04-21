-- Seed: Dialysis Consumables Catalog
-- Usage: psql -d dms -v hospital_id='<uuid>' -f 006_consumables.sql

INSERT INTO consumables (hospital_id, item_name, category, manufacturer, catalog_number, unit_of_measure, cost_per_unit, reorder_level, deleted_at) VALUES
-- Dialyzers
(:hospital_id, 'Dialyzer FX80 (High Flux)', 'dialyzer', 'Fresenius', 'FX80', 'piece', 25.00, 20, NULL),
(:hospital_id, 'Dialyzer F8HPS (High Performance)', 'dialyzer', 'Fresenius', 'F8HPS', 'piece', 30.00, 15, NULL),
(:hospital_id, 'Dialyzer Polyflux 170H', 'dialyzer', 'Gambro', 'PF170H', 'piece', 28.00, 20, NULL),
(:hospital_id, 'Dialyzer Revaclear 400', 'dialyzer', 'Gambro', 'RC400', 'piece', 32.00, 10, NULL),

-- Bloodlines
(:hospital_id, 'Bloodline Set Adult (Standard)', 'bloodline', 'Fresenius', 'BLS-A', 'set', 8.00, 50, NULL),
(:hospital_id, 'Bloodline Set Pediatric', 'bloodline', 'Fresenius', 'BLS-P', 'set', 9.00, 20, NULL),

-- AV Needles
(:hospital_id, 'AV Fistula Needle 15G', 'av_needle', 'Nipro', 'AVN-15G', 'piece', 1.50, 100, NULL),
(:hospital_id, 'AV Fistula Needle 16G', 'av_needle', 'Nipro', 'AVN-16G', 'piece', 1.50, 100, NULL),
(:hospital_id, 'AV Fistula Needle 17G', 'av_needle', 'Nipro', 'AVN-17G', 'piece', 1.50, 80, NULL),

-- Syringes
(:hospital_id, 'Syringe 5ml Luer Lock', 'syringe', 'BD', 'SYR-5ML', 'piece', 0.15, 500, NULL),
(:hospital_id, 'Syringe 10ml Luer Lock', 'syringe', 'BD', 'SYR-10ML', 'piece', 0.20, 500, NULL),
(:hospital_id, 'Syringe 20ml Luer Lock', 'syringe', 'BD', 'SYR-20ML', 'piece', 0.25, 300, NULL),
(:hospital_id, 'Syringe 50ml Luer Lock', 'syringe', 'BD', 'SYR-50ML', 'piece', 0.50, 200, NULL),

-- Gauze & Dressing
(:hospital_id, 'Gauze Swabs 5x5cm Sterile', 'gauze', 'Generic', 'GS-5X5', 'pack of 100', 2.00, 50, NULL),
(:hospital_id, 'Gauze Swabs 10x10cm Sterile', 'gauze', 'Generic', 'GS-10X10', 'pack of 100', 3.00, 30, NULL),
(:hospital_id, 'Micropore Tape 2.5cm', 'gauze', '3M', 'MPT-25', 'roll', 1.50, 100, NULL),

-- Gloves
(:hospital_id, 'Surgical Gloves Size 7 Sterile', 'gloves', 'Ansell', 'SG-7', 'pair', 0.50, 200, NULL),
(:hospital_id, 'Surgical Gloves Size 7.5 Sterile', 'gloves', 'Ansell', 'SG-7.5', 'pair', 0.50, 200, NULL),
(:hospital_id, 'Examination Gloves Medium Non-Sterile', 'gloves', 'Generic', 'EG-M', 'box of 100', 8.00, 50, NULL),
(:hospital_id, 'Examination Gloves Large Non-Sterile', 'gloves', 'Generic', 'EG-L', 'box of 100', 8.00, 50, NULL),

-- Masks & PPE
(:hospital_id, 'Surgical Face Mask 3-ply', 'mask', 'Generic', 'SFM-3P', 'box of 50', 5.00, 100, NULL),
(:hospital_id, 'N95 Respirator Mask', 'mask', '3M', 'N95-1860', 'piece', 2.00, 50, NULL),

-- Disinfectants
(:hospital_id, 'Chlorhexidine 4% Solution 500ml', 'disinfectant', 'Hibiscrub', 'CHX-4', 'bottle', 8.00, 30, NULL),
(:hospital_id, 'Isopropyl Alcohol 70% 1L', 'disinfectant', 'Generic', 'IPA-70', 'bottle', 5.00, 50, NULL),
(:hospital_id, 'Hydrogen Peroxide 6% 1L', 'disinfectant', 'Generic', 'H2O2-6', 'bottle', 6.00, 30, NULL),

-- Saline & Fluids
(:hospital_id, 'Normal Saline 0.9% 1000ml Bag', 'saline', 'Fresenius', 'NS-1000', 'bag', 2.00, 200, NULL),
(:hospital_id, 'Normal Saline 0.9% 500ml Bag', 'saline', 'Fresenius', 'NS-500', 'bag', 1.50, 150, NULL),
(:hospital_id, 'Normal Saline 0.9% 10ml Ampoule', 'saline', 'Generic', 'NS-10', 'ampoule', 0.20, 500, NULL),

-- Heparin
(:hospital_id, 'Heparin 5000 IU/ml 5ml Vial', 'heparin', 'Generic', 'HEP-5000', 'vial', 3.00, 100, NULL),
(:hospital_id, 'Heparin 1000 IU/ml 10ml Vial', 'heparin', 'Generic', 'HEP-1000', 'vial', 2.50, 80, NULL),

-- Other Essential Items
(:hospital_id, 'IV Cannula 18G', 'other', 'BD', 'IVC-18G', 'piece', 0.80, 100, NULL),
(:hospital_id, 'IV Cannula 20G', 'other', 'BD', 'IVC-20G', 'piece', 0.80, 100, NULL),
(:hospital_id, 'Blood Collection Tube EDTA 5ml', 'other', 'Vacuette', 'BCT-EDTA', 'piece', 0.50, 200, NULL),
(:hospital_id, 'Blood Collection Tube Plain 5ml', 'other', 'Vacuette', 'BCT-PLAIN', 'piece', 0.50, 200, NULL)

ON CONFLICT (hospital_id, item_name, manufacturer) DO NOTHING;

DO $$
BEGIN
  RAISE NOTICE 'Consumables catalog seeded successfully. Total items: 35';
END $$;
