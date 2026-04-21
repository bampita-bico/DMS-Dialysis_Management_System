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

type HospitalizationsHandler struct {
	pool *pgxpool.Pool
}

func NewHospitalizationsHandler(pool *pgxpool.Pool) *HospitalizationsHandler {
	return &HospitalizationsHandler{pool: pool}
}

// Create creates a new hospitalization record
// POST /api/v1/hospitalizations
func (h *HospitalizationsHandler) Create(c *gin.Context) {
	var req struct {
		PatientID            string  `json:"patient_id" binding:"required"`
		AdmissionDate        string  `json:"admission_date" binding:"required"`
		AdmissionTime        string  `json:"admission_time"`
		DischargeDate        string  `json:"discharge_date"`
		DischargeTime        string  `json:"discharge_time"`
		LengthOfStayDays     *int32  `json:"length_of_stay_days"`
		AdmissionReason      string  `json:"admission_reason" binding:"required"`
		AdmissionDiagnosis   string  `json:"admission_diagnosis"`
		Icd10Codes           string  `json:"icd10_codes"`
		AdmittingFacility    string  `json:"admitting_facility"`
		WardName             string  `json:"ward_name"`
		DialysisRelated      bool    `json:"dialysis_related"`
		AccessRelated        bool    `json:"access_related"`
		InfectionRelated     bool    `json:"infection_related"`
		TreatmentGiven       string  `json:"treatment_given"`
		ProceduresPerformed  string  `json:"procedures_performed"`
		Outcome              string  `json:"outcome"`
		DischargeDestination string  `json:"discharge_destination"`
		FollowUpRequired     bool    `json:"follow_up_required"`
		FollowUpDate         string  `json:"follow_up_date"`
		Notes                string  `json:"notes"`
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

	admissionDate, err := time.Parse("2006-01-02", req.AdmissionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admission_date format"})
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

	// Prepare parameters
	var admissionTime pgtype.Time
	if req.AdmissionTime != "" {
		t, err := time.Parse("15:04:05", req.AdmissionTime)
		if err == nil {
			admissionTime = pgtype.Time{
				Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000,
				Valid:        true,
			}
		}
	}

	var dischargeDate pgtype.Date
	if req.DischargeDate != "" {
		dDate, err := time.Parse("2006-01-02", req.DischargeDate)
		if err == nil {
			dischargeDate = pgtype.Date{Time: dDate, Valid: true}
		}
	}

	var dischargeTime pgtype.Time
	if req.DischargeTime != "" {
		t, err := time.Parse("15:04:05", req.DischargeTime)
		if err == nil {
			dischargeTime = pgtype.Time{
				Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000,
				Valid:        true,
			}
		}
	}

	var lengthOfStayDays pgtype.Int4
	if req.LengthOfStayDays != nil {
		lengthOfStayDays = pgtype.Int4{Int32: *req.LengthOfStayDays, Valid: true}
	}

	var outcome sqlc.NullHospitalizationOutcome
	if req.Outcome != "" {
		outcome = sqlc.NullHospitalizationOutcome{
			HospitalizationOutcome: sqlc.HospitalizationOutcome(req.Outcome),
			Valid:                  true,
		}
	}

	var followUpDate pgtype.Date
	if req.FollowUpDate != "" {
		fDate, err := time.Parse("2006-01-02", req.FollowUpDate)
		if err == nil {
			followUpDate = pgtype.Date{Time: fDate, Valid: true}
		}
	}

	queries := sqlc.New(tx)
	hospitalization, err := queries.CreateHospitalization(ctx, sqlc.CreateHospitalizationParams{
		HospitalID:           hospitalID,
		PatientID:            patientID,
		AdmissionDate:        pgtype.Date{Time: admissionDate, Valid: true},
		AdmissionTime:        admissionTime,
		DischargeDate:        dischargeDate,
		DischargeTime:        dischargeTime,
		LengthOfStayDays:     lengthOfStayDays,
		AdmissionReason:      req.AdmissionReason,
		AdmissionDiagnosis:   pgtype.Text{String: req.AdmissionDiagnosis, Valid: req.AdmissionDiagnosis != ""},
		Icd10Codes:           pgtype.Text{String: req.Icd10Codes, Valid: req.Icd10Codes != ""},
		AdmittingFacility:    pgtype.Text{String: req.AdmittingFacility, Valid: req.AdmittingFacility != ""},
		WardName:             pgtype.Text{String: req.WardName, Valid: req.WardName != ""},
		DialysisRelated:      req.DialysisRelated,
		AccessRelated:        req.AccessRelated,
		InfectionRelated:     req.InfectionRelated,
		TreatmentGiven:       pgtype.Text{String: req.TreatmentGiven, Valid: req.TreatmentGiven != ""},
		ProceduresPerformed:  pgtype.Text{String: req.ProceduresPerformed, Valid: req.ProceduresPerformed != ""},
		Outcome:              outcome,
		DischargeDestination: pgtype.Text{String: req.DischargeDestination, Valid: req.DischargeDestination != ""},
		FollowUpRequired:     req.FollowUpRequired,
		FollowUpDate:         followUpDate,
		RecordedBy:           userID,
		Notes:                pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hospitalization"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, hospitalization)
}

// Get retrieves a specific hospitalization by ID
// GET /api/v1/hospitalizations/:id
func (h *HospitalizationsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalization ID"})
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
	hospitalization, err := queries.GetHospitalization(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hospitalization not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalization)
}

// ListByPatient lists all hospitalizations for a patient
// GET /api/v1/patients/:patient_id/hospitalizations
func (h *HospitalizationsHandler) ListByPatient(c *gin.Context) {
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
	hospitalizations, err := queries.ListHospitalizationsByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list hospitalizations"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalizations)
}

// ListByPeriod lists hospitalizations within a date range
// GET /api/v1/hospitalizations?start_date=2024-01-01&end_date=2024-12-31
func (h *HospitalizationsHandler) ListByPeriod(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
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
	hospitalizations, err := queries.ListHospitalizationsByPeriod(ctx, sqlc.ListHospitalizationsByPeriodParams{
		HospitalID:    hospitalID,
		AdmissionDate: pgtype.Date{Time: startDate, Valid: true},
		AdmissionDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list hospitalizations"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalizations)
}

// ListDialysisRelated lists dialysis-related hospitalizations
// GET /api/v1/hospitalizations/dialysis-related
func (h *HospitalizationsHandler) ListDialysisRelated(c *gin.Context) {
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
	hospitalizations, err := queries.ListDialysisRelatedHospitalizations(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list dialysis-related hospitalizations"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalizations)
}

// ListAccessRelated lists vascular access-related hospitalizations
// GET /api/v1/hospitalizations/access-related
func (h *HospitalizationsHandler) ListAccessRelated(c *gin.Context) {
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
	hospitalizations, err := queries.ListAccessRelatedHospitalizations(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list access-related hospitalizations"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalizations)
}

// UpdateDischarge updates discharge information for a hospitalization
// PATCH /api/v1/hospitalizations/:id/discharge
func (h *HospitalizationsHandler) UpdateDischarge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalization ID"})
		return
	}

	var req struct {
		DischargeDate        string `json:"discharge_date" binding:"required"`
		DischargeTime        string `json:"discharge_time"`
		LengthOfStayDays     int32  `json:"length_of_stay_days" binding:"required"`
		Outcome              string `json:"outcome"`
		DischargeDestination string `json:"discharge_destination"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dischargeDate, err := time.Parse("2006-01-02", req.DischargeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discharge_date format"})
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

	var dischargeTime pgtype.Time
	if req.DischargeTime != "" {
		t, err := time.Parse("15:04:05", req.DischargeTime)
		if err == nil {
			dischargeTime = pgtype.Time{
				Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000,
				Valid:        true,
			}
		}
	}

	var outcome sqlc.NullHospitalizationOutcome
	if req.Outcome != "" {
		outcome = sqlc.NullHospitalizationOutcome{
			HospitalizationOutcome: sqlc.HospitalizationOutcome(req.Outcome),
			Valid:                  true,
		}
	}

	queries := sqlc.New(tx)
	hospitalization, err := queries.UpdateHospitalizationDischarge(ctx, sqlc.UpdateHospitalizationDischargeParams{
		ID:                   id,
		DischargeDate:        pgtype.Date{Time: dischargeDate, Valid: true},
		DischargeTime:        dischargeTime,
		LengthOfStayDays:     pgtype.Int4{Int32: req.LengthOfStayDays, Valid: true},
		Outcome:              outcome,
		DischargeDestination: pgtype.Text{String: req.DischargeDestination, Valid: req.DischargeDestination != ""},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update discharge info"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, hospitalization)
}
