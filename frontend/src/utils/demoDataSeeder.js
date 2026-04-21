import db from '../db/schema';

/**
 * Demo Data Seeder for Investor Presentations
 * Creates realistic African patient data in IndexedDB
 * DO NOT use in production - Kiruddu will have real patients
 */

const DEMO_PATIENTS = [
  {
    full_name: 'Nakato Sarah',
    mrn: 'KRD001',
    national_id: 'CM90012345678',
    date_of_birth: '1978-03-15',
    sex: 'female',
    blood_type: 'O+',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Teacher',
    education_level: 'tertiary',
  },
  {
    full_name: 'Okello James',
    mrn: 'KRD002',
    national_id: 'CM89012345679',
    date_of_birth: '1965-07-22',
    sex: 'male',
    blood_type: 'A+',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Protestant',
    occupation: 'Business Owner',
    education_level: 'secondary',
  },
  {
    full_name: 'Nambi Grace',
    mrn: 'KRD003',
    national_id: 'CM95012345680',
    date_of_birth: '1985-11-08',
    sex: 'female',
    blood_type: 'B+',
    marital_status: 'single',
    nationality: 'Ugandan',
    religion: 'Muslim',
    occupation: 'Nurse',
    education_level: 'tertiary',
  },
  {
    full_name: 'Musoke Daniel',
    mrn: 'KRD004',
    national_id: 'CM82012345681',
    date_of_birth: '1972-05-30',
    sex: 'male',
    blood_type: 'AB+',
    marital_status: 'widowed',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Accountant',
    education_level: 'postgraduate',
  },
  {
    full_name: 'Auma Betty',
    mrn: 'KRD005',
    national_id: 'CM92012345682',
    date_of_birth: '1980-09-12',
    sex: 'female',
    blood_type: 'O-',
    marital_status: 'divorced',
    nationality: 'Ugandan',
    religion: 'Protestant',
    occupation: 'Social Worker',
    education_level: 'tertiary',
  },
  {
    full_name: 'Ssemakula Patrick',
    mrn: 'KRD006',
    national_id: 'CM88012345683',
    date_of_birth: '1968-01-25',
    sex: 'male',
    blood_type: 'A-',
    marital_status: 'married',
    nationality: 'Ugandan',
    religion: 'Catholic',
    occupation: 'Engineer',
    education_level: 'tertiary',
  },
];

const DEMO_STAFF = [
  {
    full_name: 'Dr. Kiwanuka Robert',
    email: 'r.kiwanuka@kiruddu.go.ug',
    phone_number: '+256701234567',
    cadre: 'nephrologist',
    specialty: 'Nephrology',
    license_number: 'UMC-12345',
    is_active: true,
  },
  {
    full_name: 'Sr. Nalongo Mary',
    email: 'm.nalongo@kiruddu.go.ug',
    phone_number: '+256702345678',
    cadre: 'nurse',
    specialty: 'Dialysis Nursing',
    license_number: 'UNMC-23456',
    is_active: true,
  },
  {
    full_name: 'Sr. Akello Christine',
    email: 'c.akello@kiruddu.go.ug',
    phone_number: '+256703456789',
    cadre: 'nurse',
    specialty: 'Critical Care',
    license_number: 'UNMC-23457',
    is_active: true,
  },
];

const DEMO_MACHINES = [
  {
    machine_number: 'HD-001',
    manufacturer: 'Fresenius',
    model: '5008S',
    serial_number: 'FMC-5008-2023-001',
    operational_status: 'operational',
    location: 'Dialysis Unit A',
  },
  {
    machine_number: 'HD-002',
    manufacturer: 'Fresenius',
    model: '5008S',
    serial_number: 'FMC-5008-2023-002',
    operational_status: 'operational',
    location: 'Dialysis Unit A',
  },
  {
    machine_number: 'HD-003',
    manufacturer: 'Gambro',
    model: 'AK 200',
    serial_number: 'GMB-AK200-2022-003',
    operational_status: 'operational',
    location: 'Dialysis Unit B',
  },
];

export async function seedDemoData() {
  console.log('🌱 Starting demo data seeding...');

  try {
    // Check if already seeded
    const existingPatients = await db.patients.toArray();
    if (existingPatients.length > 0) {
      console.log('✅ Demo data already exists. Skipping seed.');
      return;
    }

    const hospitalId = 'demo_hospital';
    const timestamp = new Date().toISOString();

    // 1. Seed Staff (needed for doctor assignments)
    const staffIds = [];
    for (const staff of DEMO_STAFF) {
      const staffId = `staff_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      await db.staff_profiles.put({
        id: staffId,
        ...staff,
        hospital_id: hospitalId,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
      staffIds.push(staffId);
    }

    // 2. Seed Machines
    const machineIds = [];
    for (const machine of DEMO_MACHINES) {
      const machineId = `machine_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      await db.dialysis_machines.put({
        id: machineId,
        ...machine,
        hospital_id: hospitalId,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });
      machineIds.push(machineId);
    }

    // 3. Seed Patients with full data
    for (let i = 0; i < DEMO_PATIENTS.length; i++) {
      const patient = DEMO_PATIENTS[i];
      const patientId = `patient_${Date.now()}_${i}`;

      // Create patient
      await db.patients.put({
        id: patientId,
        ...patient,
        hospital_id: hospitalId,
        registered_by: staffIds[0], // First doctor
        primary_doctor_id: staffIds[0],
        registration_date: new Date().toISOString().split('T')[0],
        primary_language: 'English',
        interpreter_needed: false,
        is_active: true,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add contacts
      await db.patient_contacts.put({
        id: `contact_phone_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'phone',
        value: `+25670${7000000 + i}`,
        label: 'Primary Phone',
        is_primary: true,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      await db.patient_contacts.put({
        id: `contact_address_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'address',
        value: `Kampala, Nakawa Division, Plot ${100 + i}`,
        label: 'Home Address',
        is_primary: true,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add emergency contact
      await db.patient_contacts.put({
        id: `contact_emergency_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        contact_type: 'emergency',
        value: `Emergency Contact ${i}|+25677${1000000 + i}|Spouse`,
        label: 'Emergency Contact',
        is_primary: false,
        is_verified: false,
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add vascular access
      const accessTypes = ['avf', 'avg', 'tunneled_catheter'];
      const accessType = accessTypes[i % 3];
      await db.vascular_access.put({
        id: `access_${patientId}`,
        patient_id: patientId,
        hospital_id: hospitalId,
        access_type: accessType,
        access_site: accessType === 'tunneled_catheter' ? 'internal_jugular' : 'forearm',
        site_side: i % 2 === 0 ? 'left' : 'right',
        insertion_date: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        first_use_date: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        status: 'active',
        is_primary_access: true,
        inserted_by: staffIds[0],
        synced: false,
        created_at: timestamp,
        updated_at: timestamp,
      });

      // Add 2-3 recent sessions for each patient
      for (let j = 0; j < 3; j++) {
        const sessionDate = new Date(Date.now() - j * 2 * 24 * 60 * 60 * 1000);
        const sessionId = `session_${patientId}_${j}`;
        
        await db.dialysis_sessions.put({
          id: sessionId,
          patient_id: patientId,
          hospital_id: hospitalId,
          machine_id: machineIds[i % machineIds.length],
          scheduled_date: sessionDate.toISOString().split('T')[0],
          scheduled_start_time: '08:00',
          shift: 'morning',
          prescribed_duration_mins: 240,
          modality: 'hd',
          status: j === 0 ? 'scheduled' : 'completed',
          primary_nurse_id: staffIds[1], // First nurse
          supervising_doctor_id: staffIds[0], // Doctor
          was_patient_reviewed: j === 0 ? false : true,
          synced: false,
          created_at: timestamp,
          updated_at: timestamp,
        });
      }
    }

    console.log('✅ Demo data seeded successfully!');
    console.log(`   - ${DEMO_PATIENTS.length} patients`);
    console.log(`   - ${DEMO_STAFF.length} staff members`);
    console.log(`   - ${DEMO_MACHINES.length} machines`);
    console.log(`   - ${DEMO_PATIENTS.length * 3} contacts per patient`);
    console.log(`   - ${DEMO_PATIENTS.length} vascular access records`);
    console.log(`   - ${DEMO_PATIENTS.length * 3} dialysis sessions`);

  } catch (error) {
    console.error('❌ Error seeding demo data:', error);
    throw error;
  }
}

// Function to clear demo data (for re-seeding)
export async function clearDemoData() {
  console.log('🗑️  Clearing demo data...');
  await db.patients.clear();
  await db.patient_contacts.clear();
  await db.vascular_access.clear();
  await db.dialysis_sessions.clear();
  await db.staff_profiles.clear();
  await db.dialysis_machines.clear();
  console.log('✅ Demo data cleared');
}
