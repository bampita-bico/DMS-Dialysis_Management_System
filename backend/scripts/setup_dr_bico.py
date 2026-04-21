#!/usr/bin/env python3
"""
Setup Dr. Mujjabi Steve Bico as admin and mark Kiruddu as demo hospital
"""
import os
import sys
import psycopg2
import bcrypt

def setup_dr_bico():
    """Update admin user and hospital settings"""

    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    hospital_id = 'a64ad314-4a24-4d5e-bbde-87776a5aea54'

    print("=" * 70)
    print("Setting up Dr. Mujjabi Steve Bico - Kiruddu Hospital")
    print("=" * 70)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Hash password
        password = "DrBico123!"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')

        # Update admin user
        cursor.execute("""
            UPDATE users
            SET email = %s,
                full_name = %s,
                password_hash = %s,
                is_active = TRUE,
                is_verified = TRUE
            WHERE hospital_id = %s
            AND email = 'admin@kiruddu.go.ug'
        """, ('msbico@gmail.com', 'Dr. Mujjabi Steve Bico', password_hash, hospital_id))

        # Mark Kiruddu as demo hospital and set enterprise plan
        cursor.execute("""
            UPDATE hospitals
            SET subscription_plan = 'enterprise',
                settings = jsonb_set(
                    COALESCE(settings, '{}'::jsonb),
                    '{is_demo}',
                    'true'::jsonb
                ),
                enabled_modules = '{
                  "lab_management": true,
                  "full_pharmacy": true,
                  "hr_management": true,
                  "inventory_tracking": true,
                  "advanced_billing": true,
                  "offline_sync": true,
                  "chw_program": true,
                  "imaging_integration": true,
                  "outcomes_reporting": true
                }'::jsonb
            WHERE id = %s
        """, (hospital_id,))

        conn.commit()

        print("✓ Updated admin user to Dr. Mujjabi Steve Bico")
        print("✓ Marked Kiruddu as demo hospital (Enterprise plan)")
        print()
        print("=" * 70)
        print("NEW LOGIN CREDENTIALS")
        print("=" * 70)
        print(f"Email:    msbico@gmail.com")
        print(f"Password: {password}")
        print()
        print("Hospital: Kiruddu National Referral Hospital (DEMO)")
        print("Plan:     Enterprise (All features enabled)")
        print()
        print("URL: http://localhost:5173")
        print("=" * 70)

        cursor.close()
        conn.close()

    except Exception as e:
        print(f"\n✗ Failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    setup_dr_bico()
