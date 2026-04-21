package handlers

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type LabResultsHandler struct {
	pool *pgxpool.Pool
}

func NewLabResultsHandler(pool *pgxpool.Pool) *LabResultsHandler {
	return &LabResultsHandler{pool: pool}
}

// Create creates a new lab result
// POST /api/v1/lab-results
func (h *LabResultsHandler) Create(c *gin.Context) {
	var req struct {
		OrderItemID    string              `json:"order_item_id" binding:"required"`
		ValueText      string              `json:"value_text"`
		ValueNumeric   *float64            `json:"value_numeric"`
		Unit           string              `json:"unit"`
		ReferenceRange string              `json:"reference_range"`
		IsAbnormal     bool                `json:"is_abnormal"`
		IsCritical     bool                `json:"is_critical"`
		Status         sqlc.ResultStatus   `json:"status" binding:"required"`
		ResultDate     string              `json:"result_date" binding:"required"` // YYYY-MM-DD
		ResultTime     string              `json:"result_time" binding:"required"` // HH:MM:SS
		EnteredBy      string              `json:"entered_by" binding:"required"`
		Notes          string              `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Parse UUIDs
	orderItemIDParsed, err := uuid.Parse(req.OrderItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_item_id format"})
		return
	}

	enteredByParsed, err := uuid.Parse(req.EnteredBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entered_by format"})
		return
	}

	// Parse result date
	resultDate, err := time.Parse("2006-01-02", req.ResultDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid result_date format, expected YYYY-MM-DD"})
		return
	}

	// Parse result time
	resultTime, err := time.Parse("15:04:05", req.ResultTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid result_time format, expected HH:MM:SS"})
		return
	}

	// Convert time to microseconds since midnight
	resultTimeMicros := int64(resultTime.Hour()*3600+resultTime.Minute()*60+resultTime.Second()) * 1000000

	// Handle optional fields
	var valueText pgtype.Text
	if req.ValueText != "" {
		valueText = pgtype.Text{String: req.ValueText, Valid: true}
	}

	var valueNumeric pgtype.Numeric
	if req.ValueNumeric != nil {
		dec := decimal.NewFromFloat(*req.ValueNumeric)
		valueNumeric = pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	var unit pgtype.Text
	if req.Unit != "" {
		unit = pgtype.Text{String: req.Unit, Valid: true}
	}

	var referenceRange pgtype.Text
	if req.ReferenceRange != "" {
		referenceRange = pgtype.Text{String: req.ReferenceRange, Valid: true}
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
	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)

	params := sqlc.CreateLabResultParams{
		HospitalID:     uuid.MustParse(hospitalID),
		OrderItemID:    orderItemIDParsed,
		ValueText:      valueText,
		ValueNumeric:   valueNumeric,
		Unit:           unit,
		ReferenceRange: referenceRange,
		IsAbnormal:     req.IsAbnormal,
		IsCritical:     req.IsCritical,
		Status:         req.Status,
		ResultDate:     pgtype.Date{Time: resultDate, Valid: true},
		ResultTime:     pgtype.Time{Microseconds: resultTimeMicros, Valid: true},
		EnteredBy:      enteredByParsed,
		Notes:          notes,
	}

	labResult, err := queries.CreateLabResult(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lab result"})
		return
	}

	// Auto-generate critical alert if numeric value is outside critical thresholds
	var criticalAlert *sqlc.LabCriticalAlert
	if req.ValueNumeric != nil {
		criticalAlert = h.checkAndCreateCriticalAlert(ctx, queries, labResult, uuid.MustParse(hospitalID))
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// Return the result with alert info if one was created
	if criticalAlert != nil {
		c.JSON(http.StatusCreated, gin.H{
			"lab_result":     labResult,
			"critical_alert": criticalAlert,
		})
		return
	}
	c.JSON(http.StatusCreated, labResult)
}

// checkAndCreateCriticalAlert compares a numeric lab result against reference range
// critical thresholds and auto-generates a critical alert if violated.
func (h *LabResultsHandler) checkAndCreateCriticalAlert(
	ctx context.Context, queries *sqlc.Queries,
	result sqlc.LabResult, hospitalID uuid.UUID,
) *sqlc.LabCriticalAlert {
	// Get the order item to find the test_id and order_id
	orderItem, err := queries.GetLabOrderItem(ctx, result.OrderItemID)
	if err != nil {
		return nil
	}

	// Get the order to find the patient_id
	order, err := queries.GetLabOrder(ctx, orderItem.OrderID)
	if err != nil {
		return nil
	}

	// Get the test catalog entry for the name
	test, err := queries.GetLabTest(ctx, orderItem.TestID)
	if err != nil {
		return nil
	}

	// Get the default reference range for this test
	refRange, err := queries.GetDefaultReferenceRange(ctx, orderItem.TestID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil // no reference range defined — can't evaluate
		}
		return nil
	}

	// Compare the numeric value against critical thresholds
	if !result.ValueNumeric.Valid {
		return nil
	}

	resultVal := numericToDecimal(result.ValueNumeric)
	if resultVal == nil {
		return nil
	}

	var severity string
	isCritical := false

	if refRange.CriticalLow.Valid {
		critLow := numericToDecimal(refRange.CriticalLow)
		if critLow != nil && resultVal.LessThan(*critLow) {
			severity = "critical_low"
			isCritical = true
		}
	}

	if !isCritical && refRange.CriticalHigh.Valid {
		critHigh := numericToDecimal(refRange.CriticalHigh)
		if critHigh != nil && resultVal.GreaterThan(*critHigh) {
			severity = "critical_high"
			isCritical = true
		}
	}

	if !isCritical {
		return nil
	}

	// Build reference range string for the alert
	refRangeStr := pgtype.Text{}
	if refRange.ReferenceText.Valid {
		refRangeStr = refRange.ReferenceText
	}

	alert, err := queries.CreateLabCriticalAlert(ctx, sqlc.CreateLabCriticalAlertParams{
		HospitalID:     hospitalID,
		ResultID:       result.ID,
		PatientID:      order.PatientID,
		TestName:       test.Name,
		CriticalValue:  fmt.Sprintf("%s", resultVal.String()),
		ReferenceRange: refRangeStr,
		Severity:       severity,
	})
	if err != nil {
		return nil
	}

	return &alert
}

// numericToDecimal converts a pgtype.Numeric to a shopspring decimal.
func numericToDecimal(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid || n.NaN {
		return nil
	}
	bigInt := n.Int
	if bigInt == nil {
		bigInt = big.NewInt(0)
	}
	d := decimal.NewFromBigInt(bigInt, n.Exp)
	return &d
}

// Get retrieves a lab result by ID
// GET /api/v1/lab-results/:id
func (h *LabResultsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lab result ID format"})
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
	labResult, err := queries.GetLabResult(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lab result not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResult)
}

// GetByOrderItem retrieves lab result for a specific order item
// GET /api/v1/lab-order-items/:order_item_id/result
func (h *LabResultsHandler) GetByOrderItem(c *gin.Context) {
	orderItemIDStr := c.Param("order_item_id")
	orderItemID, err := uuid.Parse(orderItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order item ID format"})
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
	labResult, err := queries.GetLabResultByOrderItem(ctx, orderItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lab result not found for this order item"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResult)
}

// ListByOrder lists all lab results for a specific order
// GET /api/v1/lab-orders/:order_id/results
func (h *LabResultsHandler) ListByOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
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
	labResults, err := queries.ListLabResultsByOrder(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list lab results"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResults)
}

// ListPendingVerification lists all lab results pending verification
// GET /api/v1/lab-results/pending-verification
func (h *LabResultsHandler) ListPendingVerification(c *gin.Context) {
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
	labResults, err := queries.ListPendingVerificationResults(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pending verification results"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResults)
}

// ListCritical lists all critical lab results
// GET /api/v1/lab-results/critical
func (h *LabResultsHandler) ListCritical(c *gin.Context) {
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
	labResults, err := queries.ListCriticalResults(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list critical results"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResults)
}

// Verify verifies a lab result
// POST /api/v1/lab-results/:id/verify
func (h *LabResultsHandler) Verify(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lab result ID format"})
		return
	}

	var req struct {
		VerifiedBy string `json:"verified_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verifiedByParsed, err := uuid.Parse(req.VerifiedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verified_by format"})
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
	labResult, err := queries.VerifyLabResult(ctx, sqlc.VerifyLabResultParams{
		ID:         id,
		VerifiedBy: pgtype.UUID{Bytes: verifiedByParsed, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify lab result"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResult)
}

// Update updates a lab result
// PATCH /api/v1/lab-results/:id
func (h *LabResultsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lab result ID format"})
		return
	}

	var req struct {
		ValueText      string            `json:"value_text"`
		ValueNumeric   *float64          `json:"value_numeric"`
		Unit           string            `json:"unit"`
		ReferenceRange string            `json:"reference_range"`
		IsAbnormal     bool              `json:"is_abnormal"`
		IsCritical     bool              `json:"is_critical"`
		Status         sqlc.ResultStatus `json:"status" binding:"required"`
		Notes          string            `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Handle optional fields
	var valueText pgtype.Text
	if req.ValueText != "" {
		valueText = pgtype.Text{String: req.ValueText, Valid: true}
	}

	var valueNumeric pgtype.Numeric
	if req.ValueNumeric != nil {
		dec := decimal.NewFromFloat(*req.ValueNumeric)
		valueNumeric = pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	var unit pgtype.Text
	if req.Unit != "" {
		unit = pgtype.Text{String: req.Unit, Valid: true}
	}

	var referenceRange pgtype.Text
	if req.ReferenceRange != "" {
		referenceRange = pgtype.Text{String: req.ReferenceRange, Valid: true}
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

	params := sqlc.UpdateLabResultParams{
		ID:             id,
		ValueText:      valueText,
		ValueNumeric:   valueNumeric,
		Unit:           unit,
		ReferenceRange: referenceRange,
		IsAbnormal:     req.IsAbnormal,
		IsCritical:     req.IsCritical,
		Status:         req.Status,
		Notes:          notes,
	}

	labResult, err := queries.UpdateLabResult(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lab result"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, labResult)
}

// Delete soft deletes a lab result
// DELETE /api/v1/lab-results/:id
func (h *LabResultsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lab result ID format"})
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
	err = queries.DeleteLabResult(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete lab result"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lab result deleted successfully"})
}
