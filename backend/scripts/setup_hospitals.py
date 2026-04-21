#!/usr/bin/env python3
"""
Setup Uganda dialysis hospitals
"""
import os
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values

# Uganda dialysis hospitals with contact information
UGANDA_HOSPITALS = [
    {
        'name': 'Kiruddu National Referral Hospital',
        'address': 'Namuwongo, Makindye Division, Kampala',
        'phone': '+256 414 270 251',
        'email': 'info@kiruddu.go.ug',
        'region': 'Central',
        'level': 'Referral Hospital',
    },
    {
        'name': 'Mulago National Referral Hospital',
        'address': 'Hill Road, Kampala',
        'phone': '+256 414 554 600',
        'email': 'info@mulago.go.ug',
        'region': 'Central',
        'level': 'National Referral Hospital',
    },
    {
        'name': 'Uganda Heart Institute',
        'address': 'Mulago Hospital Complex, Kampala',
        'phone': '+256 414 540 524',
        'email': 'info@uhi.go.ug',
        'region': 'Central',
        'level': 'Specialized Institute',
    },
    {
        'name': 'Mengo Hospital',
        'address': 'Namirembe Hill, Kampala',
        'phone': '+256 414 270 021',
        'email': 'info@mengohospital.org',
        'region': 'Central',
        'level': 'Private Hospital',
    },
    {
        'name': 'Nsambya Hospital',
        'address': 'Nsambya Hill, Makindye, Kampala',
        'phone': '+256 414 267 051',
        'email': 'info@nsambya.org',
        'region': 'Central',
        'level': 'Private Not-For-Profit',
    },
    {
        'name': 'International Hospital Kampala (IHK)',
        'address': 'Namuwongo, Kampala',
        'phone': '+256 312 200 400',
        'email': 'info@ihk.co.ug',
        'region': 'Central',
        'level': 'Private Hospital',
    },
    {
        'name': 'Case Hospital',
        'address': 'Plot 69-71 Buganda Road, Kampala',
        'phone': '+256 414 250 362',
        'email': 'info@casehospital.org',
        'region': 'Central',
        'level': 'Private Hospital',
    },
    {
        'name': 'Lubaga Hospital',
        'address': 'Lubaga Hill, Kampala',
        'phone': '+256 414 270 584',
        'email': 'info@lubagahospital.org',
        'region': 'Central',
        'level': 'Private Not-For-Profit',
    },
    {
        'name': 'Mbarara Regional Referral Hospital',
        'address': 'Mbarara City, Western Uganda',
        'phone': '+256 485 420 706',
        'email': 'info@mbarararrh.go.ug',
        'region': 'Western',
        'level': 'Regional Referral Hospital',
    },
    {
        'name': 'Kampala Hospital',
        'address': 'Kololo, Kampala',
        'phone': '+256 312 200 000',
        'email': 'info@kampalahospital.com',
        'region': 'Central',
        'level': 'Private Hospital',
    },
    {
        'name': 'Gulu Regional Referral Hospital',
        'address': 'Gulu City, Northern Uganda',
        'phone': '+256 471 432 501',
        'email': 'info@gulurrh.go.ug',
        'region': 'Northern',
        'level': 'Regional Referral Hospital',
    },
    {
        'name': 'Jinja Regional Referral Hospital',
        'address': 'Jinja City, Eastern Uganda',
        'phone': '+256 434 120 165',
        'email': 'info@jinjarrh.go.ug',
        'region': 'Eastern',
        'level': 'Regional Referral Hospital',
    },
    {
        'name': 'Soroti Regional Referral Hospital',
        'address': 'Soroti City, Eastern Uganda',
        'phone': '+256 454 461 006',
        'email': 'info@sorotirh.go.ug',
        'region': 'Eastern',
        'level': 'Regional Referral Hospital',
    },
    {
        'name': 'Lacor Hospital',
        'address': 'Lacor, Gulu District',
        'phone': '+256 471 432 082',
        'email': 'info@lacorhospital.org',
        'region': 'Northern',
        'level': 'Private Not-For-Profit',
    },
    {
        'name': 'Mbale Regional Referral Hospital',
        'address': 'Mbale City, Eastern Uganda',
        'phone': '+256 454 433 260',
        'email': 'info@mbaleh.go.ug',
        'region': 'Eastern',
        'level': 'Regional Referral Hospital',
    },
]

def setup_hospitals():
    """Create Uganda dialysis hospitals in database"""

    # Database connection
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    print("=" * 70)
    print("DMS Hospital Setup - Uganda Dialysis Centers")
    print("=" * 70)
    print(f"Database: {db_config['database']}@{db_config['host']}")
    print()

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Check existing hospitals
        cursor.execute("SELECT name FROM hospitals")
        existing = {row[0] for row in cursor.fetchall()}

        hospitals_to_insert = []
        kiruddu_id = None

        for hospital in UGANDA_HOSPITALS:
            if hospital['name'] not in existing:
                hospital_id = str(uuid.uuid4())

                # Save Kiruddu ID for later
                if 'Kiruddu' in hospital['name']:
                    kiruddu_id = hospital_id

                # Map level to tier
                tier_map = {
                    'National Referral Hospital': 'national',
                    'Referral Hospital': 'regional',
                    'Regional Referral Hospital': 'regional',
                    'Specialized Institute': 'national',
                    'Private Hospital': 'private',
                    'Private Not-For-Profit': 'private',
                }
                tier = tier_map.get(hospital['level'], 'private')

                # Generate short code from name
                short_code = ''.join([word[0] for word in hospital['name'].split()[:3]]).upper()

                hospitals_to_insert.append((
                    hospital_id,
                    hospital['name'],
                    short_code,
                    tier,
                    hospital['region'],
                    hospital['address'],
                    hospital['phone'],
                    hospital['email'],
                ))

        if hospitals_to_insert:
            query = """
                INSERT INTO hospitals (
                    id, name, short_code, tier, region, address, phone, email
                ) VALUES %s
                ON CONFLICT (short_code) DO NOTHING
            """
            execute_values(cursor, query, hospitals_to_insert)
            conn.commit()
            print(f"✓ Created {len(hospitals_to_insert)} hospitals:")
            for h in hospitals_to_insert:
                print(f"  - {h[1]} ({h[3]}) - {h[4]}")
        else:
            print("✓ All hospitals already exist")

        # Get Kiruddu ID if it wasn't just created
        if not kiruddu_id:
            cursor.execute("SELECT id FROM hospitals WHERE name ILIKE '%kiruddu%' LIMIT 1")
            result = cursor.fetchone()
            if result:
                kiruddu_id = result[0]

        print()
        print("=" * 70)
        print("Hospital Setup Complete!")
        print("=" * 70)
        print(f"Total hospitals in database: {len(existing) + len(hospitals_to_insert)}")

        if kiruddu_id:
            print()
            print("=" * 70)
            print(f"Kiruddu Hospital ID: {kiruddu_id}")
            print("=" * 70)
            print()
            print("To import patient data, run:")
            print(f"  cd /home/bampita/Projects/My-apps/DMS-Dialysis_Management_System/backend/scripts")
            print(f"  python3 import_access_data.py {kiruddu_id}")

        cursor.close()
        conn.close()

        return kiruddu_id

    except Exception as e:
        print(f"\n✗ Setup failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    kiruddu_id = setup_hospitals()
