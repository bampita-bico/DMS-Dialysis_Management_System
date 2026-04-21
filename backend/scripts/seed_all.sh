#!/bin/bash
# Seed all reference data for a hospital
# Usage: ./seed_all.sh <hospital-uuid>

set -e  # Exit on error

if [ -z "$1" ]; then
  echo "Error: Hospital UUID required"
  echo "Usage: ./seed_all.sh <hospital-uuid>"
  echo "Example: ./seed_all.sh 123e4567-e89b-12d3-a456-426614174000"
  exit 1
fi

HOSPITAL_ID=$1
SEEDS_DIR="../seeds"
DB_NAME="${DB_NAME:-dms}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"

echo "========================================="
echo "DMS Reference Data Seeding"
echo "========================================="
echo "Hospital ID: $HOSPITAL_ID"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "========================================="
echo ""

# Check if hospital exists
HOSPITAL_EXISTS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM hospitals WHERE id='$HOSPITAL_ID' AND deleted_at IS NULL;")

if [ "$HOSPITAL_EXISTS" -eq 0 ]; then
  echo "Error: Hospital with ID $HOSPITAL_ID not found or is deleted"
  exit 1
fi

echo "✓ Hospital found in database"
echo ""

# Seed files in order
SEED_FILES=(
  "001_lab_tests.sql"
  "002_lab_panels.sql"
  "003_lab_reference_ranges.sql"
  "004_medications_uganda.sql"
  "005_drug_interactions.sql"
  "006_consumables.sql"
  "007_insurance_schemes.sql"
  "008_price_lists.sql"
)

for SEED_FILE in "${SEED_FILES[@]}"; do
  SEED_PATH="$SEEDS_DIR/$SEED_FILE"

  if [ ! -f "$SEED_PATH" ]; then
    echo "Warning: Seed file not found: $SEED_PATH"
    continue
  fi

  echo "Running: $SEED_FILE"
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
    -v hospital_id="'$HOSPITAL_ID'" \
    -f "$SEED_PATH" \
    -q  # Quiet mode, only show notices and errors

  if [ $? -eq 0 ]; then
    echo "✓ $SEED_FILE completed"
  else
    echo "✗ $SEED_FILE failed"
    exit 1
  fi
  echo ""
done

echo "========================================="
echo "Seeding completed successfully!"
echo "========================================="
echo ""
echo "Verification queries:"
echo "---------------------"
echo "Lab tests:         SELECT COUNT(*) FROM lab_test_catalog WHERE hospital_id='$HOSPITAL_ID';"
echo "Lab panels:        SELECT COUNT(*) FROM lab_panels WHERE hospital_id='$HOSPITAL_ID';"
echo "Medications:       SELECT COUNT(*) FROM medications WHERE hospital_id='$HOSPITAL_ID';"
echo "Consumables:       SELECT COUNT(*) FROM consumables WHERE hospital_id='$HOSPITAL_ID';"
echo "Insurance schemes: SELECT COUNT(*) FROM insurance_schemes WHERE hospital_id='$HOSPITAL_ID';"
echo "Price list:        SELECT COUNT(*) FROM price_lists WHERE hospital_id='$HOSPITAL_ID';"
echo ""
echo "To verify, run: psql -d $DB_NAME -c \"<query>\""
