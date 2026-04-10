package handlers

import (
	"encoding/json"
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

type SessionComplicationsHandler struct {
	pool *pgxpool.Pool
}

func NewSessionComplicationsHandler(pool *pgxpool.Pool) *SessionComplicationsHandler {
	return &SessionComplicationsHandler{pool: pool}
}

// Create creates a new session complication record
// POST /api/v1/session-complications
func (h *SessionComplicationsHandler) Create(c *gin.Context) {
	var req struct {
		SessionID               string                      `json:"session_id" binding:"required"`
		PatientID               string                      `json:"patient_id" binding:"required"`
		ReportedBy              string                      `json:"reported_by" binding:"required"`
		OccurredAt              string                      `json:"occurred_at" binding:"required"` // ISO8601 timestamp
		ComplicationType        string                      `json:"complication_type" binding:"required"`
		Severity                sqlc.ComplicationSeverity   `json:"severity" binding:"required"`
		Symptoms                string                      `json:"symptoms" binding:"required"`
		VitalSignsAtEvent       map[string]interface{}      `json:"vital_signs_at_event"`
		ImmediateActionTaken    string                      `json:"immediate_action_taken"`
		Outcome                 string                      `json:"outcome"`
		RequiredHospitalization bool                        `json:"required_hospitalization"`
		WasSessionTerminated    bool                        `json:"was_session_terminated"`
		DoctorNotified          bool                        `json:"doctor_notified"`
		DoctorID                string                      `json:"doctor_id"`
		FamilyNotified          bool                        `json:"family_notified"`
		Notes                   string                      `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Parse UUIDs
	sessionIDParsed, err := uuid.Parse(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id format"})
		return
	}

	patientIDParsed, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id format"})
		return
	}

	reportedByParsed, err := uuid.Parse(req.ReportedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reported_by format"})
		return
	}

	// Parse occurred_at timestamp
	occurredAt, err := time.Parse(time.RFC3339, req.OccurredAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid occurred_at format, expected ISO8601"})
		return
	}

	// Marshal vital signs to JSONB
	vitalSignsJSON, err := json.Marshal(req.VitalSignsAtEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vital_signs_at_event format"})
		return
	}

	// Handle optional fields
	var immediateActionTaken pgtype.Text
	if req.ImmediateActionTaken != "" {
		immediateActionTaken = pgtype.Text{String: req.ImmediateActionTaken, Valid: true}
	}

	var outcome pgtype.Text
	if req.Outcome != "" {
		outcome = pgtype.Text{String: req.Outcome, Valid: true}
	}

	var doctorID pgtype.UUID
	if req.DoctorID != "" {
		doctorIDParsed, err := uuid.Parse(req.DoctorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id format"})
			return
		}
		doctorID = pgtype.UUID{Bytes: doctorIDParsed, Valid: true}
	}

	var notes pgtype.Text
	if req.Notes != "" {
		notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	// Begin transaction
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Set tenant context
	hospitalIDStr := hospitalID
	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)

	params := sqlc.CreateSessionComplicationParams{
		HospitalID:              uuid.MustParse(hospitalID),
		SessionID:               sessionIDParsed,
		PatientID:               patientIDParsed,
		ReportedBy:              reportedByParsed,
		OccurredAt:              pgtype.Timestamptz{Time: occurredAt, Valid: true},
		ComplicationType:        req.ComplicationType,
		Severity:                req.Severity,
		Symptoms:                req.Symptoms,
		VitalSignsAtEvent:       vitalSignsJSON,
		ImmediateActionTaken:    immediateActionTaken,
		Outcome:                 outcome,
		RequiredHospitalization: req.RequiredHospitalization,
		WasSessionTerminated:    req.WasSessionTerminated,
		DoctorNotified:          req.DoctorNotified,
		DoctorID:                doctorID,
		FamilyNotified:          req.FamilyNotified,
		Notes:                   notes,
	}

	complication, err := queries.CreateSessionComplication(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create complication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, complication)
}

// Get retrieves a session complication by ID
// GET /api/v1/session-complications/:id
func (h *SessionComplicationsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complication ID format"})
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
	complication, err := queries.GetSessionComplication(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "complication not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, complication)
}

// ListBySession lists all complications for a specific session
// GET /api/v1/sessions/:session_id/complications
func (h *SessionComplicationsHandler) ListBySession(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID format"})
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
	complications, err := queries.ListComplicationsBySession(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list complications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, complications)
}

// ListByPatient lists all complications for a specific patient
// GET /api/v1/patients/:patient_id/complications
func (h *SessionComplicationsHandler) ListByPatient(c *gin.Context) {
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
	complications, err := queries.ListComplicationsByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list complications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, complications)
}

// ListSevere lists all severe/life-threatening complications
// GET /api/v1/session-complications/severe
func (h *SessionComplicationsHandler) ListSevere(c *gin.Context) {
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
	complications, err := queries.ListSevereComplications(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list severe complications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, complications)
}

// Update updates a session complication record
// PATCH /api/v1/session-complications/:id
func (h *SessionComplicationsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complication ID format"})
		return
	}

	var req struct {
		ComplicationType        string                      `json:"complication_type" binding:"required"`
		Severity                sqlc.ComplicationSeverity   `json:"severity" binding:"required"`
		Symptoms                string                      `json:"symptoms" binding:"required"`
		VitalSignsAtEvent       map[string]interface{}      `json:"vital_signs_at_event"`
		ImmediateActionTaken    string                      `json:"immediate_action_taken"`
		Outcome                 string                      `json:"outcome"`
		RequiredHospitalization bool                        `json:"required_hospitalization"`
		WasSessionTerminated    bool                        `json:"was_session_terminated"`
		DoctorNotified          bool                        `json:"doctor_notified"`
		DoctorID                string                      `json:"doctor_id"`
		FamilyNotified          bool                        `json:"family_notified"`
		Notes                   string                      `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Marshal vital signs to JSONB
	vitalSignsJSON, err := json.Marshal(req.VitalSignsAtEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vital_signs_at_event format"})
		return
	}

	// Handle optional fields
	var immediateActionTaken pgtype.Text
	if req.ImmediateActionTaken != "" {
		immediateActionTaken = pgtype.Text{String: req.ImmediateActionTaken, Valid: true}
	}

	var outcome pgtype.Text
	if req.Outcome != "" {
		outcome = pgtype.Text{String: req.Outcome, Valid: true}
	}

	var doctorID pgtype.UUID
	if req.DoctorID != "" {
		doctorIDParsed, err := uuid.Parse(req.DoctorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id format"})
			return
		}
		doctorID = pgtype.UUID{Bytes: doctorIDParsed, Valid: true}
	}

	var notes pgtype.Text
	if req.Notes != "" {
		notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	// Begin transaction
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

	params := sqlc.UpdateSessionComplicationParams{
		ID:                      id,
		ComplicationType:        req.ComplicationType,
		Severity:                req.Severity,
		Symptoms:                req.Symptoms,
		VitalSignsAtEvent:       vitalSignsJSON,
		ImmediateActionTaken:    immediateActionTaken,
		Outcome:                 outcome,
		RequiredHospitalization: req.RequiredHospitalization,
		WasSessionTerminated:    req.WasSessionTerminated,
		DoctorNotified:          req.DoctorNotified,
		DoctorID:                doctorID,
		FamilyNotified:          req.FamilyNotified,
		Notes:                   notes,
	}

	complication, err := queries.UpdateSessionComplication(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update complication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, complication)
}

// Delete soft deletes a session complication record
// DELETE /api/v1/session-complications/:id
func (h *SessionComplicationsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complication ID format"})
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
	err = queries.DeleteSessionComplication(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete complication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "complication deleted successfully"})
}
