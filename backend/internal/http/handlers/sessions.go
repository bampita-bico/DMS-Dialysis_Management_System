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

type SessionsHandler struct {
	pool *pgxpool.Pool
}

func NewSessionsHandler(pool *pgxpool.Pool) *SessionsHandler {
	return &SessionsHandler{pool: pool}
}

func (h *SessionsHandler) CreateSession(c *gin.Context) {
	var req struct {
		PatientID              uuid.UUID  `json:"patient_id" binding:"required"`
		ScheduleID             *uuid.UUID `json:"schedule_id"`
		MachineID              uuid.UUID  `json:"machine_id" binding:"required"`
		AccessID               *uuid.UUID `json:"access_id"`
		Modality               string     `json:"modality" binding:"required"`
		Shift                  string     `json:"shift" binding:"required"`
		ScheduledDate          string     `json:"scheduled_date" binding:"required"`
		ScheduledStartTime     string     `json:"scheduled_start_time" binding:"required"`
		PrescribedDurationMins int32      `json:"prescribed_duration_mins" binding:"required"`
		PrimaryNurseID         *uuid.UUID `json:"primary_nurse_id"`
		SupervisingDoctorID    *uuid.UUID `json:"supervising_doctor_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)

	scheduledDate, _ := time.Parse("2006-01-02", req.ScheduledDate)
	scheduledTime, _ := time.Parse("15:04", req.ScheduledStartTime)

	var scheduleID, accessID, primaryNurseID, supervisingDoctorID pgtype.UUID
	if req.ScheduleID != nil {
		scheduleID = pgtype.UUID{Bytes: *req.ScheduleID, Valid: true}
	}
	if req.AccessID != nil {
		accessID = pgtype.UUID{Bytes: *req.AccessID, Valid: true}
	}
	if req.PrimaryNurseID != nil {
		primaryNurseID = pgtype.UUID{Bytes: *req.PrimaryNurseID, Valid: true}
	}
	if req.SupervisingDoctorID != nil {
		supervisingDoctorID = pgtype.UUID{Bytes: *req.SupervisingDoctorID, Valid: true}
	}

	session, err := queries.CreateDialysisSession(ctx, sqlc.CreateDialysisSessionParams{
		HospitalID:             hospitalID,
		PatientID:              req.PatientID,
		ScheduleID:             scheduleID,
		MachineID:              req.MachineID,
		AccessID:               accessID,
		Modality:               sqlc.DialysisModality(req.Modality),
		Shift:                  sqlc.ShiftType(req.Shift),
		Status:                 sqlc.SessionStatusScheduled,
		ScheduledDate:          pgtype.Date{Time: scheduledDate, Valid: true},
		ScheduledStartTime:     pgtype.Time{Microseconds: int64(scheduledTime.Hour()*3600+scheduledTime.Minute()*60) * 1000000, Valid: true},
		PrescribedDurationMins: req.PrescribedDurationMins,
		PrimaryNurseID:         primaryNurseID,
		SupervisingDoctorID:    supervisingDoctorID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *SessionsHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	session, err := queries.GetDialysisSession(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, session)
}

func (h *SessionsHandler) ListSessionsByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	sessions, err := queries.ListSessionsByDate(ctx, sqlc.ListSessionsByDateParams{
		HospitalID:    hospitalID,
		ScheduledDate: pgtype.Date{Time: date, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list sessions"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, sessions)
}

func (h *SessionsHandler) StartSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req struct {
		PreWeightKg     *float64 `json:"pre_weight_kg"`
		PreBPSystolic   *int32   `json:"pre_bp_systolic"`
		PreBPDiastolic  *int32   `json:"pre_bp_diastolic"`
		PreHR           *int32   `json:"pre_hr"`
		PreTemp         *float64 `json:"pre_temp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	session, err := queries.StartSession(ctx, sqlc.StartSessionParams{
		ID:             sessionID,
		PreWeightKg:    pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.PreWeightKg != nil},
		PreBpSystolic:  pgtype.Int4{Int32: *req.PreBPSystolic, Valid: req.PreBPSystolic != nil},
		PreBpDiastolic: pgtype.Int4{Int32: *req.PreBPDiastolic, Valid: req.PreBPDiastolic != nil},
		PreHr:          pgtype.Int4{Int32: *req.PreHR, Valid: req.PreHR != nil},
		PreTemp:        pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.PreTemp != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionsHandler) CompleteSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req struct {
		ActualDurationMins  int32     `json:"actual_duration_mins" binding:"required"`
		PostWeightKg        *float64  `json:"post_weight_kg"`
		PostBPSystolic      *int32    `json:"post_bp_systolic"`
		PostBPDiastolic     *int32    `json:"post_bp_diastolic"`
		PostHR              *int32    `json:"post_hr"`
		WasPatientReviewed  bool      `json:"was_patient_reviewed"`
		ReviewedBy          *uuid.UUID `json:"reviewed_by"`
		SessionNotes        *string   `json:"session_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	var reviewedBy pgtype.UUID
	if req.ReviewedBy != nil {
		reviewedBy = pgtype.UUID{Bytes: *req.ReviewedBy, Valid: true}
	}

	session, err := queries.CompleteSession(ctx, sqlc.CompleteSessionParams{
		ID:                 sessionID,
		ActualDurationMins: pgtype.Int4{Int32: req.ActualDurationMins, Valid: true},
		PostWeightKg:       pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.PostWeightKg != nil},
		PostBpSystolic:     pgtype.Int4{Int32: *req.PostBPSystolic, Valid: req.PostBPSystolic != nil},
		PostBpDiastolic:    pgtype.Int4{Int32: *req.PostBPDiastolic, Valid: req.PostBPDiastolic != nil},
		PostHr:             pgtype.Int4{Int32: *req.PostHR, Valid: req.PostHR != nil},
		WasPatientReviewed: req.WasPatientReviewed,
		ReviewedBy:         reviewedBy,
		SessionNotes:       pgtype.Text{String: *req.SessionNotes, Valid: req.SessionNotes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionsHandler) AbortSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req struct {
		AbortedReason string `json:"aborted_reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	session, err := queries.AbortSession(ctx, sqlc.AbortSessionParams{
		ID:            sessionID,
		AbortedReason: pgtype.Text{String: req.AbortedReason, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to abort session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, session)
}
