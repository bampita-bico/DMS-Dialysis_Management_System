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

type ShiftAssignmentsHandler struct {
	pool *pgxpool.Pool
}

func NewShiftAssignmentsHandler(pool *pgxpool.Pool) *ShiftAssignmentsHandler {
	return &ShiftAssignmentsHandler{pool: pool}
}

// Create creates a new shift assignment
// POST /api/v1/shift-assignments
func (h *ShiftAssignmentsHandler) Create(c *gin.Context) {
	var req struct {
		StaffID        string   `json:"staff_id" binding:"required"`
		ShiftDate      string   `json:"shift_date" binding:"required"`
		ShiftType      string   `json:"shift_type" binding:"required"`
		ShiftStartTime string   `json:"shift_start_time" binding:"required"`
		ShiftEndTime   string   `json:"shift_end_time" binding:"required"`
		MachineIDs     []string `json:"machine_ids"`
		Notes          string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff_id"})
		return
	}

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_date format"})
		return
	}

	startTime, err := time.Parse("15:04:05", req.ShiftStartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_start_time format (use HH:MM:SS)"})
		return
	}

	endTime, err := time.Parse("15:04:05", req.ShiftEndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_end_time format (use HH:MM:SS)"})
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

	// Convert machine IDs to JSONB
	var machineIDs []byte
	if len(req.MachineIDs) > 0 {
		machineIDs, _ = json.Marshal(req.MachineIDs)
	}

	// Convert times to pgtype.Time (microseconds since midnight)
	shiftStartTime := pgtype.Time{
		Microseconds: int64(startTime.Hour()*3600+startTime.Minute()*60+startTime.Second()) * 1000000,
		Valid:        true,
	}
	shiftEndTime := pgtype.Time{
		Microseconds: int64(endTime.Hour()*3600+endTime.Minute()*60+endTime.Second()) * 1000000,
		Valid:        true,
	}

	queries := sqlc.New(tx)
	shift, err := queries.CreateShiftAssignment(ctx, sqlc.CreateShiftAssignmentParams{
		HospitalID:     hospitalID,
		StaffID:        staffID,
		ShiftDate:      pgtype.Date{Time: shiftDate, Valid: true},
		ShiftType:      sqlc.ShiftType(req.ShiftType),
		ShiftStartTime: shiftStartTime,
		ShiftEndTime:   shiftEndTime,
		MachineIds:     machineIDs,
		AssignedBy:     pgtype.UUID{Bytes: userID, Valid: true},
		Notes:          pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shift assignment", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, shift)
}

// Get retrieves a specific shift assignment by ID
// GET /api/v1/shift-assignments/:id
func (h *ShiftAssignmentsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift assignment ID"})
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
	shift, err := queries.GetShiftAssignment(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shift assignment not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// ListByDate lists all shift assignments for a specific date
// GET /api/v1/shift-assignments/date/:shift_date
func (h *ShiftAssignmentsHandler) ListByDate(c *gin.Context) {
	shiftDateStr := c.Param("shift_date")
	shiftDate, err := time.Parse("2006-01-02", shiftDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_date format"})
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
	shifts, err := queries.ListShiftsByDate(ctx, sqlc.ListShiftsByDateParams{
		HospitalID: hospitalID,
		ShiftDate:  pgtype.Date{Time: shiftDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shifts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

// ListByStaff lists shift assignments for a staff member within a date range
// GET /api/v1/staff/:staff_id/shifts?start_date=2024-01-01&end_date=2024-01-31
func (h *ShiftAssignmentsHandler) ListByStaff(c *gin.Context) {
	staffIDStr := c.Param("staff_id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff_id"})
		return
	}

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
	shifts, err := queries.ListShiftsByStaff(ctx, sqlc.ListShiftsByStaffParams{
		StaffID:   staffID,
		ShiftDate: pgtype.Date{Time: startDate, Valid: true},
		ShiftDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shifts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

// ListUnconfirmed lists all unconfirmed future shifts
// GET /api/v1/shift-assignments/unconfirmed
func (h *ShiftAssignmentsHandler) ListUnconfirmed(c *gin.Context) {
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
	shifts, err := queries.ListUnconfirmedShifts(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list unconfirmed shifts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

// Confirm confirms a shift assignment
// POST /api/v1/shift-assignments/:id/confirm
func (h *ShiftAssignmentsHandler) Confirm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift assignment ID"})
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
	shift, err := queries.ConfirmShift(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm shift"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// ClockIn records the clock-in time for a shift
// POST /api/v1/shift-assignments/:id/clock-in
func (h *ShiftAssignmentsHandler) ClockIn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift assignment ID"})
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
	shift, err := queries.ClockIn(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clock in"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// ClockOut records the clock-out time for a shift
// POST /api/v1/shift-assignments/:id/clock-out
func (h *ShiftAssignmentsHandler) ClockOut(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift assignment ID"})
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
	shift, err := queries.ClockOut(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clock out"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, shift)
}
