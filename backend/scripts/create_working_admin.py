#!/usr/bin/env python3
"""
Create admin user with working password using bcrypt
"""
import os
import sys
import uuid
import psycopg2

def create_admin(hospital_id):
    """Create admin user with bcrypt hashed password"""

    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    # Try to use bcrypt if available, otherwise use a known hash
    try:
        import bcrypt
        password = "admin123"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
    except ImportError:
        # Use pre-computed bcrypt hash for "admin123"
        password = "admin123"
        password_hash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

    print("=" * 70)
    print("Creating Admin User with Working Password")
    print("=" * 70)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Check if admin user exists
        cursor.execute("""
            SELECT id FROM users
            WHERE hospital_id = %s
            AND email = 'admin@kiruddu.go.ug'
        """, (hospital_id,))

        result = cursor.fetchone()

        if result:
            # Update existing user's password
            user_id = result[0]
            cursor.execute("""
                UPDATE users
                SET password_hash = %s,
                    is_active = TRUE,
                    is_verified = TRUE
                WHERE id = %s
            """, (password_hash, user_id))
            print("✓ Updated existing admin user's password")
        else:
            # Create new admin user
            user_id = str(uuid.uuid4())
            cursor.execute("""
                INSERT INTO users (
                    id, hospital_id, email, password_hash, full_name, is_active, is_verified
                ) VALUES (%s, %s, %s, %s, %s, %s, %s)
            """, (
                user_id,
                hospital_id,
                'admin@kiruddu.go.ug',
                password_hash,
                'Kiruddu Administrator',
                True,
                True
            ))
            print("✓ Created new admin user")

        email = "admin@kiruddu.go.ug"

        conn.commit()

        print(f"✓ Created admin user successfully!")
        print()
        print("=" * 70)
        print("LOGIN CREDENTIALS")
        print("=" * 70)
        print(f"Email:    {email}")
        print(f"Password: {password}")
        print()
        print("URL: http://localhost:5173")
        print("=" * 70)
        print()
        print("⚠️  IMPORTANT: Change this password after first login!")

        cursor.close()
        conn.close()

    except Exception as e:
        print(f"\n✗ Failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python3 create_working_admin.py <hospital_id>")
        sys.exit(1)

    hospital_id = sys.argv[1]
    create_admin(hospital_id)
