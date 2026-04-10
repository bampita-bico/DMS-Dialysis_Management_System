package handlers

import (
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalHistoryHandler struct {
	pool *pgxpool.Pool
}

func NewMedicalHistoryHandler(pool *pgxpool.Pool) *MedicalHistoryHandler {
	return &MedicalHistoryHandler{pool: pool}
}

// ============================================================================
// DIAGNOSES
// ============================================================================

// CreateDiagnosis creates a new diagnosis for a patient
// POST /api/v1/patients/:patient_id/diagnoses
func (h *MedicalHistoryHandler) CreateDiagnosis(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	var req struct {
		Icd10Code     string  `json:"icd10_code" binding:"required"`
		Description   string  `json:"description" binding:"required"`
		DiagnosisType string  `json:"diagnosis_type" binding:"required"`
		AdmissionID   *string `json:"admission_id"`
		Notes         string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var admissionID pgtype.UUID
	if req.AdmissionID != nil {
		admID, err := uuid.Parse(*req.AdmissionID)
		if err == nil {
			admissionID = pgtype.UUID{Bytes: admID, Valid: true}
		}
	}

	queries := sqlc.New(tx)
	diagnosis, err := queries.CreateDiagnosis(ctx, sqlc.CreateDiagnosisParams{
		HospitalID:    hospitalID,
		PatientID:     patientID,
		Icd10Code:     req.Icd10Code,
		Description:   req.Description,
		DiagnosisType: sqlc.DiagnosisType(req.DiagnosisType),
		DiagnosedBy:   userID,
		AdmissionID:   admissionID,
		Notes:         pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create diagnosis", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, diagnosis)
}

// ListDiagnosesByPatient lists all diagnoses for a patient
// GET /api/v1/patients/:patient_id/diagnoses
func (h *MedicalHistoryHandler) ListDiagnosesByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	diagnoses, err := queries.ListDiagnosesByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list diagnoses"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, diagnoses)
}

// GetPrimaryDiagnosis retrieves the primary diagnosis for a patient
// GET /api/v1/patients/:patient_id/diagnoses/primary
func (h *MedicalHistoryHandler) GetPrimaryDiagnosis(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	diagnosis, err := queries.GetPrimaryDiagnosis(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no primary diagnosis found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, diagnosis)
}

// ============================================================================
// COMORBIDITIES
// ============================================================================

// CreateComorbidity creates a new comorbidity for a patient
// POST /api/v1/patients/:patient_id/comorbidities
func (h *MedicalHistoryHandler) CreateComorbidity(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	var req struct {
		Condition   string  `json:"condition" binding:"required"`
		Icd10Code   *string `json:"icd10_code"`
		Status      string  `json:"status" binding:"required"`
		DiagnosedAt *string `json:"diagnosed_at"`
		DiagnosedBy *string `json:"diagnosed_by"`
		Notes       string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var icd10Code pgtype.Text
	if req.Icd10Code != nil {
		icd10Code = pgtype.Text{String: *req.Icd10Code, Valid: true}
	}

	var diagnosedAt pgtype.Date
	if req.DiagnosedAt != nil {
		diagDate, err := time.Parse("2006-01-02", *req.DiagnosedAt)
		if err == nil {
			diagnosedAt = pgtype.Date{Time: diagDate, Valid: true}
		}
	}

	var diagnosedBy pgtype.UUID
	if req.DiagnosedBy != nil {
		diagBy, err := uuid.Parse(*req.DiagnosedBy)
		if err == nil {
			diagnosedBy = pgtype.UUID{Bytes: diagBy, Valid: true}
		}
	} else {
		diagnosedBy = pgtype.UUID{Bytes: userID, Valid: true}
	}

	queries := sqlc.New(tx)
	comorbidity, err := queries.CreateComorbidity(ctx, sqlc.CreateComorbidityParams{
		HospitalID:  hospitalID,
		PatientID:   patientID,
		Condition:   req.Condition,
		Icd10Code:   icd10Code,
		Status:      sqlc.ComorbidityStatus(req.Status),
		DiagnosedAt: diagnosedAt,
		DiagnosedBy: diagnosedBy,
		Notes:       pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comorbidity", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, comorbidity)
}

// ListComorbiditiesByPatient lists all comorbidities for a patient
// GET /api/v1/patients/:patient_id/comorbidities
func (h *MedicalHistoryHandler) ListComorbiditiesByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	comorbidities, err := queries.ListComorbiditiesByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comorbidities"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, comorbidities)
}

// UpdateComorbidityStatus updates the status of a comorbidity
// PATCH /api/v1/comorbidities/:id/status
func (h *MedicalHistoryHandler) UpdateComorbidityStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comorbidity ID"})
		return
	}

	var req struct {
		Status     string  `json:"status" binding:"required"`
		ResolvedAt *string `json:"resolved_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var resolvedAt pgtype.Date
	if req.ResolvedAt != nil {
		resolvedDate, err := time.Parse("2006-01-02", *req.ResolvedAt)
		if err == nil {
			resolvedAt = pgtype.Date{Time: resolvedDate, Valid: true}
		}
	}

	queries := sqlc.New(tx)
	comorbidity, err := queries.UpdateComorbidityStatus(ctx, sqlc.UpdateComorbidityStatusParams{
		ID:         id,
		Status:     sqlc.ComorbidityStatus(req.Status),
		ResolvedAt: resolvedAt,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update comorbidity status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, comorbidity)
}

// ============================================================================
// ALLERGIES
// ============================================================================

// CreateAllergy creates a new allergy record for a patient
// POST /api/v1/patients/:patient_id/allergies
func (h *MedicalHistoryHandler) CreateAllergy(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	var req struct {
		Allergen string  `json:"allergen" binding:"required"`
		Category string  `json:"category" binding:"required"`
		Reaction string  `json:"reaction" binding:"required"`
		Severity string  `json:"severity" binding:"required"`
		OnsetDate *string `json:"onset_date"`
		Notes    string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var onsetDate pgtype.Date
	if req.OnsetDate != nil {
		onset, err := time.Parse("2006-01-02", *req.OnsetDate)
		if err == nil {
			onsetDate = pgtype.Date{Time: onset, Valid: true}
		}
	}

	queries := sqlc.New(tx)
	allergy, err := queries.CreateAllergy(ctx, sqlc.CreateAllergyParams{
		HospitalID: hospitalID,
		PatientID:  patientID,
		Allergen:   req.Allergen,
		Category:   sqlc.AllergyCategory(req.Category),
		Reaction:   sqlc.AllergyReaction(req.Reaction),
		Severity:   sqlc.SeverityLevel(req.Severity),
		OnsetDate:  onsetDate,
		Notes:      pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
		RecordedBy: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create allergy", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, allergy)
}

// GetActiveAllergies retrieves all active allergies for a patient
// GET /api/v1/patients/:patient_id/allergies
func (h *MedicalHistoryHandler) GetActiveAllergies(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	allergies, err := queries.GetActiveAllergies(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active allergies"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, allergies)
}

// CheckDrugAllergy checks if a patient has a drug allergy
// GET /api/v1/patients/:patient_id/allergies/check?drug=drug_name
func (h *MedicalHistoryHandler) CheckDrugAllergy(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	drugName := c.Query("drug")
	if drugName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drug parameter is required"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	hasAllergy, err := queries.CheckDrugAllergy(ctx, sqlc.CheckDrugAllergyParams{
		PatientID: patientID,
		Allergen:  drugName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check drug allergy"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient_id":   patientID,
		"drug":         drugName,
		"has_allergy":  hasAllergy,
	})
}
