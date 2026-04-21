#!/usr/bin/env python3
"""
Create admin user for a hospital
"""
import os
import sys
import uuid
import psycopg2

def create_admin_user(hospital_id, hospital_name):
    """Create admin user for hospital"""

    # Database connection
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    print("=" * 70)
    print(f"Creating Admin User for {hospital_name}")
    print("=" * 70)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Check if user already exists
        cursor.execute("""
            SELECT id, email FROM users
            WHERE hospital_id = %s
            LIMIT 1
        """, (hospital_id,))

        result = cursor.fetchone()
        if result:
            print(f"✓ User already exists:")
            print(f"  Email: {result[1]}")
            print(f"  User ID: {result[0]}")
            cursor.close()
            conn.close()
            return result[0]

        # Create new admin user
        user_id = str(uuid.uuid4())
        email = "admin@kiruddu.go.ug"
        password_hash = "$2a$10$DummyHashForImport"  # Will need to reset password
        full_name = "System Administrator"

        # Insert user
        cursor.execute("""
            INSERT INTO users (
                id, hospital_id, email, password_hash, full_name, is_active
            ) VALUES (%s, %s, %s, %s, %s, %s)
        """, (
            user_id,
            hospital_id,
            email,
            password_hash,
            full_name,
            True
        ))

        conn.commit()

        print(f"✓ Created user:")
        print(f"  Email: {email}")
        print(f"  User ID: {user_id}")
        print()
        print("⚠️  Password needs to be set via password reset!")

        cursor.close()
        conn.close()

        return user_id

    except Exception as e:
        print(f"\n✗ Failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python3 create_admin_user.py <hospital_id>")
        sys.exit(1)

    hospital_id = sys.argv[1]
    create_admin_user(hospital_id, "Kiruddu National Referral Hospital")
