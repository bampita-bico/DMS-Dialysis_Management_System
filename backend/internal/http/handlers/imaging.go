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

type ImagingHandler struct {
	pool *pgxpool.Pool
}

func NewImagingHandler(pool *pgxpool.Pool) *ImagingHandler {
	return &ImagingHandler{pool: pool}
}

func (h *ImagingHandler) CreateOrder(c *gin.Context) {
	var req struct {
		PatientID          uuid.UUID  `json:"patient_id" binding:"required"`
		SessionID          *uuid.UUID `json:"session_id"`
		Modality           string     `json:"modality" binding:"required"`
		BodyPart           string     `json:"body_part" binding:"required"`
		Laterality         *string    `json:"laterality"`
		ClinicalIndication string     `json:"clinical_indication" binding:"required"`
		Priority           string     `json:"priority" binding:"required"`
		Notes              *string    `json:"notes"`
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

	order, err := queries.CreateImagingOrder(ctx, sqlc.CreateImagingOrderParams{
		HospitalID:         hospitalID,
		PatientID:          req.PatientID,
		SessionID:          sessionID,
		OrderedBy:          orderedBy,
		OrderDate:          pgtype.Date{Time: now, Valid: true},
		OrderTime:          pgtype.Time{Microseconds: int64(now.Hour()*3600+now.Minute()*60+now.Second()) * 1000000, Valid: true},
		Modality:           sqlc.ImagingModality(req.Modality),
		BodyPart:           req.BodyPart,
		Laterality:         pgtype.Text{String: *req.Laterality, Valid: req.Laterality != nil},
		ClinicalIndication: req.ClinicalIndication,
		Priority:           sqlc.LabPriority(req.Priority),
		Notes:              pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create imaging order"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *ImagingHandler) GetOrder(c *gin.Context) {
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
	order, err := queries.GetImagingOrder(ctx, orderID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, order)
}

func (h *ImagingHandler) ListOrders(c *gin.Context) {
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
	orders, err := queries.ListPendingImagingOrders(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, orders)
}

func (h *ImagingHandler) AddReport(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		ReportText      string   `json:"report_text" binding:"required"`
		Impression      *string  `json:"impression"`
		Recommendations *string  `json:"recommendations"`
		ImageURLs       []string `json:"image_urls"`
		IsAbnormal      bool     `json:"is_abnormal"`
		IsCritical      bool     `json:"is_critical"`
		Notes           *string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	reportedBy, _ := uuid.Parse(userIDStr.(string))

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
	imageURLsJSON := []byte("[]")
	if len(req.ImageURLs) > 0 {
		imageURLsJSON = []byte(`["` + req.ImageURLs[0] + `"]`)
	}

	result, err := queries.CreateImagingResult(ctx, sqlc.CreateImagingResultParams{
		HospitalID:      hospitalID,
		OrderID:         orderID,
		ReportText:      req.ReportText,
		Impression:      pgtype.Text{String: *req.Impression, Valid: req.Impression != nil},
		Recommendations: pgtype.Text{String: *req.Recommendations, Valid: req.Recommendations != nil},
		ReportedBy:      reportedBy,
		ReportDate:      pgtype.Date{Time: now, Valid: true},
		ReportTime:      pgtype.Time{Microseconds: int64(now.Hour()*3600+now.Minute()*60+now.Second()) * 1000000, Valid: true},
		ImageCount:      pgtype.Int4{Int32: int32(len(req.ImageURLs)), Valid: true},
		ImageUrls:       imageURLsJSON,
		IsAbnormal:      req.IsAbnormal,
		IsCritical:      req.IsCritical,
		Notes:           pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add report"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, result)
}
