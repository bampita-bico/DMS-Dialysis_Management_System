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

type ClinicalOutcomesHandler struct {
	pool *pgxpool.Pool
}

func NewClinicalOutcomesHandler(pool *pgxpool.Pool) *ClinicalOutcomesHandler {
	return &ClinicalOutcomesHandler{pool: pool}
}

// Create creates a new clinical outcome record
// POST /api/v1/clinical-outcomes
func (h *ClinicalOutcomesHandler) Create(c *gin.Context) {
	var req struct {
		PatientID        string   `json:"patient_id" binding:"required"`
		AssessmentDate   string   `json:"assessment_date" binding:"required"`
		PeriodStart      string   `json:"period_start" binding:"required"`
		PeriodEnd        string   `json:"period_end" binding:"required"`
		HemoglobinAvg    *float64 `json:"hemoglobin_avg"`
		KtVAvg           *float64 `json:"ktv_avg"`
		UrrAvg           *float64 `json:"urr_avg"`
		WeightChange     *float64 `json:"weight_change"`
		HospitalDays     *int32   `json:"hospital_days"`
		InfectionEpisodes *int32   `json:"infection_episodes"`
		AccessSurvival   *string  `json:"access_survival"`
		Notes            string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)
	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	assessmentDate, err := time.Parse("2006-01-02", req.AssessmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assessment_date format"})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start format"})
		return
	}

	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end format"})
		return
	}

	// Validate period_end > period_start
	if !periodEnd.After(periodStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_end must be after period_start"})
		return
	}

	// Validate assessment_date >= period_start
	if assessmentDate.Before(periodStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assessment_date must be >= period_start"})
		return
	}

	// Start transaction with RLS
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

	// Prepare parameters - match actual CreateClinicalOutcomeParams structure
	var hemoglobin, ktv, urr pgtype.Numeric

	if req.HemoglobinAvg != nil {
		hemoglobin = pgtype.Numeric{Valid: true}
		hemoglobin.Scan(*req.HemoglobinAvg)
	}

	if req.KtVAvg != nil {
		ktv = pgtype.Numeric{Valid: true}
		ktv.Scan(*req.KtVAvg)
	}

	if req.UrrAvg != nil {
		urr = pgtype.Numeric{Valid: true}
		urr.Scan(*req.UrrAvg)
	}

	// Wrap userID in pgtype.UUID
	assessedBy := pgtype.UUID{Bytes: userID, Valid: true}

	queries := sqlc.New(tx)
	outcome, err := queries.CreateClinicalOutcome(ctx, sqlc.CreateClinicalOutcomeParams{
		HospitalID:            hospitalID,
		PatientID:             patientID,
		AssessmentDate:        pgtype.Date{Time: assessmentDate, Valid: true},
		PeriodStart:           pgtype.Date{Time: periodStart, Valid: true},
		PeriodEnd:             pgtype.Date{Time: periodEnd, Valid: true},
		Hemoglobin:            hemoglobin,
		HemoglobinTargetMin:   pgtype.Numeric{}, // Can be added to request if needed
		HemoglobinTargetMax:   pgtype.Numeric{},
		KtV:                   ktv,
		KtVTarget:             pgtype.Numeric{},
		Urr:                   urr,
		SystolicBpAvg:         pgtype.Numeric{},
		DiastolicBpAvg:        pgtype.Numeric{},
		BpControlled:          pgtype.Bool{},
		WeightGainPercent:     pgtype.Numeric{},
		Albumin:               pgtype.Numeric{},
		Phosphate:             pgtype.Numeric{},
		Calcium:               pgtype.Numeric{},
		Pth:                   pgtype.Numeric{},
		QualityOfLifeScore:    pgtype.Numeric{},
		FunctionalStatus:      pgtype.Text{},
		AdverseEventsCount:    0,
		HospitalizationsCount: int32(*req.HospitalDays),
		MissedSessionsCount:   0,
		OutcomeTrend:          sqlc.NullOutcomeTrend{},
		AssessedBy:            assessedBy,
		Notes:                 pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create clinical outcome", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, outcome)
}

// ListByPatient lists all clinical outcomes for a patient
// GET /api/v1/patients/:patient_id/clinical-outcomes
func (h *ClinicalOutcomesHandler) ListByPatient(c *gin.Context) {
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
	outcomes, err := queries.ListOutcomesByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clinical outcomes"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, outcomes)
}

// GetLatest retrieves the latest clinical outcome for a patient
// GET /api/v1/patients/:patient_id/clinical-outcomes/latest
func (h *ClinicalOutcomesHandler) GetLatest(c *gin.Context) {
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
	outcome, err := queries.GetLatestOutcome(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no clinical outcomes found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, outcome)
}

// ListDeclining lists patients with declining clinical indicators
// GET /api/v1/clinical-outcomes/declining
func (h *ClinicalOutcomesHandler) ListDeclining(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	// Default to last 90 days
	daysBack := c.DefaultQuery("days", "90")
	var assessmentDate time.Time
	if daysBack == "90" {
		assessmentDate = time.Now().AddDate(0, 0, -90)
	} else {
		// Could parse custom date if needed
		assessmentDate = time.Now().AddDate(0, 0, -90)
	}

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
	declining, err := queries.ListDecliningPatients(ctx, sqlc.ListDecliningPatientsParams{
		HospitalID:     hospitalID,
		AssessmentDate: pgtype.Date{Time: assessmentDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list declining patients"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, declining)
}

// ListByTrend lists outcomes grouped by trend (improving, stable, declining)
// GET /api/v1/clinical-outcomes/by-trend?trend=declining
func (h *ClinicalOutcomesHandler) ListByTrend(c *gin.Context) {
	trend := c.Query("trend")
	if trend == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trend parameter is required (improving, stable, declining, critical)"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

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
	outcomes, err := queries.ListOutcomesByTrend(ctx, sqlc.ListOutcomesByTrendParams{
		HospitalID: hospitalID,
		OutcomeTrend: sqlc.NullOutcomeTrend{
			OutcomeTrend: sqlc.OutcomeTrend(trend),
			Valid:        true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list outcomes by trend"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, outcomes)
}
