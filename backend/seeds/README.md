# DMS Reference Data Seeds

This directory contains SQL seed files to populate reference data for dialysis centers.

## Quick Start

```bash
cd backend
./scripts/seed_all.sh <hospital-uuid>
```

## Seed Files

| File | Description | Records |
|------|-------------|---------|
| `001_lab_tests.sql` | Standard dialysis lab test catalog | 54 tests |
| `002_lab_panels.sql` | Pre-built lab panels (Pre-dialysis, Monthly, etc.) | 12 panels |
| `003_lab_reference_ranges.sql` | Normal lab value ranges by age/sex | ~30 ranges |
| `004_medications.sql` | Common dialysis medications | 68 drugs |
| `005_drug_interactions.sql` | Critical drug interactions | ~15 interactions |
| `006_consumables.sql` | Dialysis consumables (dialyzers, bloodlines, etc.) | 35 items |
| `007_insurance_schemes.sql` | East African insurance providers | 16 schemes |
| `008_price_lists.sql` | Standard pricing for dialysis services | 33 services |

## Usage

### Seed All Data for a Hospital

```bash
# Get hospital UUID from database
psql -d dms -c "SELECT id, name FROM hospitals WHERE deleted_at IS NULL;"

# Run all seeds
export DB_NAME=dms
export DB_USER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_PASSWORD=your_password

./scripts/seed_all.sh <hospital-uuid>
```

### Seed Individual Files

```bash
psql -d dms -v hospital_id='<uuid>' -f seeds/001_lab_tests.sql
```

## Data Included

### Lab Tests (001)
- **Renal Function**: Urea, Creatinine, eGFR, BUN
- **Electrolytes**: Na, K, Cl, Bicarb, Ca, Phosphate
- **Hematology**: CBC (Hb, WBC, PLT), Iron studies
- **Bone/Mineral**: PTH, Vitamin D, ALP
- **Liver Function**: Albumin, ALT, AST
- **Infectious Disease**: Hep B, Hep C, HIV
- **Cardiac**: Troponin, BNP
- **Coagulation**: PT, INR, aPTT

### Lab Panels (002)
- Pre-Dialysis Panel (Urea, K, Ca, Phos)
- Post-Dialysis Panel (Urea for URR)
- Monthly Monitoring (Full metabolic + CBC)
- Quarterly Comprehensive (Includes PTH, lipids, HbA1c)
- Anemia Workup (CBC + Iron studies)
- Bone Panel (CKD-MBD assessment)
- Adequacy Panel (URR, Kt/V)
- Emergency Panels (Cardiac, Sepsis)

### Medications (004)
- **ESAs**: Eprex, Aranesp, Recormon (various strengths)
- **Iron**: Venofer, Imferon, oral iron
- **Phosphate Binders**: Calcium acetate, Sevelamer, Lanthanum
- **Vitamin D**: Calcitriol, Alfacalcidol
- **Antihypertensives**: 
  - ACE-I (Enalapril, Lisinopril)
  - ARBs (Losartan, Valsartan)
  - CCBs (Amlodipine, Nifedipine)
  - Beta-blockers (Atenolol, Carvedilol)
- **Anticoagulants**: Heparin, Warfarin, Enoxaparin
- **Antibiotics**: Ceftriaxone, Vancomycin, Gentamicin
- **Emergency Drugs**: Adrenaline, Hydrocortisone, Calcium gluconate

### Consumables (006)
- **Dialyzers**: FX80, Polyflux, Revaclear (various models)
- **Bloodlines**: Adult and pediatric sets
- **AV Needles**: 15G, 16G, 17G
- **PPE**: Gloves, masks, gauze
- **Fluids**: Saline, Heparin

### Insurance Schemes (007)
- **Kenya**: NHIF, Britam, AAR, CIC, Jubilee
- **Uganda**: NHIS, IAA, AAR
- **Tanzania**: NHIF-TZ, Jubilee
- **Self-Pay**: Cash, Corporate credit

### Price List (008)
- **Dialysis**: HD (KES 8,000), HDF (KES 10,000), PD
- **Procedures**: CVC insertion, AVF creation
- **Consultations**: Nephrologist, Doctor, Nurse
- **Lab Tests**: Individual test pricing
- **Administrative**: Registration, reports

## Customization

### Add Custom Items

After seeding, use the APIs to add hospital-specific items:

```bash
# Add custom medication
curl -X POST http://localhost:8080/api/v1/medications \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "generic_name": "Custom Drug",
    "strength": "100",
    "medication_form": "tablet"
  }'

# Add custom lab test
curl -X POST http://localhost:8080/api/v1/lab/tests \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "code": "CUSTOM",
    "name": "Custom Test",
    "category": "chemistry"
  }'
```

### Modify Pricing

Update price lists via API or directly in database:

```sql
UPDATE price_lists 
SET unit_price = 10000.00 
WHERE service_code = 'HD-SESS' 
AND hospital_id = '<uuid>';
```

## Verification

After seeding, verify data was inserted:

```bash
# Check counts
psql -d dms -c "SELECT COUNT(*) FROM lab_test_catalog WHERE hospital_id='<uuid>';"
psql -d dms -c "SELECT COUNT(*) FROM medications WHERE hospital_id='<uuid>';"
psql -d dms -c "SELECT COUNT(*) FROM consumables WHERE hospital_id='<uuid>';"
psql -d dms -c "SELECT COUNT(*) FROM insurance_schemes WHERE hospital_id='<uuid>';"

# Sample data
psql -d dms -c "SELECT name, category FROM lab_test_catalog WHERE hospital_id='<uuid>' LIMIT 10;"
psql -d dms -c "SELECT generic_name, strength FROM medications WHERE hospital_id='<uuid>' LIMIT 10;"
```

## Notes

- **Idempotent**: All seeds use `ON CONFLICT DO NOTHING` - safe to re-run
- **Multi-tenant**: Each hospital gets independent copies
- **Currency**: All prices in KES (Kenyan Shillings) - adjust as needed
- **Regional**: Data focused on East African context (Kenya/Uganda/Tanzania)

## Troubleshooting

### Hospital not found
```
Error: Hospital with ID xxx not found or is deleted
```
**Solution**: Verify hospital exists in database first

### Permission denied
```
Error: permission denied for table medications
```
**Solution**: Ensure database user has INSERT permissions

### Duplicate key error
```
Error: duplicate key value violates unique constraint
```
**Solution**: This is expected if data already exists - seeds will skip duplicates

## Support

For questions or issues:
- Check `/IMPLEMENTATION_PROGRESS.md` for overview
- Review table schemas in `/backend/internal/db/migrations/`
- Contact: Dr. Bampita Steve Bico
