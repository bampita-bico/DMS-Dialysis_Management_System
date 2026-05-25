import { useState } from 'react';
import FormField from './FormField';
import useOfflineData from '../../hooks/useOfflineData';
import offlineService from '../../services/offlineService';
import { authService } from '../../services/auth';

export default function PatientForm({ patient = null, onSuccess, onCancel }) {
  const { data: staff } = useOfflineData('staff_profiles');

  const [formData, setFormData] = useState({
    // Basic Information
    full_name: patient?.full_name || '',
    preferred_name: patient?.preferred_name || '',
    mrn: patient?.mrn || '',
    national_id: patient?.national_id || '',
    date_of_birth: patient?.date_of_birth?.split('T')[0] || '',
    sex: patient?.sex || '',
    blood_type: patient?.blood_type || 'unknown',

    // Demographics
    marital_status: patient?.marital_status || 'unknown',
    nationality: patient?.nationality || 'Ugandan',
    religion: patient?.religion || '',
    occupation: patient?.occupation || '',
    education_level: patient?.education_level || '',

    // Language & Communication
    primary_language: patient?.primary_language || 'English',
    interpreter_needed: patient?.interpreter_needed || false,

    // Medical
    primary_doctor_id: patient?.primary_doctor_id || '',
    renal_course: patient?.renal_course || '',
    kidney_disease_cause: patient?.kidney_disease_cause || '',
    dialysis_indication: patient?.dialysis_indication || '',
    symptom_onset_date: patient?.symptom_onset_date?.split('T')[0] || '',
    diagnosis_date: patient?.diagnosis_date?.split('T')[0] || '',
    dialysis_start_date: patient?.dialysis_start_date?.split('T')[0] || '',
    aki_stage: patient?.aki_stage || '',
    baseline_creatinine_mg_dl: patient?.baseline_creatinine_mg_dl || '',
    highest_creatinine_mg_dl: patient?.highest_creatinine_mg_dl || '',
    urine_output_ml_day: patient?.urine_output_ml_day || '',
    residual_kidney_function: patient?.residual_kidney_function || '',
    reversible_cause: patient?.reversible_cause || '',
    recovery_plan: patient?.recovery_plan || '',
    stop_dialysis_trial_date: patient?.stop_dialysis_trial_date?.split('T')[0] || '',
    aki_to_ckd_reassessment_due: patient?.aki_to_ckd_reassessment_due?.split('T')[0] || '',
    malaria_aki_phenotype: patient?.malaria_aki_phenotype || 'not_applicable',
    malaria_symptom_onset_date: patient?.malaria_symptom_onset_date?.split('T')[0] || '',
    fever_onset_date: patient?.fever_onset_date?.split('T')[0] || '',
    dark_urine_onset_date: patient?.dark_urine_onset_date?.split('T')[0] || '',
    malaria_test_date: patient?.malaria_test_date?.split('T')[0] || '',
    parasitemia: patient?.parasitemia || '',
    first_antimalarial_dose_at: patient?.first_antimalarial_dose_at ? patient.first_antimalarial_dose_at.slice(0, 16) : '',
    definitive_antimalarial_regimen: patient?.definitive_antimalarial_regimen || '',
    parasite_clearance_date: patient?.parasite_clearance_date?.split('T')[0] || '',
    pregnancy_related: patient?.pregnancy_related || false,
    gestational_age_weeks: patient?.gestational_age_weeks || '',
    expected_delivery_date: patient?.expected_delivery_date?.split('T')[0] || '',
    gravida: patient?.gravida || '',
    para: patient?.para || '',
    anc_visits: patient?.anc_visits || '',
    bp_proteinuria_summary: patient?.bp_proteinuria_summary || '',
    preeclampsia_severity: patient?.preeclampsia_severity || '',
    delivery_date: patient?.delivery_date?.split('T')[0] || '',
    fetal_outcome: patient?.fetal_outcome || '',
    magnesium_sulfate_given: patient?.magnesium_sulfate_given || false,
    pregnancy_antihypertensive_regimen: patient?.pregnancy_antihypertensive_regimen || '',
    postpartum_renal_recovery: patient?.postpartum_renal_recovery || '',
    maternal_fetal_follow_up: patient?.maternal_fetal_follow_up || '',
    transplant_referral_status: patient?.transplant_referral_status || '',
    modality_plan: patient?.modality_plan || '',
    conservative_care_flag: patient?.conservative_care_flag || false,
    tracking_notes: patient?.tracking_notes || '',

    // Contact Info (will be saved separately)
    phone_number: '',
    email: '',
    physical_address: '',

    // Emergency Contact
    emergency_contact_name: '',
    emergency_contact_phone: '',
    emergency_contact_relationship: ''
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const doctors = (staff || []).filter(s =>
    (s.cadre === 'doctor' || s.cadre === 'nephrologist') && s.is_active
  );

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    setErrors(prev => ({ ...prev, [name]: '' }));
  };

  const validate = () => {
    const newErrors = {};
    if (!formData.full_name?.trim()) newErrors.full_name = 'Full name is required';
    if (!formData.mrn?.trim()) newErrors.mrn = 'MRN is required';
    if (!formData.date_of_birth) newErrors.date_of_birth = 'Date of birth is required';
    if (new Date(formData.date_of_birth) > new Date()) newErrors.date_of_birth = 'Cannot be in future';
    if (!formData.sex) newErrors.sex = 'Sex is required';
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const newErrors = validate();
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setLoading(true);
    try {
      const currentUser = authService.getCurrentUser();

      // Patient main data
      const patientPayload = {
        full_name: formData.full_name,
        preferred_name: formData.preferred_name || null,
        mrn: formData.mrn,
        national_id: formData.national_id || null,
        date_of_birth: formData.date_of_birth,
        sex: formData.sex,
        blood_type: formData.blood_type,
        marital_status: formData.marital_status,
        nationality: formData.nationality,
        religion: formData.religion || null,
        occupation: formData.occupation || null,
        education_level: formData.education_level || null,
        primary_language: formData.primary_language,
        interpreter_needed: formData.interpreter_needed,
        primary_doctor_id: formData.primary_doctor_id || null,
        renal_course: formData.renal_course || null,
        kidney_disease_cause: formData.kidney_disease_cause || null,
        dialysis_indication: formData.dialysis_indication || null,
        symptom_onset_date: formData.symptom_onset_date || null,
        diagnosis_date: formData.diagnosis_date || null,
        dialysis_start_date: formData.dialysis_start_date || null,
        aki_stage: formData.aki_stage || null,
        baseline_creatinine_mg_dl: formData.baseline_creatinine_mg_dl || null,
        highest_creatinine_mg_dl: formData.highest_creatinine_mg_dl || null,
        urine_output_ml_day: formData.urine_output_ml_day || null,
        residual_kidney_function: formData.residual_kidney_function || null,
        reversible_cause: formData.reversible_cause || null,
        recovery_plan: formData.recovery_plan || null,
        stop_dialysis_trial_date: formData.stop_dialysis_trial_date || null,
        aki_to_ckd_reassessment_due: formData.aki_to_ckd_reassessment_due || null,
        malaria_aki_phenotype: formData.malaria_aki_phenotype,
        malaria_symptom_onset_date: formData.malaria_symptom_onset_date || null,
        fever_onset_date: formData.fever_onset_date || null,
        dark_urine_onset_date: formData.dark_urine_onset_date || null,
        malaria_test_date: formData.malaria_test_date || null,
        parasitemia: formData.parasitemia || null,
        first_antimalarial_dose_at: formData.first_antimalarial_dose_at || null,
        definitive_antimalarial_regimen: formData.definitive_antimalarial_regimen || null,
        parasite_clearance_date: formData.parasite_clearance_date || null,
        pregnancy_related: formData.pregnancy_related,
        gestational_age_weeks: formData.gestational_age_weeks || null,
        expected_delivery_date: formData.expected_delivery_date || null,
        gravida: formData.gravida || null,
        para: formData.para || null,
        anc_visits: formData.anc_visits || null,
        bp_proteinuria_summary: formData.bp_proteinuria_summary || null,
        preeclampsia_severity: formData.preeclampsia_severity || null,
        delivery_date: formData.delivery_date || null,
        fetal_outcome: formData.fetal_outcome || null,
        magnesium_sulfate_given: formData.magnesium_sulfate_given,
        pregnancy_antihypertensive_regimen: formData.pregnancy_antihypertensive_regimen || null,
        postpartum_renal_recovery: formData.postpartum_renal_recovery || null,
        maternal_fetal_follow_up: formData.maternal_fetal_follow_up || null,
        transplant_referral_status: formData.transplant_referral_status || null,
        modality_plan: formData.modality_plan || null,
        conservative_care_flag: formData.conservative_care_flag,
        tracking_notes: formData.tracking_notes || null,
        hospital_id: currentUser?.hospital_id,
        registered_by: currentUser?.id,
        registration_date: new Date().toISOString().split('T')[0],
        is_active: true,
      };

      let patientId;
      if (patient) {
        await offlineService.update('patients', patient.id, patientPayload);
        patientId = patient.id;
      } else {
        const createdPatient = await offlineService.create('patients', patientPayload, 10);
        patientId = createdPatient.id;
      }

      // Save contacts
      if (formData.phone_number) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'phone',
          value: formData.phone_number,
          label: 'Primary Phone',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      if (formData.email) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'email',
          value: formData.email,
          label: 'Primary Email',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      if (formData.physical_address) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'address',
          value: formData.physical_address,
          label: 'Home Address',
          is_primary: true,
          is_verified: false
        }, 10);
      }

      // Save emergency contact
      if (formData.emergency_contact_name) {
        await offlineService.create('patient_contacts', {
          patient_id: patientId,
          hospital_id: currentUser?.hospital_id,
          contact_type: 'emergency',
          value: `${formData.emergency_contact_name}|${formData.emergency_contact_phone}|${formData.emergency_contact_relationship}`,
          label: 'Emergency Contact',
          is_primary: false,
          is_verified: false
        }, 10);
      }

      onSuccess?.();
    } catch (error) {
      setErrors({ submit: error.message || 'Failed to save patient' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 max-h-[70vh] overflow-y-auto px-1">
      {errors.submit && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {errors.submit}
        </div>
      )}

      {/* Basic Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Basic Information</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Full Name" name="full_name" value={formData.full_name}
            onChange={handleChange} error={errors.full_name} required />

          <FormField label="Preferred Name" name="preferred_name" value={formData.preferred_name}
            onChange={handleChange} placeholder="Nickname or preferred name" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="MRN (Patient Number)" name="mrn" value={formData.mrn}
            onChange={handleChange} error={errors.mrn} required />

          <FormField label="National ID" name="national_id" value={formData.national_id}
            onChange={handleChange} placeholder="National ID or Passport" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Date of Birth" name="date_of_birth" type="date"
            value={formData.date_of_birth} onChange={handleChange} error={errors.date_of_birth} required />

          <FormField label="Sex" name="sex" type="select" value={formData.sex}
            onChange={handleChange} error={errors.sex} required
            options={[
              { value: 'male', label: 'Male' },
              { value: 'female', label: 'Female' },
              { value: 'intersex', label: 'Intersex' },
              { value: 'unknown', label: 'Unknown' }
            ]} />
        </div>

        <FormField label="Blood Type" name="blood_type" type="select"
          value={formData.blood_type} onChange={handleChange} required
          options={[
            { value: 'unknown', label: 'Unknown' },
            { value: 'A+', label: 'A+' }, { value: 'A-', label: 'A-' },
            { value: 'B+', label: 'B+' }, { value: 'B-', label: 'B-' },
            { value: 'AB+', label: 'AB+' }, { value: 'AB-', label: 'AB-' },
            { value: 'O+', label: 'O+' }, { value: 'O-', label: 'O-' }
          ]} />
      </div>

      {/* Demographics */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Demographics</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Marital Status" name="marital_status" type="select"
            value={formData.marital_status} onChange={handleChange}
            options={[
              { value: 'unknown', label: 'Unknown' },
              { value: 'single', label: 'Single' },
              { value: 'married', label: 'Married' },
              { value: 'divorced', label: 'Divorced' },
              { value: 'widowed', label: 'Widowed' }
            ]} />

          <FormField label="Nationality" name="nationality" value={formData.nationality}
            onChange={handleChange} />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Religion" name="religion" value={formData.religion}
            onChange={handleChange} placeholder="e.g., Catholic, Muslim, etc." />

          <FormField label="Occupation" name="occupation" value={formData.occupation}
            onChange={handleChange} />
        </div>

        <FormField label="Education Level" name="education_level" type="select"
          value={formData.education_level} onChange={handleChange}
          options={[
            { value: '', label: 'Not specified' },
            { value: 'none', label: 'None' },
            { value: 'primary', label: 'Primary' },
            { value: 'secondary', label: 'Secondary' },
            { value: 'tertiary', label: 'Tertiary/University' },
            { value: 'postgraduate', label: 'Postgraduate' }
          ]} />
      </div>

      {/* Contact Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Contact Information</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Phone Number" name="phone_number" value={formData.phone_number}
            onChange={handleChange} placeholder="+256..." />

          <FormField label="Email Address" name="email" type="email" value={formData.email}
            onChange={handleChange} />
        </div>

        <FormField label="Physical Address" name="physical_address" type="textarea"
          value={formData.physical_address} onChange={handleChange} rows={2}
          placeholder="Village, Parish, District..." />
      </div>

      {/* Language & Communication */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Language & Communication</h3>

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Primary Language" name="primary_language" value={formData.primary_language}
            onChange={handleChange} />

          <div className="flex items-center pt-8">
            <input type="checkbox" name="interpreter_needed" checked={formData.interpreter_needed}
              onChange={handleChange} className="h-4 w-4 text-sky-600 rounded mr-2" />
            <label className="text-sm text-gray-700">Interpreter Needed</label>
          </div>
        </div>
      </div>

      {/* Medical Information */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Medical Information</h3>

        <FormField label="Primary Doctor" name="primary_doctor_id" type="select"
          value={formData.primary_doctor_id} onChange={handleChange}
          options={doctors.map(d => ({
            value: d.id,
            label: d.full_name
          }))} />

        <FormField label="Renal Course" name="renal_course" type="select"
          value={formData.renal_course} onChange={handleChange}
          options={[
            { value: 'AKI requiring dialysis', label: 'AKI requiring dialysis' },
            { value: 'AKI on CKD watch', label: 'AKI on CKD watch' },
            { value: 'CKD on maintenance hemodialysis', label: 'CKD on maintenance hemodialysis' },
            { value: 'ESKD on maintenance hemodialysis', label: 'ESKD on maintenance hemodialysis' },
            { value: 'Pregnancy-associated AKI', label: 'Pregnancy-associated AKI' },
            { value: 'Recovery trial / dialysis reduction', label: 'Recovery trial / dialysis reduction' }
          ]} />

        <FormField label="Kidney Disease Cause" name="kidney_disease_cause" type="select"
          value={formData.kidney_disease_cause} onChange={handleChange}
          options={[
            { value: 'Diabetic kidney disease', label: 'Diabetic kidney disease' },
            { value: 'Hypertensive nephrosclerosis', label: 'Hypertensive nephrosclerosis' },
            { value: 'Lupus nephritis', label: 'Lupus nephritis' },
            { value: 'HIV-associated kidney disease', label: 'HIV-associated kidney disease' },
            { value: 'Hepatitis-related kidney disease', label: 'Hepatitis-related kidney disease' },
            { value: 'Post-infectious glomerulonephritis', label: 'Post-infectious glomerulonephritis' },
            { value: 'Sepsis-associated AKI', label: 'Sepsis-associated AKI' },
            { value: 'Severe malaria with blackwater fever', label: 'Severe malaria with blackwater fever' },
            { value: 'Severe malaria without blackwater fever', label: 'Severe malaria without blackwater fever' },
            { value: 'Malaria without AKI', label: 'Malaria without AKI' },
            { value: 'Pregnancy-associated AKI / preeclampsia', label: 'Pregnancy-associated AKI / preeclampsia' },
            { value: 'Eclampsia / HELLP-associated AKI', label: 'Eclampsia / HELLP-associated AKI' },
            { value: 'Postpartum hemorrhage-associated AKI', label: 'Postpartum hemorrhage-associated AKI' },
            { value: 'Obstructive uropathy', label: 'Obstructive uropathy' },
            { value: 'Stones / recurrent obstruction', label: 'Stones / recurrent obstruction' },
            { value: 'Prostate disease', label: 'Prostate disease' },
            { value: 'Sickle cell nephropathy', label: 'Sickle cell nephropathy' },
            { value: 'Myeloma kidney', label: 'Myeloma kidney' },
            { value: 'Drug/toxin-associated AKI', label: 'Drug/toxin-associated AKI' },
            { value: 'NSAID-associated AKI', label: 'NSAID-associated AKI' },
            { value: 'Aminoglycoside-associated AKI', label: 'Aminoglycoside-associated AKI' },
            { value: 'Snakebite-associated AKI', label: 'Snakebite-associated AKI' },
            { value: 'Primary glomerulonephritis', label: 'Primary glomerulonephritis' },
            { value: 'Congenital/urologic disease', label: 'Congenital/urologic disease' },
            { value: 'Biopsy-proven kidney disease', label: 'Biopsy-proven kidney disease' },
            { value: 'Unknown / under investigation', label: 'Unknown / under investigation' }
          ]} />

        <FormField label="Dialysis Indication" name="dialysis_indication"
          value={formData.dialysis_indication} onChange={handleChange}
          placeholder="e.g., hyperkalemia, volume overload, uremic symptoms" />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="Symptom Onset" name="symptom_onset_date" type="date"
            value={formData.symptom_onset_date} onChange={handleChange} />

          <FormField label="Diagnosis Date" name="diagnosis_date" type="date"
            value={formData.diagnosis_date} onChange={handleChange} />

          <FormField label="Dialysis Start" name="dialysis_start_date" type="date"
            value={formData.dialysis_start_date} onChange={handleChange} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="AKI Stage" name="aki_stage" type="select"
            value={formData.aki_stage} onChange={handleChange}
            options={[
              { value: 'stage_1', label: 'KDIGO AKI stage 1' },
              { value: 'stage_2', label: 'KDIGO AKI stage 2' },
              { value: 'stage_3', label: 'KDIGO AKI stage 3' },
              { value: 'not_applicable', label: 'Not applicable' }
            ]} />
          <FormField label="Baseline creatinine" name="baseline_creatinine_mg_dl" type="number" step="0.01"
            value={formData.baseline_creatinine_mg_dl} onChange={handleChange} />
          <FormField label="Highest creatinine" name="highest_creatinine_mg_dl" type="number" step="0.01"
            value={formData.highest_creatinine_mg_dl} onChange={handleChange} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="Urine output ml/day" name="urine_output_ml_day" type="number"
            value={formData.urine_output_ml_day} onChange={handleChange} />
          <FormField label="Stop-dialysis trial" name="stop_dialysis_trial_date" type="date"
            value={formData.stop_dialysis_trial_date} onChange={handleChange} />
          <FormField label="AKI-to-CKD review due" name="aki_to_ckd_reassessment_due" type="date"
            value={formData.aki_to_ckd_reassessment_due} onChange={handleChange} />
        </div>

        <FormField label="Residual Kidney Function" name="residual_kidney_function"
          value={formData.residual_kidney_function} onChange={handleChange}
          placeholder="e.g., improving urine output, residual urea clearance, anuric" />

        <FormField label="Reversible Cause" name="reversible_cause"
          value={formData.reversible_cause} onChange={handleChange}
          placeholder="e.g., sepsis treated, obstruction relieved, malaria treated, postpartum recovery" />

        <FormField label="Reversible Cause / Recovery Plan" name="recovery_plan" type="textarea"
          value={formData.recovery_plan} onChange={handleChange} rows={2}
          placeholder="Cause corrected, dialysis reduction plan, relapse criteria, consultant review date..." />

        <FormField label="Malaria AKI Phenotype" name="malaria_aki_phenotype" type="select"
          value={formData.malaria_aki_phenotype} onChange={handleChange}
          options={[
            { value: 'not_applicable', label: 'Not applicable' },
            { value: 'blackwater_fever', label: 'Blackwater fever' },
            { value: 'non_blackwater_severe_malaria', label: 'Severe malaria without blackwater fever' },
            { value: 'malaria_without_aki', label: 'Malaria without AKI' },
            { value: 'other_infection_related_aki', label: 'Other infection-related AKI' }
          ]} />

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <FormField label="Malaria symptoms" name="malaria_symptom_onset_date" type="date"
            value={formData.malaria_symptom_onset_date} onChange={handleChange} />
          <FormField label="Fever onset" name="fever_onset_date" type="date"
            value={formData.fever_onset_date} onChange={handleChange} />
          <FormField label="Dark urine onset" name="dark_urine_onset_date" type="date"
            value={formData.dark_urine_onset_date} onChange={handleChange} />
          <FormField label="Malaria test date" name="malaria_test_date" type="date"
            value={formData.malaria_test_date} onChange={handleChange} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="Parasitemia" name="parasitemia"
            value={formData.parasitemia} onChange={handleChange} placeholder="%, +/+++, or not available" />
          <FormField label="First antimalarial dose" name="first_antimalarial_dose_at" type="datetime-local"
            value={formData.first_antimalarial_dose_at} onChange={handleChange} />
          <FormField label="Parasite clearance" name="parasite_clearance_date" type="date"
            value={formData.parasite_clearance_date} onChange={handleChange} />
        </div>

        <FormField label="Definitive Antimalarial Regimen" name="definitive_antimalarial_regimen"
          value={formData.definitive_antimalarial_regimen} onChange={handleChange}
          placeholder="Drug, dose, dates, transfusion if relevant" />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="flex items-center pt-8">
            <input type="checkbox" name="pregnancy_related" checked={formData.pregnancy_related}
              onChange={handleChange} className="h-4 w-4 text-sky-600 rounded mr-2" />
            <label className="text-sm text-gray-700">Pregnancy Related</label>
          </div>

          <FormField label="Gravida" name="gravida" value={formData.gravida}
            onChange={handleChange} placeholder="G3" />

          <FormField label="Para" name="para" value={formData.para}
            onChange={handleChange} placeholder="P2" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="Gestational age weeks" name="gestational_age_weeks" type="number" step="0.1"
            value={formData.gestational_age_weeks} onChange={handleChange} />
          <FormField label="Expected delivery date" name="expected_delivery_date" type="date"
            value={formData.expected_delivery_date} onChange={handleChange} />
          <FormField label="ANC Visits" name="anc_visits" type="number"
            value={formData.anc_visits} onChange={handleChange} />
        </div>

        <FormField label="BP / Proteinuria Summary" name="bp_proteinuria_summary"
          value={formData.bp_proteinuria_summary} onChange={handleChange}
          placeholder="Highest BP, proteinuria, severe features" />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <FormField label="Preeclampsia Severity" name="preeclampsia_severity" type="select"
            value={formData.preeclampsia_severity} onChange={handleChange}
            options={[
              { value: 'none', label: 'None' },
              { value: 'preeclampsia_without_severe_features', label: 'Preeclampsia without severe features' },
              { value: 'preeclampsia_with_severe_features', label: 'Preeclampsia with severe features' },
              { value: 'eclampsia', label: 'Eclampsia' },
              { value: 'hellp', label: 'HELLP' }
            ]} />
          <FormField label="Delivery date" name="delivery_date" type="date"
            value={formData.delivery_date} onChange={handleChange} />
          <FormField label="Fetal outcome" name="fetal_outcome"
            value={formData.fetal_outcome} onChange={handleChange} />
        </div>

        <div className="flex items-center py-2">
          <input type="checkbox" name="magnesium_sulfate_given" checked={formData.magnesium_sulfate_given}
            onChange={handleChange} className="h-4 w-4 text-sky-600 rounded mr-2" />
          <label className="text-sm text-gray-700">Magnesium sulfate given when clinically indicated</label>
        </div>

        <FormField label="Pregnancy Antihypertensive Regimen" name="pregnancy_antihypertensive_regimen"
          value={formData.pregnancy_antihypertensive_regimen} onChange={handleChange}
          placeholder="Drug, dose, changes, dates, clinician" />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="Postpartum Renal Recovery" name="postpartum_renal_recovery" type="textarea"
            value={formData.postpartum_renal_recovery} onChange={handleChange} rows={2} />
          <FormField label="Maternal-Fetal Follow-Up" name="maternal_fetal_follow_up" type="textarea"
            value={formData.maternal_fetal_follow_up} onChange={handleChange} rows={2} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField label="Transplant Referral Status" name="transplant_referral_status"
            value={formData.transplant_referral_status} onChange={handleChange} />
          <FormField label="Modality Plan" name="modality_plan"
            value={formData.modality_plan} onChange={handleChange}
            placeholder="HD, PD eligibility, home HD, conservative care, palliative plan" />
        </div>

        <div className="flex items-center py-2">
          <input type="checkbox" name="conservative_care_flag" checked={formData.conservative_care_flag}
            onChange={handleChange} className="h-4 w-4 text-sky-600 rounded mr-2" />
          <label className="text-sm text-gray-700">Conservative kidney care / palliative pathway active</label>
        </div>

        <FormField label="Tracking Notes" name="tracking_notes" type="textarea"
          value={formData.tracking_notes} onChange={handleChange} rows={3}
          placeholder="Dialysis recovery plan, drug changes, access concerns, complications, follow-up plan..." />
      </div>

      {/* Emergency Contact */}
      <div className="border-b pb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Emergency Contact</h3>

        <FormField label="Emergency Contact Name" name="emergency_contact_name"
          value={formData.emergency_contact_name} onChange={handleChange} />

        <div className="grid grid-cols-2 gap-4">
          <FormField label="Emergency Contact Phone" name="emergency_contact_phone"
            value={formData.emergency_contact_phone} onChange={handleChange} />

          <FormField label="Relationship" name="emergency_contact_relationship"
            value={formData.emergency_contact_relationship} onChange={handleChange}
            placeholder="e.g., Spouse, Parent, Sibling" />
        </div>
      </div>

      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 sticky bottom-0 bg-white">
        <button type="button" onClick={onCancel} disabled={loading}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50">
          Cancel
        </button>
        <button type="submit" disabled={loading}
          className="px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50">
          {loading ? 'Saving...' : patient ? 'Update Patient' : 'Create Patient'}
        </button>
      </div>
    </form>
  );
}
