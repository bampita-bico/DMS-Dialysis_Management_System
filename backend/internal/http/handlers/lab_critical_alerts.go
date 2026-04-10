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

type LabCriticalAlertsHandler struct {
	pool *pgxpool.Pool
}

func NewLabCriticalAlertsHandler(pool *pgxpool.Pool) *LabCriticalAlertsHandler {
	return &LabCriticalAlertsHandler{pool: pool}
}

// Create creates a new lab critical alert
// POST /api/v1/lab-critical-alerts
func (h *LabCriticalAlertsHandler) Create(c *gin.Context) {
	var req struct {
		ResultID       string `json:"result_id" binding:"required"`
		PatientID      string `json:"patient_id" binding:"required"`
		TestName       string `json:"test_name" binding:"required"`
		CriticalValue  string `json:"critical_value" binding:"required"`
		ReferenceRange string `json:"reference_range"`
		Severity       string `json:"severity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Parse UUIDs
	resultIDParsed, err := uuid.Parse(req.ResultID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid result_id format"})
		return
	}

	patientIDParsed, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id format"})
		return
	}

	// Handle optional fields
	var referenceRange pgtype.Text
	if req.ReferenceRange != "" {
		referenceRange = pgtype.Text{String: req.ReferenceRange, Valid: true}
	}

	// Begin transaction
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Set tenant context
	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)

	params := sqlc.CreateLabCriticalAlertParams{
		HospitalID:     uuid.MustParse(hospitalID),
		ResultID:       resultIDParsed,
		PatientID:      patientIDParsed,
		TestName:       req.TestName,
		CriticalValue:  req.CriticalValue,
		ReferenceRange: referenceRange,
		Severity:       req.Severity,
	}

	alert, err := queries.CreateLabCriticalAlert(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create critical alert"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// Get retrieves a lab critical alert by ID
// GET /api/v1/lab-critical-alerts/:id
func (h *LabCriticalAlertsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alert, err := queries.GetLabCriticalAlert(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "critical alert not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// ListByPatient lists all critical alerts for a specific patient
// GET /api/v1/patients/:patient_id/critical-alerts
func (h *LabCriticalAlertsHandler) ListByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alerts, err := queries.ListLabCriticalAlertsByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list critical alerts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// ListUnacknowledged lists all unacknowledged critical alerts
// GET /api/v1/lab-critical-alerts/unacknowledged
func (h *LabCriticalAlertsHandler) ListUnacknowledged(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alerts, err := queries.ListUnacknowledgedCriticalAlerts(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list unacknowledged alerts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// ListByDateRange lists critical alerts within a date range
// GET /api/v1/lab-critical-alerts?start_date=2024-01-01&end_date=2024-12-31
func (h *LabCriticalAlertsHandler) ListByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date query parameters are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected YYYY-MM-DD"})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alerts, err := queries.ListCriticalAlertsByDateRange(ctx, sqlc.ListCriticalAlertsByDateRangeParams{
		HospitalID: uuid.MustParse(hospitalID),
		AlertedAt:  pgtype.Timestamptz{Time: startDate, Valid: true},
		AlertedAt_2: pgtype.Timestamptz{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list critical alerts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// Acknowledge acknowledges a critical alert
// POST /api/v1/lab-critical-alerts/:id/acknowledge
func (h *LabCriticalAlertsHandler) Acknowledge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID format"})
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by" binding:"required"`
		ActionTaken    string `json:"action_taken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acknowledgedByParsed, err := uuid.Parse(req.AcknowledgedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid acknowledged_by format"})
		return
	}

	var actionTaken pgtype.Text
	if req.ActionTaken != "" {
		actionTaken = pgtype.Text{String: req.ActionTaken, Valid: true}
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alert, err := queries.AcknowledgeCriticalAlert(ctx, sqlc.AcknowledgeCriticalAlertParams{
		ID:             id,
		AcknowledgedBy: pgtype.UUID{Bytes: acknowledgedByParsed, Valid: true},
		ActionTaken:    actionTaken,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to acknowledge critical alert"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// NotifyDoctor notifies the doctor of a critical alert
// POST /api/v1/lab-critical-alerts/:id/notify-doctor
func (h *LabCriticalAlertsHandler) NotifyDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	alert, err := queries.NotifyDoctorOfCriticalAlert(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to notify doctor"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// Delete soft deletes a lab critical alert
// DELETE /api/v1/lab-critical-alerts/:id
func (h *LabCriticalAlertsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	err = queries.DeleteLabCriticalAlert(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete critical alert"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "critical alert deleted successfully"})
}
