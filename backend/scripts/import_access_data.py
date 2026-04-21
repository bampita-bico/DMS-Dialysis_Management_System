#!/usr/bin/env python3
"""
Import MS Access CSV exports into DMS PostgreSQL database
Usage: python3 import_access_data.py <hospital_id>
"""

import csv
import sys
import os
from datetime import datetime
import psycopg2
from psycopg2.extras import execute_values
import uuid

def parse_date(date_str):
    """Parse Access date format to PostgreSQL timestamp"""
    if not date_str or date_str.strip() == '':
        return None

    # Try different date formats
    formats = [
        '%d-%b-%Y %H:%M:%S',  # 18-Apr-1972 0:00:00
        '%d-%b-%Y',           # 18-Apr-1972
        '%Y-%m-%d',           # 2025-01-29
        '%m/%d/%Y',           # 01/29/2025
    ]

    for fmt in formats:
        try:
            return datetime.strptime(date_str.strip(), fmt)
        except:
            continue

    print(f"Warning: Could not parse date: {date_str}")
    return None

def import_patients(conn, cursor, hospital_id, csv_path):
    """Import patients from Patients.csv"""
    print("Importing patients...")

    with open(csv_path, 'r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        patients = []

        for row in reader:
            # Map MS Access fields to DMS schema
            patient_id = str(uuid.uuid4())
            mrn = row.get('HospitalNumber', '').strip()
            full_name = row.get('FullName', '').strip() or f"Patient-{uuid.uuid4().hex[:8]}"

            # Parse sex
            sex_raw = row.get('Sex', '').strip().lower()
            sex = 'male' if 'male' in sex_raw else 'female' if 'female' in sex_raw else 'unknown'

            # Parse blood type
            blood_group_raw = row.get('BloodGroup', '').strip()
            blood_type_map = {
                'A+': 'a_positive', 'A-': 'a_negative',
                'B+': 'b_positive', 'B-': 'b_negative',
                'AB+': 'ab_positive', 'AB-': 'ab_negative',
                'O+': 'o_positive', 'O-': 'o_negative',
            }
            blood_type = blood_type_map.get(blood_group_raw, 'unknown')

            # Parse dates
            dob = parse_date(row.get('DOB', ''))
            if not dob:
                # Default to a reasonable date if missing
                dob = datetime(1970, 1, 1)

            registration_date = parse_date(row.get('DateOfAdmission', '')) or datetime.now()

            # Status
            status = row.get('CurrentPatientStatus', 'Active').strip()
            is_active = status.lower() == 'active'

            # We need a registered_by user - get first user from database
            # For now, use a placeholder UUID that we'll need to update
            placeholder_user = '00000000-0000-0000-0000-000000000001'

            patients.append((
                patient_id,
                hospital_id,
                mrn or f"IMPORT-{uuid.uuid4().hex[:8]}",  # Generate MRN if missing
                full_name,
                dob,
                sex,
                blood_type,
                'unknown',  # marital_status
                registration_date,
                placeholder_user,  # registered_by
                is_active,
            ))

        # Bulk insert
        if patients:
            # First, get a valid user ID to use as registered_by
            cursor.execute("""
                SELECT id FROM users
                WHERE hospital_id = %s
                LIMIT 1
            """, (hospital_id,))
            result = cursor.fetchone()
            if result:
                registered_by = result[0]
                # Update placeholder with real user
                patients = [(p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], registered_by, p[10])
                           for p in patients]

            query = """
                INSERT INTO patients (
                    id, hospital_id, mrn, full_name, date_of_birth, sex,
                    blood_type, marital_status, registration_date, registered_by,
                    is_active
                ) VALUES %s
                ON CONFLICT (hospital_id, mrn) DO NOTHING
            """
            execute_values(cursor, query, patients)
            conn.commit()
            print(f"✓ Imported {len(patients)} patients")
            return len(patients)

    return 0

def import_sessions(conn, cursor, hospital_id, csv_path):
    """Import dialysis sessions from DialysisSessions.csv"""
    print("Importing dialysis sessions...")

    # First, get patient mappings (MRN -> UUID)
    cursor.execute("""
        SELECT mrn, id FROM patients WHERE hospital_id = %s
    """, (hospital_id,))
    patient_map = {row[0]: row[1] for row in cursor.fetchall()}

    # Get default machine
    cursor.execute("""
        SELECT id FROM dialysis_machines
        WHERE hospital_id = %s AND is_hbv_dedicated = FALSE
        ORDER BY machine_code
        LIMIT 1
    """, (hospital_id,))
    result = cursor.fetchone()
    if not result:
        print("Error: No dialysis machines found. Run seed_machines.py first.")
        return 0
    default_machine_id = result[0]

    with open(csv_path, 'r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        sessions = []

        for row in reader:
            # Map to patient
            hospital_number = row.get('HospitalNumber', '').strip()
            patient_uuid = patient_map.get(hospital_number)

            if not patient_uuid:
                continue

            session_id = str(uuid.uuid4())
            session_date = parse_date(row.get('SessionDate', ''))
            if not session_date:
                continue

            # Duration in minutes
            try:
                duration_hours = float(row.get('DurationHours', 4.0) or 4.0)
                duration_mins = int(duration_hours * 60)
            except:
                duration_mins = 240

            # Determine shift based on time of day (default to morning)
            shift = 'morning'

            sessions.append((
                session_id,
                hospital_id,
                patient_uuid,
                default_machine_id,  # machine_id
                'hd',  # modality
                shift,  # shift
                'completed',  # status
                session_date,  # scheduled_date
                '08:00:00',  # scheduled_start_time
                duration_mins,  # prescribed_duration_mins
                duration_mins,  # actual_duration_mins
                row.get('Notes', ''),  # session_notes
            ))

        # Bulk insert
        if sessions:
            query = """
                INSERT INTO dialysis_sessions (
                    id, hospital_id, patient_id, machine_id, modality, shift,
                    status, scheduled_date, scheduled_start_time,
                    prescribed_duration_mins, actual_duration_mins, session_notes
                ) VALUES %s
                ON CONFLICT DO NOTHING
            """
            execute_values(cursor, query, sessions)
            conn.commit()
            print(f"✓ Imported {len(sessions)} dialysis sessions")
            return len(sessions)

    return 0

def import_vitals(conn, cursor, hospital_id, csv_path):
    """Import vitals from Vitals.csv and update session records"""
    print("Importing vitals...")

    # Get session mappings
    cursor.execute("""
        SELECT p.mrn, ds.id, ds.scheduled_date
        FROM dialysis_sessions ds
        JOIN patients p ON ds.patient_id = p.id
        WHERE ds.hospital_id = %s
        ORDER BY ds.scheduled_date
    """, (hospital_id,))
    session_map = {}
    for row in cursor.fetchall():
        mrn, session_id, date = row
        key = f"{mrn}_{date.strftime('%Y-%m-%d')}"
        session_map[key] = session_id

    with open(csv_path, 'r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        pre_vitals = {}
        post_vitals = {}

        for row in reader:
            hospital_number = row.get('HospitalNumber', '').strip()
            session_date_str = row.get('SessionDate', '').strip()
            session_date = parse_date(session_date_str)

            if not session_date:
                continue

            key = f"{hospital_number}_{session_date.strftime('%Y-%m-%d')}"
            session_id = session_map.get(key)

            if not session_id:
                continue

            # Parse vitals
            def safe_float(val):
                try:
                    return float(val) if val and val.strip() else None
                except:
                    return None

            timing = row.get('Timing', 'pre').strip().lower()
            vitals = {
                'sbp': safe_float(row.get('SBP')),
                'dbp': safe_float(row.get('DBP')),
                'hr': safe_float(row.get('HR')),
                'weight': safe_float(row.get('Weight')),
                'temp': safe_float(row.get('Temp')),
            }

            if 'pre' in timing:
                pre_vitals[session_id] = vitals
            elif 'post' in timing:
                post_vitals[session_id] = vitals

        # Update sessions with vitals
        updated = 0
        for session_id in set(list(pre_vitals.keys()) + list(post_vitals.keys())):
            pre = pre_vitals.get(session_id, {})
            post = post_vitals.get(session_id, {})

            cursor.execute("""
                UPDATE dialysis_sessions
                SET pre_weight_kg = %s,
                    pre_bp_systolic = %s,
                    pre_bp_diastolic = %s,
                    pre_hr = %s,
                    pre_temp = %s,
                    post_weight_kg = %s,
                    post_bp_systolic = %s,
                    post_bp_diastolic = %s,
                    post_hr = %s
                WHERE id = %s
            """, (
                pre.get('weight'),
                int(pre.get('sbp')) if pre.get('sbp') else None,
                int(pre.get('dbp')) if pre.get('dbp') else None,
                int(pre.get('hr')) if pre.get('hr') else None,
                pre.get('temp'),
                post.get('weight'),
                int(post.get('sbp')) if post.get('sbp') else None,
                int(post.get('dbp')) if post.get('dbp') else None,
                int(post.get('hr')) if post.get('hr') else None,
                session_id,
            ))
            updated += 1

        conn.commit()
        print(f"✓ Updated {updated} sessions with vital signs")
        return updated

    return 0

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 import_access_data.py <hospital_uuid>")
        sys.exit(1)

    hospital_id = sys.argv[1]

    # Database connection
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': os.getenv('DB_PORT', '5432'),
        'database': os.getenv('DB_NAME', 'dms'),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', ''),
    }

    print("=" * 60)
    print("DMS MS Access Data Import")
    print("=" * 60)
    print(f"Hospital ID: {hospital_id}")
    print(f"Database: {db_config['database']}@{db_config['host']}")
    print()

    # Find CSV directory
    csv_dir = '../MS_Access_ExportS'
    if not os.path.exists(csv_dir):
        csv_dir = 'MS_Access_ExportS'
    if not os.path.exists(csv_dir):
        print(f"Error: Could not find MS_Access_ExportS directory")
        sys.exit(1)

    try:
        conn = psycopg2.connect(**db_config)
        cursor = conn.cursor()

        # Import in order
        patients_count = import_patients(
            conn, cursor, hospital_id,
            os.path.join(csv_dir, 'Patients.csv')
        )

        sessions_count = import_sessions(
            conn, cursor, hospital_id,
            os.path.join(csv_dir, 'DialysisSessions.csv')
        )

        vitals_count = import_vitals(
            conn, cursor, hospital_id,
            os.path.join(csv_dir, 'Vitals.csv')
        )

        print()
        print("=" * 60)
        print("Import Summary:")
        print("=" * 60)
        print(f"Patients imported: {patients_count}")
        print(f"Sessions imported: {sessions_count}")
        print(f"Vitals imported: {vitals_count}")
        print()
        print("✓ Import completed successfully!")

        cursor.close()
        conn.close()

    except Exception as e:
        print(f"\n✗ Import failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == '__main__':
    main()
