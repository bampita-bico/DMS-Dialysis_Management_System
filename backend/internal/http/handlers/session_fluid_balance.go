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
	"github.com/shopspring/decimal"
)

type SessionFluidBalanceHandler struct {
	pool *pgxpool.Pool
}

func NewSessionFluidBalanceHandler(pool *pgxpool.Pool) *SessionFluidBalanceHandler {
	return &SessionFluidBalanceHandler{pool: pool}
}

// Create creates a new session fluid balance record
// POST /api/v1/session-fluid-balance
func (h *SessionFluidBalanceHandler) Create(c *gin.Context) {
	var req struct {
		SessionID          string   `json:"session_id" binding:"required"`
		PatientID          string   `json:"patient_id" binding:"required"`
		RecordedBy         string   `json:"recorded_by" binding:"required"`
		RecordedAt         string   `json:"recorded_at" binding:"required"` // ISO8601 timestamp
		UfGoalMl           *float64 `json:"uf_goal_ml"`
		UfAchievedMl       *float64 `json:"uf_achieved_ml"`
		UfRateMlPerHr      *float64 `json:"uf_rate_ml_per_hr"`
		FluidIntakeOralMl  *float64 `json:"fluid_intake_oral_ml"`
		FluidIntakeIvMl    *float64 `json:"fluid_intake_iv_ml"`
		FluidOutputUrineMl *float64 `json:"fluid_output_urine_ml"`
		FluidOutputOtherMl *float64 `json:"fluid_output_other_ml"`
		NetFluidBalanceMl  *float64 `json:"net_fluid_balance_ml"`
		WeightChangeKg     *float64 `json:"weight_change_kg"`
		Notes              string   `json:"notes"`
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

	recordedByParsed, err := uuid.Parse(req.RecordedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recorded_by format"})
		return
	}

	// Parse recorded_at timestamp
	recordedAt, err := time.Parse(time.RFC3339, req.RecordedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recorded_at format, expected ISO8601"})
		return
	}

	// Convert float64 to pgtype.Numeric
	toNumeric := func(val *float64) pgtype.Numeric {
		if val == nil {
			return pgtype.Numeric{}
		}
		dec := decimal.NewFromFloat(*val)
		return pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
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

	params := sqlc.CreateSessionFluidBalanceParams{
		HospitalID:         uuid.MustParse(hospitalID),
		SessionID:          sessionIDParsed,
		PatientID:          patientIDParsed,
		RecordedBy:         recordedByParsed,
		RecordedAt:         pgtype.Timestamptz{Time: recordedAt, Valid: true},
		UfGoalMl:           toNumeric(req.UfGoalMl),
		UfAchievedMl:       toNumeric(req.UfAchievedMl),
		UfRateMlPerHr:      toNumeric(req.UfRateMlPerHr),
		FluidIntakeOralMl:  toNumeric(req.FluidIntakeOralMl),
		FluidIntakeIvMl:    toNumeric(req.FluidIntakeIvMl),
		FluidOutputUrineMl: toNumeric(req.FluidOutputUrineMl),
		FluidOutputOtherMl: toNumeric(req.FluidOutputOtherMl),
		NetFluidBalanceMl:  toNumeric(req.NetFluidBalanceMl),
		WeightChangeKg:     toNumeric(req.WeightChangeKg),
		Notes:              notes,
	}

	fluidBalance, err := queries.CreateSessionFluidBalance(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create fluid balance record"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, fluidBalance)
}

// Get retrieves a session fluid balance record by ID
// GET /api/v1/session-fluid-balance/:id
func (h *SessionFluidBalanceHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fluid balance ID format"})
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
	fluidBalance, err := queries.GetSessionFluidBalance(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fluid balance record not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, fluidBalance)
}

// GetBySession retrieves fluid balance record for a specific session
// GET /api/v1/sessions/:session_id/fluid-balance
func (h *SessionFluidBalanceHandler) GetBySession(c *gin.Context) {
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
	fluidBalance, err := queries.GetFluidBalanceBySession(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fluid balance record not found for this session"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, fluidBalance)
}

// ListByPatient lists all fluid balance records for a specific patient
// GET /api/v1/patients/:patient_id/fluid-balance
func (h *SessionFluidBalanceHandler) ListByPatient(c *gin.Context) {
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
	fluidBalances, err := queries.ListFluidBalancesByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list fluid balance records"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, fluidBalances)
}

// Update updates a session fluid balance record
// PATCH /api/v1/session-fluid-balance/:id
func (h *SessionFluidBalanceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fluid balance ID format"})
		return
	}

	var req struct {
		UfGoalMl           *float64 `json:"uf_goal_ml"`
		UfAchievedMl       *float64 `json:"uf_achieved_ml"`
		UfRateMlPerHr      *float64 `json:"uf_rate_ml_per_hr"`
		FluidIntakeOralMl  *float64 `json:"fluid_intake_oral_ml"`
		FluidIntakeIvMl    *float64 `json:"fluid_intake_iv_ml"`
		FluidOutputUrineMl *float64 `json:"fluid_output_urine_ml"`
		FluidOutputOtherMl *float64 `json:"fluid_output_other_ml"`
		NetFluidBalanceMl  *float64 `json:"net_fluid_balance_ml"`
		WeightChangeKg     *float64 `json:"weight_change_kg"`
		Notes              string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Convert float64 to pgtype.Numeric
	toNumeric := func(val *float64) pgtype.Numeric {
		if val == nil {
			return pgtype.Numeric{}
		}
		dec := decimal.NewFromFloat(*val)
		return pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
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

	params := sqlc.UpdateSessionFluidBalanceParams{
		ID:                 id,
		UfGoalMl:           toNumeric(req.UfGoalMl),
		UfAchievedMl:       toNumeric(req.UfAchievedMl),
		UfRateMlPerHr:      toNumeric(req.UfRateMlPerHr),
		FluidIntakeOralMl:  toNumeric(req.FluidIntakeOralMl),
		FluidIntakeIvMl:    toNumeric(req.FluidIntakeIvMl),
		FluidOutputUrineMl: toNumeric(req.FluidOutputUrineMl),
		FluidOutputOtherMl: toNumeric(req.FluidOutputOtherMl),
		NetFluidBalanceMl:  toNumeric(req.NetFluidBalanceMl),
		WeightChangeKg:     toNumeric(req.WeightChangeKg),
		Notes:              notes,
	}

	fluidBalance, err := queries.UpdateSessionFluidBalance(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update fluid balance record"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, fluidBalance)
}

// Delete soft deletes a session fluid balance record
// DELETE /api/v1/session-fluid-balance/:id
func (h *SessionFluidBalanceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fluid balance ID format"})
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
	err = queries.DeleteSessionFluidBalance(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete fluid balance record"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "fluid balance record deleted successfully"})
}
