#!/usr/bin/env python3
"""
Seed Uganda-specific prices in UGX for Kiruddu Hospital
Prices based on typical Uganda dialysis center rates (2025)
"""
import os
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values
from datetime import date

def seed_prices(hospital_id):
    """Create Uganda price list in UGX for hospital"""

    # Database connection
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    # Uganda dialysis prices in UGX (Uganda Shillings)
    # Typical rates as of 2025
    prices = [
        # Dialysis Services
        {
            'name': 'Hemodialysis Session (4 hours)',
            'code': 'HD-SESSION',
            'category': 'Dialysis',
            'price': 450000,  # UGX 450,000
        },
        {
            'name': 'Hemodialysis Session (3 hours)',
            'code': 'HD-SESSION-3H',
            'category': 'Dialysis',
            'price': 380000,  # UGX 380,000
        },
        {
            'name': 'Emergency Hemodialysis',
            'code': 'HD-EMERGENCY',
            'category': 'Dialysis',
            'price': 550000,  # UGX 550,000
        },
        {
            'name': 'Peritoneal Dialysis Training (per day)',
            'code': 'PD-TRAINING',
            'category': 'Dialysis',
            'price': 200000,  # UGX 200,000
        },

        # Vascular Access
        {
            'name': 'AV Fistula Creation',
            'code': 'ACCESS-AVF',
            'category': 'Vascular Access',
            'price': 2500000,  # UGX 2,500,000
        },
        {
            'name': 'Central Venous Catheter Insertion',
            'code': 'ACCESS-CVC',
            'category': 'Vascular Access',
            'price': 800000,  # UGX 800,000
        },
        {
            'name': 'Permcath Insertion',
            'code': 'ACCESS-PERMCATH',
            'category': 'Vascular Access',
            'price': 1500000,  # UGX 1,500,000
        },

        # Laboratory Tests
        {
            'name': 'Complete Blood Count (CBC)',
            'code': 'LAB-CBC',
            'category': 'Laboratory',
            'price': 25000,  # UGX 25,000
        },
        {
            'name': 'Urea and Electrolytes',
            'code': 'LAB-UE',
            'category': 'Laboratory',
            'price': 35000,  # UGX 35,000
        },
        {
            'name': 'Creatinine',
            'code': 'LAB-CREAT',
            'category': 'Laboratory',
            'price': 15000,  # UGX 15,000
        },
        {
            'name': 'Parathyroid Hormone (PTH)',
            'code': 'LAB-PTH',
            'category': 'Laboratory',
            'price': 150000,  # UGX 150,000
        },
        {
            'name': 'Hepatitis B Surface Antigen',
            'code': 'LAB-HBS',
            'category': 'Laboratory',
            'price': 35000,  # UGX 35,000
        },
        {
            'name': 'HIV Test',
            'code': 'LAB-HIV',
            'category': 'Laboratory',
            'price': 20000,  # UGX 20,000
        },
        {
            'name': 'Calcium and Phosphate',
            'code': 'LAB-CA-PO4',
            'category': 'Laboratory',
            'price': 30000,  # UGX 30,000
        },
        {
            'name': 'Lipid Profile',
            'code': 'LAB-LIPIDS',
            'category': 'Laboratory',
            'price': 45000,  # UGX 45,000
        },

        # Imaging
        {
            'name': 'Chest X-Ray',
            'code': 'IMG-CXR',
            'category': 'Imaging',
            'price': 40000,  # UGX 40,000
        },
        {
            'name': 'Renal Ultrasound',
            'code': 'IMG-US-RENAL',
            'category': 'Imaging',
            'price': 80000,  # UGX 80,000
        },
        {
            'name': 'Echocardiography',
            'code': 'IMG-ECHO',
            'category': 'Imaging',
            'price': 150000,  # UGX 150,000
        },
        {
            'name': 'Doppler Study (AV Fistula)',
            'code': 'IMG-DOPPLER-AVF',
            'category': 'Imaging',
            'price': 120000,  # UGX 120,000
        },

        # Medications (per dose/unit)
        {
            'name': 'Erythropoietin 4000 IU',
            'code': 'MED-EPO-4000',
            'category': 'Medications',
            'price': 35000,  # UGX 35,000
        },
        {
            'name': 'Iron Sucrose 100mg',
            'code': 'MED-IRON-100',
            'category': 'Medications',
            'price': 25000,  # UGX 25,000
        },
        {
            'name': 'Heparin 5000 IU',
            'code': 'MED-HEPARIN',
            'category': 'Medications',
            'price': 8000,  # UGX 8,000
        },
        {
            'name': 'Calcium Carbonate 500mg (tablet)',
            'code': 'MED-CACO3',
            'category': 'Medications',
            'price': 500,  # UGX 500
        },

        # Consultations
        {
            'name': 'Nephrologist Consultation (Initial)',
            'code': 'CONS-NEPHRO-INIT',
            'category': 'Consultation',
            'price': 100000,  # UGX 100,000
        },
        {
            'name': 'Nephrologist Consultation (Follow-up)',
            'code': 'CONS-NEPHRO-FU',
            'category': 'Consultation',
            'price': 50000,  # UGX 50,000
        },
        {
            'name': 'Nutritionist Consultation',
            'code': 'CONS-NUTR',
            'category': 'Consultation',
            'price': 30000,  # UGX 30,000
        },
        {
            'name': 'Social Worker Consultation',
            'code': 'CONS-SW',
            'category': 'Consultation',
            'price': 20000,  # UGX 20,000
        },

        # Admission/Ward
        {
            'name': 'Dialysis Admission (per day)',
            'code': 'WARD-DIAL',
            'category': 'Admission',
            'price': 50000,  # UGX 50,000
        },
        {
            'name': 'ICU Admission (per day)',
            'code': 'WARD-ICU',
            'category': 'Admission',
            'price': 300000,  # UGX 300,000
        },
    ]

    print("=" * 70)
    print("Seeding Uganda Price List (UGX) for Kiruddu Hospital")
    print("=" * 70)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Check existing prices
        cursor.execute("SELECT service_code FROM price_lists WHERE hospital_id = %s", (hospital_id,))
        existing = {row[0] for row in cursor.fetchall()}

        prices_to_insert = []
        for price in prices:
            if price['code'] not in existing:
                prices_to_insert.append((
                    str(uuid.uuid4()),
                    hospital_id,
                    price['name'],
                    price['code'],
                    price['category'],
                    price['price'],
                    date(2025, 1, 1),  # effective_from
                    True,  # is_active
                ))

        if prices_to_insert:
            query = """
                INSERT INTO price_lists (
                    id, hospital_id, service_name, service_code, service_category,
                    unit_price, effective_from, is_active
                ) VALUES %s
            """
            execute_values(cursor, query, prices_to_insert)
            conn.commit()

            print(f"✓ Created {len(prices_to_insert)} price items in UGX:")

            # Group by category
            by_category = {}
            for p in prices_to_insert:
                cat = p[4]
                if cat not in by_category:
                    by_category[cat] = []
                by_category[cat].append((p[2], p[5]))

            for category, items in sorted(by_category.items()):
                print(f"\n{category}:")
                for name, price in items:
                    print(f"  - {name}: UGX {price:,.0f}")
        else:
            print("✓ All prices already exist")

        print()
        print("=" * 70)
        print("Price List Setup Complete!")
        print("=" * 70)
        print(f"Total services: {len(existing) + len(prices_to_insert)}")
        print("Currency: Uganda Shillings (UGX)")

        cursor.close()
        conn.close()

    except Exception as e:
        print(f"\n✗ Failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python3 seed_uganda_prices.py <hospital_id>")
        sys.exit(1)

    hospital_id = sys.argv[1]
    seed_prices(hospital_id)
