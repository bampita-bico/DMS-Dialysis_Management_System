package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientsHandler struct {
	pool *pgxpool.Pool
}

func NewPatientsHandler(pool *pgxpool.Pool) *PatientsHandler {
	return &PatientsHandler{pool: pool}
}

// Blood type mapping from display format to database format
var bloodTypeMapping = map[string]string{
	"A+":      "a_positive",
	"A-":      "a_negative",
	"B+":      "b_positive",
	"B-":      "b_negative",
	"AB+":     "ab_positive",
	"AB-":     "ab_negative",
	"O+":      "o_positive",
	"O-":      "o_negative",
	"unknown": "unknown",
	// Also accept the database format
	"a_positive":  "a_positive",
	"a_negative":  "a_negative",
	"b_positive":  "b_positive",
	"b_negative":  "b_negative",
	"ab_positive": "ab_positive",
	"ab_negative": "ab_negative",
	"o_positive":  "o_positive",
	"o_negative":  "o_negative",
}

func mapBloodType(input string) (string, bool) {
	mapped, ok := bloodTypeMapping[input]
	return mapped, ok
}

func (h *PatientsHandler) Create(c *gin.Context) {
	var req struct {
		MRN               string `json:"mrn" binding:"required"`
		NationalID        string `json:"national_id"`
		FullName          string `json:"full_name" binding:"required"`
		PreferredName     string `json:"preferred_name"`
		DateOfBirth       string `json:"date_of_birth" binding:"required"`
		Sex               string `json:"sex" binding:"required"`
		BloodType         string `json:"blood_type"`
		MaritalStatus     string `json:"marital_status"`
		Nationality       string `json:"nationality"`
		Religion          string `json:"religion"`
		Occupation        string `json:"occupation"`
		EducationLevel    string `json:"education_level"`
		PrimaryLanguage   string `json:"primary_language"`
		InterpreterNeeded bool   `json:"interpreter_needed"`
		PrimaryDoctorID   string `json:"primary_doctor_id"`

		RenalCourse                      string `json:"renal_course"`
		KidneyDiseaseCause               string `json:"kidney_disease_cause"`
		DialysisIndication               string `json:"dialysis_indication"`
		SymptomOnsetDate                 string `json:"symptom_onset_date"`
		DiagnosisDate                    string `json:"diagnosis_date"`
		DialysisStartDate                string `json:"dialysis_start_date"`
		AKIStage                         string `json:"aki_stage"`
		BaselineCreatinineMgDl           any    `json:"baseline_creatinine_mg_dl"`
		HighestCreatinineMgDl            any    `json:"highest_creatinine_mg_dl"`
		UrineOutputMlDay                 any    `json:"urine_output_ml_day"`
		ResidualKidneyFunction           string `json:"residual_kidney_function"`
		ReversibleCause                  string `json:"reversible_cause"`
		RecoveryPlan                     string `json:"recovery_plan"`
		StopDialysisTrialDate            string `json:"stop_dialysis_trial_date"`
		AKIToCKDReassessmentDue          string `json:"aki_to_ckd_reassessment_due"`
		MalariaAKIPhenotype              string `json:"malaria_aki_phenotype"`
		MalariaSymptomOnsetDate          string `json:"malaria_symptom_onset_date"`
		FeverOnsetDate                   string `json:"fever_onset_date"`
		DarkUrineOnsetDate               string `json:"dark_urine_onset_date"`
		MalariaTestDate                  string `json:"malaria_test_date"`
		Parasitemia                      string `json:"parasitemia"`
		FirstAntimalarialDoseAt          string `json:"first_antimalarial_dose_at"`
		DefinitiveAntimalarialRegimen    string `json:"definitive_antimalarial_regimen"`
		ParasiteClearanceDate            string `json:"parasite_clearance_date"`
		PregnancyRelated                 bool   `json:"pregnancy_related"`
		GestationalAgeWeeks              any    `json:"gestational_age_weeks"`
		ExpectedDeliveryDate             string `json:"expected_delivery_date"`
		Gravida                          string `json:"gravida"`
		Para                             string `json:"para"`
		AncVisits                        any    `json:"anc_visits"`
		BPProteinuriaSummary             string `json:"bp_proteinuria_summary"`
		PreeclampsiaSeverity             string `json:"preeclampsia_severity"`
		DeliveryDate                     string `json:"delivery_date"`
		FetalOutcome                     string `json:"fetal_outcome"`
		MagnesiumSulfateGiven            any    `json:"magnesium_sulfate_given"`
		PregnancyAntihypertensiveRegimen string `json:"pregnancy_antihypertensive_regimen"`
		PostpartumRenalRecovery          string `json:"postpartum_renal_recovery"`
		MaternalFetalFollowUp            string `json:"maternal_fetal_follow_up"`
		TransplantReferralStatus         string `json:"transplant_referral_status"`
		ModalityPlan                     string `json:"modality_plan"`
		ConservativeCareFlag             bool   `json:"conservative_care_flag"`
		TrackingNotes                    string `json:"tracking_notes"`
	}

	var raw map[string]any
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = c.ShouldBindBodyWith(&raw, binding.JSON)

	// Validate and map blood type
	var mappedBloodType string
	if req.BloodType != "" {
		var ok bool
		mappedBloodType, ok = mapBloodType(req.BloodType)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "invalid blood_type",
				"valid_values": []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-", "unknown"},
			})
			return
		}
	} else {
		mappedBloodType = "unknown"
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	userID, _ := uuid.Parse(userIDStr.(string))

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_of_birth format"})
		return
	}

	var primaryDoctorID pgtype.UUID
	if req.PrimaryDoctorID != "" {
		doctorUUID, err := uuid.Parse(req.PrimaryDoctorID)
		if err == nil {
			primaryDoctorID = pgtype.UUID{Bytes: doctorUUID, Valid: true}
		}
	}

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	patient, err := queries.CreatePatient(ctx, sqlc.CreatePatientParams{
		HospitalID:        hospitalID,
		Mrn:               req.MRN,
		NationalID:        pgtype.Text{String: req.NationalID, Valid: req.NationalID != ""},
		FullName:          req.FullName,
		PreferredName:     pgtype.Text{String: req.PreferredName, Valid: req.PreferredName != ""},
		DateOfBirth:       pgtype.Date{Time: dob, Valid: true},
		Sex:               sqlc.SexType(req.Sex),
		BloodType:         sqlc.BloodType(mappedBloodType),
		MaritalStatus:     sqlc.NullMaritalStatus{MaritalStatus: sqlc.MaritalStatus(req.MaritalStatus), Valid: req.MaritalStatus != ""},
		Nationality:       pgtype.Text{String: req.Nationality, Valid: req.Nationality != ""},
		Religion:          pgtype.Text{String: req.Religion, Valid: req.Religion != ""},
		Occupation:        pgtype.Text{String: req.Occupation, Valid: req.Occupation != ""},
		EducationLevel:    pgtype.Text{String: req.EducationLevel, Valid: req.EducationLevel != ""},
		PrimaryLanguage:   pgtype.Text{String: req.PrimaryLanguage, Valid: req.PrimaryLanguage != ""},
		InterpreterNeeded: req.InterpreterNeeded,
		RegistrationDate:  pgtype.Date{Time: time.Now(), Valid: true},
		RegisteredBy:      userID,
		PrimaryDoctorID:   primaryDoctorID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient"})
		return
	}

	if err := upsertPatientClinicalProfile(ctx, tx, hospitalID.String(), patient.ID.String(), userID.String(), raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save patient clinical profile"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

func (h *PatientsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	queries := sqlc.New(h.pool)
	patient, err := queries.GetPatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientsHandler) List(c *gin.Context) {
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	queries := sqlc.New(h.pool)
	patients, err := queries.ListActivePatients(c.Request.Context(), sqlc.ListActivePatientsParams{
		HospitalID: hospitalID,
		Limit:      100,
		Offset:     0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

func (h *PatientsHandler) Search(c *gin.Context) {
	query := c.Query("q")
	searchType := c.Query("type")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' required"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	queries := sqlc.New(h.pool)

	switch searchType {
	case "mrn":
		patient, err := queries.GetPatientByMRN(c.Request.Context(), sqlc.GetPatientByMRNParams{
			HospitalID: hospitalID,
			Mrn:        query,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": []sqlc.Patient{patient}})

	case "national_id":
		patient, err := queries.GetPatientByNationalID(c.Request.Context(), pgtype.Text{String: query, Valid: true})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": []sqlc.Patient{patient}})

	case "name":
		patients, err := queries.SearchPatientsByName(c.Request.Context(), sqlc.SearchPatientsByNameParams{
			HospitalID: hospitalID,
			FullName:   "%" + query + "%",
			Limit:      50,
			Offset:     0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": patients})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search type, use: mrn, national_id, name"})
	}
}

func (h *PatientsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if _, err := uuid.Parse(idStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	userID := c.GetString(middleware.CtxUserID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	rawPatient, err := updatePatientColumns(ctx, tx, hospitalID, idStr, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := upsertPatientClinicalProfile(ctx, tx, hospitalID, idStr, userID, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save patient clinical profile"})
		return
	}

	if rawPatient == nil {
		rawPatient, err = selectPatientJSON(ctx, tx, hospitalID, idStr)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", rawPatient)
}

func (h *PatientsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.SoftDeletePatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete patient"})
		return
	}

	c.Status(http.StatusNoContent)
}

var patientUpdateFields = map[string]string{
	"mrn":                "",
	"national_id":        "",
	"full_name":          "",
	"preferred_name":     "",
	"date_of_birth":      "::date",
	"sex":                "::sex_type",
	"blood_type":         "::blood_type",
	"marital_status":     "::marital_status",
	"nationality":        "",
	"religion":           "",
	"occupation":         "",
	"education_level":    "",
	"primary_language":   "",
	"interpreter_needed": "::boolean",
	"primary_doctor_id":  "::uuid",
	"is_active":          "::boolean",
}

func updatePatientColumns(ctx context.Context, tx pgx.Tx, hospitalID, patientID string, input map[string]any) ([]byte, error) {
	fields := make([]string, 0, len(input))
	for field := range input {
		if _, ok := patientUpdateFields[field]; ok {
			fields = append(fields, field)
		}
	}
	sort.Strings(fields)

	assignments := []string{}
	args := []any{}
	for _, field := range fields {
		value := input[field]
		if value == nil {
			continue
		}
		if text, ok := value.(string); ok && strings.TrimSpace(text) == "" {
			continue
		}
		if field == "blood_type" {
			text, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid blood_type")
			}
			mapped, ok := mapBloodType(text)
			if !ok {
				return nil, fmt.Errorf("invalid blood_type")
			}
			value = mapped
		}
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d%s", field, len(args), patientUpdateFields[field]))
	}

	if len(assignments) == 0 {
		return nil, nil
	}

	args = append(args, patientID, hospitalID)
	query := fmt.Sprintf(`
		WITH updated AS (
			UPDATE patients
			SET %s
			WHERE id = $%d::uuid
			  AND hospital_id = $%d::uuid
			  AND deleted_at IS NULL
			RETURNING *
		)
		SELECT row_to_json(updated) FROM updated
	`, strings.Join(assignments, ", "), len(args)-1, len(args))

	var raw []byte
	if err := tx.QueryRow(ctx, query, args...).Scan(&raw); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("patient not found")
		}
		return nil, fmt.Errorf("failed to update patient")
	}
	return raw, nil
}

func selectPatientJSON(ctx context.Context, tx pgx.Tx, hospitalID, patientID string) ([]byte, error) {
	var raw []byte
	err := tx.QueryRow(ctx, `
		SELECT row_to_json(t)
		FROM (
			SELECT * FROM patients
			WHERE id = $1::uuid
			  AND hospital_id = $2::uuid
			  AND deleted_at IS NULL
		) t
	`, patientID, hospitalID).Scan(&raw)
	return raw, err
}

func upsertPatientClinicalProfile(ctx context.Context, tx pgx.Tx, hospitalID, patientID, userID string, input map[string]any) error {
	entity := clinicalEntities["patient-clinical-profiles"]
	profileInput := map[string]any{
		"patient_id": patientID,
		"created_by": userID,
		"updated_by": userID,
	}

	hasClinicalField := false
	for field := range entity.Fields {
		if field == "patient_id" || field == "created_by" || field == "updated_by" {
			continue
		}
		if value, ok := input[field]; ok {
			profileInput[field] = value
			hasClinicalField = true
		}
	}
	if !hasClinicalField {
		return nil
	}

	columns, placeholders, args, err := buildInsertColumns(entity, profileInput, hospitalID)
	if err != nil {
		return err
	}

	assignments := make([]string, 0, len(columns))
	for _, col := range columns {
		if col == "hospital_id" || col == "patient_id" || col == "created_by" {
			continue
		}
		assignments = append(assignments, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}
	assignments = append(assignments, "updated_at = now()")

	query := fmt.Sprintf(`
		INSERT INTO patient_clinical_profiles (%s)
		VALUES (%s)
		ON CONFLICT (patient_id) DO UPDATE SET %s
	`, strings.Join(columns, ", "), strings.Join(placeholders, ", "), strings.Join(assignments, ", "))

	_, err = tx.Exec(ctx, query, args...)
	return err
}
