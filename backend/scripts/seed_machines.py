#!/usr/bin/env python3
"""
Seed dialysis machines for Kiruddu Hospital
"""
import os
import sys
import uuid
import psycopg2
from psycopg2.extras import execute_values

def seed_machines(hospital_id):
    """Create default dialysis machines for a hospital"""

    # Database connection
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    # Common machines used in Uganda dialysis centers
    machines = [
        {
            'code': 'HD-01',
            'serial': 'FRESENIUS-5008S-001',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 1',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-02',
            'serial': 'FRESENIUS-5008S-002',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 2',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-03',
            'serial': 'FRESENIUS-5008S-003',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 3',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-04',
            'serial': 'FRESENIUS-5008S-004',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 4',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-05',
            'serial': 'FRESENIUS-5008S-005',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 5',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-06',
            'serial': 'FRESENIUS-5008S-006',
            'model': 'Fresenius 5008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Main Dialysis Unit - Bay 6',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-HBV-01',
            'serial': 'FRESENIUS-4008S-HBV-001',
            'model': 'Fresenius 4008S',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Isolation Room - HBV',
            'hbv_dedicated': True,
        },
        {
            'code': 'HD-07',
            'serial': 'GAMBRO-AK200-001',
            'model': 'Gambro AK 200 Ultra S',
            'manufacturer': 'Gambro (Baxter)',
            'location': 'Main Dialysis Unit - Bay 7',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-08',
            'serial': 'GAMBRO-AK200-002',
            'model': 'Gambro AK 200 Ultra S',
            'manufacturer': 'Gambro (Baxter)',
            'location': 'Main Dialysis Unit - Bay 8',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-09',
            'serial': 'NIPRO-SURDIAL-001',
            'model': 'Nipro Surdial 55',
            'manufacturer': 'Nipro',
            'location': 'Main Dialysis Unit - Bay 9',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-10',
            'serial': 'NIPRO-SURDIAL-002',
            'model': 'Nipro Surdial 55',
            'manufacturer': 'Nipro',
            'location': 'Main Dialysis Unit - Bay 10',
            'hbv_dedicated': False,
        },
        {
            'code': 'HD-BACKUP',
            'serial': 'FRESENIUS-4008B-BACKUP',
            'model': 'Fresenius 4008B',
            'manufacturer': 'Fresenius Medical Care',
            'location': 'Storage - Backup Unit',
            'hbv_dedicated': False,
        },
    ]

    print("=" * 70)
    print("Seeding Dialysis Machines for Kiruddu Hospital")
    print("=" * 70)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Check existing machines
        cursor.execute("SELECT machine_code FROM dialysis_machines WHERE hospital_id = %s", (hospital_id,))
        existing = {row[0] for row in cursor.fetchall()}

        machines_to_insert = []
        for machine in machines:
            if machine['code'] not in existing:
                machines_to_insert.append((
                    str(uuid.uuid4()),
                    hospital_id,
                    machine['code'],
                    machine['serial'],
                    machine['model'],
                    machine['manufacturer'],
                    2023,  # manufacture_year
                    machine['location'],
                    'available',  # status
                    machine['hbv_dedicated'],
                ))

        if machines_to_insert:
            query = """
                INSERT INTO dialysis_machines (
                    id, hospital_id, machine_code, serial_number, model,
                    manufacturer, manufacture_year, location, status, is_hbv_dedicated
                ) VALUES %s
            """
            execute_values(cursor, query, machines_to_insert)
            conn.commit()

            print(f"✓ Created {len(machines_to_insert)} dialysis machines:")
            for m in machines_to_insert:
                hbv_label = " [HBV DEDICATED]" if m[9] else ""
                print(f"  - {m[2]}: {m[4]} - {m[7]}{hbv_label}")
        else:
            print("✓ All machines already exist")

        # Get machine IDs for session import
        cursor.execute("""
            SELECT id, machine_code FROM dialysis_machines
            WHERE hospital_id = %s
            ORDER BY machine_code
        """, (hospital_id,))

        machine_ids = cursor.fetchall()
        default_machine_id = machine_ids[0][0] if machine_ids else None

        print()
        print("=" * 70)
        print("Machine Setup Complete!")
        print("=" * 70)
        print(f"Total machines: {len(existing) + len(machines_to_insert)}")
        print(f"Available for regular dialysis: {len([m for m in machines_to_insert if not m[9]])}")
        print(f"HBV dedicated: {len([m for m in machines_to_insert if m[9]])}")

        if default_machine_id:
            print()
            print(f"Default Machine ID for imports: {default_machine_id}")

        cursor.close()
        conn.close()

        return default_machine_id

    except Exception as e:
        print(f"\n✗ Failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python3 seed_machines.py <hospital_id>")
        sys.exit(1)

    hospital_id = sys.argv[1]
    seed_machines(hospital_id)
