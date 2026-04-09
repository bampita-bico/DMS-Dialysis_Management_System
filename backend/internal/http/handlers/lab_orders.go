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

type LabOrdersHandler struct {
	pool *pgxpool.Pool
}

func NewLabOrdersHandler(pool *pgxpool.Pool) *LabOrdersHandler {
	return &LabOrdersHandler{pool: pool}
}

func (h *LabOrdersHandler) CreateOrder(c *gin.Context) {
	var req struct {
		PatientID      uuid.UUID   `json:"patient_id" binding:"required"`
		SessionID      *uuid.UUID  `json:"session_id"`
		Priority       string      `json:"priority" binding:"required"`
		ClinicalNotes  *string     `json:"clinical_notes"`
		DiagnosisCode  *string     `json:"diagnosis_code"`
		Tests          []uuid.UUID `json:"tests" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	orderedBy, _ := uuid.Parse(userIDStr.(string))

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

	now := time.Now()
	var sessionID pgtype.UUID
	if req.SessionID != nil {
		sessionID = pgtype.UUID{Bytes: *req.SessionID, Valid: true}
	}

	order, err := queries.CreateLabOrder(ctx, sqlc.CreateLabOrderParams{
		HospitalID:    hospitalID,
		PatientID:     req.PatientID,
		SessionID:     sessionID,
		OrderedBy:     orderedBy,
		OrderDate:     pgtype.Date{Time: now, Valid: true},
		OrderTime:     pgtype.Time{Microseconds: int64(now.Hour()*3600+now.Minute()*60+now.Second()) * 1000000, Valid: true},
		Priority:      sqlc.LabPriority(req.Priority),
		ClinicalNotes: pgtype.Text{String: *req.ClinicalNotes, Valid: req.ClinicalNotes != nil},
		DiagnosisCode: pgtype.Text{String: *req.DiagnosisCode, Valid: req.DiagnosisCode != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lab order"})
		return
	}

	// Create order items for each test
	for _, testID := range req.Tests {
		_, err := queries.CreateLabOrderItem(ctx, sqlc.CreateLabOrderItemParams{
			HospitalID:   hospitalID,
			OrderID:      order.ID,
			TestID:       testID,
			SpecimenType: sqlc.SpecimenTypeSerum,
			Notes:        pgtype.Text{},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *LabOrdersHandler) GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
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
	order, err := queries.GetLabOrder(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, order)
}

func (h *LabOrdersHandler) ListPendingOrders(c *gin.Context) {
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
	orders, err := queries.ListPendingLabOrders(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, orders)
}

func (h *LabOrdersHandler) CollectSpecimen(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		SpecimenBarcode string `json:"specimen_barcode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	collectedBy, _ := uuid.Parse(userIDStr.(string))

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
	item, err := queries.CollectSpecimen(ctx, sqlc.CollectSpecimenParams{
		ID:                  itemID,
		SpecimenCollectedBy: pgtype.UUID{Bytes: collectedBy, Valid: true},
		SpecimenBarcode:     pgtype.Text{String: req.SpecimenBarcode, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to collect specimen"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *LabOrdersHandler) AddResult(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		ValueText      *string  `json:"value_text"`
		ValueNumeric   *float64 `json:"value_numeric"`
		Unit           *string  `json:"unit"`
		ReferenceRange *string  `json:"reference_range"`
		IsAbnormal     bool     `json:"is_abnormal"`
		IsCritical     bool     `json:"is_critical"`
		Notes          *string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	enteredBy, _ := uuid.Parse(userIDStr.(string))

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

	now := time.Now()
	result, err := queries.CreateLabResult(ctx, sqlc.CreateLabResultParams{
		HospitalID:     hospitalID,
		OrderItemID:    itemID,
		ValueText:      pgtype.Text{String: *req.ValueText, Valid: req.ValueText != nil},
		ValueNumeric:   pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.ValueNumeric != nil},
		Unit:           pgtype.Text{String: *req.Unit, Valid: req.Unit != nil},
		ReferenceRange: pgtype.Text{String: *req.ReferenceRange, Valid: req.ReferenceRange != nil},
		IsAbnormal:     req.IsAbnormal,
		IsCritical:     req.IsCritical,
		Status:         sqlc.ResultStatusPreliminary,
		ResultDate:     pgtype.Date{Time: now, Valid: true},
		ResultTime:     pgtype.Time{Microseconds: int64(now.Hour()*3600+now.Minute()*60+now.Second()) * 1000000, Valid: true},
		EnteredBy:      enteredBy,
		Notes:          pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add result"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *LabOrdersHandler) VerifyResult(c *gin.Context) {
	resultID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid result ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	verifiedBy, _ := uuid.Parse(userIDStr.(string))

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
	result, err := queries.VerifyLabResult(ctx, sqlc.VerifyLabResultParams{
		ID:         resultID,
		VerifiedBy: pgtype.UUID{Bytes: verifiedBy, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify result"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *LabOrdersHandler) ListCriticalAlerts(c *gin.Context) {
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
	alerts, err := queries.ListUnacknowledgedCriticalAlerts(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list alerts"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, alerts)
}
