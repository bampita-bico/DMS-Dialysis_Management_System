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

type PaymentsHandler struct {
	pool *pgxpool.Pool
}

func NewPaymentsHandler(pool *pgxpool.Pool) *PaymentsHandler {
	return &PaymentsHandler{pool: pool}
}

// Create creates a new payment record
// POST /api/v1/payments
func (h *PaymentsHandler) Create(c *gin.Context) {
	var req struct {
		InvoiceID         *string `json:"invoice_id"`
		AccountID         string  `json:"account_id" binding:"required"`
		PatientID         string  `json:"patient_id" binding:"required"`
		PaymentDate       string  `json:"payment_date" binding:"required"`
		PaymentTime       string  `json:"payment_time"`
		Amount            float64 `json:"amount" binding:"required"`
		PaymentMethod     string  `json:"payment_method" binding:"required"`
		ReferenceNumber   string  `json:"reference_number"`
		MobileMoneyNumber string  `json:"mobile_money_number"`
		BankName          string  `json:"bank_name"`
		ChequeNumber      string  `json:"cheque_number"`
		CardLastFour      string  `json:"card_last_four"`
		Notes             string  `json:"notes"`
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

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment_date format"})
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
	var invoiceID pgtype.UUID
	if req.InvoiceID != nil {
		invID, err := uuid.Parse(*req.InvoiceID)
		if err == nil {
			invoiceID = pgtype.UUID{Bytes: invID, Valid: true}
		}
	}

	var paymentTime pgtype.Time
	if req.PaymentTime != "" {
		t, err := time.Parse("15:04:05", req.PaymentTime)
		if err == nil {
			paymentTime = pgtype.Time{
				Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000,
				Valid:        true,
			}
		}
	} else {
		now := time.Now()
		paymentTime = pgtype.Time{
			Microseconds: int64(now.Hour()*3600+now.Minute()*60+now.Second()) * 1000000,
			Valid:        true,
		}
	}

	var amount pgtype.Numeric
	amount.Scan(req.Amount)

	queries := sqlc.New(tx)
	payment, err := queries.CreatePayment(ctx, sqlc.CreatePaymentParams{
		HospitalID:        hospitalID,
		InvoiceID:         invoiceID,
		AccountID:         accountID,
		PatientID:         patientID,
		PaymentDate:       pgtype.Date{Time: paymentDate, Valid: true},
		PaymentTime:       paymentTime,
		Amount:            amount,
		PaymentMethod:     sqlc.PaymentMethod(req.PaymentMethod),
		ReferenceNumber:   pgtype.Text{String: req.ReferenceNumber, Valid: req.ReferenceNumber != ""},
		MobileMoneyNumber: pgtype.Text{String: req.MobileMoneyNumber, Valid: req.MobileMoneyNumber != ""},
		BankName:          pgtype.Text{String: req.BankName, Valid: req.BankName != ""},
		ChequeNumber:      pgtype.Text{String: req.ChequeNumber, Valid: req.ChequeNumber != ""},
		CardLastFour:      pgtype.Text{String: req.CardLastFour, Valid: req.CardLastFour != ""},
		ReceivedBy:        userID,
		Notes:             pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// Get retrieves a specific payment by ID
// GET /api/v1/payments/:id
func (h *PaymentsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
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
	payment, err := queries.GetPayment(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ListByInvoice lists all payments for an invoice
// GET /api/v1/invoices/:invoice_id/payments
func (h *PaymentsHandler) ListByInvoice(c *gin.Context) {
	invoiceIDStr := c.Param("invoice_id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice_id"})
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
	payments, err := queries.ListPaymentsByInvoice(ctx, pgtype.UUID{Bytes: invoiceID, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ListByPatient lists all payments for a patient
// GET /api/v1/patients/:patient_id/payments?limit=20&offset=0
func (h *PaymentsHandler) ListByPatient(c *gin.Context) {
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
	payments, err := queries.ListPaymentsByPatient(ctx, sqlc.ListPaymentsByPatientParams{
		PatientID: patientID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ListByDateRange lists payments within a date range
// GET /api/v1/payments?start_date=2024-01-01&end_date=2024-12-31
func (h *PaymentsHandler) ListByDateRange(c *gin.Context) {
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
	payments, err := queries.ListPaymentsByDate(ctx, sqlc.ListPaymentsByDateParams{
		HospitalID:  hospitalID,
		PaymentDate: pgtype.Date{Time: startDate, Valid: true},
		PaymentDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ListByMethod lists payments by payment method
// GET /api/v1/payments/method/:method
func (h *PaymentsHandler) ListByMethod(c *gin.Context) {
	method := c.Param("method")
	if method == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment method is required"})
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
	payments, err := queries.ListPaymentsByMethod(ctx, sqlc.ListPaymentsByMethodParams{
		HospitalID:    hospitalID,
		PaymentMethod: sqlc.PaymentMethod(method),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetTotal gets the total payment amount for a date range
// GET /api/v1/payments/total?start_date=2024-01-01&end_date=2024-12-31
func (h *PaymentsHandler) GetTotal(c *gin.Context) {
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
	total, err := queries.GetPaymentTotal(ctx, sqlc.GetPaymentTotalParams{
		HospitalID:  hospitalID,
		PaymentDate: pgtype.Date{Time: startDate, Valid: true},
		PaymentDate_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment total"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
