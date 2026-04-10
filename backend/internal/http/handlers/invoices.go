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

type InvoicesHandler struct {
	pool *pgxpool.Pool
}

func NewInvoicesHandler(pool *pgxpool.Pool) *InvoicesHandler {
	return &InvoicesHandler{pool: pool}
}

// Create creates a new invoice
// POST /api/v1/invoices
func (h *InvoicesHandler) Create(c *gin.Context) {
	var req struct {
		AccountID      string  `json:"account_id" binding:"required"`
		PatientID      string  `json:"patient_id" binding:"required"`
		SessionID      *string `json:"session_id"`
		InvoiceNumber  string  `json:"invoice_number" binding:"required"`
		InvoiceDate    string  `json:"invoice_date" binding:"required"`
		DueDate        string  `json:"due_date" binding:"required"`
		TotalAmount    float64 `json:"total_amount" binding:"required"`
		DiscountAmount float64 `json:"discount_amount"`
		TaxAmount      float64 `json:"tax_amount"`
		NetAmount      float64 `json:"net_amount" binding:"required"`
		Status         string  `json:"status"`
		Notes          string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	invoiceDate, err := time.Parse("2006-01-02", req.InvoiceDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice_date format"})
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format"})
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
	var sessionID pgtype.UUID
	if req.SessionID != nil {
		sessID, err := uuid.Parse(*req.SessionID)
		if err == nil {
			sessionID = pgtype.UUID{Bytes: sessID, Valid: true}
		}
	}

	status := sqlc.InvoiceStatus("issued")
	if req.Status != "" {
		status = sqlc.InvoiceStatus(req.Status)
	}

	var totalAmount, discountAmount, taxAmount, netAmount pgtype.Numeric
	totalAmount.Scan(req.TotalAmount)
	discountAmount.Scan(req.DiscountAmount)
	taxAmount.Scan(req.TaxAmount)
	netAmount.Scan(req.NetAmount)

	queries := sqlc.New(tx)
	invoice, err := queries.CreateInvoice(ctx, sqlc.CreateInvoiceParams{
		HospitalID:     hospitalID,
		AccountID:      accountID,
		PatientID:      patientID,
		SessionID:      sessionID,
		InvoiceNumber:  req.InvoiceNumber,
		InvoiceDate:    pgtype.Date{Time: invoiceDate, Valid: true},
		DueDate:        pgtype.Date{Time: dueDate, Valid: true},
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		TaxAmount:      taxAmount,
		NetAmount:      netAmount,
		Status:         status,
		IssuedBy:       pgtype.UUID{Bytes: userID, Valid: true},
		Notes:          pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// Get retrieves a specific invoice by ID
// GET /api/v1/invoices/:id
func (h *InvoicesHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice ID"})
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
	invoice, err := queries.GetInvoice(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// GetByNumber retrieves an invoice by its invoice number
// GET /api/v1/invoices/number/:invoice_number
func (h *InvoicesHandler) GetByNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoice_number")
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
	invoice, err := queries.GetInvoiceByNumber(ctx, sqlc.GetInvoiceByNumberParams{
		HospitalID:    hospitalID,
		InvoiceNumber: invoiceNumber,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// ListByPatient lists all invoices for a patient
// GET /api/v1/patients/:patient_id/invoices?limit=20&offset=0
func (h *InvoicesHandler) ListByPatient(c *gin.Context) {
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
	invoices, err := queries.ListInvoicesByPatient(ctx, sqlc.ListInvoicesByPatientParams{
		PatientID: patientID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invoices"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// ListByAccount lists all invoices for a billing account
// GET /api/v1/billing-accounts/:account_id/invoices
func (h *InvoicesHandler) ListByAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
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
	invoices, err := queries.ListInvoicesByAccount(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invoices"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// ListByStatus lists invoices by status
// GET /api/v1/invoices?status=overdue
func (h *InvoicesHandler) ListByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
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
	invoices, err := queries.ListInvoicesByStatus(ctx, sqlc.ListInvoicesByStatusParams{
		HospitalID: hospitalID,
		Status:     sqlc.InvoiceStatus(status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invoices"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// ListOverdue lists overdue invoices
// GET /api/v1/invoices/overdue
func (h *InvoicesHandler) ListOverdue(c *gin.Context) {
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
	invoices, err := queries.ListOverdueInvoices(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list overdue invoices"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// UpdateStatus updates the status of an invoice
// PATCH /api/v1/invoices/:id/status
func (h *InvoicesHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice ID"})
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
	invoice, err := queries.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
		ID:     id,
		Status: sqlc.InvoiceStatus(req.Status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update invoice status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// UpdatePayment updates the paid amount for an invoice
// PATCH /api/v1/invoices/:id/payment
func (h *InvoicesHandler) UpdatePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice ID"})
		return
	}

	var req struct {
		PaidAmount float64 `json:"paid_amount" binding:"required"`
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

	var paidAmount pgtype.Numeric
	paidAmount.Scan(req.PaidAmount)

	queries := sqlc.New(tx)
	invoice, err := queries.UpdateInvoicePayment(ctx, sqlc.UpdateInvoicePaymentParams{
		ID:         id,
		PaidAmount: paidAmount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update invoice payment"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}
