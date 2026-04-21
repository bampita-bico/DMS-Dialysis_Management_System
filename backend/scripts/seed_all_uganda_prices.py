#!/usr/bin/env python3
"""
Complete Uganda pricing in UGX - All services, subscriptions, and MOH purchases
"""
import os
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values
from datetime import date

def seed_complete_prices(hospital_id):
    """Create complete Uganda price list in UGX"""

    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    # Complete Uganda pricing in UGX
    prices = [
        # ============================================
        # DIALYSIS SERVICES
        # ============================================
        {
            'name': 'Hemodialysis Session (4 hours)',
            'code': 'HD-SESSION-4H',
            'category': 'Dialysis',
            'price': 450000,
        },
        {
            'name': 'Hemodialysis Session (3 hours)',
            'code': 'HD-SESSION-3H',
            'category': 'Dialysis',
            'price': 380000,
        },
        {
            'name': 'Hemodialysis Session (2 hours) - Short',
            'code': 'HD-SESSION-2H',
            'category': 'Dialysis',
            'price': 300000,
        },
        {
            'name': 'Emergency Hemodialysis',
            'code': 'HD-EMERGENCY',
            'category': 'Dialysis',
            'price': 550000,
        },
        {
            'name': 'Peritoneal Dialysis Training (per day)',
            'code': 'PD-TRAINING',
            'category': 'Dialysis',
            'price': 200000,
        },
        {
            'name': 'Home Dialysis Setup and Training',
            'code': 'HD-HOME-SETUP',
            'category': 'Dialysis',
            'price': 500000,
        },
        {
            'name': 'Continuous Ambulatory Peritoneal Dialysis (CAPD) - Monthly',
            'code': 'CAPD-MONTHLY',
            'category': 'Dialysis',
            'price': 1200000,
        },

        # ============================================
        # VASCULAR ACCESS
        # ============================================
        {
            'name': 'AV Fistula Creation',
            'code': 'ACCESS-AVF',
            'category': 'Vascular Access',
            'price': 2500000,
        },
        {
            'name': 'AV Fistula Revision',
            'code': 'ACCESS-AVF-REVISION',
            'category': 'Vascular Access',
            'price': 1800000,
        },
        {
            'name': 'Central Venous Catheter Insertion (Temporary)',
            'code': 'ACCESS-CVC-TEMP',
            'category': 'Vascular Access',
            'price': 800000,
        },
        {
            'name': 'Permcath Insertion (Tunneled)',
            'code': 'ACCESS-PERMCATH',
            'category': 'Vascular Access',
            'price': 1500000,
        },
        {
            'name': 'Permcath Removal',
            'code': 'ACCESS-PERMCATH-REMOVE',
            'category': 'Vascular Access',
            'price': 300000,
        },
        {
            'name': 'Fistulogram with Angioplasty',
            'code': 'ACCESS-FISTULOGRAM',
            'category': 'Vascular Access',
            'price': 2000000,
        },

        # ============================================
        # LABORATORY TESTS - ROUTINE
        # ============================================
        {
            'name': 'Complete Blood Count (CBC)',
            'code': 'LAB-CBC',
            'category': 'Laboratory',
            'price': 25000,
        },
        {
            'name': 'Urea',
            'code': 'LAB-UREA',
            'category': 'Laboratory',
            'price': 15000,
        },
        {
            'name': 'Creatinine',
            'code': 'LAB-CREAT',
            'category': 'Laboratory',
            'price': 15000,
        },
        {
            'name': 'Urea and Electrolytes (U&E)',
            'code': 'LAB-UE',
            'category': 'Laboratory',
            'price': 35000,
        },
        {
            'name': 'Sodium',
            'code': 'LAB-NA',
            'category': 'Laboratory',
            'price': 10000,
        },
        {
            'name': 'Potassium',
            'code': 'LAB-K',
            'category': 'Laboratory',
            'price': 10000,
        },
        {
            'name': 'Calcium',
            'code': 'LAB-CA',
            'category': 'Laboratory',
            'price': 12000,
        },
        {
            'name': 'Phosphate',
            'code': 'LAB-PO4',
            'category': 'Laboratory',
            'price': 12000,
        },
        {
            'name': 'Calcium and Phosphate',
            'code': 'LAB-CA-PO4',
            'category': 'Laboratory',
            'price': 30000,
        },
        {
            'name': 'Liver Function Tests (LFT)',
            'code': 'LAB-LFT',
            'category': 'Laboratory',
            'price': 40000,
        },
        {
            'name': 'Fasting Blood Sugar',
            'code': 'LAB-FBS',
            'category': 'Laboratory',
            'price': 8000,
        },
        {
            'name': 'Random Blood Sugar',
            'code': 'LAB-RBS',
            'category': 'Laboratory',
            'price': 6000,
        },
        {
            'name': 'HbA1c',
            'code': 'LAB-HBA1C',
            'category': 'Laboratory',
            'price': 35000,
        },

        # ============================================
        # LABORATORY TESTS - SPECIALIZED
        # ============================================
        {
            'name': 'Parathyroid Hormone (PTH)',
            'code': 'LAB-PTH',
            'category': 'Laboratory',
            'price': 150000,
        },
        {
            'name': 'Vitamin D (25-OH)',
            'code': 'LAB-VIT-D',
            'category': 'Laboratory',
            'price': 120000,
        },
        {
            'name': 'Hepatitis B Surface Antigen (HBsAg)',
            'code': 'LAB-HBS-AG',
            'category': 'Laboratory',
            'price': 35000,
        },
        {
            'name': 'Hepatitis C Antibody',
            'code': 'LAB-HCV',
            'category': 'Laboratory',
            'price': 40000,
        },
        {
            'name': 'HIV Test',
            'code': 'LAB-HIV',
            'category': 'Laboratory',
            'price': 20000,
        },
        {
            'name': 'Lipid Profile',
            'code': 'LAB-LIPIDS',
            'category': 'Laboratory',
            'price': 45000,
        },
        {
            'name': 'Ferritin',
            'code': 'LAB-FERRITIN',
            'category': 'Laboratory',
            'price': 60000,
        },
        {
            'name': 'Iron Studies (Full Panel)',
            'code': 'LAB-IRON-PANEL',
            'category': 'Laboratory',
            'price': 80000,
        },
        {
            'name': 'C-Reactive Protein (CRP)',
            'code': 'LAB-CRP',
            'category': 'Laboratory',
            'price': 30000,
        },
        {
            'name': 'Albumin',
            'code': 'LAB-ALBUMIN',
            'category': 'Laboratory',
            'price': 15000,
        },

        # ============================================
        # IMAGING
        # ============================================
        {
            'name': 'Chest X-Ray',
            'code': 'IMG-CXR',
            'category': 'Imaging',
            'price': 40000,
        },
        {
            'name': 'Renal Ultrasound',
            'code': 'IMG-US-RENAL',
            'category': 'Imaging',
            'price': 80000,
        },
        {
            'name': 'Abdominal Ultrasound',
            'code': 'IMG-US-ABD',
            'category': 'Imaging',
            'price': 70000,
        },
        {
            'name': 'Echocardiography',
            'code': 'IMG-ECHO',
            'category': 'Imaging',
            'price': 150000,
        },
        {
            'name': 'Doppler Study (AV Fistula)',
            'code': 'IMG-DOPPLER-AVF',
            'category': 'Imaging',
            'price': 120000,
        },
        {
            'name': 'CT Scan (Non-contrast)',
            'code': 'IMG-CT-PLAIN',
            'category': 'Imaging',
            'price': 250000,
        },
        {
            'name': 'CT Scan (With contrast)',
            'code': 'IMG-CT-CONTRAST',
            'category': 'Imaging',
            'price': 350000,
        },

        # ============================================
        # MEDICATIONS (Per dose/unit)
        # ============================================
        {
            'name': 'Erythropoietin (Eprex) 4000 IU',
            'code': 'MED-EPO-4000',
            'category': 'Medications',
            'price': 35000,
        },
        {
            'name': 'Erythropoietin (Eprex) 10000 IU',
            'code': 'MED-EPO-10000',
            'category': 'Medications',
            'price': 80000,
        },
        {
            'name': 'Iron Sucrose (Venofer) 100mg',
            'code': 'MED-IRON-SUCROSE',
            'category': 'Medications',
            'price': 25000,
        },
        {
            'name': 'Heparin 5000 IU (per vial)',
            'code': 'MED-HEPARIN-5000',
            'category': 'Medications',
            'price': 8000,
        },
        {
            'name': 'Calcium Carbonate 500mg (per tablet)',
            'code': 'MED-CACO3-500',
            'category': 'Medications',
            'price': 500,
        },
        {
            'name': 'Sevelamer 800mg (per tablet)',
            'code': 'MED-SEVELAMER',
            'category': 'Medications',
            'price': 2000,
        },
        {
            'name': 'Calcitriol 0.25mcg (per capsule)',
            'code': 'MED-CALCITRIOL',
            'category': 'Medications',
            'price': 1500,
        },
        {
            'name': 'Folic Acid 5mg (per tablet)',
            'code': 'MED-FOLIC-ACID',
            'category': 'Medications',
            'price': 300,
        },

        # ============================================
        # CONSULTATIONS
        # ============================================
        {
            'name': 'Nephrologist Consultation (Initial)',
            'code': 'CONS-NEPHRO-INIT',
            'category': 'Consultation',
            'price': 100000,
        },
        {
            'name': 'Nephrologist Consultation (Follow-up)',
            'code': 'CONS-NEPHRO-FU',
            'category': 'Consultation',
            'price': 50000,
        },
        {
            'name': 'Nutritionist/Dietitian Consultation',
            'code': 'CONS-NUTR',
            'category': 'Consultation',
            'price': 30000,
        },
        {
            'name': 'Social Worker Consultation',
            'code': 'CONS-SW',
            'category': 'Consultation',
            'price': 20000,
        },
        {
            'name': 'Psychologist Consultation',
            'code': 'CONS-PSYCH',
            'category': 'Consultation',
            'price': 40000,
        },
        {
            'name': 'Cardiology Consultation',
            'code': 'CONS-CARDIO',
            'category': 'Consultation',
            'price': 80000,
        },

        # ============================================
        # ADMISSION/WARD
        # ============================================
        {
            'name': 'Dialysis Admission (per day)',
            'code': 'WARD-DIAL',
            'category': 'Admission',
            'price': 50000,
        },
        {
            'name': 'ICU Admission (per day)',
            'code': 'WARD-ICU',
            'category': 'Admission',
            'price': 300000,
        },
        {
            'name': 'General Ward (per day)',
            'code': 'WARD-GENERAL',
            'category': 'Admission',
            'price': 30000,
        },
        {
            'name': 'Private Room (per day)',
            'code': 'WARD-PRIVATE',
            'category': 'Admission',
            'price': 100000,
        },

        # ============================================
        # DIALYSIS SUPPLIES (Per session)
        # ============================================
        {
            'name': 'Dialyzer (High Flux)',
            'code': 'SUPPLY-DIALYZER-HF',
            'category': 'Supplies',
            'price': 45000,
        },
        {
            'name': 'Bloodline Set',
            'code': 'SUPPLY-BLOODLINE',
            'category': 'Supplies',
            'price': 15000,
        },
        {
            'name': 'AV Fistula Needles (pair)',
            'code': 'SUPPLY-NEEDLES',
            'category': 'Supplies',
            'price': 8000,
        },
        {
            'name': 'Bicarbonate Concentrate',
            'code': 'SUPPLY-BICARB',
            'category': 'Supplies',
            'price': 12000,
        },

        # ============================================
        # PACKAGES (Monthly/Annual)
        # ============================================
        {
            'name': 'Monthly Dialysis Package (3x/week)',
            'code': 'PKG-DIALYSIS-3X-MONTH',
            'category': 'Packages',
            'price': 5400000,  # 12 sessions @ 450k
        },
        {
            'name': 'Monthly Dialysis Package (2x/week)',
            'code': 'PKG-DIALYSIS-2X-MONTH',
            'category': 'Packages',
            'price': 3600000,  # 8 sessions @ 450k
        },
        {
            'name': 'Annual Dialysis Package (3x/week)',
            'code': 'PKG-DIALYSIS-3X-YEAR',
            'category': 'Packages',
            'price': 64800000,  # 144 sessions @ 450k
        },
    ]

    print("=" * 80)
    print("Uganda Complete Price List in UGX - All Services")
    print("=" * 80)

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
                    date(2025, 1, 1),
                    True,
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

            total_value = 0
            for category, items in sorted(by_category.items()):
                print(f"\n{category}:")
                for name, price in items:
                    print(f"  - {name}: UGX {price:,.0f}")
                    total_value += price

        else:
            print("✓ All prices already exist")

        print()
        print("=" * 80)
        print("Price List Setup Complete!")
        print("=" * 80)
        print(f"Total services: {len(existing) + len(prices_to_insert)}")
        print(f"New prices added: {len(prices_to_insert)}")
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
        print("Usage: python3 seed_all_uganda_prices.py <hospital_id>")
        sys.exit(1)

    hospital_id = sys.argv[1]
    seed_complete_prices(hospital_id)
