package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DialysisSessionsHandler struct {
	pool *pgxpool.Pool
}

func NewDialysisSessionsHandler(pool *pgxpool.Pool) *DialysisSessionsHandler {
	return &DialysisSessionsHandler{pool: pool}
}

// Create creates a new dialysis session
// POST /api/v1/dialysis-sessions
func (h *DialysisSessionsHandler) Create(c *gin.Context) {
	var req struct {
		PatientID              string  `json:"patient_id" binding:"required"`
		ScheduleID             *string `json:"schedule_id"`
		MachineID              string  `json:"machine_id" binding:"required"`
		AccessID               *string `json:"access_id"`
		Modality               string  `json:"modality" binding:"required"`
		Shift                  string  `json:"shift" binding:"required"`
		Status                 string  `json:"status"`
		ScheduledDate          string  `json:"scheduled_date" binding:"required"`
		ScheduledStartTime     string  `json:"scheduled_start_time" binding:"required"`
		PrescribedDurationMins int32   `json:"prescribed_duration_mins" binding:"required"`
		PrimaryNurseID         *string `json:"primary_nurse_id"`
		SupervisingDoctorID    *string `json:"supervising_doctor_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	machineID, err := uuid.Parse(req.MachineID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine_id"})
		return
	}

	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_date format"})
		return
	}

	scheduledStartTime, err := time.Parse("15:04:05", req.ScheduledStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_start_time format (use HH:MM:SS)"})
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
	var scheduleID pgtype.UUID
	if req.ScheduleID != nil {
		schID, err := uuid.Parse(*req.ScheduleID)
		if err == nil {
			scheduleID = pgtype.UUID{Bytes: schID, Valid: true}
		}
	}

	var accessID pgtype.UUID
	if req.AccessID != nil {
		accID, err := uuid.Parse(*req.AccessID)
		if err == nil {
			accessID = pgtype.UUID{Bytes: accID, Valid: true}
		}
	}

	var primaryNurseID pgtype.UUID
	if req.PrimaryNurseID != nil {
		nurseID, err := uuid.Parse(*req.PrimaryNurseID)
		if err == nil {
			primaryNurseID = pgtype.UUID{Bytes: nurseID, Valid: true}
		}
	}

	var supervisingDoctorID pgtype.UUID
	if req.SupervisingDoctorID != nil {
		docID, err := uuid.Parse(*req.SupervisingDoctorID)
		if err == nil {
			supervisingDoctorID = pgtype.UUID{Bytes: docID, Valid: true}
		}
	}

	status := sqlc.SessionStatus("scheduled")
	if req.Status != "" {
		status = sqlc.SessionStatus(req.Status)
	}

	// Convert scheduled start time to pgtype.Time (microseconds since midnight)
	scheduledStartTimePG := pgtype.Time{
		Microseconds: int64(scheduledStartTime.Hour()*3600+scheduledStartTime.Minute()*60+scheduledStartTime.Second()) * 1000000,
		Valid:        true,
	}

	queries := sqlc.New(tx)
	session, err := queries.CreateDialysisSession(ctx, sqlc.CreateDialysisSessionParams{
		HospitalID:             hospitalID,
		PatientID:              patientID,
		ScheduleID:             scheduleID,
		MachineID:              machineID,
		AccessID:               accessID,
		Modality:               sqlc.DialysisModality(req.Modality),
		Shift:                  sqlc.ShiftType(req.Shift),
		Status:                 status,
		ScheduledDate:          pgtype.Date{Time: scheduledDate, Valid: true},
		ScheduledStartTime:     scheduledStartTimePG,
		PrescribedDurationMins: req.PrescribedDurationMins,
		PrimaryNurseID:         primaryNurseID,
		SupervisingDoctorID:    supervisingDoctorID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dialysis session", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// Get retrieves a specific dialysis session by ID
// GET /api/v1/dialysis-sessions/:id
func (h *DialysisSessionsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
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
	session, err := queries.GetDialysisSession(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dialysis session not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// ListByPatient lists dialysis sessions for a patient
// GET /api/v1/patients/:patient_id/dialysis-sessions?limit=20&offset=0
func (h *DialysisSessionsHandler) ListByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
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
	sessions, err := queries.ListSessionsByPatient(ctx, sqlc.ListSessionsByPatientParams{
		PatientID: patientID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list dialysis sessions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// ListByDate lists dialysis sessions for a specific date
// GET /api/v1/dialysis-sessions/date/:scheduled_date
func (h *DialysisSessionsHandler) ListByDate(c *gin.Context) {
	scheduledDateStr := c.Param("scheduled_date")
	scheduledDate, err := time.Parse("2006-01-02", scheduledDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_date format"})
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
	sessions, err := queries.ListSessionsByDate(ctx, sqlc.ListSessionsByDateParams{
		HospitalID:    hospitalID,
		ScheduledDate: pgtype.Date{Time: scheduledDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions by date"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// ListActiveByMachine lists active sessions on a specific machine
// GET /api/v1/machines/:machine_id/active-sessions
func (h *DialysisSessionsHandler) ListActiveByMachine(c *gin.Context) {
	machineIDStr := c.Param("machine_id")
	machineID, err := uuid.Parse(machineIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine_id"})
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
	sessions, err := queries.ListActiveSessionsByMachine(ctx, machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active sessions by machine"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// ListActive lists all active sessions
// GET /api/v1/dialysis-sessions/active
func (h *DialysisSessionsHandler) ListActive(c *gin.Context) {
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
	sessions, err := queries.ListActiveSessions(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active sessions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// Start starts a dialysis session with pre-treatment vitals
// POST /api/v1/dialysis-sessions/:id/start
func (h *DialysisSessionsHandler) Start(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req struct {
		PreWeightKg     float64 `json:"pre_weight_kg" binding:"required"`
		PreBpSystolic   int32   `json:"pre_bp_systolic" binding:"required"`
		PreBpDiastolic  int32   `json:"pre_bp_diastolic" binding:"required"`
		PreHr           int32   `json:"pre_hr" binding:"required"`
		PreTemp         float64 `json:"pre_temp"`
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

	var preWeightKg, preTemp pgtype.Numeric
	preWeightKg.Scan(req.PreWeightKg)
	if req.PreTemp > 0 {
		preTemp.Scan(req.PreTemp)
	}

	queries := sqlc.New(tx)
	session, err := queries.StartSession(ctx, sqlc.StartSessionParams{
		ID:             id,
		PreWeightKg:    preWeightKg,
		PreBpSystolic:  pgtype.Int4{Int32: req.PreBpSystolic, Valid: true},
		PreBpDiastolic: pgtype.Int4{Int32: req.PreBpDiastolic, Valid: true},
		PreHr:          pgtype.Int4{Int32: req.PreHr, Valid: true},
		PreTemp:        preTemp,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// Complete completes a dialysis session with post-treatment vitals
// POST /api/v1/dialysis-sessions/:id/complete
func (h *DialysisSessionsHandler) Complete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req struct {
		ActualDurationMins int32   `json:"actual_duration_mins" binding:"required"`
		PostWeightKg       float64 `json:"post_weight_kg" binding:"required"`
		PostBpSystolic     int32   `json:"post_bp_systolic" binding:"required"`
		PostBpDiastolic    int32   `json:"post_bp_diastolic" binding:"required"`
		PostHr             int32   `json:"post_hr" binding:"required"`
		WasPatientReviewed bool    `json:"was_patient_reviewed"`
		ReviewedBy         *string `json:"reviewed_by"`
		SessionNotes       string  `json:"session_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
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

	var postWeightKg pgtype.Numeric
	postWeightKg.Scan(req.PostWeightKg)

	var reviewedBy pgtype.UUID
	if req.ReviewedBy != nil {
		revID, err := uuid.Parse(*req.ReviewedBy)
		if err == nil {
			reviewedBy = pgtype.UUID{Bytes: revID, Valid: true}
		}
	} else if req.WasPatientReviewed {
		// Default to current user if reviewed but no reviewer specified
		reviewedBy = pgtype.UUID{Bytes: userID, Valid: true}
	}

	queries := sqlc.New(tx)
	session, err := queries.CompleteSession(ctx, sqlc.CompleteSessionParams{
		ID:                 id,
		ActualDurationMins: pgtype.Int4{Int32: req.ActualDurationMins, Valid: true},
		PostWeightKg:       postWeightKg,
		PostBpSystolic:     pgtype.Int4{Int32: req.PostBpSystolic, Valid: true},
		PostBpDiastolic:    pgtype.Int4{Int32: req.PostBpDiastolic, Valid: true},
		PostHr:             pgtype.Int4{Int32: req.PostHr, Valid: true},
		WasPatientReviewed: req.WasPatientReviewed,
		ReviewedBy:         reviewedBy,
		SessionNotes:       pgtype.Text{String: req.SessionNotes, Valid: req.SessionNotes != ""},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// Abort aborts a dialysis session with a reason
// POST /api/v1/dialysis-sessions/:id/abort
func (h *DialysisSessionsHandler) Abort(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req struct {
		AbortedReason string `json:"aborted_reason" binding:"required"`
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

	queries := sqlc.New(tx)
	session, err := queries.AbortSession(ctx, sqlc.AbortSessionParams{
		ID:            id,
		AbortedReason: pgtype.Text{String: req.AbortedReason, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to abort session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// UpdateStatus updates the status of a dialysis session
// PATCH /api/v1/dialysis-sessions/:id/status
func (h *DialysisSessionsHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
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

	queries := sqlc.New(tx)
	session, err := queries.UpdateSessionStatus(ctx, sqlc.UpdateSessionStatusParams{
		ID:     id,
		Status: sqlc.SessionStatus(req.Status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}
